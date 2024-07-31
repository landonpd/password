package main

import (
	"fmt"
	"os"
	stor "password/passwordStorage"
	uI "password/userInterface"

	//"Users/landondixon/src/goCode/passwordManager/passwordFunctions"  //github.com/landonpd/password/

	tea "github.com/charmbracelet/bubbletea" //if my code is on github I can just include it from their, might be a sneaky little way to get my file there
	// "crypto/aes" I will be useing aes
	// "crypto/rand"
)

// SavedPassword represents a password entry

func main() {
	//variables and constants
	const LIMIT = 100000
	const FILENAME = "passwordFile.txt"
	//var savedMasterPassword string //, stringtoWrite,inputtedmasterPassword, string
	//var newKey int //, oldKey, passwordCount int //choice, createPasswordChoice,
	//var passwords []pswrd.SavedPassword

	// var pswrds []pswrd.SavedPassword
	display := uI.GetInputDisplay
	//var masterPasswordToSave pswrd.SavedPassword
	//passwordCorrect:=false//,done:=false
	// config := pswrd.CharmConfig{
	// 	DatabaseName: "password-db",
	// 	AccountName:  "Landon",
	// 	AccountKey:   "my-passwords",
	// }
	// pswrds.Config = config
	//generating new encryption key

	// db, err := kv.OpenWithDefaults(config.DatabaseName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	//reading file to get all the passwords including old key and master password through charm
	// data, newAcct := stor.ReadDataCharm(config, db)
	//read data without charm
	fileData, err := stor.ReadData(FILENAME)

	if err != nil {
		fmt.Println("file doesn't exist")
		//maybe make the file, or have it ship with a blank file, this will have to change later when it is an api through google doc
	} else if len(fileData) == 0 {
		//file is empty
		display = uI.NewAcctDisplay
		//writes correct to the file
		// data := []byte("correct")
		// err := os.WriteFile(FILENAME, data, 0644) //0644 specifies the file is readable and writeable by the owner and readable by everyone else
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}
	//determines if it is a new user or not, will need to update how this works, will need to be a part of the file
	// if newAcct {
	// 	pswrds.Pswrds = append(pswrds.Pswrds, pswrd.SavedPassword{Website: "", EncryptedPswrd: ""})
	// 	display = uI.NewAcctDisplay
	// } else {
	// 	//data exists
	// 	err := yaml.Unmarshal(data, &pswrds)
	// 	if err != nil {
	// 		fmt.Printf("Error while unMarshaling. %v", err)
	// 	}
	// 	//got them now need to decrypt them
	// 	for i, password := range pswrds.Pswrds {
	// 		//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
	// 		pswrds.Pswrds[i].EncryptedPswrd = pswrd.Decrypt(password.EncryptedPswrd, len(pswrds.Pswrds[0].Website))

	// 	}

	// 	//prepping for end making new key and saving it
	// 	//should be the new key instead of the old one
	// 	//fmt.Println(pswrds.Pswrds[0].EncryptedPswrd)
	// }

	////splitting data from file into website and password pairs
	//fmt.Println(fileData)
	// fileDataLst := strings.Split(fileData, "\n")
	// fileDataLst = fileDataLst[:len(fileDataLst)-1]
	// for _, passwords := range fileDataLst {
	// 	// fmt.Println(passwords)
	// 	savedPassword := strings.Split(passwords, " : ")
	// 	tempPsswrd := pswrd.SavedPassword{Website: savedPassword[0], EncryptedPswrd: savedPassword[1]}
	// 	pswrds = append(pswrds, tempPsswrd)
	// }

	////got them now need to decrypt them
	// for i, password := range pswrds {
	// 	//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
	// 	pswrds[i].EncryptedPswrd = pswrd.Decrypt(password.EncryptedPswrd, len(pswrds[0].Website))

	// }
	//fmt.Println(pswrds.Pswrds[len(pswrds.Pswrds)-1].Website)
	//fmt.Println(fileDataLst[len(fileDataLst)-1])

	//possibly messing things up, change the key for the encryption but never put in the master key so old master key, new encryption key maybe
	// rand.Seed(time.Now().Unix())
	// newKey = rand.Int() % LIMIT
	// pswrds[0].Website = pswrd.GeneratePassword(newKey)
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

	p := tea.NewProgram(uI.InitialModel(fileData, display)) //db, pswrds.Config
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
