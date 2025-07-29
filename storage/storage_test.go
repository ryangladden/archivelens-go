package storage

// import (
// 	"bytes"
// 	"io"
// 	"mime/multipart"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"github.com/rs/zerolog/log"
// )

// var (
// 	sm *StorageManager
// )

// func TestMain(m *testing.M) {
// 	sm = NewStorageManager("http://gertrude:9000", "s3-test", "us-east-1")
// 	m.Run()
// }

// func TestUploadFile(t *testing.T) {

// 	file := createMultipartFileHeader("../cookies.txt")
// 	sm.UploadFile(file, "cookies-file-teehee")
// 	log.Debug().Msgf("Presigned URL: %s", *sm.GeneratePresignedURL("cookie-file-teehee", 30))
// }

// func createMultipartFileHeader(filePath string) *multipart.FileHeader {
// 	// open the file
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		log.Fatal().Err(err)
// 		return nil
// 	}
// 	defer file.Close()*

// 	// create a buffer to hold the file in memory
// 	var buff bytes.Buffer
// 	buffWriter := io.Writer(&buff)

// 	// create a new form and create a new file field
// 	formWriter := multipart.NewWriter(buffWriter)
// 	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
// 	if err != nil {
// 		log.Fatal().Err(err)
// 		return nil
// 	}

// 	// copy the content of the file to the form's file field
// 	if _, err := io.Copy(formPart, file); err != nil {
// 		log.Fatal().Err(err)
// 		return nil
// 	}

// 	// close the form writer after the copying process is finished
// 	// I don't use defer in here to avoid unexpected EOF error
// 	formWriter.Close()

// 	// transform the bytes buffer into a form reader
// 	buffReader := bytes.NewReader(buff.Bytes())
// 	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

// 	// read the form components with max stored memory of 1MB
// 	multipartForm, err := formReader.ReadForm(1 << 20)
// 	if err != nil {
// 		log.Fatal().Err(err)
// 		return nil
// 	}

// 	// return the multipart file header
// 	files, exists := multipartForm.File["file"]
// 	if !exists || len(files) == 0 {
// 		log.Fatal().Msg("multipart file not exists")
// 		return nil
// 	}

// 	return files[0]
// }
