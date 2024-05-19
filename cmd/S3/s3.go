package S3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"time"
)

var (
	loadEnvError       = godotenv.Load()
	awsAccessKeyID     = os.Getenv("S3_API_KEY")
	awsSecretAccessKey = os.Getenv("S3_SECRET_KEY")
	awsRegion          = os.Getenv("S3_REGION")
	bucketName         = os.Getenv("S3_BUCKET_NAME")

	sess = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	}))
)

func init() {
	if loadEnvError != nil {
		fmt.Println("Error loading .env file:", loadEnvError)

	}
}

func UploadAudioFile(filename string) error {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}
	defer f.Close()

	// Используйте базовое имя файла для ключа S3
	key := filepath.Base(filename)

	// Определите MIME-тип для MP3 файла
	mimeType := "audio/mpeg"

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
	return nil
}

func GetAllBucketObjects() []*s3.Object {
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		fmt.Println(err)
	}
	objects := resp.Contents
	return objects
}

func GetFileLink(filename string) string {
	svc := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}

	req, _ := svc.GetObjectRequest(params)

	url, err := req.Presign(15 * time.Minute) // Set link expiration time
	if err != nil {
		fmt.Println("[AWS GET LINK]:", params, err)
	}

	return url
}
