package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sawaguchitake/usage-timeline/internal/reader"
	"github.com/sawaguchitake/usage-timeline/internal/utils"
)

func GetRecordsHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "./_works" // デフォルト
	}

	var records []reader.UsageRecord
	var err error

	if file != "" {
		// 特定のファイルのみ
		records, err = utils.ReadRecord(filepath.Join(dir, file))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file %s: %v", file, err), http.StatusInternalServerError)
			return
		}
	} else {
		// 全ファイル
		records, err = readAllRecords(dir)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading records: %v", err), http.StatusInternalServerError)
			return
		}
	}

	utils.SortRecords(records)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "./_works" // デフォルト
	}

	files, err := listFiles(dir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing files: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"files": files})
}

func listFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".csv" || ext == ".xlsx" {
			relPath, _ := filepath.Rel(dir, path)
			files = append(files, relPath)
		}
		return nil
	})
	return files, err
}

func readAllRecords(dir string) ([]reader.UsageRecord, error) {
	var allRecords []reader.UsageRecord

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".csv" || ext == ".xlsx" {
			records, err := utils.ReadRecord(path)
			if err != nil {
				log.Printf("Error reading %s: %v", path, err)
				return nil // 続行
			}
			allRecords = append(allRecords, records...)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allRecords, nil
}
