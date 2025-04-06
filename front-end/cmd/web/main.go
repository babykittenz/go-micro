package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}

	// Print the current working directory to debug
	fmt.Println("Current working directory:", cwd)

	// Check if we're already in cmd/web directory
	dirName := filepath.Base(cwd)
	parentDir := filepath.Base(filepath.Dir(cwd))

	var templatesDir string
	if parentDir == "cmd" && dirName == "web" {
		// We're already in cmd/web, so templates should be directly in templates/
		templatesDir = filepath.Join(cwd, "templates")
	} else {
		// We're in the project root, so use the full path
		templatesDir = filepath.Join(cwd, "cmd", "web", "templates")
	}

	fmt.Println("Templates directory:", templatesDir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml", templatesDir)
	})

	fmt.Println("Starting front end service on port 80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string, templatesDir string) {
	partials := []string{
		filepath.Join(templatesDir, "base.layout.gohtml"),
		filepath.Join(templatesDir, "header.partial.gohtml"),
		filepath.Join(templatesDir, "footer.partial.gohtml"),
	}

	var templateSlice []string
	templateSlice = append(templateSlice, filepath.Join(templatesDir, t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
