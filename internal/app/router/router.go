package router

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/midoks/simdht/internal/conf"
	"github.com/midoks/simdht/internal/dht"
	"github.com/midoks/simdht/internal/logs"
	"github.com/midoks/simdht/internal/tools"
)

func autoMakeCustomConf(customConf string) error {

	if tools.IsExist(customConf) {
		return nil
	}

	// auto make custom conf file
	cfg := ini.Empty()
	if tools.IsFile(customConf) {
		// Keeps custom settings if there is already something.
		if err := cfg.Append(customConf); err != nil {
			return err
		}
	}

	cfg.Section("").Key("app_name").SetValue("simdht")
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("web").Key("http_port").SetValue("11010")
	cfg.Section("session").Key("provider").SetValue("memory")

	os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
	if err := cfg.SaveTo(customConf); err != nil {
		return err
	}

	return nil
}

func Init(customConf string) error {
	var err error

	if strings.EqualFold(customConf, "") {
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

	//stat DHT
	dht.Run()

	if strings.EqualFold(conf.App.RunMode, "dev") {
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}

	return nil
}
