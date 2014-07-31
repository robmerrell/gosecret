package main

import (
	"fmt"
)

var s3hostFmt = "https://%s.s3.amazonaws.com/%s"

// generateS3Url generates the URL required for the upload request to S3.
func generateS3Url(bucket, file string) string {
	return fmt.Sprintf(s3hostFmt, bucket, file)
}
