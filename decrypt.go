package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"flag"
	"io/ioutil"
	"os"
)

// flags
var decryptKeyFlag string
var decryptInFilenameArg string
var decryptOutFilenameArg string

var decryptDoc = `
Usage: decrypt [options] in-file out-file

Decrypt an input file using a keypair and write the results to an output file
`

// decryptAction is the action invoked by comandante
func decryptAction() error {
	// make sure filenames are set
	if decryptInFilenameArg == "" {
		return errors.New("Please provide a valid input file")
	}
	if decryptOutFilenameArg == "" {
		return errors.New("Please provide a valid output file")
	}

	// read the input file
	contents, err := ioutil.ReadFile(decryptInFilenameArg)
	if err != nil {
		return err
	}

	// decrypt and write to the outfile
	decrypted, err := decrypt([]byte(decryptKeyFlag), []byte(contents))
	if err != nil {
		return err
	}
	return ioutil.WriteFile(decryptOutFilenameArg, decrypted, 0644)
}

// decryptFlagInit initializes the flagset for the decrypt command
func decryptFlagInit(fs *flag.FlagSet) {
	defaultKey := os.Getenv("GOSECRET_KEY")
	fs.StringVar(&decryptKeyFlag, "key", defaultKey, "A 16, 24 or 32 byte key to use for decryption. Defaults to value in $GOSECRET_KEY")
}

// decryptFlagPostParse sets filenames from the arguments provided by the flagset
func decryptFlagPostParse(fs *flag.FlagSet) {
	// make sure the input file is reachable
	if filename := fs.Arg(0); filename != "" {
		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			decryptInFilenameArg = filename
		}
	}

	decryptOutFilenameArg = fs.Arg(1)
}

// decrypt decrypts a message using a given key
func decrypt(key, contents []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	if len(contents) < aes.BlockSize {
		return []byte{}, errors.New("File to decrypt is too small")
	}

	iv := contents[:aes.BlockSize]
	// decrypted := make([]byte, len(contents)-aes.BlockSize)
	decrypted := contents[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decrypted, decrypted)

	return decrypted, nil
}
