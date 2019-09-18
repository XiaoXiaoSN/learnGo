package elasticsearch

import (
	"log"

	es "github.com/elastic/go-elasticsearch/v6"
)

// ElasticSearch 是一個 elastic 物件
type ElasticSearch struct {
	Client *es.Client
	config *es.Config
}

// NewElasticSearch 產生一個 ElasticSearch 物件
func NewElasticSearch(cfg *Configuration) (*ElasticSearch, error) {
	esConfig := Configuration{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}
	es, err := es.NewClient(esConfig)
	if err != nil {
		log.Printf("Error creating the elasticsearch client: %s\n", err)
	}

	return &ElasticSearch{
		Client: es,
		config: &esConfig,
	}, nil
}

// GetConfig 回傳啟動這個 elastic 物件的設定檔案
func (es *ElasticSearch) GetConfig() *es.Config {
	return es.config
}
