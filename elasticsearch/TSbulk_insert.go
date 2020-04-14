package elasticsearch

import (
	"JYB_Crawler/Basics"
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
	"strconv"
	"sync/atomic"
	"time"
)

func BulkInsert() error {
	defer close(Docsc)
	bulk := Client.Bulk().Index(Index).Type(Typ)
	for d := range Docsc {
		// Simple progress
		// AddUint64(): total增加1，并返回一个新的值 类似于total++
		atomic.AddUint64(&Total, 1)

		// Enqueue the document
		bulk.Add(elastic.NewBulkIndexRequest().Id(strconv.Itoa(d.ID)).Doc(d))
		//当bulk中的doc的数量达到bulkSize时，执行一次批量插入操作
		if bulk.NumberOfActions() >= BulkSize {
			// Commit
			res, err := bulk.Do(IndexCtx)
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
		case <-IndexCtx.Done():
			return IndexCtx.Err()
		}
		//写入一次程序暂停3s
		time.Sleep(3e9)
	}
	// Commit the final batch before exiting
	if bulk.NumberOfActions() > 0 {
		_, err := bulk.Do(IndexCtx)
		if err != nil {
			return err
		}
	}
	return nil
}


//插入类型表到es
func TpBulkInsert() error {

	bulk := Client.Bulk().Index("crawler_type").Type("type_info")

	for _, d := range Basics.EveryType {
		// AddUint64(): total增加1，并返回一个新的值 类似于total++
		atomic.AddUint64(&Total, 1)
		// Enqueue the document
		bulk.Add(elastic.NewBulkIndexRequest().Id(strconv.Itoa(int(d.ID))).Doc(d))

		if bulk.NumberOfActions() >= 20 {
			// Commit
			res, err := bulk.Do(context.Background())
			if err != nil {
				return err
			}
			if res.Errors {
				// Look up the failed documents with res.Failed(), and e.g. recommit
				return errors.New("bulk commit failed")
			}
			// "bulk" is reset after Do, so you can reuse it
		}
	}
	return nil
}
