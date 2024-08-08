package passwordStorage //remake as file interaction/password storage

import (
	pswrd "password/passwordFunctions"
	//bubbleT "Users/landondixon/src/goCode/passwordManager/bubbleTeaFunctions"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	//"io/ioutil"
)

// reads data from charm cloud9L
// func ReadDataCharm(c pswrd.CharmConfig, db *kv.KV) ([]byte, bool) {
// 	//var charmData []byte
// 	display := false

// 	// Fetch updates and easily define your own syncing strategy
// 	if err := db.Sync(); err != nil {
// 		log.Fatal(err)
// 	}
// 	// _, err := db.Get(context.Background(), key)
// 	// if err == charmstore.ErrNotFound {
// 	// 	// Key does not exist
// 	// 	return false, nil
// 	// if len(charmData)==0{ //find out how to check if a key works

// 	// }else{
// 	// Quickly get a value
// 	charmData, err := db.Get([]byte(c.AccountKey)) //charmData, err := ioutil.ReadFile("passwords.yaml")

// 	if charmData == nil { //don't know what this error actually is
// 		// Connection refused error
// 		display = true
// 	} else if err != nil {
// 		panic(err)
// 	}

// 	return charmData, display
// }

// will run a command on the command line
func runCommand(command string, args ...string) { //the dots make this a variadic function, any number of string arguments can be passed in after command, I can pass in a slice if I put ... after the slice name
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		Args := []string{"/c", command}
		Args = append(Args, args...)
		cmd = exec.Command("cmd", Args...) //using the dots after allows me to pass a slice into the variadic function
	case "linux", "darwin":
		cmd = exec.Command(command, args...)
	default:
		fmt.Println("Unsupported OS")
		return
	}
	_, err := cmd.CombinedOutput() //get output here if I need it
	if err != nil {
		fmt.Printf("Error running command '%s %s': %s\n", command, args, err)
		return
	}
	//fmt.Printf("Output of '%s':\n%s\n", command, string(output))
}

// read data from a file, just reads it in and stores it all together
func ReadData(fileName string) (string, error) {
	//first pulls from github to ensure the file is up to date
	//runCommand("git", "pull") //should just work hopefully, fingers crossed
	runCommand("git", "fetch")
	runCommand("git", "checkout", "origin/main", "--", fileName)
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("error opening %s: %s", fileName, err)
		return "", err
	}
	return string(data), nil
}

// writes the data to the given file, creates the file if it doesn't exist
func WritePasswords(key []byte, fileName string, pswrds []pswrd.SavedPassword) {
	// var pswrds []pswrd.SavedPassword
	// pswrds = passwords
	// pswrds.Config = c
	//not going to be needed anymore
	// for i, password := range pswrds {
	// 	//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
	// 	pswrds[i].EncryptedPswrd = pswrd.Encrypt(password.EncryptedPswrd, len(pswrds[0].Website))

	// }
	//look into what the heck yaml is cause it seems to simplify the code a bit
	// pswrd.DisplayPasswords(pswrds)
	var allText = "correct\n"
	for _, pswrd := range pswrds {
		allText = allText + pswrd.Website + ": " + pswrd.EncryptedPswrd + "\n"
	}
	// fmt.Println(allText)
	data := []byte(allText)
	encryptedData := pswrd.EncryptAes(key, data)
	err := os.WriteFile(fileName, encryptedData, 0644) //0644 specifies the file is readable and writeable by the owner and readable by everyone else
	if err != nil {
		fmt.Println(err)
		return
	}
	//after everything is written, automatically commits the changes to the file so that it is updated on all devices
	runCommand("git", "add", fileName)
	runCommand("git", "commit", "-m", "Used password manager.")
	runCommand("git", "push")
}

// writes to charm cloud
// func WritePasswordsCharm(passwords []pswrd.SavedPassword, Key int, db *kv.KV, c pswrd.CharmConfig) {

// 	var pswrds pswrd.Passwords
// 	pswrds.Pswrds = passwords
// 	pswrds.Config = c
// 	for i, password := range pswrds.Pswrds {
// 		//have to use pswrds.Pswrds[i] to actually update the password, possibly need to change in other places
// 		pswrds.Pswrds[i].EncryptedPswrd = pswrd.Encrypt(password.EncryptedPswrd, len(pswrds.Pswrds[0].Website))

// 	}

// 	yamlData, err := yaml.Marshal(&pswrds)

// 	if err != nil {
// 		fmt.Printf("Error while Marshaling. %v", err)
// 	}
// 	//fmt.Println(string(yamlData))

// 	// fmt.Println(c.AccountKey)
// 	if err := db.Set([]byte(c.AccountKey), yamlData); err != nil { //failing now for some reason
// 		log.Fatal(err)
// 	}

// 	if err := db.Sync(); err != nil {
// 		log.Fatal(err)
// 	}

// 	//fileName := "passwords.yaml"

// 	// err = ioutil.WriteFile(fileName, yamlData, 0644)
// 	// if err != nil {
// 	//     panic("Unable to write data into the file")
// 	// }
// }
