package s3cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
	"github.com/yylq/s3tool/pkg/common"
	"github.com/yylq/s3tool/pkg/utils"
)

var (
	default_contxt = common.BuildS3Context()
	default_type   = "application/octet-stream"
	sourceFile     string
	outputFile     string
	objtype        string
	max_age        int32

	objectCmd = &cobra.Command{
		Use:   "object",
		Short: "manager s3 object",
	}
	createObjectCmd = &cobra.Command{
		Use:   "create",
		Short: "create s3 object",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createOject()
		},
	}
	deleteObjectCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete s3 object",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteOject()
		},
	}
	saveOjectCmd = &cobra.Command{
		Use:   "get",
		Short: "get s3 object",
		RunE: func(cmd *cobra.Command, args []string) error {
			return saveOject()
		},
	}

	listObjectCmd = &cobra.Command{
		Use:   "list",
		Short: "list s3 object",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOject()
		},
	}

	ftype = map[string]string{
		".json": "application/json",
		".xml":  "text/xml",
		".txt":  "text/plain",
		".html": "text/html",
		".md":   "text/plain",
	}
)

func createOject() error {
	if sourceFile == "" {
		return fmt.Errorf("invalid source file")
	}

	inputFile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	newmd5 := utils.Md5sumFile(sourceFile)
	suffix := path.Ext(sourceFile)

	v, ok := ftype[suffix]
	if ok {
		objtype = v
	}

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
		DisableSSL:  aws.Bool(true),
		Credentials: creds,
	}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Printf("NewSession")
		return err
	}

	svc := s3.New(sess)
	//if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucket}); err != nil {
	//	log.Printf("WaitUntilBucketExists")
	//	return err
	//}
	cachecontrol := fmt.Sprintf("max-age=%d", max_age)
	putobj := &s3.PutObjectInput{
		Body:         inputFile,
		Bucket:       &bucket,
		Key:          &key,
		CacheControl: &cachecontrol,
		ContentType:  &objtype,
		Metadata:     map[string]*string{"md5": &newmd5},
	}

	_, err = svc.PutObject(putobj)
	if err != nil {
		return err
	}

	log.Printf("created:[ http://%s.%s/%s ]", bucket, default_contxt.Endpoints, key)
	return nil
}

func deleteOject() error {

	bucket := default_contxt.Bucket
	key := default_contxt.Key
	creds := credentials.NewStaticCredentials(default_contxt.Ak, default_contxt.SK, "")
	_, err := creds.Get()
	if err != nil {
		log.Printf("err:%v", err)
		return err
	}

	config := &aws.Config{
		Region:      aws.String(default_contxt.RegionId),
		Endpoint:    aws.String(default_contxt.Endpoints),
		DisableSSL:  aws.Bool(true),
		Credentials: creds,
	}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Printf("NewSession")
		return err
	}

	svc := s3.New(sess)
	//if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucket}); err != nil {
	//	log.Printf("WaitUntilBucketExists")
	//	return err
	//}

	obj := &s3.DeleteObjectInput{
		Bucket: &default_contxt.Bucket,
		Key:    &default_contxt.Key}

	_, err = svc.DeleteObject(obj)
	if err != nil {
		return err
	}
	log.Printf("delete http://%s.%s/%s", bucket, default_contxt.Endpoints, key)
	return nil
}

func listOject() error {

	bucket := default_contxt.Bucket

	creds := credentials.NewStaticCredentials(default_contxt.Ak, default_contxt.SK, "")
	_, err := creds.Get()
	if err != nil {
		log.Printf("err:%v", err)
		return err
	}

	config := &aws.Config{
		Region:      aws.String(default_contxt.RegionId),
		Endpoint:    aws.String(default_contxt.Endpoints),
		DisableSSL:  aws.Bool(true),
		Credentials: creds,
	}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Printf("NewSession")
		return err
	}

	svc := s3.New(sess)
	//if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucket}); err != nil {
	//	log.Printf("WaitUntilBucketExists")
	//	return err
	//}
	obj := &s3.ListObjectsInput{
		Bucket: &bucket,
	}

	out, err := svc.ListObjects(obj)
	for _, obj := range out.Contents {
		log.Printf("%s.%s/%s", bucket, default_contxt.Endpoints, *obj.Key)
	}
	return err
}

func saveOject() error {
	if outputFile == "" {
		return fmt.Errorf("invalid outputFile file")
	}
	err := utils.DownOject(outputFile, default_contxt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DownOject err:%v", err)
		return err
	}
	return nil
}

func init() {
	//rootCmd.AddCommand(objectCmd)
	objectCmd.AddCommand(createObjectCmd, deleteObjectCmd, listObjectCmd, saveOjectCmd)
	objectCmd.PersistentFlags().StringVarP(&default_contxt.Key, "key", "k", "", "s3 store key")
	objectCmd.PersistentFlags().StringVarP(&default_contxt.Bucket, "bucket", "b", "", "s3 store bucket")
	createObjectCmd.PersistentFlags().StringVarP(&sourceFile, "file", "f", "", "s3 store src file")
	createObjectCmd.PersistentFlags().StringVarP(&objtype, "type", "t", default_type, "s3 store type")
	createObjectCmd.PersistentFlags().Int32VarP(&max_age, "maxage", "m", 3600, "cache control max age")

	saveOjectCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "s3 store output file")

}
