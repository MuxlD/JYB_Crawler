package main

import (
	"JYB_Crawler.Vn/Basics"
	"JYB_Crawler.Vn/eduData"
	"JYB_Crawler.Vn/elasticsearch"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	//初始化
	Basics.StartMySql()
	elasticsearch.InitMapping()
	//开始爬虫工作
	ctx, cancel := context.WithCancel(context.Background())
	go eduData.StartContext(ctx, 5, 15)

	<-ch

	fmt.Println("收到 ctrl+c 命令....")

	cancel()

	fmt.Println("程序完整退出....")
}
