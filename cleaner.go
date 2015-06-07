package main

//run command via go run cleaner.go --filepath=/Users/caveman/Desktop
import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"
)

type MetaData struct {
	lastAccessed  string
	lastModiefied string
}

func unixTimeToDate(lastAccess string) time.Time {
	i, err := strconv.ParseInt(lastAccess, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm
}

func fileInfo(fileName string) MetaData {
	// for time formatting in golang
	// details: http://golang.org/pkg/time/#Time.Format and
	// https://gobyexample.com/time-formatting-parsing

	const layout = "2006-01-02"
	var fileMeta MetaData

	file, error := os.Stat(fileName)

	if error != nil {
		fmt.Println(error)
	}
	lastAccess := (file.Sys().(*syscall.Stat_t).Atimespec.Sec)
	// dirty hack fix it via actively coverting or using sth else than ParseInt!!
	dt := unixTimeToDate(fmt.Sprint(lastAccess))
	fileMeta.lastAccessed = dt.Format(layout)
	fmt.Println("getting file info...")
	fileMeta.lastModiefied = file.ModTime().Format(layout)
	return fileMeta
}

func listFiles(filePath string) {
	files, _ := ioutil.ReadDir(filePath)
	for _, f := range files {
		fileName := f.Name()
		fullFilePath := filePath + fileName
		fileMeta := fileInfo(fullFilePath)
		fmt.Println("--------------")
		fmt.Println("File:", fileName, ", last accessed at:", fileMeta.lastAccessed, "and last modified at:", fileMeta.lastModiefied)
	}

}

func main() {
	filePath := flag.String("filepath", "default", "a string")
	flag.Parse()
	listFiles(*filePath)
}
