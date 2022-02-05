package main

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func addTimeToFileName(fileName string) string {
	dt := time.Now()
	dtString := dt.Format("01-02-2006 15:04:05 Monday")
	newName := dtString + "-" + fileName
	return newName
}
func printFile(pathOfFile string) {
	out := exec.Command("lp", pathOfFile)
	err := out.Run()
	fmt.Println("File being printed")
	if err != nil {
		fmt.Println(err)
	}
}
func multipleFiles(r *http.Request) {
	// setting the max size
	err := r.ParseMultipartForm(200 << 20)
	// handle err
	if err != nil {
		fmt.Println(err)
	}
	// list of files from the form
	var filesList []*multipart.FileHeader
	// iterate through them
	for _, fh := range r.MultipartForm.File["uploadFile"] {
		// append files to lise
		filesList = append(filesList, fh)
		// printing to console, fileName, size, and header
		fmt.Printf("Uploaded File: %+v\n", fh.Filename)
		fmt.Printf("FileSize: %+v\n", fh.Size)
		fmt.Printf("MIME Header: %+v\n", fh.Header)
		// add timeToFileName
		tempFileName := addTimeToFileName(fh.Filename)
		// create and Write each File
		//getting the file from file.Header
		theFile, _ := fh.Open()
		// passing the file along with the file name to save it as
		createAndWriteFile(tempFileName, theFile)
		fmt.Println("Successfully created the file")
		// printing
		go printFile(tempFileName)
		// print each one of the files
	}
}

// need to refactor code into smaller function
func pdfFileUpload(res http.ResponseWriter, r *http.Request) {
	fmt.Println("File upload endpoint")
	//If get then serves index.html
	if r.Method == "GET" {
		http.ServeFile(res, r, "static/index.html")
		return
	}
	if r.Method == "POST" {
		fmt.Fprintf(res, "Successfully uploaded File\n")
		multipleFiles(r)
	}
}
func createAndWriteFile(theFileName string, file multipart.File) {
	f, err := os.Create("./files/" + theFileName)
	if err != nil {
		fmt.Println("There is an error with making temp file")
		fmt.Print(err)
	}

	//Reading the file into byte Array
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("There was an error reading file into fileByes")
		fmt.Print(err)
	}

	defer f.Close()
	f.Write(fileBytes)
	//Write the file into temp from byte array

	fmt.Println("Successfully saved file")
}
func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	//PDF Submitting Post file
	http.HandleFunc("/upload", pdfFileUpload)
	fmt.Println("Server listening on port :3000")
	http.ListenAndServe(":3000", nil)
}
