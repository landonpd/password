package passwordFunctions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	badRand "math/rand"
	"strings"
	//"os"
	//"io/ioutil"
	// "crypto/aes"
)

// not needed
// type CharmConfig struct {
// 	DatabaseName string
// 	AccountName  string
// 	AccountKey   string
// }

type SavedPassword struct {
	Website        string
	EncryptedPswrd string //`yaml:"password"`
}

// not needed either, just gonna be savedPassword slice
// type Passwords struct {
// 	//master password here maybe
// 	//maybe key here, probably eventually but not now
// 	Pswrds []SavedPassword // `yaml:"Passwords"`
// 	Config CharmConfig
// }

const CHARACTERS = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%&*"
const NUM_CHAR = len(CHARACTERS)

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

// randomly generates passwords for the given length using the characters string and randlonly selecting from it
func GeneratePassword(length int) string {
	newPassword := ""
	for i := 0; i < length; i++ {
		newPassword += string(CHARACTERS[badRand.Intn(NUM_CHAR)])
	}
	return newPassword
}

// hashes a byte slice using sha256 and returns the 32 bytes of the result
func HashKey(key []byte) []byte {
	//hashes the key
	h := sha256.New() //new sha256 object
	h.Write(key)      //actually hashes
	return h.Sum(nil) //does the last step of the checksum
}

// encrypts text using aes, counter stream mode and the key
func EncryptAes(key, text []byte) []byte {

	//fmt.Printf("the type of this %v\n", len(hashedKey)) //h.Sum(nil) returns [32]byte of the hash

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	result := make([]byte, aes.BlockSize+len(text)) //allocating space for encryption, blocksize is added to make sure the thing is the right length

	iv := result[:aes.BlockSize]                      //makes the iv the right size
	io.ReadFull(rand.Reader, iv)                      //randomly generates iv
	stream := cipher.NewCTR(c, iv)                    //makes a CTR stream object using the aes cipher and the iv, New CFBEncrypter(c,iv)
	stream.XORKeyStream(result[aes.BlockSize:], text) //xors each byte in text with a byte from the stream cipher and stores it in result, preserves iv at the front
	return result

}

// decrypts text using aes, counter stream mode and the key
func DecryptAes(key, cipherTxt []byte) []byte {
	p, err := aes.NewCipher(key) //makes an aes cipher object
	if err != nil {              //error checking
		fmt.Println(err)
		return nil
	}
	if len(cipherTxt) < aes.BlockSize { //making sure the ciphertext is long enough, it must be at least as long as the block size, the iv is in the first iv
		fmt.Println("ciphertext is too short")
		return nil
	}
	iv := cipherTxt[:aes.BlockSize]           //getting the iv from the beginning of the ciphertext
	realCipher := cipherTxt[aes.BlockSize:]   //retrieving the actual ciphertext, without the iv
	plainTxt := make([]byte, len(realCipher)) //allocating space for plainTxt
	stream := cipher.NewCTR(p, iv)            //makes a new stream object
	stream.XORKeyStream(plainTxt, realCipher) // xors result of ctr stream with realCipher and stores in plainTxt
	return plainTxt

}

// this takes the master password and the ciphertext from the file. Checks if the first line of the ciphertext is the word correct
// this will be used to check if the master password is correct
func CheckMasterPswrd(pswrd, cipherTxt string) bool {
	key := HashKey([]byte(pswrd))
	plainTxt := string(DecryptAes(key, []byte(cipherTxt)))
	pswrds := strings.Split(plainTxt, "\n")
	return pswrds[0] == "correct"
}

// decrypts the given password
// going to redo see encrypt

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
