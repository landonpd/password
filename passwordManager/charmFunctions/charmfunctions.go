package charmFunctions
import (
	pswrd "Users/landondixon/src/goCode/passwordManager/passwordFunctions"
	//bubbleT "Users/landondixon/src/goCode/passwordManager/bubbleTeaFunctions"
	"github.com/charmbracelet/charm/kv"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	
	
	//"io/ioutil"
)

func ReadData(c pswrd.CharmConfig,db *kv.KV)  ([]byte,bool){
	//var charmData []byte
	display:=false
	


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
	
	if charmData ==nil{ //don't know what this error actually is
		// Connection refused error
		display=true
	}else if err!=nil{
		panic(err)
	}
	
	return charmData,display
}


//writes to charm cloud
func WritePasswords(passwords []pswrd.SavedPassword,Key int,db *kv.KV,c pswrd.CharmConfig){

	
	var pswrds pswrd.Passwords
	pswrds.Pswrds=passwords
	pswrds.Config=c
	for i,password:=range pswrds.Pswrds {
		//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
		pswrds.Pswrds[i].EncryptedPswrd=pswrd.Encrypt(password.EncryptedPswrd,len(pswrds.Pswrds[0].Website))
		
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