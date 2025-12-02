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
	sheet := r.URL.Query().Get("sheet")
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "./_works" // デフォルト
	}

	var records []reader.UsageRecord
	var err error

	if file != "" {
		// 特定のファイルのみ
		path := filepath.Join(dir, file)
		if strings.ToLower(filepath.Ext(path)) == ".xlsx" && sheet != "" {
			// シート指定がある場合はオプションで渡す
			records, err = utils.ReadRecord(path, reader.Options{SheetName: sheet})
		} else {
			records, err = utils.ReadRecord(path)
		}
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

// GetSheetsHandler returns the list of sheets for the specified Excel file.
// Query params: file (required), dir (optional)
func GetSheetsHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "./_works"
	}
	if file == "" {
		http.Error(w, "file parameter is required", http.StatusBadRequest)
		return
	}
	path := filepath.Join(dir, file)
	if strings.ToLower(filepath.Ext(path)) != ".xlsx" {
		http.Error(w, "sheets are only available for .xlsx files", http.StatusBadRequest)
		return
	}

	sheets, err := reader.GetSheetList(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing sheets: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"sheets": sheets})
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
