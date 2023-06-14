package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"log"
	"os"
	"path/filepath"
)

var tempFilePath = "/tmp/file"

// resizedFilePath must end with one of the known image file extensions due to the limits of "imaging" library.
// However, it does not affect the actual extension of the file.
var resizedFilePath = "/tmp/resized.jpeg"

func handler(ctx context.Context, s3Event events.S3Event) (err error) {
	s3EventRecord := s3Event.Records[0].S3
	sourceBucket := s3EventRecord.Bucket.Name
	sourceKey := s3EventRecord.Object.Key

	destinationBucket := os.Getenv("RESIZED_BUCKET_NAME")

	// Check if the file has the desired extension
	desiredExtensions := []string{".jpg", ".jpeg", ".png"} // todo consider extracting it by env var
	extension := filepath.Ext(sourceKey)
	if !contains(desiredExtensions, extension) {
		fmt.Printf("Skipping file %s as it doesn't have a valid extension", sourceKey)
		return
	}

	// Create a new AWS session
	sess := session.Must(session.NewSession())

	// Create an S3 downloader
	downloader := s3manager.NewDownloader(sess)

	// Copy the object to a local file
	file, err := os.Create(tempFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = downloader.DownloadWithContext(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(sourceBucket),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		return err
	}

	// Resize image
	resize()

	resized, err := os.Open(resizedFilePath)
	if err != nil {
		return err
	}
	defer resized.Close()

	// Upload the resized image to the destination bucket
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(destinationBucket),
		Key:    aws.String(sourceKey),
		Body:   resized,
	})
	if err != nil {
		return
	}

	fmt.Printf("File copied from %s/%s to %s/%s", sourceBucket, sourceKey, destinationBucket, sourceKey)

	return
}
func main() {
	lambda.Start(handler)
}

// Helper function to check if a slice contains a given string
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func resize() {
	fmt.Println("Resizing start")
	image, err := imaging.Open(tempFilePath)
	if err != nil {
		log.Fatalf("Error while loading image file: %v", err)
	}

	resized := imaging.Resize(image, 1600, 1200, imaging.Lanczos)

	err = imaging.Save(resized, resizedFilePath)
	if err != nil {
		log.Fatalf("Error while saving image file: %v", err)
	}
}
