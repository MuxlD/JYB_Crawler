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
	tsCapacity      = 1000 //仓库容量
)

type TsCrawler ChromeBrowser

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
	go ts.GetAllEdu()

	//多消费者
	//for i := 1; i < =chromeNum; i++ {
	go ts.Crawler(strconv.Itoa(1))
	//time.Sleep(time.Second * 10)
	//}
	go elasticsearch.BulkInsert()
}

//主要爬虫程序
func (ts *TsCrawler) Crawler(chromeId string) {
	//初始化浏览器
	chrome := NewChromedp(context.Background())
	ts.allocCtx = chrome.allocCtx
	ts.cancel = chrome.cancel
	var stop bool
	var ok bool
	var tsCraw Basics.TsUrl
	for {
		select {
		//当chrome.Close()被执行时,程序才会进入该case
		case <-ts.allocCtx.Done():
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
			ts.CrawlerByUrl(tsCraw, chrome)
			continue
		}
	}
}
