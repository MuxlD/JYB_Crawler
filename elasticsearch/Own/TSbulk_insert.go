package Own

import (
	"chromedp_test/elasticsearch"
	"errors"
	"github.com/olivere/elastic/v7"
	"sync/atomic"
	"time"
)


func BulkIndex() error {
	defer close(elasticsearch.Docsc)
	bulk := elasticsearch.Client.Bulk().Index(elasticsearch.Index).Type(elasticsearch.Typ)
	for d := range elasticsearch.Docsc {
		// Simple progress
		// AddUint64(): total增加1，并返回一个新的值 类似于total++
		atomic.AddUint64(&elasticsearch.Total, 1)
		//返回从begin开始的时间

		// Enqueue the document
		bulk.Add(elastic.NewBulkIndexRequest().Id(d.ID).Doc(d))
		//当bulk中的doc的数量达到bulkSize时，执行一次批量插入操作
		if bulk.NumberOfActions() >= elasticsearch.BulkSize {
			// Commit
			res, err := bulk.Do(elasticsearch.IndexCtx)
			if err != nil {
				return err
			}
			if res.Errors {
				// Look up the failed documents with res.Failed(), and e.g. recommit
				return errors.New("bulk commit failed")
			}
			// "bulk" is reset after Do, so you can reuse it
		}

		select {
		default:
		case <-elasticsearch.IndexCtx.Done():
			return elasticsearch.IndexCtx.Err()
		}
		//写入一次程序暂停3s
		time.Sleep(3e9)
	}
	// Commit the final batch before exiting
	if bulk.NumberOfActions() > 0 {
		_, err := bulk.Do(elasticsearch.IndexCtx)
		if err != nil {
			return err
		}
	}
	return nil
}
