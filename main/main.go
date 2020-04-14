
package main

import (
	"JYB_Crawler/Basics"
	"JYB_Crawler/eduData"

	//_ "JYB_Crawler/elasticsearch"

)

func main() {

	Basics.StartMySql()
	//开始爬虫工作
	eduData.StartContext(5,15)
}

