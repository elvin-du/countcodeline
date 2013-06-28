/*
parse the file that be defined by @conf file, and print the number ofrows, comment rows, and the number of files.
*/
package main

import(
	"os"
	"container/list"
	"log"
	"path/filepath"
	"bufio"
	"encoding/xml"
	"io/ioutil"
	"strings"
	"fmt"
)

type CONF struct{
	FILES	[]FILE	`xml:"file"`
}

type FILE struct{
	EXT		string	`xml:"ext"`
	COMMENT	string	`xml:"comment"`
}

func init(){
	log.SetFlags(log.Llongfile| log.LstdFlags)
}

func main(){
	conf,err := ParseConf()
	if nil != err{
		log.Println(err)
		return
	}
	allFiles := GetAllFilesName()
	fmt.Println("AllFiles: ", allFiles.Len())
	srcMap := GetParsedFilesByConf(allFiles, conf)
	for k, v := range srcMap{
		fmt.Println(k, "file: ", len(v))
	}
	res := Parse(srcMap, conf)
	total := 0
	for k, v := range res{
		total += v.COMMENT + v.CODE
		fmt.Println("===========================================================")
		fmt.Println(k, "file file line: ", v.COMMENT + v.CODE)
		fmt.Println(k, "file comment line: ", v.COMMENT)
		fmt.Println(k, "file code line: ", v.CODE)
		fmt.Println("===========================================================")
	}
	fmt.Println("file total line: ", total)
}

type SRCNUM struct{
	CODE	int
	COMMENT int
}

//example:map["go" or "css"]123
func Parse(files map[string][]string, conf CONF)(parseResult map[string]SRCNUM){
	parseResult = map[string]SRCNUM{}

	for k, v := range files{
		te, tc := 0, 0
		for _, v2 := range v{
			e, c := ComputeLine(v2, conf)
			te += e
			tc += c
		}
		parseResult[k] = SRCNUM{te, tc}
	}
	return
}

func GetComment(ext string, conf CONF)(comment string){
	for _, v := range conf.FILES{
		if v.EXT == ext{
			comment = v.COMMENT
			break
		}
	}
	return
}

func ComputeLine(path string, conf CONF)(code, comment int){
	f,err := os.Open(path)
	if nil != err{
		log.Println(err)
		return
	}
	defer f.Close()
	ext := filepath.Ext(path)
	strComment := GetComment(ext, conf)

	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), strComment) {
			comment += 1
			continue
		}
		code += 1
	}
	return
}

/*
parse named conf file to get extention, 
and comment into return value @conf .
*/
func ParseConf()(conf CONF, err error){
	confPath,_ := os.Getwd()
	confPath +=  "/conf.xml"
	file,err := os.Open(confPath)
	if nil != err{
		log.Println(err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil{
		log.Println(err)
		return
	}

	conf = CONF{}
	err = xml.Unmarshal(data, &conf)
	if nil != err{
		log.Println(err)
		return
	}
	//log.Println(conf)

	return
}

func GetParsedFilesByConf(lst list.List, conf CONF)(result map[string][]string){
	result = map[string][]string{}

	for e := lst.Front(); nil != e; e = e.Next(){
		ext := filepath.Ext(e.Value.(string))
		for _,file := range conf.FILES{
			if ext == file.EXT{
				result[ext] = append(result[ext], e.Value.(string))
				break
			}
		}
	}

	return
}

/*
put all file names of the project into @lst 
*/
func GetAllFilesName() (lst list.List){
	fullPath := GetSrcFullPath()
	fmt.Println("fullPath:" ,fullPath)
	filepath.Walk(fullPath,func(path string,fi os.FileInfo,err error)error{
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		lst.PushBack(path)
		return nil
	})
	return
}

/*
get full path of project which will be parsed.
*/
func GetSrcFullPath()(fullPath string){
	args := os.Args;
	parameterLen := len(args)
	if parameterLen == 1{
		fullPath,_ = os.Getwd()
	}else if(parameterLen >= 2){
		fullPath = args[1]
	}

	fullPath,_ = filepath.Abs(fullPath)
	return
}
