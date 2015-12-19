package main

import (
	//"image"
	//"image/color"
	//"runtime"

	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const UPLOADEDFILEPATH string = "uploadedImage/"
const THUMBNAILPATH string = "thumbnailImage/"
const THUMBNAIL_WIDTH int = 100
const THUMBNAIL_HEIGHT int = 100

func checkFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	filename := header.Filename
	fileExtender := strings.Split(header.Filename, ".")[len(strings.Split(header.Filename, "."))-1]

	for i := 1; checkFileExist(UPLOADEDFILEPATH + filename); i++ {
		filename = strings.Join(strings.Split(header.Filename, ".")[:len(strings.Split(header.Filename, "."))-1], "") + "_" + strconv.Itoa(i) + "." + fileExtender
	}

	out, err := os.Create(UPLOADEDFILEPATH + filename)
	if err != nil {
		//fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		fmt.Fprintln(w, err)
		return
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintln(w, filename)

	thumbnailName, err := makeThumbnail(filename)
	if err != nil {
		fmt.Println("thumbnail error")
	} else {
		fmt.Fprintln(w, THUMBNAILPATH+thumbnailName)
	}

}

func makeThumbnail(filename string) (string, error) {
	/*
	   //make thumbnail from bytes
	   			reader := bufio.NewReader(conn)

	   			var values []byte
	   			buf := make([]byte, 10)
	   			for {
	   				size, err := reader.Read(buf)
	   				if err != nil {
	   					return
	   				}

	   				values = append(values, buf...)
	   				if size < len(buf) {
	   					break
	   				}
	   			}
	   			byte_reader := bytes.NewReader(values)

	   			//load images and make 100x100 thumbnails of them
	   			var thumbnails []image.Image
	   			//	for _, file := range files {
	   			//img, err := imaging.Open(file)
	   			img, err := imaging.Decode(byte_reader)
	   			if err != nil {
	   				panic(err)
	   			}
	*/
	img, err := imaging.Open(UPLOADEDFILEPATH + filename)
	if err != nil {
		return "", errors.New("makeThumbnail fail")
	}

	thumb := imaging.Thumbnail(img, THUMBNAIL_WIDTH, THUMBNAIL_HEIGHT, imaging.CatmullRom)

	//create a new blank image
	//dst := imaging.New(THUMBNAIL_WIDTH, THUMBNAIL_HEIGHT, color.NRGBA{0, 0, 0, 0})

	//paste thumbnails into hte new image side by side
	//dst = imaging.Paste(dst, thumb, image.Pt(THUMBNAIL_WIDTH, 0))

	//save the combined image to file
	imaging.Save(thumb, THUMBNAILPATH+"thumbnail_"+filename)

	return "thumbnail_" + filename, nil
}

func uploadForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ` <html>
 <title>Go upload</title>
 <body>

 <form action="./receive" method="post" enctype="multipart/form-data">
 <label for="file">Filename:</label>
 <input type="file" name="file" id="file">
 <input type="submit" name="submit" value="Submit">
 </form>

 </body>
 </html>
`)
}

func main() {
	http.HandleFunc("/receive", uploadHandler)
	http.HandleFunc("/uploadForm", uploadForm)
	http.ListenAndServe(":8080", nil)
}
