package service

import (
	"errors"
	"github.com/gugabfigueiredo/star-wars-api/log"
)

var ErrMissingUser = errors.New("Missing User")

type IHelloService interface {
	SayHello(msisdn string) (string, error)
}

type HelloService struct {
	Logger *log.Logger
}

func (service *HelloService) SayHello(user string) (string, error) {

	if len(user) == 0 {
		service.Logger.E("Error saying hello.", "reason", "Missing user.")
		return "", ErrMissingUser
	}

	message := "hello " + user

	service.Logger.I("Hello repository is saying 'hello'.", "response", message)

	return message, nil
}