package elasticsearch

import (
	"JYB_Crawler.Vn/Basics"
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
	"sync/atomic"
)

//Consumer
func BulkInsert(indexCtx context.Context) error {
	log.Println("begin bulk insert...")
	bulk := Client.Bulk().Index(Index).Type(Typ)
	for d := range Docsc {
		// Simple progress
		// AddUint64(): total增加1，并返回一个新的值 类似于total++，可避免多协程同时修改
		atomic.AddUint64(&Total, 1)
		// Enqueue the document
		bulk.Add(elastic.NewBulkIndexRequest().Id(strconv.Itoa(d.ID)).Doc(d))
		//当bulk中的doc的数量达到bulkSize时，执行一次批量插入操作
		if bulk.NumberOfActions() >= BulkSize {
			// Commit
			log.Println(Index, Total, "articles inserted successfully...")
			res, err := bulk.Do(indexCtx)
			if err != nil {
				return err
			}
			if res.Errors {
				// Look up the failed documents with res.Failed(), and e.g. recommit
				return errors.New("bulk commit failed")
			}
			// "bulk" is reset after Do, so you can reuse it
			log.Println(Index, "intermediate item bulk inserted successfully...")
		}

		select {
		default:
		case <-indexCtx.Done():
			return indexCtx.Err()
		}
	}
	// Commit the final batch before exiting
	if bulk.NumberOfActions() > 0 {
		_, err = bulk.Do(indexCtx)
		if err != nil {
			return err
		}
	}
	log.Println(Index, "all inserted successfully...")
	return nil
}

//插入类型表到es
func TpBulkInsert() error {
	log.Println("begin type bulk insert...")
	bulk := Client.Bulk().Index("crawler_type").Type("type_info")

	for _, d := range Basics.EveryType {

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
	log.Println("crawler_type bulk inserted successfully...")
	return nil
}
