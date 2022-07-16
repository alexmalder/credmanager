package src

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// The data struct for the decoded data
// Notice that all fields must be exportable!
type Bw struct {
	Items []Item `json:"items"`
}

type Item struct {
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Uris     []struct {
			Uri string `json:"uri"`
		}
	}
	Name string `json:"name"`
    Notes string `json:"notes"`
}

func ReadJson() []Item {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(os.Getenv("BITWARDEN_BACKUP_PATH"))
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	// Now let's unmarshall the data into `payload`
	var payload Bw
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	//log.Println(payload)
    return payload.Items
}
