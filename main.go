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

// need to refactor code into smaller function
func pdfFileUpload(res http.ResponseWriter, r *http.Request) {
	fmt.Println("File upload endpoint")
	//If get then serves index.html
	if r.Method == "GET" {
		http.ServeFile(res, r, "static/index.html")
		return
	}
	if r.Method == "POST" {
		//hopefully this makes 200MB the max limit!
		r.ParseMultipartForm(200 << 20)
		//Make this a multiple file so the file will be array of files
		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Println("there is an error parsing file from Form!")
			fmt.Println(err)
			return
		}
		//Displays all the values send from the html form.
		for key, value := range r.Form {
			fmt.Print("%s = %s\n", key, value)
		}
		//Data on the files displayed on the console.
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("FileSize: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		//Creating a temp file
		//create name based on current time
		theFileName := addTimeToFileName(handler.Filename)
		createAndWriteFile(theFileName, file)
		fmt.Fprintf(res, "Successfully uploaded File\n")
		go printFile(theFileName)
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
