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
		SaramaConfig *sarama.Config
	}
	// Log models
	Log struct {
		basemodel.BaseModel
		Level    string `json:"level"`
		Messages string `json:"messages"`
	}
)

// SubmitKafkaLog func
func (n *NorthstarLib) SubmitKafkaLog(i interface{}, model string, secret string) (err error) {
	if len(model) < 1 {
		model = "log"
	}
	build := kafkaLogBuilder(i, model)

	jMarshal, _ := json.Marshal(build)

	kafkaProducer, err := sarama.NewAsyncProducer([]string{n.Host}, n.SaramaConfig)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: n.Topic,
		Value: sarama.StringEncoder(secret + ":" + model + ":" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", n.Topic)
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", n.Topic, err)
	}

	return err
}

func kafkaLogBuilder(i interface{}, model string) (payload map[string]interface{}) {
	inrec, _ := json.Marshal(i)
	json.Unmarshal(inrec, &payload)

	return payload
}
