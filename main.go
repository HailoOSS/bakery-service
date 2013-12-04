package main

import (
	log "github.com/cihub/seelog"
	"github.com/hailocab/{{REPONAME}}/handler"
	service "github.com/hailocab/go-platform-layer/server"
)

func main() {
	defer log.Flush()

	service.Name = "com.hailocab.service.{{SERVICENAME}}"
	service.Description = "Please provide a short description of what your service does. It should be about this long."
	service.Version = ServiceVersion
	service.Source = "github.com/hailocab/{{REPONAME}}"
	service.OwnerEmail = "jonathan@hailocab.com"
	service.OwnerMobile = "+4407546186424"

	service.Init()

	service.Register(&service.Endpoint{
		Name:    "foo",
		Mean:    50,
		Upper95: 100,
		Handler: handler.Foo,
		// remember to choose an appropriate authoriser for each endpoint
		Authoriser: service.RoleAuthoriser([]string{"ADMIN"}),
	})

	service.Run()
}
