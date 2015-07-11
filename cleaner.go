package main

//run command via go run cleaner.go --filepath=/Users/caveman/Desktop
import (
	// "flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
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
	fileMeta.lastModiefied = file.ModTime().Format(layout)
	return fileMeta
}

func listFiles(w http.ResponseWriter, filePath string) {
	files, _ := ioutil.ReadDir(filePath)

	fmt.Fprint(w, "File:", "| <b>last accessed at:</b>", "|and last modified at:")
	for _, f := range files {
		fileName := f.Name()
		fullFilePath := filePath + fileName
		fileMeta := fileInfo(fullFilePath)
		fmt.Println(fileName, "|", fileMeta.lastAccessed, "|", fileMeta.lastModiefied)
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, fileName, "|", fileMeta.lastAccessed, "|", fileMeta.lastModiefied)
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to MAC Cleaner!\n")
}

func prepareCleaning(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Path is: %s!\n", ps.ByName("path"))
	filePath := ps.ByName("path")
	listFiles(w, filePath)
}
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/cleaner/*path", prepareCleaning)
	log.Fatal(http.ListenAndServe(":8080", router))
}
