package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
	"logscan/pkg"
	"time"
)

// 获取search界面stacklBar图表数据
func GetBarChartData(ctx *gin.Context) {
	//es客户端
	client, err := elastic.NewClient(
		elastic.SetURL(pkg.ESConfig.URL),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Println("Failed to create elastic client:", err)
	}

	//参数处理
	startTime := ctx.Query("start")
	endTime := ctx.Query("end")
	lucene := ctx.Query("lucene")
	//根据时间间隔来决定interval
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		log.Println("Failed to parse start time:", err)
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		log.Println("Failed to parse end time:", err)
	}
	duration := end.Sub(start)
	var timeInterval string
	if duration.Hours()/24 > 365 {
		//时间跨度超过一年，按照16天为一个间隔
		timeInterval = "16d"
	} else if duration.Hours()/24 > 180 {
		//时间跨度超过半年，按照7天为一个间隔
		timeInterval = "7d"
	} else if duration.Hours()/24 > 30 {
		//时间跨度超过一个月，按照1天为一个间隔
		timeInterval = "1d"
	} else if duration.Hours()/24 > 7 {
		//时间跨度超过一周，按照12小时为一个间隔
		timeInterval = "12h"
	} else if duration.Hours()/24 > 3 {
		//时间跨度超过三天，按照1小时为一个间隔
		timeInterval = "1h"
	} else if duration.Hours()/24 > 1 {
		//时间跨度超过一天，按照1分钟为一个间隔
		timeInterval = "1m"
	} else if duration.Minutes() > 30 {
		//时间跨度超过半小时，按照30秒为一个间隔
		timeInterval = "30s"
	} else {
		//时间跨度小于半小时，按照1秒为一个间隔
		timeInterval = "1s"
	}
	var boolQuery elastic.Query
	timeRangeQuery := elastic.NewRangeQuery("@timestamp").Gte(startTime).Lte(endTime)
	if len(lucene) == 0 {
		boolQuery = timeRangeQuery
	} else {
		luceneQuery := elastic.NewQueryStringQuery(lucene)
		boolQuery = elastic.NewBoolQuery().Must(timeRangeQuery, luceneQuery)
	}

	//获得按照分钟聚合的数据，分别获取正常和异常的数据
	dateAgg := elastic.NewDateHistogramAggregation().Field("@timestamp").FixedInterval(timeInterval)
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
