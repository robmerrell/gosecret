package main

import (
	"fmt"
	"github.com/robmerrell/gosecret/vendor/github.com/robmerrell/comandante"
	"os"
)

func main() {
	bin := comandante.New("gosecret", "Manage encrypted files in an S3 bucket")
	bin.IncludeHelp()

	// decrypt
	decryptCmd := comandante.NewCommand("decrypt", "Decrypt a file", decryptAction)
	decryptCmd.Documentation = decryptDoc
	decryptCmd.FlagInit = decryptFlagInit
	decryptCmd.FlagPostParse = decryptFlagPostParse
	bin.RegisterCommand(decryptCmd)

	// encrypt
	encryptCmd := comandante.NewCommand("encrypt", "Encrypt a file", encryptAction)
	encryptCmd.Documentation = encryptDoc
	encryptCmd.FlagInit = encryptFlagInit
	encryptCmd.FlagPostParse = encryptFlagPostParse
	bin.RegisterCommand(encryptCmd)

	// download
	downloadCmd := comandante.NewCommand("download", "Download a file", downloadAction)
	downloadCmd.Documentation = downloadDoc
	downloadCmd.FlagInit = downloadFlagInit
	downloadCmd.FlagPostParse = downloadFlagPostParse
	bin.RegisterCommand(downloadCmd)

	// upload
	uploadCmd := comandante.NewCommand("upload", "Upload a file", uploadAction)
	uploadCmd.Documentation = uploadDoc
	uploadCmd.FlagInit = uploadFlagInit
	uploadCmd.FlagPostParse = uploadFlagPostParse
	bin.RegisterCommand(uploadCmd)

	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
