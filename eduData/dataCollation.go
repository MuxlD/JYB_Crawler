package eduData

import (
	"JYB_Crawler/Basics"
	"JYB_Crawler/elasticsearch"
	"context"
	"log"
)

var (
	chromedpTimeout int //爬取网络超时时间
	done            chan struct{}
	tsCh            chan Basics.TsUrl
	tsCapacity      = 1000 //仓库容量
)

type TsCrawler struct {
	Ctx context.Context
}

func (ts *TsCrawler) NewTsCrawler(ctx context.Context) *TsCrawler {
	return &TsCrawler{ctx}
}

func StartContext(goroutineNum, cralerTimeout int) {
	ts := new(TsCrawler)
	//构建类型切片，设置通用超时时间
	Basics.EveryType = make([]Basics.Type, 20)
	chromedpTimeout = cralerTimeout

	//获取所有类型链接,为变量everyType赋值
	err := ts.FindAllType(Basics.JYBFL)
	if err != nil {
		log.Println("FindAllType error, info：", err)
		return
	}
	err = elasticsearch.TpBulkInsert()
	if err != nil {
		log.Println("TpBulkInsert error,info:", err)
		return
	}

	//ts.Do(goroutineNum)
}

func (ts *TsCrawler) Do(chromeNum int) {
	//记录初始ts
	//设置仓库容量
	tsCh = make(chan Basics.TsUrl, tsCapacity)
	//生产者信号通道
	done = make(chan struct{})

	//生产者函数，对变量everyTypeInfo的值进行处理
	ts.GetAllEdu()

	//多消费者
	//for i := 0; i < chromeNum; i++ {
	//go ts.Crawler(strconv.Itoa(0))
	//time.Sleep(time.Second * 10)
	//}
}

//主要爬虫程序
func (ts *TsCrawler) Crawler(chromeId string) {
	//初始化浏览器
	chrome := NewChromedp(context.Background())
	ts.NewTsCrawler(chrome.allocCtx)
	var stop bool
	var ok bool
	var tsCraw Basics.TsUrl
	for {
		select {
		//当chrome.Close()被执行时才会进入该case
		case <-ts.Ctx.Done():
			log.Printf("收到退出信号，" + chromeId + "号协程执行退出...\n")
			//关闭浏览器
			return
		case <-done:
			//生产者结束信号
			stop = true
		case tsCraw, ok = <-tsCh:
			if !ok && stop {
				log.Printf("通道内容已消费完," + chromeId + "号协程退出...\n")
				//关闭浏览器
				chrome.Close()
				return
			}
			//消费函数
			CrawlerByUrl(tsCraw, chrome)
			continue
		}
	}
}
