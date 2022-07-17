package src

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/openpgp"
)

var (
	gpghomedir = os.Getenv("GPG_HOMEDIR")
	gpgsecring = os.Getenv("GPG_SECRING")
	gpgpubring = os.Getenv("GPG_PUBRING")
	passphrase = os.Getenv("GPG_PASSPHRASE")
)

// read file as string utility function
func FileAsString(path string) string {
	content, err := ioutil.ReadFile(path) // the file is inside the local directory
	if err != nil {
		log.Println("Err")
	}
	//log.Println(string(content))    // This is some content
	return string(content)
}

// read keyring from file, public or private
func readKeyring(keyring string) openpgp.EntityList {
	// Read in public key
	keyringFileBuffer, err := os.Open(keyring)
	checkErr(err)
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	checkErr(err)
	return entityList
}

// encode function
func EncryptString(secretString string) string {
	//log.Printf("Secret to hide: \n%s", secretString)
	// Read in public key
	entityList := readKeyring(gpghomedir + "/" + gpgpubring)
	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	checkErr(err)
	_, err = w.Write([]byte(secretString))
	checkErr(err)
	err = w.Close()
	checkErr(err)
	// Encode to base64
	bytes, err := ioutil.ReadAll(buf)
	checkErr(err)
	encStr := base64.StdEncoding.EncodeToString(bytes)
	// Output encrypted/encoded string
	//log.Println("Encrypted Secret:", encStr)
	return encStr
}

// decode function
func DecryptString(encString string) string {
	//log.Println("Passphrase:", passphrase)
	// Open the private key file
	var entity *openpgp.Entity
	entityList := readKeyring(gpghomedir + "/" + gpgsecring)
	if len(entityList) != 1 {
		log.Fatal("Entity list length is not 1")
	}
	entity = entityList[0]
	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	passphraseByte := []byte(passphrase)
	//log.Println("Decrypting private key using passphrase")
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}
	//log.Println("Finished decrypting private key using passphrase")

	// Decode the base64 string
	dec, err := base64.StdEncoding.DecodeString(encString)
	checkErr(err)
	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
	checkErr(err)
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	checkErr(err)
	decStr := string(bytes)
	return decStr
}
