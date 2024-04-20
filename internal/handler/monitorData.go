package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"logscan/pkg"
)

// 获取监控面版的数据，返回每个host的异常日志数量
func GetMonitorData(ctx *gin.Context) {
	client, err := elastic.NewClient(
		elastic.SetURL(pkg.ESConfig.URL),
		elastic.SetSniff(false),
	)
	if err != nil {
		fmt.Println("Failed to create elastic client:", err)
		return
	}

	start := ctx.Query("start")
	end := ctx.Query("end")

	timeRangeQuery := elastic.NewRangeQuery("@timestamp").Gte(start).Lte(end)
	anomalyQuery := elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("level", "Normal"))
	boolQuery := elastic.NewBoolQuery().Must(timeRangeQuery, anomalyQuery)
	hostQuery := elastic.NewTermsAggregation().Field("host.keyword").Size(1000).Missing(0)

	response, err := client.Search().
		Index("logscan").
		Query(boolQuery).
		Aggregation("host", hostQuery).
		Size(0).
		Do(ctx)
	if err != nil {
		fmt.Println("Failed to search log from es:", err)
		return
	}

	agg, found := response.Aggregations.Terms("host")
	if !found {
		fmt.Println("host aggregation not found")
		return
	}
	result := make(map[string]int64)
	for _, bucket := range agg.Buckets {
		host, err := bucket.Key.(string)
		if !err {
			fmt.Println("Failed to convert key to string")
			return
		}
		result[host] = bucket.DocCount
	}
	ctx.JSON(200, result)
}
