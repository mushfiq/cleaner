package main

//run command via go run cleaner.go --filepath=/Users/caveman/Desktop
import (
	// "flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
	"html/template"
	"github.com/julienschmidt/httprouter"
	// "github.com/gorilla/schema"|
)

type MetaData struct {
	LastAccessed  string
	LastModiefied string
	FileName 	  string
}

var fileInfoPageTmpl, err = template.ParseFiles("index.html")

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
	fileMeta.LastAccessed = dt.Format(layout)
	fileMeta.LastModiefied = file.ModTime().Format(layout)
	fileMeta.FileName = file.Name()
	return fileMeta
}

func listFiles(w http.ResponseWriter, filePath string) {
	files, _ := ioutil.ReadDir(filePath)

	// fmt.Fprint(w, "Filename :	", "| last accessed at:", "|and last modified at")
	
	allFiles := []MetaData{}
	 
	for _, f := range files {
		fileName := f.Name()
		fullFilePath := filePath + fileName
		fileMeta := fileInfo(fullFilePath)
		singleFileInfo := MetaData{fileMeta.LastAccessed, fileMeta.LastModiefied, fileMeta.FileName}
		allFiles = append(allFiles, singleFileInfo)
		fmt.Println(allFiles)
	}
	
	if err := fileInfoPageTmpl.Execute(w, allFiles); err != nil {
		fmt.Println("Failed to build page", err)
	}
	
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to MAC Cleaner!\n")
}

func prepareCleaning(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// fmt.Fprintf(w, "Path is: %s!\n", ps.ByName("path"))
	filePath := ps.ByName("path")
	listFiles(w, filePath)
}
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/cleaner/*path", prepareCleaning)
	log.Fatal(http.ListenAndServe(":8080", router))
}
