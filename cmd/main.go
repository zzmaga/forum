package main

import (
	"log"

	"forum/architecture/repository"
	"forum/architecture/service"
	wh "forum/architecture/web/handler"
	"forum/architecture/web/server"
	database "forum/internal"
)

func main() {
	// init db with architecture migrations
	db, err := database.InitDB("forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// build repo and service
	repo := repository.NewRepo(db)
	svc := service.NewService(repo)

	// web handler configs
	hcfg := &wh.Configs{TemplatesDir: "ui/html", StaticFilesDir: "ui/static"}
	mh, err := wh.NewMainHandler(svc, hcfg)
	if err != nil {
		log.Fatal(err)
	}
	handler := mh.InitRoutes(hcfg)

	// server configs
	srvCfg := &server.Configs{
		Port:           ":8080",
		ReadTimeout:    5000,
		WriteTimeout:   5000,
		IdleTimeout:    60000,
		MaxHeaderBytes: 1 << 20,
	}

	srv := &server.Server{}
	if err := srv.Run(srvCfg, handler); err != nil {
		log.Fatal(err)
	}
}
