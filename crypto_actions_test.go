package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

var testKey = []byte("1234123412341234")

func TestEncryptDecrypt(t *testing.T) {
	message := []byte("test message")

	encrypted, _ := encrypt(testKey, message)
	decrypted, _ := decrypt(testKey, encrypted)

	if string(decrypted) != string(message) {
		t.Error("Couldn't decrypt encrypted message correctly")
	}
}

func TestEncryptShouldFailWithBadKey(t *testing.T) {
	_, err := encrypt([]byte("bad key"), []byte("adsf"))
	if err == nil {
		t.Error("Expected an error, but didn't recive one")
	}
}

func TestDecryptShouldFailWithBadKey(t *testing.T) {
	_, err := decrypt([]byte("bad key"), []byte("adsf"))
	if err == nil {
		t.Error("Expected an error, but didn't recive one")
	}
}

func TestEncryptPostFlagParsing(t *testing.T) {
	infile := "testdata/encrypted"
	outfile := "out-file"

	fs := flag.NewFlagSet("name", flag.ExitOnError)
	fs.Parse([]string{infile, "out-file"})

	encryptFlagPostParse(fs)

	if encryptInFilenameArg != infile {
		t.Errorf("Got %s for input file, but expected %s", encryptInFilenameArg, infile)
	}

	if encryptOutFilenameArg != outfile {
		t.Errorf("Got %s for output file, but expected %s", encryptOutFilenameArg, outfile)
	}
}

func TestDecryptPostFlagParsing(t *testing.T) {
	infile := "testdata/encrypted"
	outfile := "out-file"

	fs := flag.NewFlagSet("name", flag.ExitOnError)
	fs.Parse([]string{infile, "out-file"})

	decryptFlagPostParse(fs)

	if decryptInFilenameArg != infile {
		t.Errorf("Got %s for input file, but expected %s", decryptInFilenameArg, infile)
	}

	if decryptOutFilenameArg != outfile {
		t.Errorf("Got %s for output file, but expected %s", decryptOutFilenameArg, outfile)
	}
}

func TestEncryptAction(t *testing.T) {
	outfile := "testdata/test_encrypt_action"

	// since we are generating a file make sure it doesn't already exist
	_, err := os.Stat(outfile)
	if err == nil {
		os.Remove(outfile)
	}

	encryptKeyFlag = string(testKey)
	encryptInFilenameArg = "testdata/plain"
	encryptOutFilenameArg = outfile

	encryptAction()

	contents, _ := ioutil.ReadFile(outfile)
	os.Remove(outfile)
	if len(contents) < 1 {
		t.Error("No contents where written to the outfile when encrypting")
	}
}

func TestDecryptAction(t *testing.T) {
	outfile := "testdata/test_decrypt_action"

	// since we are generating a file make sure it doesn't already exist
	_, err := os.Stat(outfile)
	if err == nil {
		os.Remove(outfile)
	}

	decryptKeyFlag = string(testKey)
	decryptInFilenameArg = "testdata/encrypted"
	decryptOutFilenameArg = outfile

	decryptAction()

	contents, _ := ioutil.ReadFile(outfile)
	os.Remove(outfile)
	message := string(contents)
	if message != "This is a test file" {
		t.Errorf("Expected This is a test file but got %s", message)
	}
}
