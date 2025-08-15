package handler

import (
	"forum/architecture/service"
	"forum/architecture/web/handler/view"
)

// TODO add configs from env.configs
type Configs struct {
	TemplatesDir   string // `env:"templates"`
	StaticFilesDir string // `env:"static_files`
}
type MainHandler struct {
	view    view.View
	service service.Service
}

func NewMainHandler(service *service.Service, configs *Configs) (*MainHandler, error) {
	mh := &MainHandler{
		view:    *view.NewView{configs.TemplatesDir},
		service: service,
	}
	return mh, nil
}
