package passwordStorage //remake as file interaction/password storage

import (
	"fmt"
	"os"
	"os/exec"
	pswrd "password/passwordFunctions"
	"runtime"
)

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
}

// read data from a file, just reads it in and stores it all together
func ReadData(fileName string) (string, error) {
	//first pulls from github to ensure the file is up to date, using fetch and checkout to only 'pull' the .bin file instead of everything
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
	var allText = "correct\n"
	for _, pswrd := range pswrds {
		allText = allText + pswrd.Website + ": " + pswrd.EncryptedPswrd + "\n"
	}
	data := []byte(allText)
	encryptedData := pswrd.EncryptAes(key, data)       //encrypting all the data at once
	err := os.WriteFile(fileName, encryptedData, 0644) //0644 specifies the file is readable and writeable by the owner and readable by everyone else
	if err != nil {
		fmt.Println(err)
		return
	}
	//after everything is written, automatically commits the changes to the file so that it is updated on all devices
	runCommand("git", "add", fileName)
	runCommand("git", "commit", "-m", "Used password manager.")
	runCommand("git", "push", "--force") //forces the push to accept the local changes
}
