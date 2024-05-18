package S3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
)

func UploadAudioFile(filename string) error {
	// Замените these на ваши учетные данные AWS
	awsAccessKeyID := os.Getenv("S3_API_KEY")
	awsSecretAccessKey := os.Getenv("S3_SECRET_KEY")
	awsRegion := "us-east-1"

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	}))

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
		Bucket:      aws.String("discord-audio-records"),
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
