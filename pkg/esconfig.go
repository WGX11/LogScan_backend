package pkg

import "os"

type ElasticSearchConfig struct {
	URL string
}

var ESConfig ElasticSearchConfig

func InitConfig() {
	ESConfig.URL = os.Getenv("ELASTICSEARCH_URL")
	if ESConfig.URL == "" {
		ESConfig.URL = "http://localhost:9200"
	}
}

func init() {
	InitConfig()
}
