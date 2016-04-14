package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

var contentPath = flag.String("path", "./files", "directory to serve")

const (
	commonHtmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_HREF_TARGET_BLANK

	commonExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
)

type File struct {
	Path string
	Name string
}

type Note struct {
	Url         string
	Dirs        []*File
	Files       []*File
	ContentHTML template.HTML
	Content     string
	Edit        bool
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requested: %s", r.URL.Path)
	url := r.URL.Path

	// Ignore favicon
	if url == "/favicon.ico" {
		return
	}

	// Tailing "/" = directory, otherwise serve a .md-File (TODO: only, if corresponding file not found)
	if strings.HasSuffix(url, "/") {
		url += "index.md"
	} else if !strings.HasSuffix(url, ".md") {
		url += ".md"
	}

	// If file not found and directory exists, serve directories index
	// otherwise recurse until directory exists
	if _, err := os.Stat(*contentPath + url); os.IsNotExist(err) {
		urlDir, _ := path.Split(url)
		if _, err := os.Stat(*contentPath + urlDir); err == nil {
			url = urlDir + "index.md"
		} else {
			for len(urlDir) > 1 {
				urlDir, _ = path.Split(urlDir[:len(urlDir)-1])
				if _, err := os.Stat(*contentPath + urlDir); err == nil {
					url = urlDir + "index.md"
					break
				}
			}
			if len(urlDir) == 1 {
				url = "/index.md"
			}
		}
	}

	var edit bool
	if r.FormValue("edit") == "true" {
		edit = true
	}
	content := r.FormValue("content")
	file := sanitizeFilename(r.FormValue("file"))

	switch r.FormValue("action") {
	default:
	case "save_edit":
		edit = true
		fallthrough
	case "save":
		err := ioutil.WriteFile(*contentPath+url, []byte(content), 0644)
		if err != nil {
			log.Println("WriteFile error:", err)
		}
	case "create_folder":
		urlDir, _ := path.Split(url)
		err := os.Mkdir(*contentPath+urlDir+file, 0755)
		if err != nil {
			log.Println("Unable to create dir:", err)
		}
		url = urlDir + file + "/"

	case "create_file":
		if !strings.HasSuffix(file, ".md") {
			file += ".md"
		}
		urlDir, _ := path.Split(url)
		err := ioutil.WriteFile(*contentPath+urlDir+file, []byte{}, 0644)
		if err != nil {
			log.Println("CreateFile error:", err)
		}
		url = urlDir + file
	}

	if r.URL.Path != url {
		log.Printf("Redirect to: %s", url)
		http.Redirect(w, r, url, 302)
	}

	// We already tried hard to find a regular file, if still no file is found, serve empty content
	fileContent, _ := ioutil.ReadFile(*contentPath + url)

	renderer := blackfriday.HtmlRenderer(commonHtmlFlags, "", "")
	contentHTML := template.HTML(string(blackfriday.MarkdownOptions(fileContent, renderer, blackfriday.Options{Extensions: commonExtensions})))

	dirs := make([]*File, 0)
	files := make([]*File, 0)
	urlDir, _ := path.Split(url)
	if urlDir != "/" {
		dirs = append(dirs, &File{Path: ".." + "/", Name: ".."})
	}

	fs, err := filepath.Glob(*contentPath + urlDir + "*")
	if err != nil {
		log.Println("Glob error:", err)
	}

	replacer := strings.NewReplacer("_", " ", ".", " ")

	for _, f := range fs {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		if info.IsDir() {
			dirs = append(dirs, &File{Path: info.Name() + "/", Name: info.Name()})
		} else {
			if strings.HasSuffix(info.Name(), ".md") {
				files = append(files, &File{Path: info.Name(), Name: replacer.Replace(strings.TrimSuffix(info.Name(), ".md"))})
			}
		}
	}

	note := &Note{Url: url, Files: files, Dirs: dirs, ContentHTML: contentHTML, Content: string(fileContent), Edit: edit}

	renderTemplate(w, note)
	log.Printf("Served %s", url)
}

func renderTemplate(w http.ResponseWriter, note *Note) {
	t := template.New("notes")
	var err error

	t, err = template.ParseFiles("templates/content.tpl")
	if err != nil {
		log.Print("Could not parse template", err)
	}

	t.ParseFiles("templates/header.tpl", "templates/footer.tpl")
	err = t.Execute(w, note)
	if err != nil {
		log.Print("Could not execute template: ", err)
	}
}

func sanitizeFilename(filename string) string {
	return filename
}

func main() {
	http.HandleFunc("/", handler)

	// Static resources
	http.Handle("/_static/", http.StripPrefix("/_static/", http.FileServer(http.Dir("./static/"))))

	var bind = flag.String("bind", ":8000", "ip/port to bin, example: 0.0.0.0:8000")

	flag.Parse()
	var err error

	if *bind != "" {
		err = http.ListenAndServe(*bind, nil)
	}
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
