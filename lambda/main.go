package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strings"
)

func handler(ctx context.Context, s3Event events.S3Event) (err error) {
	s3EventRecord := s3Event.Records[0].S3
	sourceBucket := s3EventRecord.Bucket.Name
	sourceKey := s3EventRecord.Object.Key

	destinationBucket := os.Getenv("RESIZED_BUCKET_NAME")

	fmt.Printf("Copy %s", sourceKey)

	// Create a new AWS session
	sess := session.Must(session.NewSession())

	// Create an S3 client
	s3client := s3.New(sess)

	// Copy the object to the destination bucket
	_, err = s3client.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		CopySource: aws.String(strings.Join([]string{sourceBucket, sourceKey}, "/")),
		Bucket:     aws.String(destinationBucket),
		Key:        aws.String(sourceKey),
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
