package mailhook

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Mailhook class
type Mailhook struct {
	ext      PdfExtractor
	filesSub chan PdfFile
	authKeys []string
}

// GetPdfs method
func (h *Mailhook) GetPdfs() chan PdfFile {
	return h.filesSub
}

// EnableWebHook method
func (h *Mailhook) EnableWebHook(endpoint string, host string, destFolder string, authKeys []string) {
	fmt.Println("Starting server...")
	h.authKeys = authKeys

	r := gin.Default()
	r.Use(h.authMiddleware)

	r.POST(endpoint, func(c *gin.Context) {
		file, _ := c.FormFile("file")

		if file == nil || file.Size == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "No attachment provided!",
			})
			return
		}

		log.Println(file.Filename)
		log.Println(file.Header.Get("Content-Type"))

		if strings.Compare(file.Header.Get("Content-Type"), "application/zip") == 0 {
			go h.ext.ParseZip(file, destFolder, h.filesSub)
		} else if strings.Compare(file.Header.Get("Content-Type"), "application/pdf") == 0 {
			path := filepath.Join(destFolder, file.Filename)
			err := c.SaveUploadedFile(file, path)
			if err != nil {
				fmt.Printf("Something went wrong writing file %s\n", err.Error())
			} else {
				h.filesSub <- PdfFile{path: path}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})

	r.Run(host)
}

func (h *Mailhook) authMiddleware(c *gin.Context) {
	authkey := c.GetHeader("authkey")
	authorized := false
	for i := range h.authKeys {
		if h.authKeys[i] == authkey {
			authorized = true
		}
	}
	if authorized {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}

}

// MailhookBuilder method
func MailhookBuilder() Mailhook {
	return Mailhook{PdfExtractor{}, make(chan PdfFile), []string{}}
}
