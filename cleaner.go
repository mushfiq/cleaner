package main

//run command via go run cleaner.go --filepath=/Users/caveman/Desktop
import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

type MetaData struct {
	LastAccessed, LastModiefied, FileName, FullFilePath string
}

var fileInfoPageTmpl, err = template.ParseFiles("templates/index.html")

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
	fileMeta.FullFilePath = fileName
	return fileMeta
}

func listFiles(w http.ResponseWriter, filePath string, r *http.Request) {
	files, _ := ioutil.ReadDir(filePath)

	allFiles := []MetaData{}

	for _, f := range files {
		fileName := f.Name()
		FullFilePath := filePath + fileName
		fileMeta := fileInfo(FullFilePath)
		singleFileInfo := MetaData{fileMeta.LastAccessed, fileMeta.LastModiefied, fileMeta.FileName, fileMeta.FullFilePath}
		allFiles = append(allFiles, singleFileInfo)
	}

	if err := fileInfoPageTmpl.Execute(w, allFiles); err != nil {
		fmt.Println("Failed to build page", err)
	}

}

func deleteFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// dirty?
	if r.FormValue("filename") != "" {
		for _, value := range r.Form {
			for i := 0; i < len(value); i++ {
				if value[i] != "" {
					filePath := ps.ByName("path")
					fmt.Println("Delteing", value[i])
					err := os.Remove(value[i])
					if err != nil {
						fmt.Println(err)
						return
					}
					// after successful file deletion redirecting to file listing
					listFiles(w, filePath, r)
				}

			}
		}
	}

}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to MAC Cleaner!\n")
}

func prepareCleaning(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filePath := ps.ByName("path")
	listFiles(w, filePath, r)
}

func main() {
	fmt.Println("Mac cleaner running on port 8080")
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/cleaner/*path", prepareCleaning)
	router.POST("/cleaner/*path", deleteFiles)
	log.Fatal(http.ListenAndServe(":8080", router))
}
