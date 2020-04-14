
package main

import (
	"chromedp_test/Basics"
	"chromedp_test/eduData"
)

func main() {

	Basics.StartMySql()
	//开始爬虫工作
	eduData.InitTsCrawler().StartContext(5,15)
}

