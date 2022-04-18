package render

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TplFile implements TemplateFile interface.
type TplFile struct {
	name string
	data []byte
	ext  string
}

func (f *TplFile) Name() string {
	return f.name
}

func (f *TplFile) Data() []byte {
	return f.data
}

func (f *TplFile) Ext() string {
	return f.ext
}

// NewTplFile cerates new template file with given name and data.
func NewTplFile(name string, data []byte, ext string) *TplFile {
	return &TplFile{name, data, ext}
}

// TplFileSystem implements TemplateFileSystem interface.
type TplFS struct {
	files []TemplateFile
}

func (fs TplFS) ListFiles() []TemplateFile {
	return fs.files
}

func (fs TplFS) Get(name string) (io.Reader, error) {
	for i := range fs.files {
		if fs.files[i].Name()+fs.files[i].Ext() == name {
			return bytes.NewReader(fs.files[i].Data()), nil
		}
	}
	return nil, fmt.Errorf("file '%s' not found", name)
}

// NewTemplateFileSystem creates new template file system with given options.
func NewTemplateFS(opt Options, omitData bool) TplFS {
	var err error
	fs := TplFS{}
	fs.files = make([]TemplateFile, 0, 10)

	// // Directories are composed in reverse order because later one overwrites previous ones,
	// // so once found, we can directly jump out of the loop.
	dirs := make([]string, 0, len(opt.AppendDirectories)+1)
	for i := len(opt.AppendDirectories) - 1; i >= 0; i-- {
		dirs = append(dirs, opt.AppendDirectories[i])
	}
	dirs = append(dirs, opt.Directory)

	// var err error
	for i := range dirs {
		// Skip ones that does not exists for symlink test,
		// but allow non-symlink ones added after start.
		if !IsExist(dirs[i]) {
			continue
		}

		dirs[i], err = filepath.EvalSymlinks(dirs[i])
		if err != nil {
			panic("EvalSymlinks(" + dirs[i] + "): " + err.Error())
		}
	}

	lastDir := dirs[len(dirs)-1]

	// We still walk the last (original) directory because it's non-sense we load templates not exist in original directory.
	if err := filepath.Walk(lastDir, func(path string, info os.FileInfo, _ error) error {

		r, err := filepath.Rel(lastDir, path)
		if err != nil {
			return err
		}

		fmt.Println("R", r, lastDir, path)

		ext := GetExt(r)

		for _, extension := range opt.Extensions {
			if ext != extension {
				continue
			}

			var data []byte
			if !omitData {
				// Loop over candidates of directory, break out once found.
				// The file always exists because it's inside the walk function,
				// and read original file is the worst case.
				// for i := range dirs {
				// path = filepath.Join(dirs[i], r)

				fmt.Println("path", path)
				if !IsFile(path) {
					continue
				}

				data, err = ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				// break
				// }
			}

			name := filepath.ToSlash((r[0 : len(r)-len(ext)]))

			fmt.Println("ddd:", name, data, ext)
			fs.files = append(fs.files, NewTplFile(name, data, ext))
		}

		return nil
	}); err != nil {
		panic("NewTemplateFileSystem: " + err.Error())
	}

	return fs
}
