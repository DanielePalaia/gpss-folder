package main

import (
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type kafkaEngine struct {
	kafkaIp   []string
	topicJson string
	topicCsv  string
	topicText string
	delimiter string
	wJson     *kafka.Writer
	wCsv      *kafka.Writer
	//wTxt      *kafka.Writer
}

var writer *kafka.Writer

func makeKafkaEngine(kafkaIp string, topicJson string, topicCsv string, topicText string) *kafkaEngine {
	engine := new(kafkaEngine)
	engine.kafkaIp = make([]string, 1)
	engine.kafkaIp[0] = kafkaIp
	engine.topicJson = topicJson
	engine.topicCsv = topicCsv
	//engine.topicText = topicText

	return engine
}

func (engine *kafkaEngine) createConnection(topic string) *kafka.Writer {

	var w *kafka.Writer

	// Connection to the first topic
	if topic != "" {
		config := kafka.WriterConfig{
			Brokers:      engine.kafkaIp,
			Topic:        engine.topicJson,
			BatchTimeout: 1 * time.Nanosecond,
		}

		w = kafka.NewWriter(config)
		if w == nil {
			fmt.Println("Error creating topic")
		}
	}

	return w
}

func (engine *kafkaEngine) Configure() {

	if engine.topicJson != "" {
		engine.wJson = engine.createConnection(engine.topicJson)
	}
	if engine.topicCsv != "" {
		engine.wCsv = engine.createConnection(engine.topicCsv)
	}
	/*if engine.topicText != "" {
		engine.wTxt = engine.createConnection(engine.topicText)
	}*/

}

func (engine *kafkaEngine) PushJson(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return engine.wJson.WriteMessages(parent, message)
}

func (engine *kafkaEngine) PushCsv(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return engine.wCsv.WriteMessages(parent, message)
}

/*func (engine *kafkaEngine) PushTxt(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return engine.wTxt.WriteMessages(parent, message)
}*/
