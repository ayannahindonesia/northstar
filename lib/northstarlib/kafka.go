package northstarlib

import (
	"encoding/json"
	"log"

	"github.com/ayannahindonesia/basemodel"

	"github.com/Shopify/sarama"
)

type (
	// NorthstarLib main type
	NorthstarLib struct {
		Host         string
		Topic        string
		Secret       string
		Send         bool
		SaramaConfig *sarama.Config
	}
	// Log models
	Log struct {
		basemodel.BaseModel
		Level    string `json:"level"`
		Tag      string `json:"tag"`
		Messages string `json:"messages"`
	}
)

// SubmitKafkaLog func
func (n *NorthstarLib) SubmitKafkaLog(l Log, model string) (err error) {
	if !n.Send {
		return nil
	}
	if len(model) < 1 {
		model = "log"
	}
	build := kafkaLogBuilder(l, model)

	jMarshal, _ := json.Marshal(build)

	kafkaProducer, err := sarama.NewAsyncProducer([]string{n.Host}, n.SaramaConfig)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: n.Topic,
		Value: sarama.StringEncoder(n.Secret + ":" + model + ":" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", n.Topic)
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", n.Topic, err)
	}

	return err
}

func kafkaLogBuilder(l Log, model string) (payload map[string]interface{}) {
	inrec, _ := json.Marshal(l)
	json.Unmarshal(inrec, &payload)

	return payload
}
