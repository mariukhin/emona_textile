package cmd

import (
	"backend/app/logger"
	"backend/app/model"
	api "backend/app/rest"
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type ServerCommand struct {
	AppSecret      string   `long:"app-secret" env:"APP_SECRET" description:"App server secret key"`
	AllowedOrigins []string `long:"allowed-origins" env:"ALLOWED_ORIGINS" env-delim:"," description:"List of allowed origins. Default value is '*'"`

	SecureCookiesDisabled bool `long:"secure-cookies-disabled" env:"SECURE_COOKIES_DISABLED" description:"Disable secure cookies"`

	AppLogLevel  string `long:"app-log" env:"APP_LOG" choice:"debug" choice:"info" choice:"warn" choice:"error" default:"info" description:"Application log level"`
	AppLogFormat string `long:"app-log-fmt" env:"APP_LOG_FMT" choice:"text" choice:"json" default:"text" description:"Application log format"`

	Port int `short:"p" long:"app-port" env:"APP_PORT" default:"4000" description:"App server port"`

	MongoDBCommand
}

type serverApp struct {
	*ServerCommand

	restSrv    *api.Rest
	appLogger  *logrus.Logger
	terminated chan struct{}
}

func (c *ServerCommand) Execute(args []string) error {
	logLevel, err := logrus.ParseLevel(c.AppLogLevel)
	if err != nil {
		return errors.Wrap(err, "fail to parse app log level")
	}

	logFormat, err := logger.ParseLogFormat(c.AppLogFormat)
	if err != nil {
		return errors.Wrap(err, "fail to parse app log format")
	}

	appLogger := logrus.New()
	appLogger.Out = os.Stdout
	appLogger.Level = logLevel
	appLogger.SetFormatter(logFormat.Formatter())

	appLogger.Infof("Starting EmonaTextile server")

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		appLogger.Warn("Interrupt signal")
		cancel()
	}()

	app, err := c.newServerApp(appLogger)
	if err != nil {
		appLogger.Panicf("Failed to setup EmonaTextile server, %+v", err)
		return err
	}
	if err = app.run(ctx); err != nil {
		appLogger.Errorf("EmonaTextile server terminated with error %+v", err)
		return err
	}
	appLogger.Info("EmonaTextile server terminated")

	return nil
}

func (c *ServerCommand) newServerApp(appLogger *logrus.Logger) (*serverApp, error) {
	//ctx := logger.WithLogger(context.Background(), appLogger.WithFields(logrus.Fields{
	//	"init": "ensure_indexes",
	//}))

	storageOpts := model.NewStorageOptions()
	storageOpts.MongoURI = c.MongoURI
	storageOpts.MongoDatabase = c.MongoDatabase

	storage, err := storageOpts.Storage()
	if err != nil {
		return nil, errors.Wrap(err, "fail to create storage")
	}

	carouselStore, err := model.NewCarouselStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create carousel storage")
	}

	srv := &api.Rest{
		AllowedOrigins: c.AllowedOrigins,

		SecureCookiesDisabled: c.SecureCookiesDisabled,
		CarouselStore:         carouselStore,

		AppLog: appLogger,
	}

	return &serverApp{
		ServerCommand: c,
		restSrv:       srv,
		appLogger:     appLogger,
		terminated:    make(chan struct{}),
	}, nil
}

func (s *serverApp) run(ctx context.Context) error {
	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		s.appLogger.Info("Shutdown initiated")
		s.restSrv.Shutdown()
	}()

	s.restSrv.Run(s.Port)

	close(s.terminated)
	return nil
}
