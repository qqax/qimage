package qimage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"
)

const MaxUploadSize = 1024 * 1024 // 1MB

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
// each time Write is called
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}

	fmt.Printf("File upload in progress: %d\n", pr.BytesRead)
}

func multipleUploadHandler(w http.ResponseWriter, r *http.Request) {
	//if r.Method != "POST" {
	//	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//	return
	//}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {

		err := uploadHandler(fileHeader, []string{jpegType, pngType}, saveFile)
		if err != nil {
			return
		}
	}

	fmt.Fprintf(w, "Upload successful")
}

//func uploadHandler(w http.ResponseWriter, r *http.Request) {
//	//if r.Method != "POST" {
//	//	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//	//	return
//	//}
//
//	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
//	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
//		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
//		return
//	}
//
//	// The argument to FormFile must match the name attribute
//	// of the file input on the frontend
//	file, fileHeader, err := r.FormFile("file")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	defer file.Close()
//
//	buff := make([]byte, 512)
//	_, err = file.Read(buff)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	filetype := http.DetectContentType(buff)
//	if filetype != "image/jpeg" && filetype != "image/png" {
//		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
//		return
//	}
//
//	_, err = file.Seek(0, io.SeekStart)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Create the uploads folder if it doesn't
//	// already exist
//	err = os.MkdirAll("./uploads", os.ModePerm)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Create a new file in the uploads directory
//	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	defer dst.Close()
//
//	// Copy the uploaded file to the filesystem
//	// at the specified destination
//	_, err = io.Copy(dst, file)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	fmt.Fprintf(w, "Upload successful")
//}

const jpegType = "image/jpeg"
const pngType = "image/png"

type saveMultipartFile func(file *multipart.File, filename string, size int64) error

func uploadHandler(fileHeader *multipart.FileHeader, allowedFileTypes []string, saveMultipartFile saveMultipartFile) error {
	// Restrict the size of each uploaded file to 1MB.
	// To prevent the aggregate size from exceeding
	// a specified value, use the http.MaxBytesReader() method
	// before calling ParseMultipartForm()
	if fileHeader.Size > MaxUploadSize {
		return errors.New(fmt.Sprintf("the uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename))
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return err
	}

	filetype := http.DetectContentType(buff)
	if slices.Contains(allowedFileTypes, filetype) {
		return errors.New("the provided file format is not allowed. Please upload a JPEG or PNG image")
	}

	err = saveMultipartFile(&file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		return err
	}

	return nil
}

func saveFile(file *multipart.File, filename string, size int64) error {
	_, err := (*file).Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	//_, err = io.Copy(f, *file)
	//if err != nil {
	//	return err
	//}

	pr := &Progress{
		TotalSize: size,
	}

	_, err = io.Copy(f, io.TeeReader(*file, pr))
	if err != nil {
		return err
	}

	return nil
}
