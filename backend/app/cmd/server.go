package cmd

import (
	"backend/app/logger"
	"backend/app/model"
	api "backend/app/rest"
	"backend/app/template"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerCommand struct {
	AppSecret      string   `long:"app-secret" env:"APP_SECRET" required:"true" description:"App server secret key"`
	AllowedOrigins []string `long:"allowed-origins" env:"ALLOWED_ORIGINS" env-delim:"," description:"List of allowed origins. Default value is '*'"`

	MediaDir  string `long:"media-root" env:"MEDIA_ROOT" required:"true" description:"Media files root directory"`
	MediaHost string `long:"media-host" env:"MEDIA_HOST" default:"http://localhost:8000/media" description:"Media files host"`

	SecureCookiesDisabled bool `long:"secure-cookies-disabled" env:"SECURE_COOKIES_DISABLED" description:"Disable secure cookies"`

	AppLogLevel  string `long:"app-log" env:"APP_LOG" choice:"debug" choice:"info" choice:"warn" choice:"error" default:"info" description:"Application log level"`
	AppLogFormat string `long:"app-log-fmt" env:"APP_LOG_FMT" choice:"text" choice:"json" default:"text" description:"Application log format"`

	BackOfficeHost string `long:"back-office-host" env:"BACK_OFFICE_HOST" default:"http://localhost:8000" description:"Back-office host"`
	AdminPanelHost string `long:"admin-panel-host" env:"ADMIN_PANEL_HOST" default:"http://localhost:8000" description:"Admin Panel host"`

	MaxEmailConfirmationCount  int           `long:"max-email-confirmation-count" env:"MAX_EMAIL_CONFIRMATION_COUNT" default:"3" description:"Max email confirmation count"`
	MaxEmailConfirmationPeriod time.Duration `long:"max-email-confirmation-period" env:"MAX_EMAIL_CONFIRMATION_PERIOD" default:"12h" description:"Max email confirmation period"`
	ContactImportFileMaxSizeMB int           `long:"contact-import-max-size-mb" env:"CONTACT_IMPORT_MAX_SIZE_MB" default:"10" description:"Max size of contact import file, in megabytes"`

	AccessTokenDuration            string `long:"access-token-duration" env:"ACCESS_TOKEN_DURATION" default:"5m" description:"Rest API access token duration, valid units 's', 'm', 'h'. Example: 10m - 10 minutes"`
	RefreshTokenDuration           string `long:"refresh-token-duration" env:"REFRESH_TOKEN_DURATION" default:"24h" description:"Rest API refresh token duration, valid units 's', 'm', 'h'. Example: 24h - 24 hours"`
	EmailConfirmationTokenDuration string `long:"email-confirm-token-duration" env:"EMAIL_CONFIRMATION_TOKEN_DURATION" default:"12h" description:"Email confirmation token duration, valid units 's', 'm', 'h'. Example: 12h - 12 hours"`
	InvitationTokenDuration        string `long:"invitation-token-duration" env:"INVITATION_TOKEN_DURATION" default:"24h" description:"User invitation token duration, valid units 's', 'm', 'h'. Example: 24h - 24 hours"`
	PasswordResetTokenDuration     string `long:"pass-reset-token-duration" env:"PASS_RESET_TOKEN_DURATION" default:"10m" description:"Password reset token duration, valid units 's', 'm', 'h'. Example: 10m - 10 minutes"`

	Port int `short:"p" long:"app-port" env:"APP_PORT" default:"8000" description:"App server port"`

	Mailgun MailgunCommand `group:"Mailgun config" namespace:"mailgun"`

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

	appLogger.Infof("Starting CotonSMS server")

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
		appLogger.Panicf("Failed to setup CotonSMS server, %+v", err)
		return err
	}
	if err = app.run(ctx); err != nil {
		appLogger.Errorf("CotonSMS server terminated with error %+v", err)
		return err
	}
	appLogger.Info("CotonSMS server terminated")

	return nil
}

func (c *ServerCommand) newServerApp(appLogger *logrus.Logger) (*serverApp, error) {
	ctx := logger.WithLogger(context.Background(), appLogger.WithFields(logrus.Fields{
		"init": "ensure_indexes",
	}))

	storageOpts := model.NewStorageOptions()
	storageOpts.MongoURI = c.MongoURI
	storageOpts.MongoDatabase = c.MongoDatabase

	storage, err := storageOpts.Storage()
	if err != nil {
		return nil, errors.Wrap(err, "fail to create storage")
	}

	refBookStore, err := model.NewRefBookStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create refbook storage")
	}

	err = model.InitPhoneService(ctx, refBookStore)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create init phone service")
	}

	staffStore, err := model.NewStaffStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create staff storage")
	}

	err = staffStore.EnsureIndexes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to ensure staff storage indexes")
	}

	userStore, err := model.NewUserStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create user storage")
	}

	err = userStore.EnsureIndexes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to ensure user storage indexes")
	}

	accountsStore, err := model.NewAccountStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create accounts storage")
	}

	err = accountsStore.EnsureIndexes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to ensure accounts storage indexes")
	}

	contactsStore, err := model.NewContactsStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create contacts storage")
	}

	err = contactsStore.EnsureIndexes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to ensure contacts storage indexes")
	}

	pricesStore, err := model.NewPricesStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create prices storage")
	}

	err = pricesStore.EnsureIndexes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to ensure prices storage indexes")
	}

	campaignStore, err := model.NewCampaignStore(storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create campaign storage")
	}

	err = campaignStore.EnsureIndexes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to ensure campaign storage indexes")
	}

	linkBuilder := api.NewLinkBuilder(c.BackOfficeHost, c.AdminPanelHost)

	//mailgunSenderOpts := mailgun.NewSenderOptions()
	//mailgunSenderOpts.DomainName = c.Mailgun.DomainName
	//mailgunSenderOpts.ApiKey = c.Mailgun.ApiKey
	//mailgunSenderOpts.FromEmail = c.Mailgun.FromEmail
	//mailgunSenderOpts.APIBase = c.Mailgun.APIBase()

	//sender, err := mailgunSenderOpts.Sender()
	//if err != nil {
	//	return nil, errors.Wrap(err, "fail to create email sender")
	//}

	tokenAuthorityConf := api.NewTokenAuthorityConf(c.AppSecret)
	tokenAuthorityConf.AccessTokenDuration, err = time.ParseDuration(c.AccessTokenDuration)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse access token duration")
	}
	tokenAuthorityConf.RefreshTokenDuration, err = time.ParseDuration(c.RefreshTokenDuration)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse refresh token duration")
	}
	tokenAuthorityConf.EmailConfirmationTokenDuration, err = time.ParseDuration(c.EmailConfirmationTokenDuration)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse email confirmation token duration")
	}
	tokenAuthorityConf.InvitationTokenDuration, err = time.ParseDuration(c.InvitationTokenDuration)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse email invitation token duration")
	}
	tokenAuthorityConf.PasswordResetTokenDuration, err = time.ParseDuration(c.PasswordResetTokenDuration)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse email password reset token duration")
	}
	tokenAuthority, err := api.NewTokenAuthority(tokenAuthorityConf, storage)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create token authority")
	}

	emailTemplates, err := template.NewEmailTemplates(c.BackOfficeHost, c.AdminPanelHost)
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize email templates")
	}

	srv := &api.Rest{
		AllowedOrigins: c.AllowedOrigins,

		MediaDir:  c.MediaDir,
		MediaHost: c.MediaHost,

		AdminPanelHost: c.AdminPanelHost,
		BackOfficeHost: c.BackOfficeHost,

		SecureCookiesDisabled: c.SecureCookiesDisabled,

		ContactImportFileMaxSize: int64(c.ContactImportFileMaxSizeMB * 1024 * 1024),

		MaxEmailConfirmationCount:  c.MaxEmailConfirmationCount,
		MaxEmailConfirmationPeriod: c.MaxEmailConfirmationPeriod,

		StaffStore:    staffStore,
		UserStore:     userStore,
		AccountsStore: accountsStore,
		ContactsStore: contactsStore,
		PricesStore:   pricesStore,
		CampaignStore: campaignStore,

		PhoneService: model.SharedPhoneService,

		RefBookStore: refBookStore,
		// MailSender:   sender,
		LinkBuilder: linkBuilder,

		TokenAuthority: tokenAuthority,
		EmailTemplates: emailTemplates,

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
