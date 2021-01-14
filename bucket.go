package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
	"log"
)

var (

	atrr_str string
	BucketCmd = &cobra.Command{
		Use:   "bucket",
		Short: "manager s3 bucket",
	}
	createBucketCmd = &cobra.Command{
		Use:   "create",
		Short: "create s3 bucket",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", default_contxt)
			err:= createBucketOject()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}
			log.Printf("Successfully create bucket %s in  %s\n", default_contxt.Bucket, default_contxt.Endpoints)
		},
	}
	deleteBucketCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete s3 bucket",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", default_contxt)
			err:= deleteBucketOject()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}
			log.Printf("Successfully delete bucket %s in  %s\n", default_contxt.Bucket, default_contxt.Endpoints)
		},
	}
	listBucketCmd = &cobra.Command{
		Use:   "list",
		Short: "list s3 bucket",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", default_contxt)
			err:= ListBucketOject()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}

		},
	}
	corsBucketCmd = &cobra.Command{
		Use:   "cors",
		Short: "cors s3 bucket",
	}

	setBucketCorsCmd = &cobra.Command{
		Use:   "set",
		Short: "set s3 bucket cors",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", default_contxt)
			err:= setBucketCors()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}

		},
	}
	listBucketCorsCmd = &cobra.Command{
		Use:   "list",
		Short: "list s3 bucket cors",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", default_contxt)
			err:= ListBucketCors()
			if err != nil {
				log.Printf("err:%v",err)
				return;
			}

		},
	}


)

func createBucketOject()  error {


	bucket := default_contxt.Bucket
	log.Printf("bucket:%v",bucket)

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
	private :="public-read"
	obj:= &s3.CreateBucketInput{Bucket:&bucket,ACL:&private}

	out, err := svc.CreateBucket(obj)
	log.Printf("out:%v",out)
	log.Printf("curl http://%s bucket:%s",default_contxt.Endpoints, bucket)
	return err
}
func deleteBucketOject()  error {


	bucket := default_contxt.Bucket
	log.Printf("bucket:%v",bucket)

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
	obj:= &s3.DeleteBucketInput{Bucket:&bucket}

	out, err := svc.DeleteBucket(obj)
	log.Printf("out:%v",out)
	log.Printf("curl http://%s bucket:%s",default_contxt.Endpoints, bucket)
	return err
}
func ListBucketOject()  error {


	bucket := default_contxt.Bucket
	log.Printf("bucket:%v",bucket)
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
	obj:= &s3.ListBucketsInput{}

	out, err := svc.ListBuckets(obj)
	log.Printf("out:%v",out)



	return err
}
func setBucketCors()  error {

	bucket := default_contxt.Bucket
	log.Printf("bucket:%v",bucket)


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



	method :=[]string{"PUT","GET", "POST", "DELETE","OPTIONS"}
	origin :=[]string{"*",}
	conf:= &s3.CORSConfiguration{

		CORSRules: []*s3.CORSRule{
			&s3.CORSRule{
				AllowedHeaders:[]*string{&origin[0]} ,
				AllowedMethods:[]*string{ &method[0],&method[1], &method[2],&method[3]},
				AllowedOrigins: []*string{&origin[0]},

			},
		},
	}
	bytes,_ := json.Marshal(conf)
	log.Printf("%s",string(bytes))
	obj:= &s3.PutBucketCorsInput{
		Bucket: &bucket,
		CORSConfiguration: conf,
	}

	out, err := svc.PutBucketCors(obj)
	log.Printf("out:%v",out)
	return err
}
func ListBucketCors()  error {

	bucket := default_contxt.Bucket
	log.Printf("bucket:%v",bucket)


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


	obj:= &s3.GetBucketCorsInput{
		Bucket: &bucket,
	}

	out, err := svc.GetBucketCors(obj)
	log.Printf("out:%v",out)
	return err
}
func init() {
	rootCmd.AddCommand(BucketCmd)
	BucketCmd.AddCommand(createBucketCmd, deleteBucketCmd,listBucketCmd, corsBucketCmd)
	corsBucketCmd.AddCommand(setBucketCorsCmd, listBucketCorsCmd)
	BucketCmd.PersistentFlags().StringVarP(&default_contxt.Bucket, "bucket", "b", "", "s3 store bucket")
	corsBucketCmd.PersistentFlags().StringVarP(&atrr_str, "atrr", "a", "", "s3 bucket attr")
}
