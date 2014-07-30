package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
)

// flags
var encryptKeyFlag string
var encryptInFilenameArg string
var encryptOutFilenameArg string

var encryptDoc = `
Usage: encrypt [options] in-file out-file

Encrypt an input file using a keypair and write the results to an output file
`

// encryptAction is the action invoked by comandante
func encryptAction() error {
	// make sure filenames are set
	if encryptInFilenameArg == "" {
		return errors.New("Please provide a valid input file")
	}
	if encryptOutFilenameArg == "" {
		return errors.New("Please provide a valid output file")
	}

	// read the input file
	contents, err := ioutil.ReadFile(encryptInFilenameArg)
	if err != nil {
		return err
	}

	// encrypt and write to the outfile
	encrypted, err := encrypt([]byte(encryptKeyFlag), []byte(contents))
	if err != nil {
		return err
	}
	return ioutil.WriteFile(encryptOutFilenameArg, encrypted, 0644)
}

// encryptFlagInit initializes the flagset for the encrypt command
func encryptFlagInit(fs *flag.FlagSet) {
	defaultKey := os.Getenv("GOSECRET_KEY")
	fs.StringVar(&encryptKeyFlag, "key", defaultKey, "A 16, 24 or 32 byte key to use for encryption. Defaults to value in $GOSECRET_KEY")
}

// encryptFlagPostParse sets filenames from the arguments provided by the flagset
func encryptFlagPostParse(fs *flag.FlagSet) {
	// make sure the input file is reachable
	if filename := fs.Arg(0); filename != "" {
		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			encryptInFilenameArg = filename
		}
	}

	encryptOutFilenameArg = fs.Arg(1)
}

// encrypt encryptes a message using a given key
func encrypt(key, contents []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	iv := contents[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(contents))
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], contents)

	return ciphertext, nil
}

// func decrypt(key, text []byte) []byte {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if len(text) < aes.BlockSize {
// 		panic("ciphertext too short")
// 	}
// 	iv := text[:aes.BlockSize]
// 	text = text[aes.BlockSize:]
// 	cfb := cipher.NewCFBDecrypter(block, iv)
// 	cfb.XORKeyStream(text, text)
// 	return decodeBase64(text)
// }
