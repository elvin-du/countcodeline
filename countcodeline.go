package main

import(
	"os"
	"io"
	"log"
	"path/filepath"
	"bufio"
)

func init(){

}


func main(){
	fullPath := GetSrcFullPath()

	ParseConf()
	log.Println(fullPath)
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

	rd := bufio.NewReader(file)
	var line string
	for{
		line, err = rd.ReadString('\n')
		if io.EOF == err || nil != err{
			break
		}
		conf = append(conf, line)
	}
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
