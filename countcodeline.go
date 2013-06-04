package main

import(
	"os"
	"container/list"
	"io"
	"log"
	"path/filepath"
	"bufio"
)

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
	srcList := GetParsedFilesByConf(allFiles, conf)
	res := Parse(srcList)
	total := 0
	for k,v := range res{
		total += v
		println(k,":",v)
	}
	println("Total: ", total)
}

//example:map["go" or "css"]123
func Parse(files list.List)(parseResult map[string]int){
	parseResult = map[string]int{}
	for e := files.Front(); nil != e; e = e.Next(){
		v := e.Value.(string)
		ext := filepath.Ext(v)
		parseResult[ext] += ComputeLine(v)
	}
	return
}

func ComputeLine(path string)(num int){
	f,err := os.Open(path)
	if nil != err{
		log.Println(err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for{
		_,err := r.ReadString('\n')
		if io.EOF == err || nil != err{
			break
		}
		num += 1
	}
	return
}

func ParseConf()(conf []string, err error){
	confPath,_ := os.Getwd()
	confPath +=  "/conf"
	file,err := os.Open(confPath)
	if nil != err{
		log.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		conf = append(conf, scanner.Text())
	}
	//log.Println(conf)
	return
}

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

func GetParsedFilesByConf(lst list.List,conf []string)(l list.List){
	for e := lst.Front(); nil != e; e = e.Next(){
		ext := filepath.Ext(e.Value.(string))
		for _,extension := range conf{
			if ext == extension{
				l.PushBack(e.Value.(string))
				continue
			}
		}
	}

	for e := l.Front(); nil != e; e = e.Next(){
		println("matched string:" ,e.Value.(string))
	}
	return
}

func GetAllFilesName() (lst list.List){
	fullPath := GetSrcFullPath()
	log.Println("fullPath:" ,fullPath)
	filepath.Walk(fullPath,func(path string,fi os.FileInfo,err error)error{
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		//println(fi.Name())
		//println(path)
		//name := fi.Name()
		lst.PushBack(path)
		return nil
	})
	return
}
