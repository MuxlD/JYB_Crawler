
package main

import (
	"chromedp_test/Basics"
	"chromedp_test/eduData"
	"context"
	"log"
)

func main() {

	Basics.StartMySql()

	parentCtx:=context.Background()
	log.Println("parentCtx:",parentCtx)

	//生成一个可取消的context
	ctx,cancel := context.WithCancel(parentCtx)
	defer cancel()
	log.Println("The first cancelable context:",ctx)

	//开始爬虫工作
	eduData.StartContext(ctx,5,15)
}

