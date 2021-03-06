package router

import (
	"fmt"
	"os"
	"path/filepath"

	"net/http"
	_ "net/http/pprof"

	"gopkg.in/ini.v1"

	"github.com/midoks/simdht/internal/app/template"
	"github.com/midoks/simdht/internal/assets/templates"
	"github.com/midoks/simdht/internal/conf"
	"github.com/midoks/simdht/internal/dht"
	"github.com/midoks/simdht/internal/logs"
	"github.com/midoks/simdht/internal/mgdb"
	"github.com/midoks/simdht/internal/render"
	"github.com/midoks/simdht/internal/tools"
)

func autoMakeCustomConf(customConf string) error {

	if tools.IsExist(customConf) {
		return nil
	}

	// auto make custom conf file
	cfg := ini.Empty()
	if tools.IsFile(customConf) {
		if err := cfg.Append(customConf); err != nil {
			return err
		}
	}

	cfg.Section("").Key("app_name").SetValue("simdht")
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("web").Key("http_port").SetValue("11010")
	cfg.Section("session").Key("provider").SetValue("memory")

	cfg.Section("mongodb").Key("addr").SetValue("127.0.0.1:27017")
	cfg.Section("mongodb").Key("db").SetValue("simdht")

	os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
	if err := cfg.SaveTo(customConf); err != nil {
		return err
	}

	return nil
}

func Init(customConf string) error {
	var err error

	if customConf == "" {
		customConf = filepath.Join(conf.CustomDir(), "conf", "app.conf")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return fmt.Errorf("custom conf file get absolute path: %s", err)
		}
	}

	err = autoMakeCustomConf(customConf)
	if err != nil {
		return err
	}

	conf.Init(customConf)
	logs.Init()
	mgdb.Init()

	renderOpt := render.Options{
		Directory:         filepath.Join(conf.WorkDir(), "templates"),
		AppendDirectories: []string{filepath.Join(conf.CustomDir(), "templates")},
		Funcs:             template.FuncMap(),
		IndentJSON:        true,
	}

	if !conf.Web.LoadAssetsFromDisk {
		renderOpt.FileSystem = templates.NewTemplateFileSystem("", renderOpt.AppendDirectories[0])
	}

	render.Renderer(renderOpt)

	//stat DHT
	go dht.Run()

	// go tool pprof --seconds 30 http://localhost:6060/debug/pprof/profile
	// go tool pprof -http=:11011 --seconds 30 http://localhost:6060/debug/pprof/profile

	// go tool pprof -pdf profile.out > cpu.pdf
	// go tool pprof -pdf memprofile.out > mem.pdf
	if conf.App.RunMode != "prod" {
		go func() {
			port := fmt.Sprintf(":%s", conf.Debug.Port)
			http.ListenAndServe(port, nil)
		}()
	}

	return nil
}
