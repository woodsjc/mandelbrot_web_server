package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const validPathStr = `^/(:?\w+[.]png)?`

var validPath = regexp.MustCompile(validPathStr)

func checkInit() {
	if _, err := os.Stat("mandel.png"); os.IsNotExist(err) {
		generateMandelbrot()
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request, title string) {
	if len(title) > 0 && title[0] == '/' {
		title = "." + title
	} else {
		title = "./" + title
	}

	file, err := os.Stat(title)
	if err != nil {
		log.Printf("imageHandler unable to find %s", title)
		http.NotFound(w, r)
		return
	}

	size := strconv.FormatInt(file.Size(), 10)
	log.Printf("imageHandler serving: %s @ %s", title, size)
	http.ServeFile(w, r, title)
}

func baseHandler(w http.ResponseWriter, r *http.Request, title string) {
	content := `<html>
<head></head>
<body>
<p>Mandelbrot!!!
<br>
<br>
</p>
<img src="mandel.png" />
</body>
</html>`

	fmt.Fprintf(w, content)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Printf("URL path that doesn't exist: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(m[0], ".png") {
			fileHandler(w, r, m[0])
		} else {
			fn(w, r, m[0])
		}
	}
}

func main() {
	checkInit()

	http.HandleFunc("/", makeHandler(baseHandler))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
