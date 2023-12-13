package common

import "os"

type S3Context struct {
	Endpoints string
	RegionId  string
	Bucket    string
	Ak        string
	SK        string
	Key       string
}

func BuildS3Context() *S3Context {
	return &S3Context{
		Endpoints: os.Getenv("S3_ENDPOINTS"),
		RegionId:  os.Getenv("S3_REGIONID"),
		Bucket:    os.Getenv("S3_BUCKET"),
		Ak:        os.Getenv("S3_AK"),
		SK:        os.Getenv("S3_SK"),
	}
}
