package main

import(
	"os"
	"log"
	"os/exec"
	"path/filepath"
)

func init(){

}

func main(){
	fullPath := GetSrcFullPath()
	log.Println(fullPath)
}

func GetSrcFullPath()(fullPath string){
	args := os.Args;
	parameterLen := len(args)
	if parameterLen == 1{
		fullPath,_ = exec.LookPath(args[0])
	}else if(parameterLen >= 2){
		fullPath = args[1]
	}

	fullPath,_ = filepath.Abs(fullPath)
	return
}
