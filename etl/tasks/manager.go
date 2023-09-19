package tasks

import (
	"context"
	"errors"

	"github.com/juliotorresmoreno/iot/etl/kafka"
)

type Task interface {
	Context() context.Context
	Name() string
	Channel() string
	UUID() string
}

var DefaultTaskManager *TaskManager

func init() {
	DefaultTaskManager, _ = NewTaskManager()
}

type TaskManager struct {
	channel  string
	tasks    map[string]Task
	kafkaCli *kafka.KafkaClient
}

func NewTaskManager() (*TaskManager, error) {
	t := &TaskManager{
		channel: "jobs",
		tasks:   make(map[string]Task),
	}

	kafkaCli, err := kafka.NewKafkaClient("jobs")
	t.kafkaCli = kafkaCli

	return t, err
}

func (t *TaskManager) Add(task Task) error {
	if task.UUID() == "" {
		return errors.New("uuid is not provided")
	}
	t.tasks[task.UUID()] = task

	err := t.kafkaCli.Pub(map[string]interface{}{
		"uuid":    task.UUID(),
		"name":    task.Name(),
		"channel": task.Channel(),
	})

	return err
}

func (t *TaskManager) Del(uuid string) error {
	if uuid == "" {
		return errors.New("uuid is not provided")
	}
	delete(t.tasks, uuid)
	return nil
}

func (t *TaskManager) List() []Task {
	tasks := make([]Task, 0)
	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (t *TaskManager) Subscribe() error {
	go t.kafkaCli.Sub(func(ch string, data any) error {
		return nil
	})
	return nil
}
