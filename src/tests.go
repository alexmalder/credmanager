package src

import "log"

// test encode and decode functions
func TestGPG() {
	rawData := FileAsString("test.yml")
	encStr := EncryptString(rawData)
	decStr := DecryptString(encStr)
	log.Printf("Decrypted Secret: \n%s", decStr)
	if rawData == decStr {
		log.Println("Test successfull")
	}
}
