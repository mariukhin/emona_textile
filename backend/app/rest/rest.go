package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"backend/app/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Rest struct {
	AllowedOrigins []string

	SecureCookiesDisabled bool

	CarouselStore model.CarouselStore
	LinkBuilder   LinkBuilder

	AppLog *logrus.Logger

	carouselManagement *carouselManagement

	httpServer *http.Server
	lock       sync.Mutex
}

func (s *Rest) Run(port int) {
	s.AppLog.Infof("Start EmonaTextile REST API server on port %d", port)

	// set app log as default logger
	if s.AppLog.IsLevelEnabled(logrus.WarnLevel) {
		log.SetOutput(s.AppLog.WriterLevel(logrus.WarnLevel))
		log.SetFlags(log.Lmsgprefix)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	s.lock.Lock()
	s.httpServer = s.makeHTTPServer(port, s.routes())

	// set app log as default error logger for http server
	if s.AppLog.IsLevelEnabled(logrus.ErrorLevel) {
		w := s.AppLog.WriterLevel(logrus.ErrorLevel)
		defer w.Close()
		s.httpServer.ErrorLog = log.New(w, "http", log.Lmsgprefix)
	}

	s.lock.Unlock()

	err := s.httpServer.ListenAndServe()
	s.AppLog.Warnf("CotonSMS REST API server terminated, %v", err)
}

func (s *Rest) Shutdown() {
	s.AppLog.Warn("Shutdown EmonaTextile API server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.AppLog.Errorf("shutdown error, %v", err)
		}
		s.AppLog.Debug("Shutdown EmonaTextile API server completed")
	}

	s.lock.Unlock()
}

func (s *Rest) makeHTTPServer(port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		// WriteTimeout:      120 * time.Second,
		IdleTimeout: 30 * time.Second,
	}
}

func (s *Rest) routes() chi.Router {
	rootRouter := chi.NewRouter()

	// A good base middleware stack
	rootRouter.Use(middleware.RequestID)
	rootRouter.Use(middleware.RealIP)
	//rootRouter.Use(RequestLogger(s.AppLog))
	//rootRouter.Use(HttpLogger(s.AppLog))
	// TODO make custom recoverer
	//rootRouter.Use(middleware.Recoverer)

	allowedOrigins := s.AllowedOrigins
	if allowedOrigins == nil || len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	rootRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "X-Auth-Token", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	rootRouter.Use(middleware.Timeout(60 * time.Second))

	rootRouter.Route("/", func(r chi.Router) {
		r.Get("/carousel", s.carouselManagement.list)
	})

	return rootRouter
}
