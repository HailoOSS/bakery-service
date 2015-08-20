package main

import (
	log "github.com/cihub/seelog"

	"github.com/hailocab/bakery-service/handler"
	service "github.com/hailocab/go-platform-layer/server"
)

func main() {
	defer log.Flush()

	service.Name = "com.hailocab.infrastructure.bakery"
	service.Description = "Makes artefacts"
	service.Version = ServiceVersion
	service.Source = "github.com/hailocab/bakery-service"
	service.OwnerEmail = "platform@hailocab.com"
	service.OwnerTeam = "Platform"

	service.Init()

	service.Register(&service.Endpoint{
		Name:       "foo",
		Mean:       50,
		Upper95:    100,
		Handler:    handler.Foo,
		Authoriser: service.RoleAuthoriser([]string{"ADMIN", "PLATFORM"}),
	})

	service.Run()
}
