package src

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Queries []struct {
		Query string
	}
	InsertSecret  string `yaml:"insert_secret"`
	SelectSecrets string `yaml:"select_secrets"`
	UpdateSecret  string `yaml:"update_secret"`
	DeleteSecret  string `yaml:"delete_secret"`
}

func ReadConfig() Config {
	yfile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	data := Config{}
	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		log.Fatal(err2)
	}
	//for _, v := range data.Queries {fmt.Printf("%s\n", v.Query)}
	return data
}
