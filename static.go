package Form_Mailly_Go

import (
	"embed"
	"net/http"
	"strconv"
)

//go:embed public/index.html
var publicFiles embed.FS

var indexHTML []byte

func init() {
	b, err := publicFiles.ReadFile("public/index.html")
	if err != nil {
		panic(err) // Fail fast if the resource isn't embedded
	}
	indexHTML = b
}

func HomeHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.Header().Set("Content-Length", strconv.Itoa(len(indexHTML)))
	response.WriteHeader(http.StatusOK)
	_, err := response.Write(indexHTML)
	if err != nil {
		return
	}
}
