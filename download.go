package main

import (
	"errors"
	"flag"
	"github.com/robmerrell/gosecret/vendor/github.com/kr/s3"
	"github.com/robmerrell/gosecret/vendor/github.com/kr/s3/s3util"
	"io"
	"net/http"
	"os"
)

// flags and args
var downloadBucketNameFlag string
var downloadAccessKeyFlag string
var downloadSecretKeyFlag string
var downloadFilenameArg string
var downloadDestinationFilenameArg string

var downloadDoc = `
Usage: download [options] file [destination file]

Download a file from an s3 bucket.
If a destination file isn't specified the downloaded file will be named the same as the soruce file.
`

func downloadAction() error {
	// make sure that we have all of the required data
	if downloadFilenameArg == "" {
		return errors.New("Please provide a valid filename to download")
	}
	if downloadDestinationFilenameArg == "" {
		downloadDestinationFilenameArg = downloadFilenameArg
	}
	if downloadBucketNameFlag == "" {
		return errors.New("Please provide an S3 bucket name with --bucket or $GOSECRET_BUCKET")
	}
	if downloadAccessKeyFlag == "" {
		return errors.New("Please provide an AWS access key with --access-key or $GOSECRET_ACCESS_KEY")
	}
	if downloadSecretKeyFlag == "" {
		return errors.New("Please provide an AWS secrety key with --secret-key or $GOSECRET_SECRET_KEY")
	}

	// create the config needed for the downloader
	config := &s3util.Config{
		Keys: &s3.Keys{
			AccessKey: downloadAccessKeyFlag,
			SecretKey: downloadSecretKeyFlag,
		},
		Service: s3.DefaultService,
	}

	return download(downloadBucketNameFlag, downloadFilenameArg, downloadDestinationFilenameArg, config)
}

// downloadFlagInit initializes the flagset for the download command
func downloadFlagInit(fs *flag.FlagSet) {
	defaultBucket := os.Getenv("GOSECRET_BUCKET")
	fs.StringVar(&downloadBucketNameFlag, "bucket", defaultBucket, "S3 bucket to download from. Defaults to value in $GOSECRET_BUCKET")

	defaultAccessKey := os.Getenv("GOSECRET_ACCESS_KEY")
	fs.StringVar(&downloadAccessKeyFlag, "access-key", defaultAccessKey, "S3 Access Key. Defaults to value in $GOSECRET_ACCESS_KEY")

	defaultSecretKey := os.Getenv("GOSECRET_SECRET_KEY")
	fs.StringVar(&downloadSecretKeyFlag, "secret-key", defaultSecretKey, "S3 Secret Key. Defaults to value in $GOSECRET_SECRET_KEY")
}

// downloadFlagPostParse sets the downloadable filename from the arguments provided by the flagset
func downloadFlagPostParse(fs *flag.FlagSet) {
	if filename := fs.Arg(0); filename != "" {
		downloadFilenameArg = filename
	}

	if destFilename := fs.Arg(1); destFilename != "" {
		downloadDestinationFilenameArg = destFilename
	}
}

// download downloads a file from an s3 bucket.
func download(bucket, sourceFile, destFile string, config *s3util.Config) error {
	headers := http.Header{}
	headers.Add("x-amz-acl", "private")
	s3File, err := s3util.Open(generateS3Url(bucket, sourceFile), config)
	if err != nil {
		return err
	}
	defer s3File.Close()

	// open the local file to download to
	localFile, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// copy the file
	_, err = io.Copy(localFile, s3File)
	return err
}
