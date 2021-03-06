package readconfig

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

// Config - структура для считывания конфигурационного файла
type Config struct {
	Db    string `yaml:"db" json:"db"`
	Port  uint   `yaml:"port" json:"port" `
	Host  string `yaml:"dbhost" json:"dbhost"`
	User  string `yaml:"dbuser" json:"dbuser"`
	Pass  string `yaml:"dbpass" json:"dbpass"`
	Limit uint   `yaml:"limit" json:"limit"`
}

func (e *Config) Validate() error {
	var err error
	/*
		err = e.CheckUrl(e.Db_url)
		if err != nil {
			return err
		}

		err = e.CheckUrl(e.Jaeger_url)
		if err != nil {
			return err
		}

		err = e.CheckUrl(e.Sentry_url)
		if err != nil {
			return err
		}

		err = e.CheckUrl(e.Kafka_broker)
		if err != nil {
			return err
		}
	*/
	return err
}

func (e *Config) SetPort(p uint) error {

	e.Port = p
	return nil

}

/*
func (e *Config) SetDb(p string) error {

	err := e.CheckUrl(p)
	if err == nil {
		e.Db_url = p
		return nil
	} else {
		return err

	}

}
*/

func (e *Config) CheckUrl(path string) error {

	_, err := url.ParseRequestURI(path)

	if err != nil {
		return err
	}
	return nil

}

// GetBaseFile возвращаем имя программы
func GetBaseFile() string {
	filename := os.Args[0] // get command line first parameter
	return strings.Split(filepath.Base(filename), ".")[0]
}

// GetDefaultConfigFile - возвращаем полное имя конфиг файла
func GetDefaultConfigFile() string {

	return GetBaseFile() + ".json"
}

func ReadConfig(ConfigName string) (x *Config, err error) {
	var file []byte
	if file, err = ioutil.ReadFile(ConfigName); err != nil {
		return nil, err
	}
	x = new(Config)
	switch strings.ToLower(path.Ext(ConfigName)) {

	case ".yaml", ".yml":
		err = yaml.Unmarshal(file, &x)
	case ".json":
		err = json.Unmarshal(file, &x)

	case ".ini":
		cfg, err := ini.Load(ConfigName)
		if err == nil {
			x.Port = cfg.Section("").Key("port").MustUint()
			//x.Port2 = cfg.Section("").Key("port2").MustUint()
			//x.Db_url = cfg.Section("").Key("db_url").String()
			//x.Jaeger_url = cfg.Section("").Key("jaeger_url").String()
			//x.Sentry_url = cfg.Section("").Key("sentry_url").String()
			//x.Kafka_broker = cfg.Section("").Key("kafka_broker").String()
			//x.Some_app_id = cfg.Section("").Key("some_app_id").String()
			//x.Some_app_key = cfg.Section("").Key("some_app_key").String()
		}

	}

	if err != nil {
		return nil, err
	}

	if err = x.Validate(); err != nil {
		return nil, err
	}

	return x, nil
}
