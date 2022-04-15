package render

import (
	"html/template"
	"io"
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
		// ListFiles() []TemplateFile
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
		// Delims Delims
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

func Html(status int, name string) {

}

func Renderer(op Options) {

}
