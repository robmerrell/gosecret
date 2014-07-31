package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func replaceUrl(newUrl string, refVar *string, testFunc func()) {
	oldUrl := *refVar
	*refVar = newUrl

	testFunc()

	*refVar = oldUrl
}

func TestUpload(t *testing.T) {
	multiPartCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("ETag", "faketag")
		uploadRes, _ := ioutil.ReadFile("testdata/upload_res")
		multiPartCount++
		fmt.Fprint(w, string(uploadRes))
	}))
	defer server.Close()

	replaceUrl(server.URL+"/%s/%s", &s3hostFmt, func() {
		err := upload("testbucket", "testdata/plain", nil)
		if err != nil {
			t.Errorf("Couldn't upload file: %s", err)
		}
		if multiPartCount < 1 {
			t.Errorf("No file was uploaded")
		}
	})
}

func TestUploadFlagPostParse(t *testing.T) {
	filename := "testdata/plain"

	fs := flag.NewFlagSet("name", flag.ExitOnError)
	fs.Parse([]string{filename})

	uploadFlagPostParse(fs)

	if uploadFilenameArg != filename {
		t.Errorf("Got %s for filename, but expected %s", encryptInFilenameArg, filename)
	}
}

func TestUploadAction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("ETag", "faketag")
		uploadRes, _ := ioutil.ReadFile("testdata/upload_res")
		fmt.Fprint(w, string(uploadRes))
	}))
	defer server.Close()

	filename := "testdata/plain"
	uploadBucketNameFlag = "testbucket"
	uploadAccessKeyFlag = "testaccess"
	uploadSecretKeyFlag = "testsecret"
	uploadFilenameArg = filename

	replaceUrl(server.URL+"/%s/%s", &s3hostFmt, func() {
		err := uploadAction()
		if err != nil {
			t.Errorf("Couldn't upload file: %s", err)
		}
	})
}

func TestGenerateS3Url(t *testing.T) {
	url := generateS3Url("bucket", "file")
	validUrl := "https://bucket.s3.amazonaws.com/file"
	if url != validUrl {
		t.Errorf("Didn't generate a correct URL. Expected %s, got %s", validUrl, url)
	}
}
