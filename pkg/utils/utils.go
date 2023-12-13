package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/yylq/s3tool/pkg/common"
)

func DownOject(out string, context *common.S3Context) error {
	if out == "" {
		return fmt.Errorf("invalid outputFile file")
	}
	var mataMd5 string
	bucket := context.Bucket
	key := context.Key
	creds := credentials.NewStaticCredentials(context.Ak, context.SK, "")
	_, err := creds.Get()
	if err != nil {
		return err
	}

	config := &aws.Config{
		Region:      aws.String(context.RegionId),
		Endpoint:    aws.String(context.Endpoints),
		DisableSSL:  aws.Bool(true),
		Credentials: creds,
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	getobj := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	outObject, err := svc.GetObject(getobj)
	if err != nil {
		return err
	}
	meta := outObject.Metadata
	if meta != nil {
		mataMd5 = *meta["Md5"]
		fmt.Println(mataMd5)
	}
	fmt.Println(outObject)
	tmpout, _ := ioutil.TempDir("/tmp", "s3_")
	tmpFile := tmpout + "/" + path.Base(out)
	{
		outFile, err := os.Create(tmpFile)
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, outObject.Body)
		if err != nil {
			return err
		}
	}
	newmd5 := Md5sumFile(tmpFile)
	if mataMd5 != "" && newmd5 == mataMd5 {
		os.Rename(tmpFile, out)
		return nil
	}
	os.Remove(tmpout)
	return nil
}
func Md5sumFile(filePath string) string {

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Println(err)
		return ""
	}

	md5Hash := hash.Sum(nil)
	md5String := fmt.Sprintf("%x", md5Hash)
	return md5String
}
