package eduData

import (
	"chromedp_test/Basics"
	"context"
	"fmt"
	"log"
)

var (
	chromedpTimeout int //爬取网络超时时间
	done            chan struct{}
	tsCh            chan Basics.TsUrl
	tsCapacity      = 1000 //仓库容量
)

//用于记录当前程序流信息
type TsCrawler struct {
	FlowTo     string
	CurrentCtx context.Context
}

func InitTsCrawler() *TsCrawler {
	return &TsCrawler{
		"begin",
		nil,
	}
}

func (ts *TsCrawler) CurrentFlow(flowname string, ctx context.Context) *TsCrawler {
	ts.FlowTo = ts.FlowTo + "->" + flowname
	ts.CurrentCtx = ctx
	return ts
}

func (ts *TsCrawler) StartContext(goroutineNum, cralerTimeout int) {

	//构建类型切片，设置通用超时时间
	everyType = make([]Basics.Type, 20)
	chromedpTimeout = cralerTimeout

	//建立用于类型爬取的context
	chrome := NewChromedp(context.Background())
	ctx, _ := chrome.NewTab()
	//获取所有类型链接,为变量everyType赋值
	err := ts.CurrentFlow("FindAllType", ctx).FindAllType(Basics.JYBFL)
	//类型爬取结束取消该ctx
	chrome.Close()
	if err != nil {
		log.Println("FindAllType error, info：", err)
		return
	}

	ts.Do(goroutineNum)
}

func (ts *TsCrawler) Do(chromeNum int) {

	//设置仓库容量
	tsCh = make(chan Basics.TsUrl, tsCapacity)
	//生产者信号通道
	done = make(chan struct{})

	//生产者函数，对变量everyTypeInfo的值进行处理
	ts.GetAllEdu()
	fmt.Println(len(tsCh))
	//多消费者
	//for i := 0; i < chromeNum; i++ {
	//go ts.Crawler(strconv.Itoa(0))
	//time.Sleep(time.Second * 10)
	//}
}

//主要爬虫程序
func (ts *TsCrawler) Crawler(chromeId string) {
	//初始化浏览器
	chrome := NewChromedp(ts.Ctx)
	var stop bool
	var ok bool
	var tsCraw Basics.TsUrl
	for {
		select {
		case <-ts.Ctx.Done():
			log.Printf("收到退出信号，" + chromeId + "号协程执行退出...\n")
			//关闭浏览器
			chrome.Close()
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
