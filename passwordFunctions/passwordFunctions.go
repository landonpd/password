package passwordFunctions

import (
	"fmt"
	"math/rand"
	//"os"
	//"io/ioutil"
)

type CharmConfig struct {
	DatabaseName string
	AccountName  string
	AccountKey   string
}

type SavedPassword struct {
	Website        string
	EncryptedPswrd string //`yaml:"password"`
}

type Passwords struct {
	//master password here maybe
	//maybe key here, probably eventually but not now
	Pswrds []SavedPassword // `yaml:"Passwords"`
	Config CharmConfig
}

const CHARACTERS = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%&*"
const NUM_CHAR = len(CHARACTERS)

// randomly generates passwords for the given length using the characters string and randlonly selecting from it
func GeneratePassword(length int) string {
	newPassword := ""
	for i := 0; i < length; i++ {
		newPassword += string(CHARACTERS[rand.Intn(NUM_CHAR)])
	}
	return newPassword
}

// encrypts the given password using a ceasur cipher.
// Going to redo, just need to put the whole thing together then encrypt the entire file.
func Encrypt(password string, key int) string {
	encryptedPassword := ""
	for _, char := range password {
		place := 0
		for j, ch := range CHARACTERS {
			if char == ch {
				place = j
				break
			}
		}
		encryptedPassword += string(CHARACTERS[(place+key)%NUM_CHAR])
	}

	return encryptedPassword
}

// decrypts the given password
// going to redo see encrypt
func Decrypt(encryptedPassword string, key int) string {
	decryptedPassword := ""
	adjustedKey := key % NUM_CHAR
	for _, char := range encryptedPassword {
		place := 0
		for j, ch := range CHARACTERS {
			if char == ch {
				place = j
				break
			}
		}

		adjustedPlace := place - adjustedKey
		if adjustedPlace < 0 {
			adjustedPlace = NUM_CHAR + adjustedPlace
		}
		decryptedPassword += string(CHARACTERS[adjustedPlace])
	}

	return decryptedPassword
}

// printys out the list of passwords, numbered, with website: password
func DisplayPasswords(passwords []SavedPassword) { //(passwords []SavedPassword, oldKey int)
	fmt.Println()
	for i := 1; i < len(passwords); i++ {
		// fmt.Printf("%d. %s: %s\n", i, passwords[i].Website, decrypt(passwords[i].password, oldKey))
		fmt.Printf("%d. %s: %s\n", i, passwords[i].Website, passwords[i].EncryptedPswrd)
	}
	fmt.Println()

}

// func WritePasswords(passwords []SavedPassword,Key int) {
// 	stringtoWrite:=""

// 	for _,password:= range passwords{
// 		stringtoWrite+=password.Website
// 		stringtoWrite+=" "
// 		stringtoWrite+=Encrypt(password.EncryptedPswrd,Key)
// 		stringtoWrite+=" "

// 	}
// 	os.WriteFile("passwordFile.txt",[]byte(stringtoWrite) , 0644)
// }
