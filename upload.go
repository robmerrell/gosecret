package main

import (
	"errors"
	"flag"
	"github.com/robmerrell/gosecret/vendor/github.com/kr/s3"
	"github.com/robmerrell/gosecret/vendor/github.com/kr/s3/s3util"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// flags and args
var uploadBucketNameFlag string
var uploadAccessKeyFlag string
var uploadSecretKeyFlag string
var uploadFilenameArg string

var uploadDoc = `
Usage: upload [options] file

Upload a file to an s3 bucket
`

func uploadAction() error {
	// make sure that we have all of the required data
	if uploadFilenameArg == "" {
		return errors.New("Please provide a valid filename to upload")
	}
	if uploadBucketNameFlag == "" {
		return errors.New("Please provide an S3 bucket name with --bucket or $GOSECRET_BUCKET")
	}
	if uploadAccessKeyFlag == "" {
		return errors.New("Please provide an AWS access key with --access-key or $GOSECRET_ACCESS_KEY")
	}
	if uploadSecretKeyFlag == "" {
		return errors.New("Please provide an AWS secrety key with --secret-key or $GOSECRET_SECRET_KEY")
	}

	// create the config needed for the uploader
	config := &s3util.Config{
		Keys: &s3.Keys{
			AccessKey: uploadAccessKeyFlag,
			SecretKey: uploadSecretKeyFlag,
		},
		Service: s3.DefaultService,
	}

	return upload(uploadBucketNameFlag, uploadFilenameArg, config)
}

// uploadFlagInit initializes the flagset for the upload command
func uploadFlagInit(fs *flag.FlagSet) {
	defaultBucket := os.Getenv("GOSECRET_BUCKET")
	fs.StringVar(&uploadBucketNameFlag, "bucket", defaultBucket, "S3 bucket to upload into. Defaults to value in $GOSECRET_BUCKET")

	defaultAccessKey := os.Getenv("GOSECRET_ACCESS_KEY")
	fs.StringVar(&uploadAccessKeyFlag, "access-key", defaultAccessKey, "S3 Access Key. Defaults to value in $GOSECRET_ACCESS_KEY")

	defaultSecretKey := os.Getenv("GOSECRET_SECRET_KEY")
	fs.StringVar(&uploadSecretKeyFlag, "secret-key", defaultSecretKey, "S3 Secret Key. Defaults to value in $GOSECRET_SECRET_KEY")
}

// uploadFlagPostParse sets the uploadable filename from the arguments provided by the flagset
func uploadFlagPostParse(fs *flag.FlagSet) {
	// make sure the input file is reachable
	if filename := fs.Arg(0); filename != "" {
		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			uploadFilenameArg = filename
		}
	}
}

// upload uploads a file to an s3 bucket.
func upload(bucket, file string, config *s3util.Config) error {
	// open the local file to upload
	localFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer localFile.Close()

	headers := http.Header{}
	headers.Add("x-amz-acl", "private")
	s3File, err := s3util.Create(generateS3Url(bucket, filepath.Base(file)), headers, config)
	if err != nil {
		return err
	}
	defer s3File.Close()

	// copy the file
	_, err = io.Copy(s3File, localFile)
	return err
}
