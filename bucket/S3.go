package bucket

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CreateSession(config *aws.Config) (*session.Session, error) {
	return session.NewSession(config)
}

func CreateBucket(svc *s3.S3, name string) error {
	if _, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(name)}); err != nil {
		return err
	}
	return nil
}

func UploadFile(sess *session.Session, path string, bucket string) error {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Error on open file")
	}
	defer file.Close()
	uploader := s3manager.NewUploader(sess)
	if _, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		Body:   file,
	}); err != nil {
		return err
	}

	return nil
}

func DownloadFile(sess *session.Session, filename string, bucket string) (file *os.File, err error) {
	item, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Unable to open file %q, %v", item.Name(), err)
	}
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(item,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		log.Fatal("Error on download file")
		return nil, err
	}
	fmt.Println("file: ", item.Name(), numBytes)
	return item, nil
}

func ListBuckets(svc *s3.S3) error {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		log.Fatal("Unable to list buckets", err)
		return err
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
	return nil
}

func DeleteBucket(svs *s3.S3, bucket string) error {
	_, err := svs.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(bucket)})

	if err != nil {
		return err
	}

	return nil
}
