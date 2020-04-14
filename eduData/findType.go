package eduDate

import (
	"chromedp_test/Basics"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"time"
)

var everyType []Basics.Type

func (ts *TsCrawler) FindAllType(url string) (err error) {

	start := time.Now()

	ctx1, cancel1 := context.WithTimeout(ts.Ctx, time.Duration(chromedpTimeout)*time.Second)
	defer cancel1()

	err = chromedp.Run(ctx1,
		chromedp.Navigate(url),
		// 存在类型组，说明成功进入
		chromedp.WaitVisible(`.class-groups`),
	)
	if err != nil {
		//没有进入全部分类页面
		log.Println("链接：", url)
		return
	}

	for i := 0; i < 20; i++ {
		everyType[i].ID = uint(i)
		err = chromedp.Run(ctx1, eduType(i))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//fmt.Println(everyTypeInfo)
	log.Printf("类型抓取成功,链接：%v，爬取耗时：%v\n", url, time.Since(start))
	return
}

//将抓取到的类型链接放入切片,并跳转到该类型下
func eduType(i int) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.JavascriptAttribute(TypeSel(i+1), "href", &everyType[i].TypeUrl),
		chromedp.Text(TypeSel(i+1), &everyType[i].Name),
	}
}

//返回抓取类型链接的sel
func TypeSel(i int, ) (sel string) {
	if i <= 10 {
		sel = `.class-groups .class-groups-left div:nth-child(` + strconv.Itoa(i) + `) a`
	} else {
		sel = `.class-groups .class-groups-right div:nth-child(` + strconv.Itoa(i-10) + `) a`
	}
	return
}
