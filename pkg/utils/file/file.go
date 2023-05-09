package file

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadAndSaveFile(fileURL, savePath string) error {
	// Создание директории, если она не существует
	dir := filepath.Dir(savePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
