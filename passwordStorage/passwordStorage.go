package passwordStorage //remake as file interaction/password storage

import (
	pswrd "password/passwordFunctions"
	//bubbleT "Users/landondixon/src/goCode/passwordManager/bubbleTeaFunctions"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/charm/kv"
	"gopkg.in/yaml.v3"
	//"io/ioutil"
)

// reads data from charm cloud9L
func ReadDataCharm(c pswrd.CharmConfig, db *kv.KV) ([]byte, bool) {
	//var charmData []byte
	display := false

	// Fetch updates and easily define your own syncing strategy
	if err := db.Sync(); err != nil {
		log.Fatal(err)
	}
	// _, err := db.Get(context.Background(), key)
	// if err == charmstore.ErrNotFound {
	// 	// Key does not exist
	// 	return false, nil
	// if len(charmData)==0{ //find out how to check if a key works

	// }else{
	// Quickly get a value
	charmData, err := db.Get([]byte(c.AccountKey)) //charmData, err := ioutil.ReadFile("passwords.yaml")

	if charmData == nil { //don't know what this error actually is
		// Connection refused error
		display = true
	} else if err != nil {
		panic(err)
	}

	return charmData, display
}

// read data from a file, just reads it in and stores it all together
func ReadData(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("error opening %s: %s", fileName, err)
		return ""
	}
	return string(data)
}

// writes the data to the given file, creates the file if it doesn't exist
func WritePasswords(fileName string, passwords []pswrd.SavedPassword) {
	var pswrds pswrd.Passwords
	pswrds.Pswrds = passwords
	// pswrds.Config = c
	for i, password := range pswrds.Pswrds {
		//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
		pswrds.Pswrds[i].EncryptedPswrd = pswrd.Encrypt(password.EncryptedPswrd, len(pswrds.Pswrds[0].Website))

	}
	var allText string
	for _, pswrd := range pswrds.Pswrds {
		allText = allText + pswrd.Website + " : " + pswrd.EncryptedPswrd + "\n"
	}
	data := []byte(allText)
	err := os.WriteFile(fileName, data, 0644) //0644 specifies the file is readable and writeable by the owner and readable by everyone else
	if err != nil {
		fmt.Println(err)
		return
	}
}

// writes to charm cloud
func WritePasswordsCharm(passwords []pswrd.SavedPassword, Key int, db *kv.KV, c pswrd.CharmConfig) {

	var pswrds pswrd.Passwords
	pswrds.Pswrds = passwords
	pswrds.Config = c
	for i, password := range pswrds.Pswrds {
		//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
		pswrds.Pswrds[i].EncryptedPswrd = pswrd.Encrypt(password.EncryptedPswrd, len(pswrds.Pswrds[0].Website))

	}

	yamlData, err := yaml.Marshal(&pswrds)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	//fmt.Println(string(yamlData))

	// fmt.Println(c.AccountKey)
	if err := db.Set([]byte(c.AccountKey), yamlData); err != nil { //failing now for some reason
		log.Fatal(err)
	}

	if err := db.Sync(); err != nil {
		log.Fatal(err)
	}

	//fileName := "passwords.yaml"

	// err = ioutil.WriteFile(fileName, yamlData, 0644)
	// if err != nil {
	//     panic("Unable to write data into the file")
	// }
}