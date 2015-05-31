package main
//run command via go run cleaner.go --filepath=/Users/caveman/Desktop
import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type MetaData struct {
	fileName string
	fileSize int64
}

func fileInfo(fileName string) (string) {
	// for time formatting in golang
	// details: http://golang.org/pkg/time/#Time.Format and 
 	// https://gobyexample.com/time-formatting-parsing
	
	const layout = "2006-01-02"
	
	file, error := os.Stat(fileName)
	if error != nil {
		fmt.Println(error)
	}
	
	return file.ModTime().Format(layout)
}

func listFiles(filePath string) {
	files, _ := ioutil.ReadDir(filePath)
	for _, f := range files {
		fileName := f.Name()
		fullFilePath := filePath+fileName
		fLastModified := fileInfo(fullFilePath)
		fmt.Println("File:", filePath+fileName, "and last modified at:", fLastModified)
	}

}

func main() {
	filePath := flag.String("filepath", "default", "a string")
	flag.Parse()
	listFiles(*filePath)
}
