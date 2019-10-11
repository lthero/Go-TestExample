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
	"strconv"
)

const (
	cityList = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`
	//LOGPATH  LOGPATH/time.Now().Format(FORMAT)/*.log
	LOGPATH = `C:\Go\GO_Work\spiderlog\`
	//FORMAT .
	FORMAT = "20060102"
	//LineFeed 换行
	LineFeed = "\r\n"
	// 用户信息正则
)
var path = LOGPATH + time.Now().Format(FORMAT) + `\`
var ageRe =regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var hunRe =regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)

type Request struct {
	Url string
	ParserFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items []interface{}
}

type Profile struct {
	Name 		string
	Gender 		string
	Age 		int
	Height 		int
	Weight 		int
	Income 		string
	Marriage 	string
	Education 	string
	Occupation 	string
	Hukou 		string
	Xinzuo 		string
	House		string
	Car 		string
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

//运行函数
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
			switch item.(type) {

			case string:
				writelog(`spider.log`,item.(string))
				break
			case int:
				writelog(`spider.log`, item.(string))
				break
			case float64:
				writelog(`spider.log`, item.(string))
				break
			}
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
		log.Printf("%s \n",m[2])
		result.Items = append(result.Items,string(m[2]))
		result.Requests = append(result.Requests, Request{
			Url: string(m[1]),
			ParserFunc: func(c []byte) ParseResult {
				return ParseProfile(c,string(m[2]))
			},
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
//获取用户具体信息
func ParseProfile(contents []byte,name string) ParseResult{
	profile := Profile{}
	profile.Name = name
	age,err := strconv.Atoi(string(extractString(contents,ageRe)))

	if err != nil {
		profile.Age = age
	}
	profile.Marriage = extractString(contents,hunRe)

	result := ParseResult{
		Items:[]interface{}{profile},
	}
	return result
}

func extractString (contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >=2 {
		return string(match[1])
	}else{
		return ""
	}
}