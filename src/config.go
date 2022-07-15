package src

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// config yaml for sql
type Config struct {
	Queries []struct {
		Query string
	}
	InsertSecret  string `yaml:"insert_secret"`
	SelectSecrets string `yaml:"select_secrets"`
	SelectSecret  string `yaml:"select_secret"`
	UpdateSecret  string `yaml:"update_secret"`
	DeleteSecret  string `yaml:"delete_secret"`
}

// read config yaml and return Config object
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
	//log.Println(data)
	return data
}
