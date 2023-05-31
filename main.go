package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// UploadedFile represents the information of an uploaded file
type UploadedFile struct {
	Name          string
	Size          int64
	SizeFormatted string
	Progress      float64
}

var (
	uploadedFiles = make(map[string]*UploadedFile)
	filesMutex    = &sync.Mutex{}
)

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create("uploads/" + file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()

	// Copy the file data and track progress
	progressWriter := &ProgressWriter{Writer: dst, Size: file.Size}
	_, err = io.Copy(progressWriter, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store the uploaded file information
	filesMutex.Lock()
	uploadedFiles[file.Filename] = &UploadedFile{
		Name:     file.Filename,
		Size:     file.Size,
		Progress: 100.0,
	}
	filesMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

// Function to format file size
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return strconv.FormatInt(size, 10) + " B"
	}
	div := int(math.Log(float64(size)) / math.Log(unit))
	sizeFormatted := float64(size) / math.Pow(unit, float64(div))
	unitString := []string{"KB", "MB", "GB"}[div-1]
	return strconv.FormatFloat(sizeFormatted, 'f', 2, 64) + " " + unitString
}

func listFiles(c *gin.Context) {
	filesMutex.Lock()
	defer filesMutex.Unlock()

	files, err := os.ReadDir("uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fileData []UploadedFile
	for _, file := range files {
		if !file.IsDir() {
			filePath := fmt.Sprintf("uploads/%s", file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Failed to get file info for %s: %v\n", filePath, err)
				continue
			}

			fileSize := fileInfo.Size()
			fileSizeFormatted := formatFileSize(fileSize)

			fileData = append(fileData, UploadedFile{
				Name:          file.Name(),
				Size:          fileSize,
				SizeFormatted: fileSizeFormatted,
				Progress:      100.0,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"files": fileData})
}

func downloadFile(c *gin.Context) {
	filename := c.Query("filename")
	filePath := fmt.Sprintf("uploads/%s", filename)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set the appropriate headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

type ProgressWriter struct {
	io.Writer
	Size      int64
	BytesRead int64
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)
	pw.BytesRead += int64(n)
	progress := float64(pw.BytesRead) / float64(pw.Size) * 100.0
	fmt.Printf("Upload progress: %.2f%%\n", progress)
	return
}

func main() {
	r := gin.Default()

	r.POST("/upload", uploadFile)
	r.GET("/files", listFiles)
	r.GET("/download", downloadFile)

	// Serve frontend files
	r.Static("/static", "./fe/static")
	r.LoadHTMLFiles("./fe/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	fmt.Println("Server running on port 8000...")
	r.Run(":8000")
}
