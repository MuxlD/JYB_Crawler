package eduData

import (
	"JYB_Crawler/Basics"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
)

func CrawlerByUrl(tsCraw Basics.TsUrl, chr *ChromeBrowser) {

	goCtx, cancel := chr.NewTab()
	defer cancel()

	var bts Basics.TrainingSchool
	log.Println("当前学校链接：", tsCraw.Url)

	//进入具体的信息爬取
	var err error
	bts, err = EveryEdu(goCtx, tsCraw.Url)
	if err != nil {
		log.Println("信息爬取失败，失败信息：", tsCraw.Url, err)
	}
	//bts.TypeId = tsCraw.TypeID
	bts.Url = tsCraw.Url
	fmt.Println(bts)

}

func EveryEdu(ctx context.Context, url string) (bts Basics.TrainingSchool, err error) {

	fmt.Println("entry:", url)
	//start := time.Now()
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		//等待机构详细出现
		chromedp.WaitVisible(`index-agency-intro-container`),
	)
	if err != nil {
		//页面加载不成功
		log.Println( "页面加载不成功,页面模板可能存在差距，链接：", url,err)
		return
	}

	var courseHtml string
	var bsHtml string
	//写入基础信息

	err = chromedp.Run(ctx, dpCrawl(&bts, &courseHtml, &bsHtml))
	if err != nil {
		log.Println("信息爬取失败.")
		return
	}

	bts.Course = Splice(courseHtml, `title="(.*)">`)
	bts.BrightSpot = Splice(bsHtml, `.png" alt="">(.*)</span>`)

	//log.Printf("抓取成功,爬取耗时：%v\n", time.Since(start))
	log.Println("exit:",url)
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
