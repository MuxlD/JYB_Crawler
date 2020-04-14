package eduData

import (
	"JYB_Crawler/Basics"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
)

var maxID uint

//生产者函数
func (ts *TsCrawler) GetAllEdu() {
	db := Basics.GetDB()
	db.Exec("TRUNCATE TABLE types;")

	ctx := context.Background()
	chrome := NewChromedp(ctx)

	//便利出所有的类型
	for id, ets := range Basics.EveryType {
		var count int
		_ = db.Table("ts_urls").Select("max(id)").Row().Scan(&maxID)
		start := time.Now()

		maxPa, mul := maxPage(chrome, ets)
		db.Model(&Basics.Type{}).Where("id = ?", id+1).Update("max_page", maxPa)

		ctx, cancel := chrome.NewTab()
		ctx0, _ := context.WithTimeout(ctx, time.Duration(mul*chromedpTimeout)*time.Second)

		fmt.Println("最大页码为", maxPa)

		var oneUrl Basics.TsUrl
		oneUrl.TypeID = ets.ID

		//分页加载
		for n := 1; n <= maxPa; n++ {
			//学校链接爬取,及通道的写入
			var urlHtml string
			nowPageUrl := ets.TypeUrl + "p" + strconv.Itoa(n) + ".html"
			err := chromedp.Run(ctx0,
				chromedp.Navigate(nowPageUrl),
				chromedp.WaitVisible(`.mt10`),
				chromedp.OuterHTML(`.mt10 .office-result-list`, &urlHtml),
			)
			if err != nil {
				log.Println("html提取失败", err)
				break
			}
			fmt.Println("当前链接：", nowPageUrl)
			selectUrl(urlHtml, `href="(.*)" target="_blank" class="office-rlist-name"`, oneUrl, db, &count)
		}
		db.Model(&Basics.Type{}).Where("id = ?", id+1).Update("count", count)
		log.Printf("抓取成功:%v，爬取耗时：%v\n", ets.TypeUrl, time.Since(start))
		cancel()
	}

	//关闭通道，通知所有类目下的商品获取完成
	//close(done)
	log.Println("所有学校url提取完成...")
}

//提取url
func selectUrl(html string, reg string, tst Basics.TsUrl, db *gorm.DB, count *int) {

	result := SelfReg(html, reg)
	for i := range result {
		tst.Url = Basics.JYB + result[i][1]
		*count++
		maxID++
		tst.ID = maxID
		db.Create(&tst)
		//tsCh <- tst
	}
}

//最大页码查询
func maxPage(chrome *ChromeBrowser, ets Basics.Type) (maxPa, mul int) {
	ctx, c := chrome.NewTab()
	defer c()
	ctx, _ = context.WithTimeout(ctx, time.Duration(chromedpTimeout)*time.Second)
	err := chromedp.Run(ctx,
		//页面跳转
		chromedp.Navigate(ets.TypeUrl),
		// 存在类型组，说明成功进入
		chromedp.WaitVisible(`.mt10`),
	)
	if err != nil {
		log.Fatalln("类型首页加载失败...", ets.TypeUrl)
	}
	//页面验证，找出最大页码
	var xPage string //最大页码
	err = chromedp.Run(ctx,
		//获取最大页码
		chromedp.Text(`.mt10 .pagination li:nth-last-child(2) a`, &xPage),
	)
	if err != nil {
		log.Println(err)
	}
	maxPa, _ = strconv.Atoi(xPage)
	if maxPa == 0 {
		maxPa = 1
		//页面加载不成功
		log.Println("该类型的最大页码可能为：", maxPa)
	}
	//mul为multipler,用于设置超时倍数
	mul = (maxPa + 3) / 3
	return
}
