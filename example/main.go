package main

import (
	"github.com/corneldamian/gowebservice"
	"github.com/kardianos/service"

	"log"
	"os"

	"github.com/corneldamian/gowebservice/example/config"
	"github.com/corneldamian/gowebservice/example/routing"
)

func main() {
	cfg := config.GetConfig()

	svcConfig := &service.Config{
		Name:        "GoWebService",
		DisplayName: "Go Web Service Example",
		Description: "This is an example Go Web service.",
	}

	s, err := gowebservice.DoInit(cfg, svcConfig, routing.MiddlewaresFactory, routing.GetRoutesFactory())
	if err != nil {
		return
	}

	if len(os.Args) > 1 {
		if err := service.Control(s, os.Args[1]); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
