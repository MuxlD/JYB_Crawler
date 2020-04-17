package eduData

import (
	"JYB_Crawler/Basics"
	"JYB_Crawler/elasticsearch"
	"context"
	"log"
	"strconv"
)

var (
	chromedpTimeout int //爬取网络超时时间
	done            chan struct{}
	tsCh            chan Basics.TsUrl
	tsCapacity      = 1000   //仓库容量
	PendLink        []string //正常爬取出错的待处理链接
)

type TsCrawler struct{}

func StartContext(goroutineNum, cralerTimeout int) {
	ts := new(TsCrawler)
	//构建类型切片，设置通用超时时间
	chromedpTimeout = cralerTimeout

	//获取所有类型链接,为变量everyType赋值
	err := ts.FindAllType(Basics.JYBFL)
	if err != nil {
		log.Println("FindAllType error, info：", err)
		return
	}

	ts.Do(goroutineNum)
}

func (ts *TsCrawler) Do(goroutineNum int) {
	//记录初始ts
	//设置仓库容量
	tsCh = make(chan Basics.TsUrl, tsCapacity)
	//生产者信号通道
	done = make(chan struct{})

	//生产者函数，对变量everyType的值进行处理,完善type对象空白字段
	go ts.AllLink()
	//建立批量插入任务
	indexCtx, indexCancel := context.WithCancel(context.Background())
	//多消费者，es生产者
	for i := 1; i <= goroutineNum; i++ {
		go ts.Crawler(strconv.Itoa(i), indexCtx, indexCancel)
		//time.Sleep(time.Second * 10)
	}
	err := elasticsearch.BulkInsert(indexCtx)
	if err != nil {
		log.Println("BulkInsert error, info:", err)
		//出错时，关闭批量插入es的任务,主goroutine结束
		indexCancel()
	}
}

//主要爬虫程序
func (ts *TsCrawler) Crawler(chromeId string, indexCtx context.Context, indexCancel context.CancelFunc) {
	log.Println("chrome ID", chromeId, ":Successfully entered goroutine...")
	//新的浏览器
	chrome := NewChromedp(context.Background())
	var stop bool
	var ok bool
	var tsCraw Basics.TsUrl
	for {
		select {
		//当chrome.Close()(即:chrome.cancel())被执行时,程序才会进入该case//异常退出
		case <-chrome.allocCtx.Done():
			log.Printf("收到退出信号，" + chromeId + "号协程执行退出...\n")
			return
		//indexCancel(),当es插入异常时关闭
		case <-indexCtx.Done():
			log.Printf("bulk insert or CrawlerByUrl error...\n")
			//关闭浏览器
			chrome.Close()
			return
		case <-done:
			//生产者结束信号
			stop = true
		case tsCraw, ok = <-tsCh:
			if !ok && stop {
				log.Printf("通道内容已消费完," + chromeId + "号协程退出...\n")
				//所有内容写入完成关闭es写入通道
				close(elasticsearch.Docsc)
				//关闭浏览器
				chrome.Close()
				return
			}
			//消费函数
			reUrl := ts.CrawlerByUrl(tsCraw, chrome)
			if reUrl != "" {
				PendLink = append(PendLink, reUrl)
			}
			continue
		}
	}
}
