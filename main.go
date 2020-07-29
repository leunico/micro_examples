package main

import (
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/go-micro/v2"
	"myauth/handler"
	"myauth/client"

	myauth "myauth/proto/myauth"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.myauth"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the Myauth service client
		micro.WrapHandler(client.MyauthWrapper(service)),
	)

	// Register Handler
	myauth.RegisterMyauthHandler(service.Server(), new(handler.Myauth))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
