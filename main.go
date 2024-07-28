package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	stor "password/passwordStorage"
	bubbleT "password/userInterface"
	"time"

	pswrd "password/passwordFunctions" //"Users/landondixon/src/goCode/passwordManager/passwordFunctions"  //github.com/landonpd/password/

	tea "github.com/charmbracelet/bubbletea" //if my code is on github I can just include it from their, might be a sneaky little way to get my file there
	"github.com/charmbracelet/charm/kv"
	"gopkg.in/yaml.v3"
	// "crypto/aes" I will be useing aes
	// "crypto/rand"
)

// SavedPassword represents a password entry

func main() {
	//variables and constants
	const LIMIT = 100000

	//var savedMasterPassword string //, stringtoWrite,inputtedmasterPassword, string
	var newKey int //, oldKey, passwordCount int //choice, createPasswordChoice,
	//var passwords []pswrd.SavedPassword
	var pswrds pswrd.Passwords
	display := bubbleT.GetInputDisplay
	//var masterPasswordToSave pswrd.SavedPassword
	//passwordCorrect:=false//,done:=false
	config := pswrd.CharmConfig{
		DatabaseName: "password-db",
		AccountName:  "Landon",
		AccountKey:   "my-passwords",
	}
	pswrds.Config = config
	//generating new encryption key

	db, err := kv.OpenWithDefaults(config.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//reading file to get all the passwords including old key and master password through charm
	data, newAcct := stor.ReadDataCharm(config, db)

	if newAcct {
		pswrds.Pswrds = append(pswrds.Pswrds, pswrd.SavedPassword{Website: "", EncryptedPswrd: ""})
		display = bubbleT.NewAcctDisplay
	} else {
		//data exists
		err := yaml.Unmarshal(data, &pswrds)
		if err != nil {
			fmt.Printf("Error while unMarshaling. %v", err)
		}
		//got them now need to decrypt them
		for i, password := range pswrds.Pswrds {
			//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
			pswrds.Pswrds[i].EncryptedPswrd = pswrd.Decrypt(password.EncryptedPswrd, len(pswrds.Pswrds[0].Website))

		}

		//prepping for end making new key and saving it
		//should be the new key instead of the old one
		//fmt.Println(pswrds.Pswrds[0].EncryptedPswrd)
	}
	rand.Seed(time.Now().Unix())
	newKey = rand.Int() % LIMIT
	pswrds.Pswrds[0].Website = pswrd.GeneratePassword(newKey)
	// db, err := kv.OpenWithDefaults(config.DatabaseName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// // Fetch updates and easily define your own syncing strategy
	// if err := db.Sync(); err != nil {
	// 	log.Fatal(err)
	// }

	// // Quickly get a value
	// charmData, err := db.Get([]byte(config.AccountKey))
	// if err != nil {
	// 	panic(err)
	// }
	// if len(charmData)==0{
	// 	panic("no charm data")//need to go to new user
	// }
	//fmt.Println(string(charmData)) //comes out in bytes

	//use the below code to get the data out of the yaml file, it's still saved their
	//reading from yaml file

	//getting the data, parsing it, then decrypting it to read

	//---------------------------------------------------------------BubbleTea code below----------------------------------------------------------------------------------------\\

	//fmt.Println(bubbleT.DisplayType(display))
	//running bubble tea wich gives really nice terminal interface, see extra file for details
	var allText string
	for _, pswrd := range pswrds.Pswrds {
		allText = allText + "\n" + pswrd.Website + " : " + pswrd.EncryptedPswrd
	}
	//stor.WriteData("testFile.txt", allText) allText would be everything in the file,
	p := tea.NewProgram(bubbleT.InitialModel(pswrds.Pswrds, db, pswrds.Config, display))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
