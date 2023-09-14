package config

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/yaml.v2"
)

type KafkaConfig struct {
	Producer *kafka.ConfigMap `yaml:"producer"`
	Consumer *kafka.ConfigMap `yaml:"consumer"`
}

type Config struct {
	Limit    int               `yaml:"limit"`
	Addr     string            `yaml:"addr"`
	Env      string            `yaml:"env"`
	Logger   bool              `yaml:"logger"`
	Database map[string]string `yaml:"database"`
	Kaftka   *KafkaConfig      `yaml:"kafka"`
}

var config interface{}
var configPath string = ""

func getConfigArgs() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	configPathDefault := path.Join(dir, "config.yml")
	flag.StringVar(&configPath, "c", configPathDefault, "config path")
	flag.Parse()
}

func GetConfig() (Config, error) {
	if config != nil {
		return config.(Config), nil
	}

	if configPath == "" {
		getConfigArgs()
	}
	result := Config{
		Kaftka: &KafkaConfig{
			Consumer: &kafka.ConfigMap{},
			Producer: &kafka.ConfigMap{},
		},
	}
	f, err := os.Open(configPath)
	if err != nil {
		return result, err
	}
	buff, err := io.ReadAll(f)
	if err != nil {
		return result, err
	}
	fmt.Println(string(buff))
	err = yaml.Unmarshal(buff, &result)
	if err != nil {
		return result, err
	}
	config = result
	return config.(Config), nil
}
