package pkg

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"time"
)

func SearchLogFromEs(startTime, endTime, lucene string) *elastic.SearchResult {
	//创建Elasticsearch客户端配置
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Println("Failed to create es client:", err)
	}
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		log.Println("Failed to parse start time:", err)
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		log.Println("Failed to parse end time:", err)
	}

	//配置查询query
	var query elastic.Query
	query = elastic.NewRangeQuery("@timestamp").Gte(start).Lte(end)
	//Lucene风格的条件过滤
	if len(lucene) != 0 {
		luceneQuery := elastic.NewQueryStringQuery(lucene)
		query = elastic.NewBoolQuery().Must(query, luceneQuery)
	}
	response, err := client.Search().
		Index("logscan").
		Query(query).
		Size(500).
		Sort("@timestamp", true).
		Do(context.Background())
	if err != nil {
		log.Println("Failed to search log from es:", err)
	}
	fmt.Println(response)
	return response
}
