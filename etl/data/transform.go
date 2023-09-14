package data

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/juliotorresmoreno/iot/etl/db"
	"github.com/juliotorresmoreno/iot/etl/entity"
	"github.com/juliotorresmoreno/iot/etl/kafka"
	"github.com/juliotorresmoreno/iot/etl/tasks"
)

type ETL struct {
	channel  string
	dbCli    *db.Manager
	kafkaCli *kafka.KafkaClient
}

func MakeETL(source entity.Source) (*ETL, error) {
	e := &ETL{
		channel: "import",
	}

	dbCli, err := db.MakeManager(source)
	if err != nil {
		return e, err
	}
	e.dbCli = dbCli

	return e, err
}

func (ETL *ETL) SetKafkaClient(kafkaCli *kafka.KafkaClient) {
	ETL.kafkaCli = kafkaCli
}

type Task struct {
	uuid    string
	name    string
	channel string
	context context.Context
}

func (t *Task) Context() context.Context {
	return t.context
}

func (t *Task) UUID() string {
	return t.uuid
}

func (t *Task) Name() string {
	return t.name
}

func (t *Task) Channel() string {
	return t.channel
}

func (t *Task) MarshalJSON() ([]byte, error) {
	buff := bytes.NewBufferString("")
	json.NewEncoder(buff).Encode(map[string]string{
		"uuid":    t.UUID(),
		"name":    t.Name(),
		"channel": t.Channel(),
	})
	return buff.Bytes(), nil
}

func (c *ETL) Run() tasks.Task {
	ctx := context.Background()

	task := Task{
		context: ctx,
		uuid:    uuid.New().String(),
		name:    "data",
		channel: c.channel,
	}

	return &task
}
