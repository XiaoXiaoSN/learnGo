package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	brokers = "localhost:32768"
	topic   = "topic001"
)

func main() {
	broker := sarama.NewBroker(brokers)
	err := broker.Open(nil)
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	// 取得 broker 的資訊
	{
		request := sarama.MetadataRequest{Topics: []string{topic}}
		response, err := broker.GetMetadata(&request)
		if err != nil {
			panic(err)
		}

		fmt.Println("There are", len(response.Topics), "topics active in the cluster.")
		fmt.Printf("%+v\n\n", response)
	}

	// 發資料囉
	{
		request := sarama.ProduceRequest{
			RequiredAcks: sarama.NoResponse,
		}
		request.AddMessage(topic, 0, &sarama.Message{
			Key:   []byte("key"),
			Value: []byte("hello testing"),
		})
		response, err := broker.Produce(&request)
		if err != nil {
			panic(err)
		}

		fmt.Printf("發資料囉 %+v\n", response)
	}

	return
}
