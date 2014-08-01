package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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
		t.Errorf("Got %s for filename, but expected %s", uploadFilenameArg, filename)
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

func TestDownload(t *testing.T) {
	downloadRes, _ := ioutil.ReadFile("testdata/download_res")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(downloadRes))
	}))
	defer server.Close()

	replaceUrl(server.URL+"/%s/%s", &s3hostFmt, func() {
		// make sure the file doesn't already exist
		testfile := "test_download_func"
		_, err := os.Stat(testfile)
		if err == nil {
			os.Remove(testfile)
		}

		err = download("testbucket", testfile, nil)
		if err != nil {
			t.Errorf("Couldn't download file: %s", err)
		}

		downloadedFile, _ := ioutil.ReadFile(testfile)
		if string(downloadedFile) != string(downloadRes) {
			t.Error("Downloaded file doesn't match the test download file")
		}
		os.Remove(testfile)
	})
}

func TestDownloadFlagPostParse(t *testing.T) {
	filename := "plain"

	fs := flag.NewFlagSet("name", flag.ExitOnError)
	fs.Parse([]string{filename})

	downloadFlagPostParse(fs)

	if downloadFilenameArg != filename {
		t.Errorf("Got %s for filename, but expected %s", downloadFilenameArg, filename)
	}
}

func TestDownloadAction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downloadRes, _ := ioutil.ReadFile("testdata/download_res")
		fmt.Fprint(w, string(downloadRes))
	}))
	defer server.Close()

	filename := "test_download_action"
	downloadBucketNameFlag = "testbucket"
	downloadAccessKeyFlag = "testaccess"
	downloadSecretKeyFlag = "testsecret"
	downloadFilenameArg = filename

	replaceUrl(server.URL+"/%s/%s", &s3hostFmt, func() {
		err := downloadAction()
		if err != nil {
			t.Errorf("Couldn't download file: %s", err)
		}

		os.Remove(filename)
	})
}
