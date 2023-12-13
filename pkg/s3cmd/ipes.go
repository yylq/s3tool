package s3cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/yylq/s3tool/pkg/utils"
)

var (
	IpesCmd = &cobra.Command{
		Use:   "ipes",
		Short: "ipes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return downIpes()
		},
	}
)

type IpesMeta struct {
	Obj string `json:"obj"`
	Md5 string `json:"md5"`
}

func downIpes() error {

	default_contxt.Bucket = "x86-ipes"
	default_contxt.Key = "ipes/agent_meta.json"
	outputFile = "agent_meta.json"
	err := utils.DownOject(outputFile, default_contxt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DownOject :%s err:%v", default_contxt.Key, err)
		return err
	}
	var meta []IpesMeta
	data, err := os.ReadFile(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadFile :%s err:%v", outputFile, err)
		return err
	}
	err = json.Unmarshal(data, &meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unmarshal :%s err:%v", outputFile, err)
		return err
	}
	fmt.Println(meta)
	for _, one := range meta {
		default_contxt.Key = one.Obj
		outputFile = path.Base(one.Obj)
		err := utils.DownOject(outputFile, default_contxt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DownOject :%s err:%v", one.Obj, err)
			return err
		}
		fmt.Printf("%s is download\n", outputFile)
	}
	fmt.Println("down ipes softs is ok")
	return nil
}
