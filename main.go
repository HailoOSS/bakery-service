package main

import (
	"github.com/hailocab/{{REPONAME}}/handler"
	service "github.com/hailocab/go-platform-layer/server"
)

func main() {
	service.Name = "com.hailocab.service.{{SERVICENAME}}"
	service.Description = "Please provide a short description of what your service does. It should be about this long."
	service.Version = ServiceVersion
	service.Source = "github.com/hailocab/{{REPONAME}}"
	service.OwnerEmail = "youremail@hailocab.com"
	service.OwnerMobile = "+44123412341234"

	service.Init()

	service.Register(&service.Endpoint{
		Name:    "foo",
		Mean:    10,
		Upper95: 20,
		Handler: handler.Foo,
	})

	service.Run()
}
