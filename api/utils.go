package api

import (
	"net/http"
	"os"
	"path"
)

// ShowError showing error message as json
func ShowError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8;")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"message": "` + msg + `"}`))
}

// Redirect404 redirect into 404 page
func Redirect404(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+r.Host+"/error/404", http.StatusFound)
}

// ResponseJSON response as JSON
func ResponseJSON(w http.ResponseWriter, status int, response []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(response)
}

func GetBaseDir() (string, error) {
	return os.Getwd()
}

func PrepareUploadDir(dirName string) error {
	dir, err := GetBaseDir()
	if err != nil {
		return err
	}

	dirPath := path.Join(dir, dirName)

	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dirPath, 0777)
}
