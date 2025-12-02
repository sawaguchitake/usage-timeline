package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/sawaguchitake/usage-timeline/internal/web"
)

func main() {
	// 静的ファイルの提供
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API
	http.HandleFunc("/api/records", web.GetRecordsHandler)
	http.HandleFunc("/api/files", web.GetFilesHandler)
	http.HandleFunc("/api/sheets", web.GetSheetsHandler)

	// ルートで静的ファイルを扱う（index.htmlが自動で提供される）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
