package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

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
		file, handler, err := r.FormFile("pdfFile")
		if err != nil {
			fmt.Println("there is an error retreving file!")
			fmt.Println(err)
			return
		}
		//Displayes all the values send from the html form.
		for key, value := range r.Form {
			fmt.Print("%s = %s\n", key, value)
		}
		//Data on the files displayed on the console.
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("FileSize: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		//Creating a temp file
		tempFile, err := ioutil.TempFile("mainFiles", "upload-*.pdf")
		if err != nil {
			fmt.Println("There is an error with making temp file")
			fmt.Print(err)
		}

		//Reading the file into byte Array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("There was an error reading file into fileByes")
			fmt.Print(err)
		}
		//Write the file into temp from byte array
		defer tempFile.Close()
		tempFile.Write(fileBytes)
		fmt.Println("Successfully saved file")
		fmt.Fprintf(res, "Successfully uploaded File\n")
	}

}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/find&replace", func(w http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			//fmt.Fprintf(w, "Cannot GET find&replace "+request.RemoteAddr)
			return
		}
		fmt.Fprintf(w, "Find and replace, Will be fixed")
	})
	//PDF Submitting Post file
	http.HandleFunc("/pdf", pdfFileUpload)
	fmt.Println("Server listening on port :3000")

	http.ListenAndServe(":3000", nil)
}
