package src

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// config yaml for sql
type Config struct {
	/*
		Queries []struct {
			Query string
		}
		Drops []struct {
			Query string
		}
		InsertSecret   string `yaml:"insert_secret"`
		SelectSecrets  string `yaml:"select_secrets"`
		SelectSecret   string `yaml:"select_secret"`
		UpdateSecret   string `yaml:"update_secret"`
		DeleteSecret   string `yaml:"delete_secret"`
		InsertRevision string `yaml:"insert_revision"`
	*/
}

var (
	insertSecret       = "INSERT INTO secrets (key, revision, value, username, uri, notes, type) VALUES($1, $2, $3, $4, $5, $6, $7)"
	selectSecrets      = "SELECT * FROM secrets WHERE is_deleted=false"
	selectSecret       = "SELECT * FROM secrets WHERE key=$1 and is_deleted=false"
	updateSecret       = "UPDATE secrets SET revision=$2, value=$3, username=$4, uri=$5, notes=$6, is_deleted=$7 WHERE key=$1"
	insertRevision     = "INSERT INTO revision (key, revision, value, username, uri, notes, type, is_deleted) VALUES($1, $2, $3, $4, $5, $6, $7, $8)"
	createTableSecrets = `
	CREATE TABLE IF NOT EXISTS secrets (
        key VARCHAR(255) UNIQUE PRIMARY KEY NOT NULL,
        revision INTEGER NOT NULL,
        value VARCHAR(4096) NOT NULL,
        username VARCHAR(255),
        uri VARCHAR(1024),
        notes VARCHAR(4096),
        type VARCHAR(8),
        is_deleted BOOLEAN DEFAULT FALSE
    )`
	createTableRevision = `
	CREATE TABLE IF NOT EXISTS revision (
        id SERIAL PRIMARY KEY,
        key VARCHAR(255) NOT NULL,
        revision INTEGER NOT NULL,
        value VARCHAR(4096) NOT NULL,
        username VARCHAR(255),
        uri VARCHAR(1024),
        notes VARCHAR(4096),
        type VARCHAR(8),
        is_deleted BOOLEAN DEFAULT FALSE
      )
	`
	dropTableSecrets  = `DROP TABLE secrets`
	dropTableRevision = `DROP TABLE revision`
)

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
