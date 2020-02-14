package go_conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/spf13/viper"
	"strings"
)

type Confhub struct {
	filename   string
	filepath   string
	core       *viper.Viper
	filesuffix string
	fullpath   string
}

func New(file string, defConfFile ...string) *Confhub {
	var (
		tmp []string
		l   int
	)
	name := file
	path := "./"
	suffix := ""
	if strings.Contains(file, "/") {
		tmp := strings.Split(file, "/")
		l := len(tmp) - 1
		path = strings.Join(tmp[0:l], "/")
		name = tmp[l]
	}

	tmp = strings.Split(name, ".")
	l = len(tmp) - 1
	if l >= 1 {
		name = strings.Join(tmp[0:l], ".")
		suffix = tmp[l]
	}

	core := viper.New()
	core.SetConfigName(name)
	core.AddConfigPath(path)
	if suffix != "" {
		core.SetConfigType(suffix)
	} else {
		suffix = "toml"
	}
	if len(defConfFile) > 0 {
		defConf, err := New(defConfFile[0]).ReadConf()
		if err == nil {
			for k, v := range defConf {
				core.SetDefault(k, v)
			}
		}
	}
	fullpath := zfile.RealPath(path + "/" + name + "." + suffix)
	return &Confhub{filename: name, filepath: path, filesuffix: suffix, core: core, fullpath: fullpath}
}

func (c *Confhub) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return c.core.Unmarshal(rawVal, opts...)
}

func (c *Confhub) Object() *viper.Viper {
	return c.core
}

func (c *Confhub) ReadConf() (data map[string]interface{}, err error) {
	err = c.core.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// _ = c.core.WriteConfigAs("." + c.filesuffix)
			return
		}
		err = nil
	}
	data = c.core.AllSettings()
	return
}

func (c *Confhub) Exist() bool {
	return zfile.FileExist(c.fullpath)
}

func (c *Confhub) Set(key string, value interface{}) {
	c.core.Set(key, value)
}

func (c *Confhub) Get(key string) (value interface{}) {
	return c.core.Get(key)
}

func (c *Confhub) ConfigChange(fn func(e fsnotify.Event)) {
	c.core.WatchConfig()
	c.core.OnConfigChange(fn)
}

func (c *Confhub) GetAll() map[string]interface{} {
	return c.core.AllSettings()
}

func (c *Confhub) WriteConf(filepath ...string) error {
	if len(filepath) > 0 {
		return c.core.WriteConfigAs(filepath[0])
	}
	return c.core.WriteConfigAs(c.fullpath)
}
