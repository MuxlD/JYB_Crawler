package eduData

import (
	"context"
	"fmt"
	"log"
	"testing"
)

//检查每次新建的context是否相同
func TestNewChromedp(t *testing.T) {
	var ctx1, ctx2 context.Context
	for i := 1; i < 3; i++ {
		chr := NewChromedp(context.Background())
		if i == 1 {
			ctx1 = chr.allocCtx
		}
		ctx2 = chr.allocCtx
		log.Println(chr.allocCtx)
	}
	if ctx1 == ctx2 {
		fmt.Println(true)
		return
	}
	fmt.Println(false)
	//函数输出为false
}

func TestNewTab(t *testing.T) {
	chr := NewChromedp(context.Background())
	var ctx1, ctx2 context.Context
	for i := 1; i < 3; i++ {
		ctx, _ := chr.NewTab()
		if i == 1 {
			ctx1 = ctx
		}
		ctx2 = ctx
		fmt.Println(ctx)
	}
	if ctx1 == ctx2 {
		fmt.Println(true)
		return
	}
	fmt.Println(false)
	//函数输出为false
}
