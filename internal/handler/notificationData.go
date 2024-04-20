package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
	"logscan/pkg"
	"time"
)

// 获取持续更新的报警日志数据
func GetNotificationData(ctx *gin.Context) {
	//创建Elasticsearch客户端配置
	client, err := elastic.NewClient(
		elastic.SetURL(pkg.ESConfig.URL),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Println("Failed to create es client:", err)
	}

	startTime := ctx.Query("start")
	endTime := ctx.Query("end")
	lucene := ctx.Query("lucene")
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
	//过滤异常日志
	anomalyTermQuery := elastic.NewTermQuery("level.keyword", "Normal")
	query = elastic.NewBoolQuery().Must(query).MustNot(anomalyTermQuery)
	response, err := client.Search().
		Index("logscan").
		Query(query).
		Size(500).
		Sort("@timestamp", true).
		Do(context.Background())
	if err != nil {
		log.Println("Failed to search log from es:", err)
	}
	ctx.JSON(200, response.Hits.Hits)
}
