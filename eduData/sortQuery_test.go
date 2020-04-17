package eduData

import (
	"JYB_Crawler/Basics"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"testing"
	"time"
)

var date time.Time

func TestCrawlerByUrl(t *testing.T) {
	Basics.StartMySql()
	log.SetFlags(log.Lshortfile)
	db := Basics.GetDB()
	db.Model(Basics.Type{}).Find(&Basics.EveryType)
	var dbAllTs []Basics.TsUrl
	var total int
	db = db.Model(Basics.TsUrl{}).Count(&total)
	//验证数据库中是否有数据
	if total <= 0 {
		log.Println("table ts_urls is empty")
	}
	db.Find(&dbAllTs)
	chrome := NewChromedp(context.Background())
	defer chrome.Close()
	//parentCtx, cancel := chrome.NewTab()
	//defer cancel()
	ts := new(TsCrawler)
	date = time.Now()
	for _, ats := range dbAllTs {

		err := ts.CrawlerByUrl(ats, chrome)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	fmt.Println("总用时", time.Since(date))
}

func TestTsCrawler_CrawlerByUrl(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	//chrome := NewChromedp(context.Background())
	//defer chrome.Close()
	//goCtx, cancel := chrome.NewTab()
	////goCtx ,cancel := chromedp.NewContext(context.Background())
	//defer cancel()
	url := "https://product.pconline.com.cn/"
	goCtx, cancel := chromedp.NewContext(context.Background())
	ctx, cancel := context.WithTimeout(goCtx, 15*time.Second)
	defer cancel()

	//start := time.Now()
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		//等待机构详细出现
		//chromedp.WaitVisible(`index-agency-intro-container`),
	)
	if err != nil {
		//页面加载不成功
		log.Println("页面加载不成功,页面模板可能存在差距，链接：", url, err)
		return
	}
}
