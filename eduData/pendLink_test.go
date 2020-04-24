package eduData

import (
	"JYB_Crawler/Basics"
	"context"
	"testing"
)

var testTsUrl Basics.TsUrl

func Test_crawlerByPendUrl(t *testing.T) {
	chrome := NewChromedp(context.Background())
	testTsUrl = Basics.TsUrl{
		TypeID: 1,
		Url:    "https://cs.jiaoyubao.cn/edu/zsrcjykjyxgs/",
	}
	_ = crawlerByPendUrl(testTsUrl, chrome)
}
