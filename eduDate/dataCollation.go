package eduDate

import (
	"chromedp_test/Basics"
	"context"
	"fmt"
	"log"
)

var (
	BrowserMap      map[string]*ChromeBrowser
	chromedpTimeout int //爬取网络超时时间
	done            chan struct{}
	tsCh            chan TsCrawler
	tsCapacity      = 1000 //仓库容量
)


//用于类型及学校url信息的爬取
type TsCrawler struct {
	Type    string
	TypeId  int
	TypeUrl string
	Url     string
	Ctx     context.Context
}

func NewTsCrawler(ctx context.Context) *TsCrawler {
	return &TsCrawler{
		Ctx: ctx,
	}
}

//ctx 为 parent context (context.Background)
func StartContext(firstCtx context.Context, goroutineNum, cralerTimeout int) {

	BrowserMap = make(map[string]*ChromeBrowser, goroutineNum+2)
	everyType = make([]TsCrawler, 20)
	chromedpTimeout = cralerTimeout

	//建立用于类型爬取的context
	ctx, cancel := NewChromeCtx(firstCtx, "-2")
	log.Println("typeCrawler used ", ctx) //InitChromedp's context.
	//获取所有类型链接,为变量everyType赋值
	typeCrawler := NewTsCrawler(ctx)
	err := typeCrawler.FindAllType(Basics.JYBFL)
	//类型爬取结束取消该ctx
	cancel()
	log.Println(ctx, " cancelled...")
	if err != nil {
		log.Println("FindAllType error, info：", err)
		return
	}

	log.Println("TsCrawler's parent context is", firstCtx) //context.Background.WithCancel
	tsCrawler := NewTsCrawler(firstCtx)
	tsCrawler.Do(goroutineNum)
}

func (ts *TsCrawler) Do(chromeNum int) {

	//设置仓库容量
	tsCh = make(chan TsCrawler, tsCapacity)
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
	var tsCraw TsCrawler
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
