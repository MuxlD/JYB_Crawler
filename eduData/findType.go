package eduData

import (
	"JYB_Crawler.Vn/Basics"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"time"
)

func (ts *TsCrawler) FindAllType(url string) (err error) {

	db := Basics.GetDB()

	db.Model(Basics.Type{}).Find(&Basics.EveryType)
	if len(Basics.EveryType) < 20 {
		//清除表格数据
		db.Exec("TRUNCATE TABLE types;")
		//清空切片
		Basics.EveryType = nil
		Basics.EveryType = make([]Basics.Type, 20)
	} else if Basics.EveryType[0].TypeName != "" && Basics.EveryType[19].TypeName != "" {
		log.Println("everyType form database...")
		return nil
	}

	start := time.Now()
	//建立用于类型爬取的context
	chrome := NewChromedp(context.Background())
	defer chrome.Close()
	ctx0, cancel0 := chrome.NewTab()
	defer cancel0()
	ctx, cancel := context.WithTimeout(ctx0, time.Duration(chromedpTimeout)*time.Second)
	defer cancel()

	err = chromedp.Run(ctx,
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
		Basics.EveryType[i].ID = uint(i + 1)
		err = chromedp.Run(ctx, eduType(i))
		if err != nil {
			fmt.Println(err)
			return
		}
		db.Create(&Basics.EveryType[i])
	}
	//fmt.Println(everyTypeInfo)
	log.Printf("类型抓取成功,链接：%v，爬取耗时：%v\n", url, time.Since(start))
	return
}

//将抓取到的类型链接放入切片,并跳转到该类型下
func eduType(i int) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.JavascriptAttribute(TypeSel(i+1), "href", &Basics.EveryType[i].TypeUrl),
		chromedp.Text(TypeSel(i+1), &Basics.EveryType[i].TypeName),
	}
}

//返回抓取类型链接的sel
func TypeSel(i int) (sel string) {
	if i <= 10 {
		sel = `.class-groups .class-groups-left div:nth-child(` + strconv.Itoa(i) + `) a`
	} else {
		sel = `.class-groups .class-groups-right div:nth-child(` + strconv.Itoa(i-10) + `) a`
	}
	return
}
