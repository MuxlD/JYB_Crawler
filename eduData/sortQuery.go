package eduData

import (
	"JYB_Crawler/Basics"
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"time"
)

func (ts *TsCrawler) CrawlerByUrl(tsCraw Basics.TsUrl, chr *ChromeBrowser) (err error) {
	log.Println("当前学校链接：", tsCraw.Url)

	start := time.Now()

	var bts Basics.TrainingSchool
	var retry bool

	//进入具体的信息爬取
	for i := 0; i < 3; i++ {
		bts, retry = ts.EveryEdu(chr, tsCraw.Url)
		if !retry { //retry == false 时跳出循环
			break
		}
	}
	//成功赋值字段name,course,brightSpot,info,campus,phoneNumber
	if bts.Name == "" {
		log.Println("获取学校信息失败", tsCraw.Url)
		err = errors.New("多次尝试无效")
		return
	}

	bts.TypeID = tsCraw.TypeID
	bts.TypeUrl = Basics.EveryType[tsCraw.TypeID-1].TypeUrl
	bts.TypeName = Basics.EveryType[tsCraw.TypeID-1].TypeName
	bts.Url = tsCraw.Url
	//成功赋值字段name,course,brightSpot,info,campus,phoneNumber,type_id,type_url,type_name,url
	fmt.Println(bts)

	//elasticsearch.Docsc <- bts

	log.Printf("抓取成功,爬取耗时：%v\n", time.Since(start))

	return
}

func (ts *TsCrawler) EveryEdu(chr *ChromeBrowser, url string) (bts Basics.TrainingSchool, retry bool) {
	// retry == true 时触发循环
	goCtx, cancel := chr.NewTab()
	defer cancel()
	ctx, cancel := context.WithTimeout(goCtx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		//等待链接页面
		chromedp.WaitVisible(`.page-container`),
	)
	if err != nil {
		//页面加载不成功
		log.Println("验证页面第一条class失败", err)
		retry = true
		return
	}

	err = chromedp.Run(ctx,
		//验证是否存在机构详细
		chromedp.WaitVisible(`.index-agency-intro-container`),
	)
	if err != nil {
		log.Println("等待机构详细显示失败", err)
		retry = true
		return
	}

	var courseHtml string
	var bsHtml string
	//写入基础信息
	err = chromedp.Run(ctx, dpCrawl(&bts, &courseHtml, &bsHtml))
	if err != nil {
		log.Println("获取机构信息失败...")
		retry = true
		return
	}

	if courseHtml != "" {
		bts.Course = Splice(courseHtml, `title="(.*)">`)
	} else {
		log.Println(bts.Name, "获取到的课程为空...")
	}

	if bsHtml != "" {
		bts.BrightSpot = Splice(bsHtml, `.png" alt="">(.*)</span>`)
	} else {
		log.Println(bts.Name, "获取到的特色为空...")
	}

	//成功赋值字段name,course,brightSpot,info,campus,phoneNumber
	return
}

//直接爬取
func dpCrawl(bts *Basics.TrainingSchool, courseHtml, bsHtml *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.JavascriptAttribute(`.header-agency-outer .header-agency-logo img`, "title", &bts.Name),

		//chromedp.OuterHTML(`.header-agency-outer`, &ts.Course),
		chromedp.OuterHTML(`.agency-nav .agency-nav-toggle-list`, courseHtml),
		chromedp.OuterHTML(`.index-agency-intro-right`, bsHtml),

		chromedp.Text(`.index-agency-intro-jj p`, &bts.Info),
		chromedp.Text(`.index-agency-intro-xq p`, &bts.Campus),
		chromedp.Text(`.index-agency-intro-tel p`, &bts.PhoneNumber),
	}
}

//字符串拼接
func Splice(bodystr string, rxp string) (c []string) {
	result := SelfReg(bodystr, rxp)
	for i := range result {
		c = append(c, result[i][1])
	}
	return
}

//传入body与正则公式，返回筛选结果
func SelfReg(bodystr string, rxp string) [][]string {
	reg := regexp.MustCompile(rxp)
	result := reg.FindAllStringSubmatch(bodystr, -1)
	return result
}
