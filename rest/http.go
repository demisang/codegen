package rest

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"time"

	"github.com/demisang/codegen"
	"github.com/sirupsen/logrus"
)

type Server struct {
	app    app
	server *http.Server
	log    *logrus.Logger
}

type app interface {
	GetTemplatesList(ctx context.Context) ([]codegen.Template, error)
	RawList(ctx context.Context, templateID, targetRelativeRootDir string, placeholders []codegen.Placeholder) ([]codegen.PreviewListItem, error)
	PreviewList(ctx context.Context, templateID, targetRelativeRootDir string, placeholders []codegen.Placeholder) ([]codegen.PreviewListItem, error)
	Generate(ctx context.Context, templateID, targetRelativeRootDir string, placeholders []codegen.Placeholder) (string, error)
	GetDirectories(ctx context.Context, selectedDir string) ([]string, error)
}

const (
	readHeaderTimeout = 30 * time.Second
)

//go:embed public
var public embed.FS

func NewServer(app app, log *logrus.Logger, host string, port int) *Server {
	s := Server{
		app: app,
		log: log,
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	mux := http.NewServeMux()

	// fs := http.FileServer(http.Dir("./public"))
	// mux.Handle("/", http.FileServer(http.Dir("./public")))
	sub, _ := fs.Sub(public, "public")
	mux.Handle("/", http.FileServer(http.FS(sub)))
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	// mux.HandleFunc("/", fs.ServeHTTP)
	// mux.HandleFunc("/", s.loggingMiddleware(http.HandlerFunc(s.linkGet)).ServeHTTP)
	mux.HandleFunc("/templates", s.loggingMiddleware(s.enableCors(
		http.HandlerFunc(s.templates),
	)).ServeHTTP)
	mux.HandleFunc("/raw-list", s.loggingMiddleware(s.enableCors(
		http.HandlerFunc(s.rawList),
	)).ServeHTTP)
	mux.HandleFunc("/preview-list", s.loggingMiddleware(s.enableCors(
		http.HandlerFunc(s.previewList),
	)).ServeHTTP)
	mux.HandleFunc("/generate", s.loggingMiddleware(s.enableCors(
		http.HandlerFunc(s.generate),
	)).ServeHTTP)
	mux.HandleFunc("/directories", s.loggingMiddleware(s.enableCors(
		http.HandlerFunc(s.directories),
	)).ServeHTTP)

	s.server = &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readHeaderTimeout}

	return &s
}

func (s *Server) enableCors(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// "OPTIONS" request just OK response
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s.log.Infof("request %s", r.RequestURI)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.server.Shutdown(ctx)
		if err != nil {
			s.log.Errorf("server shutdown: %v", err)
		}
	}()

	s.log.Infof("starting server %s", s.server.Addr)

	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		s.log.Infof("server closed")
	} else if err != nil {
		s.log.Infof("error starting server: %v", err)
	}

	return err
}

func errResponse(w http.ResponseWriter, _ *http.Request, statusCode int, err error) {
	http.Error(w, err.Error(), statusCode)
}

func response(w http.ResponseWriter, r *http.Request, statusCode int, content any) {
	body, err := json.Marshal(content)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(statusCode)
	_, _ = io.WriteString(w, string(body))
}
