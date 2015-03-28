package main

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

func fileInfo(fileName string) (string, int64) {
	file, error := os.Stat(fileName)
	if error != nil {
		fmt.Println(error)
	}
	return file.Name(), file.Size()
}

func listFiles(filePath string) {
	files, _ := ioutil.ReadDir(filePath)
	for _, f := range files {
		fileName := f.Name()
		fname, fsize := fileInfo(fileName)
		fmt.Println("File:", fname, "and size is:", float64(float64(fsize)/(1024*1024)))
	}

}

func main() {
	filePath := flag.String("filepath", "default", "a string")
	flag.Parse()
	listFiles(*filePath)
}
