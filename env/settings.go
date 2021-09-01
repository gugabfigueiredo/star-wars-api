package env

import (
	"encoding/json"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/repository"
	"time"
)

type settings struct {
	Log *log.Config

	Server struct {
		Port              string        `default:"8080"`
		Context           string        `default:"sw-api"`
		UpdateRefsTimeout time.Duration `default:"4h"`
	}

	Database *repository.Config
}


var Settings settings

func (s settings) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}