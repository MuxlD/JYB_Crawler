package elasticsearch

import (
	"context"
	"log"
)

const (
	TsMapping = `{
  "mappings": {
    "training_school": {
      "properties": {
        "id": {
          "type": "long"
        },
        "url": {
          "type": "keyword"
        },
        "name": {
          "type": "keyword"
        },
        "type_id": {
          "type": "integer"
        },
        "type_url": {
          "type": "keyword"
        },
        "type_name": {
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

	TpMapping = `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "type_info": {
      "properties": {
        "type_id": {
          "type": "integer"
        },
        "type_url": {
          "type": "keyword"
        },
        "type_name": {
          "type": "keyword"
        },
        "max_page": {
          "type": "integer"
        },
        "count": {
          "type": "integer"
        }
      }
    }
  }
}`
)

func Mapping(index, mapping string) {
	ctx := context.Background()
	//验证索引index是否存在
	//如果index存在，则将其删除
	exists, err := Client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Fatal(index, "IndexExists error:", err)
	}
	if exists {
		_, err := Client.DeleteIndex(index).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	//按照mapping新建索引
	_, err = Client.CreateIndex(index).Body(mapping).Do(ctx)
	if err != nil {
		log.Fatal(index, "CreateIndex error:", err)
	}
	log.Println(index, "mapping created successfully...")
}
