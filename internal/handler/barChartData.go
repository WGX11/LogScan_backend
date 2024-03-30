package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
)

// 获取search界面stacklBar图表数据
func GetBarChartData(ctx *gin.Context) {
	//es客户端
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Println("Failed to create elastic client:", err)
	}

	//参数处理
	startTime := ctx.Query("start")
	endTime := ctx.Query("end")
	lucene := ctx.Query("lucene")
	var boolQuery elastic.Query
	timeRangeQuery := elastic.NewRangeQuery("@timestamp").Gte(startTime).Lte(endTime)
	if len(lucene) == 0 {
		boolQuery = timeRangeQuery
	} else {
		luceneQuery := elastic.NewQueryStringQuery(lucene)
		boolQuery = elastic.NewBoolQuery().Must(timeRangeQuery, luceneQuery)
	}

	//获得按照分钟聚合的数据，分别获取正常和异常的数据
	dateAgg := elastic.NewDateHistogramAggregation().Field("@timestamp").CalendarInterval("minute")
	normalAgg := elastic.NewFiltersAggregation().Filters(elastic.NewTermQuery("level.keyword", "Normal"))
	anomalyAgg := elastic.NewFiltersAggregation().Filters(elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("level.keyword", "Normal")))
	dateAgg.SubAggregation("normal", normalAgg)
	dateAgg.SubAggregation("anomaly", anomalyAgg)

	//es查询
	response, err := client.Search().
		Index("logscan").
		Query(boolQuery).
		Aggregation("documents_over_time", dateAgg).
		Size(0).
		Do(context.Background())
	if err != nil {
		log.Println("Failed to search documents:", err)
	}

	result, found := response.Aggregations.DateHistogram("documents_over_time")
	if !found {
		log.Println("Failed to find aggregation:", err)
	}
	ctx.JSON(200, result.Buckets)

}
