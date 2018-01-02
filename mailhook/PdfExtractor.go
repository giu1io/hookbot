package mailhook

// Sizer type
import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
)

type Sizer interface {
	Size() int64
}

// PdfExtractor class
type PdfExtractor struct {
}

// ParseZip method
func (e *PdfExtractor) ParseZip(file *multipart.FileHeader, destinationFolder string, c chan PdfFile) {
	var err error
	testPdf, _ := regexp.Compile("\\.pdf")

	openFile, err := file.Open()
	defer openFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	sr := openFile.(Sizer).Size()
	r, err := zip.NewReader(openFile, sr)

	if err != nil {
		fmt.Println(sr)
		fmt.Println(err)
		return
	}

	for _, zf := range r.File {
		if testPdf.FindStringIndex(zf.Name) != nil {
			path, err := e.writeFileOnDisk(zf, destinationFolder)
			if err != nil {
				fmt.Printf("Something went wrong writing unzipped file %s\n", err.Error())
			}
			c <- PdfFile{path: path}
		}
	}
}

// writeFileOnDisk method
func (e *PdfExtractor) writeFileOnDisk(file *zip.File, destination string) (string, error) {
	fpath := filepath.Join(destination, file.Name)
	memFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer memFile.Close()

	diskFile, err := os.OpenFile(
		fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return "", err
	}
	defer diskFile.Close()

	_, err = io.Copy(diskFile, memFile)
	if err != nil {
		return "", err
	}

	fmt.Printf("File %s wrote on disk successfully\n", fpath)
	return fpath, nil
}
