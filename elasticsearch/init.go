package elasticsearch

import (
	"JYB_Crawler/Basics"
	"flag"
	"github.com/olivere/elastic/v7"
	"log"
)

var (
	url      string
	Index    string
	Typ      string
	sniff    bool
	BulkSize int
)

var (
	Client *elastic.Client
	err    error
	Docsc  chan Basics.TrainingSchool
	Total  uint64
)

func InitMapping() {
	flag.StringVar(&url, "url", "http://localhost:9200", "Elasticsearch URL")
	flag.StringVar(&Index, "index", "crawler_training_school", "Elasticsearch index name")
	flag.StringVar(&Typ, "type", "training_school", "Elasticsearch type name")
	flag.BoolVar(&sniff, "sniff", true, "Enable or disable sniffing")
	//每50条批量插入一次
	flag.IntVar(&BulkSize, "size", 50, "Number of documents to collect before committing")

	flag.PrintDefaults()
	//解析os.Args[1:]中的命令行标志
	flag.Parse()
	//设置log.Println的打印格式(文件+行号)  example: elasticsearch/simple_demo/log/SetFlags.go
	log.SetFlags(log.Lshortfile)

	if url == "" {
		log.Fatal("missing url parameter")
	}
	if Index == "" {
		log.Fatal("missing index parameter")
	}
	if Typ == "" {
		log.Fatal("missing type parameter")
	}
	if BulkSize <= 0 {
		log.Fatal("bulk-size must be a positive number")
	}

	Docsc = make(chan Basics.TrainingSchool, BulkSize)
	// Create an Elasticsearch client
	Client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(sniff))
	if err != nil {
		log.Fatal(err)
	}

	//创建索引结构
	Mapping(Index, TsMapping)
	Mapping("crawler_type", TpMapping)
}
