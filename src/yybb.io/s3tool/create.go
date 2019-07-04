package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	default_contxt = buildDefaultContext()
	sourceFile string

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create s3 object",
		Run: func(cmd *cobra.Command, args []string) {
			err:= createOject()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}
			log.Printf("Successfully created bucket %s and uploaded data with key %s\n", default_contxt.Bucket, default_contxt.Key)
		},
	}

)
type CreateContxt struct {
	Endpoints string
	RegionId string
	Bucket string
	Ak string
	SK string
	Key string
}

func buildDefaultContext() CreateContxt {
	return CreateContxt{
		Endpoints:os.Getenv("S3_Endpoints"),
		RegionId: os.Getenv("S3_RegionId"),
		Bucket:os.Getenv("S3_Bucket"),
		Ak: os.Getenv("S3_AK"),
		SK: os.Getenv("S3_SK"),
	}
}
func createOject()  error {
	if sourceFile =="" {
		return fmt.Errorf("invalid source file")
	}
	inputFile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer inputFile.Close()


	bucket := default_contxt.Bucket
	key := default_contxt.Key
	creds := credentials.NewStaticCredentials(default_contxt.Ak, default_contxt.SK, "")
	_, err = creds.Get()
	if err != nil {
		return err
	}

	config := &aws.Config{
		Region:      aws.String(default_contxt.RegionId),
		Endpoint:    aws.String(default_contxt.Endpoints),
		DisableSSL:  aws.Bool(false),
		Credentials: creds,
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucket}); err != nil {

		return err
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:   inputFile,
		Bucket: &bucket,
		Key:    &key,
	})

	return err
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVarP(&default_contxt.Key, "key", "k","", "s3 store key")
	createCmd.PersistentFlags().StringVarP(&default_contxt.Bucket, "bucket", "b", "","s3 store buket")
	createCmd.PersistentFlags().StringVarP(&sourceFile, "file", "f", "","s3 store src file")
}
