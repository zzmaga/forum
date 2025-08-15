package handler

import (
	"net/http"

	"forum/architecture/service"
	"forum/architecture/web/handler/view"
)

type Configs struct {
	TemplatesDir   string `cenv:"templates"`
	StaticFilesDir string `cenv:"static_files"`
}

type MainHandler struct {
	// templates    *template.Template
	view    view.View
	service *service.Service
}

func NewMainHandler(service *service.Service, configs *Configs) (*MainHandler, error) {
	mh := &MainHandler{
		view:    *view.NewView(configs.TemplatesDir),
		service: service,
	}
	return mh, nil
}

func (m *MainHandler) InitRoutes(configs *Configs) http.Handler {
	mux := http.NewServeMux()
	// HERE IS ALL ROUTES
	fsStatic := http.FileServer(http.Dir(configs.StaticFilesDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	// AnyRoutes
	mux.HandleFunc("/", m.IndexHandler)
	mux.HandleFunc("/signup", m.SignUpHandler)
	mux.HandleFunc("/signin", m.SignInHandler)
	mux.HandleFunc("/signout", m.SignOutHandler)

	mux.Handle("/post/get", http.HandlerFunc(m.PostViewHandler))
	mux.Handle("/post/create", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostCreateHandler)))
	mux.Handle("/post/edit", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostEditHandler)))
	mux.Handle("/post/vote", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostVoteHandler)))
	mux.Handle("/post/delete", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostDeleteHandler)))

	mux.Handle("/posts/own", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostsOwnHandler)))
	mux.Handle("/posts/voted", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostsVotedHandler)))

	mux.Handle("/post/comment/create", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostCommentCreateHandler)))
	mux.Handle("/post/comment/delete", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostCommentDeleteHandler)))
	mux.Handle("/post/comment/vote", m.MiddlewareSessionChecker(http.HandlerFunc(m.PostCommentVoteHandler)))

	mux.Handle("/categories/posts", http.HandlerFunc(m.CategoriesPostsHandler))

	return mux
}
