package Own

import (
	"chromedp_test/elasticsearch"
	"context"
	"log"
)

const (
	mapping = `{
  "mappings": {
    "training_school": {
      "properties": {
        "id": {
          "type": "keyword"
        },
        "type_name": {
          "type": "keyword"
        },
        "type_url": {
          "type": "keyword"
        },
        "type_id": {
          "type": "integer"
        },
        "name": {
          "type": "keyword"
        },
        "url": {
          "type": "keyword"
        },
        "bright_spot": {
          "type": "keyword"
        },
        "info": {
          "type": "text"
        },
        "course": {
          "type": "keyword"
        },
        "campus": {
          "type": "keyword",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        },
        "phone_number": {
          "type": "keyword"
        }
      }
    }
  },
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 0
  }
}`
)

func Mapping(index string) {
	ctx := context.Background()
	//验证索引index是否存在
	//如果index存在，则将其删除
	exists, err := elasticsearch.Client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		_, err := elasticsearch.Client.DeleteIndex(index).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	//按照mapping新建索引
	_, err = elasticsearch.Client.CreateIndex(index).Body(mapping).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
