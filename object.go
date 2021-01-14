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
	"path"
)

var (
	default_contxt = buildDefaultContext()
	default_type = "application/octet-stream"
	sourceFile string
	objtype string

	objectCmd = &cobra.Command{
		Use:   "object",
		Short: "manager s3 object",
	}
	createObjectCmd = &cobra.Command{
		Use:   "create",
		Short: "create s3 object",
		Run: func(cmd *cobra.Command, args []string) {
		        log.Printf("%v", default_contxt)
			err:= createOject()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}
			log.Printf("Successfully created bucket %s and uploaded data with key %s\n", default_contxt.Bucket, default_contxt.Key)
		},
	}
	deleteObjectCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete s3 object",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", default_contxt)
			err:= deleteOject()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}

		},
	}

	listObjectCmd = &cobra.Command{
		Use:   "list",
		Short: "list s3 object",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", default_contxt)
			err:= listOject()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}

		},
	}

	ftype = map[string]string{
		".json":"application/json",
		".xml": "text/xml",
		".txt": "text/plain",
		".html":"text/html",
		".md": "text/plain",
	}


)
type s3Context struct {
	Endpoints string
	RegionId string
	Bucket string
	Ak string
	SK string
	Key string
}

func buildDefaultContext() s3Context {
	return s3Context{
		Endpoints: os.Getenv("S3_ENDPOINTS"),
		RegionId: os.Getenv("S3_REGIONID"),
		Bucket: os.Getenv("S3_BUCKET"),
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
    suffix := path.Ext(sourceFile)
	log.Printf("suffix:%v", suffix)
    v, ok := ftype[suffix]
    if ok {
		objtype = v
	}

	bucket := default_contxt.Bucket
	key := default_contxt.Key
	creds := credentials.NewStaticCredentials(default_contxt.Ak, default_contxt.SK, "")
	_, err = creds.Get()
	if err != nil {
		log.Printf("err:%v",err)
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
	cachecontrol :="no-cache"
    putobj := &s3.PutObjectInput{
		Body:   inputFile,
		Bucket: &bucket,
		Key:    &key,
		CacheControl:&cachecontrol,
		ContentType : &objtype,
	}
    log.Printf("putobj:%v", putobj)

	out, err := svc.PutObject(putobj)
	log.Printf("out:%v",out)
	log.Printf("curl http://%s.%s/%s",bucket,default_contxt.Endpoints,key)
	return err
}

func deleteOject()  error {


	bucket := default_contxt.Bucket
	key := default_contxt.Key
	creds := credentials.NewStaticCredentials(default_contxt.Ak, default_contxt.SK, "")
	_, err := creds.Get()
	if err != nil {
		log.Printf("err:%v",err)
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

	obj:= &s3.DeleteObjectInput{
		Bucket:&default_contxt.Bucket,
		Key:&default_contxt.Key}

	out, err := svc.DeleteObject(obj)
	log.Printf("out:%v",out)
	log.Printf("delete http://%s.%s/%s",bucket,default_contxt.Endpoints,key)
	return err
}

func listOject()  error {


	bucket := default_contxt.Bucket

	creds := credentials.NewStaticCredentials(default_contxt.Ak, default_contxt.SK, "")
	_, err := creds.Get()
	if err != nil {
		log.Printf("err:%v",err)
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
	obj:= &s3.ListObjectsInput{
		Bucket:&bucket,
		}

	out, err := svc.ListObjects(obj)
	log.Printf("out:%v",out)
	for _, obj := range out.Contents {
		log.Printf("curl http://%s.%s/%s",bucket,default_contxt.Endpoints,*obj.Key)
	}
	return err
}



func init() {
	rootCmd.AddCommand(objectCmd)
	objectCmd.AddCommand(createObjectCmd, deleteObjectCmd, listObjectCmd)
	objectCmd.PersistentFlags().StringVarP(&default_contxt.Key, "key", "k", "", "s3 store key")
	objectCmd.PersistentFlags().StringVarP(&default_contxt.Bucket, "bucket", "b", "", "s3 store bucket")
	createObjectCmd.PersistentFlags().StringVarP(&sourceFile, "file", "f", "","s3 store src file")
	createObjectCmd.PersistentFlags().StringVarP(&objtype, "type", "t", default_type,"s3 store type")

}
