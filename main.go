package main

import (
	"github.com/hailocab/{{REPONAME}}/handler"
	service "github.com/hailocab/go-platform-layer/server"
)

func main() {
	server.Name = "com.hailocab.service.{{REPONAME}}"
	server.Description = "Please provide a short description of what your service does. It should be about this long."
	server.Version = 20130524110011
	server.Source = "github.com/hailocab/{{REPONAME}}"
	server.OwnerEmail = "youremail@hailocab.com"
	server.OwnerMobile = "+44123412341234"

	service.Register(&service.Endpoint{
		Name:    "foo",
		Mean:    10,
		Upper95: 20,
		Handler: handler.Foo,
	})

	service.Run()
}
