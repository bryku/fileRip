package main

import(
	"fmt"
	"github.com/bryku/prams"
	"os"
	"strings"
	"encoding/hex"
)
func PNG(data []byte) []string{
	fileStart := "89504e47"
	fileEnd := "49454e44ae426082"
	fileList := strings.Split(hex.EncodeToString(data),fileStart)
	for key, _ := range fileList {
		if(strings.Contains(fileList[key], fileEnd) == true){
			s := strings.Split(fileList[key], fileEnd)
			fileList[key] = fileStart + s[0] + fileEnd
		}else{
			fileList[key] = ""
		}
	}
	return fileList
}
func JPG(data []byte) []string{
	fileStart := "ffd8ffe000104a46494600010101012c012c0000ffe1"
	fileEnd := "ff00f00fffd9"
	fileList := strings.Split(hex.EncodeToString(data),fileStart)
	for key, _ := range fileList {
		if(strings.Contains(fileList[key], fileEnd) == true){
			s := strings.Split(fileList[key], fileEnd)
			fileList[key] = fileStart + s[0] + fileEnd
		}else{
			fileList[key] = ""
		}
	}
	return fileList
}
func CreateFile(files []string, filePath string, fileType string) string{
	ret := ""
	fp := strings.Split(filePath, ".")
	filePath = fp[0]
	fileName := ""

	for key, _ := range files {
		fileName = fmt.Sprintf("%s%s%d%s%s", filePath, "-", key, ".", fileType)
		if(len(files[key]) > 0){

			newdata, _ := hex.DecodeString(files[key])
			newfile, err := os.Create(fileName) // create new file
			if err != nil {ret += fmt.Sprintf("%s%s%s", fileName," - fail","\n")
			}else{
				newfile.Write(newdata)
				newfile.Close()
				ret += fmt.Sprintf("%s%s%s", fileName, " - success", "\n")
			}
		}else{
			ret += fmt.Sprintf("%s%s%s", fileName, " - No Data", "\n")
		}
	}
	return ret
}
func main(){
	status := ""

	filePath, isPath := prams.Get("-path")
	fileType, isType := prams.Get("-type")
	_,	isReport := prams.Get("report")
	_, isHelp := prams.Get("help")

	if(isPath == true && isType == true){
		if(len(filePath) < 1){
			status = "File Path: Required"
		}else if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			file, err := os.Open(filePath)
			if err != nil {status = "File Path: Unknown"}
			fileStat, err := file.Stat()
			if err != nil {status = "File: Unknown Error"}
			fileSize := fileStat.Size()
			data := make([]byte, fileSize)
			file.Read(data)

			if(len(data) > 0){
				switch(fileType){
					case "png":
						files := PNG(data)
						status = CreateFile(files, filePath, fileType)
					default:
						status = "File Type: Unknown [png, jpg]"
				}
			}else{
				status = "File Data: None"
			}
		}else{
			status = "File Path: Not Found"
		}
	}else{
		isHelp = true
	}
	if(isHelp == true){
		isReport = true
		status = `
fileRip: 2018-12-06
Author: drburnett
Version: 2.03
Summary: Extracts specified file types and create those as separate files.
Examples:
fileRip -path /home/$USER/Downloads/file -type png
fileRip -path /home/$USER/Downloads/file -type png -report true
`
	}
	if(isReport == true){
		fmt.Println("File Path:", filePath)
		fmt.Println("File Type:", fileType)
		fmt.Println("Report:", isReport)
	}
	fmt.Println(status)
}



