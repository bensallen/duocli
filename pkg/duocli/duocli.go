package duocli

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/admin"
)

// Config is a struct resprenting the JSON configuration file for duocli.
// UserAgent and Timeout are optional fields.
type Config struct {
	Ikey      string `json:"ikey"`
	Skey      string `json:"skey"`
	APIHost   string `json:"api_host"`
	UserAgent string `json:"user_agent"`
	Timeout   int64  `json:"timeout"`
}

// LoadAdminConfig loads a JSON configuration file from path and returns a populated *admin.Client from duo_api_golang.
func LoadAdminConfig(path string) (*admin.Client, error) {
	if path == "" {
		path = ".duocli.json"
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	duoapi, err := loadConfig(file)
	if err != nil {
		return nil, err
	}
	return admin.New(*duoapi), nil
}

// LoadConfig loads a JSON configuration file from path and returns a populated *duoapi.DuoApi.
func LoadConfig(path string) (*duoapi.DuoApi, error) {
	if path == "" {
		path = ".duocli.json"
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return loadConfig(file)
}

func loadConfig(file io.Reader) (*duoapi.DuoApi, error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	if conf.Ikey == "" {
		return nil, fmt.Errorf("ikey not specified in configuration file")
	}
	if conf.Skey == "" {
		return nil, fmt.Errorf("ikey not specified in configuration file")
	}
	if conf.APIHost == "" {
		return nil, fmt.Errorf("api_host not specified in configuration file")
	}
	if conf.UserAgent == "" {
		conf.UserAgent = "duocli"
	}
	if conf.Timeout == 0 {
		conf.Timeout = 10
	}

	return duoapi.NewDuoApi(
		conf.Ikey,
		conf.Skey,
		conf.APIHost,
		conf.UserAgent,
		duoapi.SetTimeout(time.Duration(conf.Timeout)*time.Second),
	), nil
}
