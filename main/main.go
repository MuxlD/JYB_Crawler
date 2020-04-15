package main

import (
	"JYB_Crawler/Basics"
	"JYB_Crawler/eduData"
	"JYB_Crawler/elasticsearch"
	//_ "JYB_Crawler/elasticsearch"
)

func main() {

	Basics.StartMySql()
	elasticsearch.InitMapping()
	//开始爬虫工作
	eduData.StartContext(5, 15)
}
