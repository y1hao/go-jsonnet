package trace

import (
	_ "embed"

	"log"
	"net/http"
	"text/template"

	"github.com/google/go-jsonnet"
)

//go:embed index.html
var templ string

type Server struct {
	filename string
	result   string
	frames   map[int][]*jsonnet.TraceItem
}

func NewServer(filename string, result string, frames map[int][]*jsonnet.TraceItem) *Server {
	return &Server{
		filename: filename,
		result:   result,
		frames:   frames,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Filename string
		Result   string
		Frames   map[int][]*jsonnet.TraceItem
	}{
		Filename: s.filename,
		Result:   s.result,
		Frames:   s.frames,
	}

	tmpl, err := template.New("root").Parse(templ)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "root", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
