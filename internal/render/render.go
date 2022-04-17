package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"sync"
)

const (
	_TMPL_DIR = "templates"

	_CONTENT_TYPE    = "Content-Type"
	_CONTENT_BINARY  = "application/octet-stream"
	_CONTENT_JSON    = "application/json"
	_CONTENT_HTML    = "text/html"
	_CONTENT_PLAIN   = "text/plain"
	_CONTENT_XHTML   = "application/xhtml+xml"
	_CONTENT_XML     = "text/xml"
	_DEFAULT_CHARSET = "UTF-8"
)

var (
	// Provides a temporary buffer to execute templates into and catch errors.
	bufpool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}

	// Included helper functions for use when rendering html
	helperFuncs = template.FuncMap{
		"yield": func() (string, error) {
			return "", fmt.Errorf("yield called with no layout defined")
		},
		"current": func() (string, error) {
			return "", nil
		},
	}
)

type (

	// TemplateFile represents a interface of template file that has name and can be read.
	TemplateFile interface {
		Name() string
		Data() []byte
		Ext() string
	}
	// TemplateFileSystem represents a interface of template file system that able to list all files.
	TemplateFileSystem interface {
		ListFiles() []TemplateFile
		Get(string) (io.Reader, error)
	}

	// Delims represents a set of Left and Right delimiters for HTML template rendering
	Delims struct {
		// Left delimiter, defaults to {{
		Left string
		// Right delimiter, defaults to }}
		Right string
	}

	// RenderOptions represents a struct for specifying configuration options for the Render middleware.
	Options struct {
		// Directory to load templates. Default is "templates".
		Directory string
		// Addtional directories to overwite templates.
		AppendDirectories []string
		// Layout template name. Will not render a layout if "". Default is to "".
		Layout string
		// Extensions to parse template files from. Defaults are [".tmpl", ".html"].
		Extensions []string
		// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Default is [].
		Funcs []template.FuncMap
		// Delims sets the action delimiters to the specified strings in the Delims struct.
		Delims Delims
		// Appends the given charset to the Content-Type header. Default is "UTF-8".
		Charset string
		// Outputs human readable JSON.
		IndentJSON bool
		// Outputs human readable XML.
		IndentXML bool
		// Prefixes the JSON output with the given bytes.
		PrefixJSON []byte
		// Prefixes the XML output with the given bytes.
		PrefixXML []byte
		// Allows changing of output to XHTML instead of HTML. Default is "text/html"
		HTMLContentType string
		// TemplateFileSystem is the interface for supporting any implmentation of template file system.
		TemplateFileSystem
	}
)

var renderOption Options

func init() {
	renderOption = Options{
		IndentJSON:      true,
		HTMLContentType: _CONTENT_HTML,
		Directory:       _TMPL_DIR,
		Extensions:      []string{".tmpl", ".html"},
	}
}

func Renderer(op Options) {

}

func HTML(status int, name string) {
	t1, err := template.ParseFiles(name)
	if err != nil {
		// fmt.Println(err)
	}
	fmt.Println(t1, err)

	fmt.Println(renderOption)
	// t1.Execute(w, "hello world")
}
