package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"northstar/application"
	"northstar/models"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

type (
	// Consumer type
	Consumer struct {
		KafkaConsumer     sarama.Consumer
		PartitionConsumer sarama.PartitionConsumer
	}
)

var wg sync.WaitGroup

func init() {
	var err error
	topic := application.App.Config.GetString(fmt.Sprintf("%s.kafka.topic", application.App.ENV))

	kafka := &Consumer{}
	kafka.KafkaConsumer, err = sarama.NewConsumer([]string{application.App.Kafka.Host}, application.App.Kafka.Config)
	if err != nil {
		log.Printf("error while creating new kafka consumer : %v", err)
	}

	kafka.SetPartitionConsumer(topic)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer kafka.KafkaConsumer.Close()
		for {
			message, err := kafka.Listen()
			if err != nil {
				log.Printf("error occured when listening kafka : %v", err)
			}
			if message != nil {
				err := processMessage(message)
				if err != nil {
					log.Printf("%v . message : %v", err, string(message))
				}
			}
		}
	}()
}

// SetPartitionConsumer func
func (k *Consumer) SetPartitionConsumer(topic string) (err error) {
	k.PartitionConsumer, err = k.KafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetOldest)

	return err
}

// Listen to kafka
func (k *Consumer) Listen() ([]byte, error) {
	select {
	case err := <-k.PartitionConsumer.Errors():
		return nil, err
	case msg := <-k.PartitionConsumer.Messages():
		return msg.Value, nil
	}
}

func processMessage(kafkaMessage []byte) (err error) {
	splitKafkaString := strings.SplitN(string(kafkaMessage), ":", 2)

	switch splitKafkaString[0] {
	default:
		return fmt.Errorf("cannot process kafka message : %v", string(kafkaMessage))
	case "log":
		mod := models.Log{}

		json.Unmarshal([]byte(splitKafkaString[1]), &mod)
		err = mod.Create()
		break
	}

	return err
}
