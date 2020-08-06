package main

import (
	"net/http"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"

	"myauth/client"
	"myauth/handler"

	"github.com/micro/go-micro/v2"

	myauth "myauth/proto/myauth"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.myauth"),
		micro.Version("latest"),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// Initialise service
	service.Init(
		// create wrap for the Myauth service client
		micro.WrapHandler(client.MyauthWrapper(service)),
	)
	go PrometheusBoot()
	// Register Handler
	myauth.RegisterMyauthHandler(service.Server(), new(handler.Myauth))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func PrometheusBoot() {
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("localhost:8085", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
