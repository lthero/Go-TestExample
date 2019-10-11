package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"net/http"
	"log"
	"os"
	"time"
	"io"
)

const (
	cityList = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`
	//LOGPATH  LOGPATH/time.Now().Format(FORMAT)/*.log
	LOGPATH = `C:\Go\GO_Work\spiderlog\`
	//FORMAT .
	FORMAT = "20060102"
	//LineFeed 换行
	LineFeed = "\r\n"
)
var path = LOGPATH + time.Now().Format(FORMAT) + `\`

type Request struct {
	Url string
	ParserFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items []interface{}
}

func NilParser([]byte) ParseResult{
	return ParseResult{}
}
//主函数
func main(){


	url := "http://www.zhenai.com/zhenghun"
	run(Request{
		Url:url,
		ParserFunc:printCityList,
	})
}

func run(seeds ...Request) {
	CreateDir(path)

	var requests []Request
	for _,r := range seeds {
		requests = append(requests,r)
	}

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		go writelog(`spider.log`,r.Url)
		body, err := fetch(r.Url)
		if err != nil {
			log.Printf("Error info: %s\n",err)
			continue
		}

		ParseResult := r.ParserFunc(body)
		requests = append(requests,ParseResult.Requests...)
		for _, item := range ParseResult.Items {
			//
			//
			// adminlog.Printf("%s\n",item)
			writelog(`spider.log`,item.(string))
		}
	}
}

func fetch(url string) ([]byte,error){
	resp,err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status Code : d%",resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func printCityList(contents []byte) ParseResult{
	re:=regexp.MustCompile(cityList)
	content := re.FindAllSubmatch(contents,-1)
	result := ParseResult{}
	for _,m := range content {
		result.Items = append(result.Items,string(m[2]))
		result.Requests = append(result.Requests, Request{
			Url: string(m[1]),
			ParserFunc: NilParser,
		})
	}
	return result
}
//write log
func writelog(fileName, msg string) error {
	if !IsExist(path) {
		return CreateDir(path)
	}
	var (
		err error
		f   *os.File
	)

	f, err = os.OpenFile(path+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	_, err = io.WriteString(f, LineFeed+msg)

	defer f.Close()
	return err
}
//CreateDir  文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

//IsExist  判断文件夹/文件是否存在  存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
