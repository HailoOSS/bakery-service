package main

import (
	"time"

	protoBuild "github.com/hailocab/bakery-service/proto/build"

	"github.com/hailocab/bakery-service/aws"
	"github.com/hailocab/bakery-service/elastic"
	"github.com/hailocab/bakery-service/handler"

	log "github.com/cihub/seelog"
	service "github.com/hailocab/platform-layer/server"
	"github.com/hailocab/service-layer/config"
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
		Authoriser:       service.RoleAuthoriser([]string{"ADMIN", "PLATFORM"}),
		Handler:          handler.Build,
		Mean:             50,
		Name:             "build",
		RequestProtocol:  new(protoBuild.Request),
		ResponseProtocol: new(protoBuild.Response),
		Upper95:          100,
	})

	config.WaitUntilLoaded(time.Second * 2)

	aws.Init()
	elastic.Init()

	service.Run()
}
