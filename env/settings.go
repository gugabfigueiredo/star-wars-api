package env

import (
	"encoding/json"
	"github.com/gugabfigueiredo/star-wars-api/log"
)

type settings struct {
	Log *log.Config

	Server struct {
		Port string `default:"8080"`
		Context string `default:"sw-api"`
	}

	Database struct {
		UserName        string `default:"user"`
		Password        string `default:"pass"`
		Host            string `default:"localhost"`
		Port            string `default:"5432"`
		DatabaseName    string `default:"sw-api"`
		SetMaxOpenConns int    `default:"100"`
	}
}


var Settings settings

func (s settings) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}