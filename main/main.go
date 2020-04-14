
package main

import (
	"chromedp_test/eduDate"
	"context"
	"log"
)

func main() {
	//连接数据库
	//Basics.StartMySql()

	//context.Background()为空context，不可取消
	parentCtx:=context.Background()
	log.Println("parentCtx:",parentCtx)

	//生成一个可取消的context
	ctx,cancel := context.WithCancel(parentCtx)
	defer cancel()
	log.Println("The first cancelable context:",ctx)

	//开始爬虫工作
	eduDate.StartContext(ctx,5,15)
}

