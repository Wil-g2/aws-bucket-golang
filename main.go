package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/wil-g2/aws-bucket/bucket"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Erro to load env file", err)
	}
}

func main() {
	S3Endpoint := os.Getenv("ENDPOINT")
	sess, err := bucket.CreateSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("S3_ID"), os.Getenv("S3_SECRET"), os.Getenv("S3_TOKEN")),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.UsWest2RegionID),
		Endpoint:         aws.String(S3Endpoint),
	},
	)

	if err != nil {
		log.Fatal(err)
	}
	// Create S3 service client
	svc := s3.New(sess)

	err = bucket.CreateBucket(svc, "test")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Erro to create a bucket")
	}

	fileBucket, err := bucket.DownloadFile(sess, "image.png", "testgo")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Erro to create a file")
	}

	fmt.Println(*fileBucket)
	bucket.ListBuckets(svc)

	if err = bucket.DeleteBucket(svc, "test"); err != nil {
		fmt.Println(err)
	}
}
