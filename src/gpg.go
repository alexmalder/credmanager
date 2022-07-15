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
	gpghomedir     = os.Getenv("GPG_HOMEDIR") + "/"
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
func readKeyring(keyring string) (openpgp.EntityList, error) {
	// Read in public key
	keyringFileBuffer, err := os.Open(keyring)
	if err != nil {
		return nil, err
	}
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}
	return entityList, nil
}

// encode function
func EncTest(secretString string) (string, error) {
	log.Printf("Secret to hide: \n%s", secretString)
	log.Printf("Public Keyring: %s", gpghomedir+"pubring.gpg")

	// Read in public key
	entityList, err := readKeyring(gpghomedir + "pubring.gpg")

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(secretString))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	// Encode to base64
	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}
	encStr := base64.StdEncoding.EncodeToString(bytes)

	// Output encrypted/encoded string
	log.Println("Encrypted Secret:", encStr)

	return encStr, nil
}

// decode function
func DecTest(encString string) (string, error) {
	log.Println("Secret Keyring:", gpghomedir+"secring.gpg")
	log.Println("Passphrase:", passphrase)

	// Open the private key file
	var entity *openpgp.Entity
	entityList, err := readKeyring(gpghomedir + "secring.gpg")
	if len(entityList) != 1 {
		return "", err
	}
	entity = entityList[0]

	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	passphraseByte := []byte(passphrase)
	log.Println("Decrypting private key using passphrase")
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}
	log.Println("Finished decrypting private key using passphrase")

	// Decode the base64 string
	dec, err := base64.StdEncoding.DecodeString(encString)
	if err != nil {
		return "", err
	}

	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	decStr := string(bytes)

	return decStr, nil
}

// test encode and decode functions
func TestGPG() {
	rawData := FileAsString("test.yml")
	encStr, err := EncTest(rawData)
	if err != nil {
		log.Fatal(err)
	}
	decStr, err := DecTest(encStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Decrypted Secret: \n%s", decStr)
	if rawData == decStr {
		log.Println("Test successfull")
	}
}
