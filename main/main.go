package main

import (
	"JYB_Crawler.Vn/Basics"
	"JYB_Crawler.Vn/eduData"
	"JYB_Crawler.Vn/elasticsearch"
	//_ "JYB_Crawler.Vn/elasticsearch"
)

func main() {

	Basics.StartMySql()
	elasticsearch.InitMapping()
	//开始爬虫工作
	eduData.StartContext(5, 15)
}
