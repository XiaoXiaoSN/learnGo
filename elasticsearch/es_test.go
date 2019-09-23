package elasticsearch_test

import (
	"testing"

	es "github.com/XiaoXiaoSN/elasticsearch"
)

func TestConnection(t *testing.T) {
	t.Run("", func(t *testing.T) {
		cfg := es.Configuration{
			Addresses: []string{
				"http://localhost:9200",
			},
			Username: "foo",
			Password: "bar",
		}

		es, _ := es.NewElasticSearch(cfg)
		es.GetConfig()
	})
}
