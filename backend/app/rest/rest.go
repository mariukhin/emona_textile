package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	contacts_import "amifactory.team/sequel/coton-app-backend/app/import"
	"amifactory.team/sequel/coton-app-backend/app/mail"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"amifactory.team/sequel/coton-app-backend/app/template"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Rest struct {
	AllowedOrigins []string

	MediaDir  string
	MediaHost string

	BackOfficeHost string
	AdminPanelHost string

	SecureCookiesDisabled bool

	ContactImportFileMaxSize int64

	MaxEmailConfirmationCount  int
	MaxEmailConfirmationPeriod time.Duration

	StaffStore    model.StaffStore
	UserStore     model.UserStore
	AccountsStore model.AccountStore
	ContactsStore model.ContactsStore
	PricesStore   model.PricesStore
	RefBookStore  model.RefBookStore
	CampaignStore model.CampaignStore

	PhoneService model.PhoneService

	MailSender  mail.Sender
	LinkBuilder LinkBuilder

	TokenAuthority TokenAuthority
	EmailTemplates template.EmailTemplates

	AppLog *logrus.Logger

	adminCommon *adminCommon

	// Admin Panel (Staff)
	staffAuth       *staffAuth
	staffRole       *staffRole
	staffProfile    *staffProfile
	staffManagement *staffManagement
	accounts        *staffAccounts
	usersManagement *usersManagement
	staffPrices     *staffPrices

	// Back-office (Customers)
	userAuth             *userAuth
	userCommon           *userCommon
	userProfile          *userProfile
	userAccounts         *userAccounts
	userAccountRoles     *userAccountRoles
	userContactGroups    *userAccountContactGroups
	userContactGroupVars *userAccountContactGroupVars
	userContacts         *userAccountContacts
	userCampaigns        *userCampaigns
	userPrices           *userPrices

	httpServer *http.Server
	lock       sync.Mutex
}

func (s *Rest) Run(port int) {
	s.AppLog.Infof("Start CotonSMS REST API server on port %d", port)

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
	s.AppLog.Warn("Shutdown CotonSMS REST API server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.AppLog.Errorf("shutdown error, %v", err)
		}
		s.AppLog.Debug("Shutdown CotonSMS REST API server completed")
	}

	s.lock.Unlock()
}

func (s *Rest) makeHTTPServer(port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		// WriteTimeout:      120 * time.Second, // TODO: such a long timeout needed for blocking export (backup) request
		IdleTimeout: 30 * time.Second,
	}
}

func (s *Rest) routes() chi.Router {
	rootRouter := chi.NewRouter()

	// A good base middleware stack
	rootRouter.Use(middleware.RequestID)
	rootRouter.Use(middleware.RealIP)
	rootRouter.Use(RequestLogger(s.AppLog))
	rootRouter.Use(HttpLogger(s.AppLog))
	// TODO make custom recoverer
	rootRouter.Use(middleware.Recoverer)

	s.adminCommon = s.adminControllers()

	s.staffAuth, s.staffRole, s.staffProfile, s.staffManagement,
		s.accounts, s.usersManagement, s.staffPrices = s.staffControllers()

	s.userAuth, s.userCommon, s.userProfile, s.userAccounts,
		s.userAccountRoles, s.userContactGroups, s.userContactGroupVars,
		s.userContacts, s.userCampaigns, s.userPrices = s.userControllers()

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

	// Staff routes
	rootRouter.Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Route("/api/v1", func(r chi.Router) {
			// public routes
			r.With(AllowContentTypeJson()).Group(func(r chi.Router) {
				r.Post("/sign-in/email", s.staffAuth.signInEmail)
				r.Post("/sign-out", s.staffAuth.signOut)
				r.Post("/confirm-email", s.staffAuth.confirmEmail)
				r.Post("/refresh", s.staffAuth.refreshToken)
			})

			// routes that require authentication
			r.Group(func(r chi.Router) {
				// Check staff access token
				r.Use(s.staffAuth.authenticator)

				// Staff API: https://docs.coton.amifactory.network/#/Staff%20API
				r.Route("/profile", func(r chi.Router) {
					r.Get("/", s.staffProfile.getStaffProfile)
					r.With(AllowContentTypeJson()).Put("/", s.staffProfile.updateStaffProfile)
					r.Post("/photo", s.staffProfile.updateStaffProfilePhoto)
					r.Delete("/photo", s.staffProfile.deleteStaffProfilePhoto)
				})

				// Staff Common: https://docs.coton.amifactory.network/#/Staff%20Common
				r.Get("/staff-roles/", s.staffRole.getStaffRoles)
				r.Get("/account-member-roles/", s.adminCommon.getAccountMemberRoles)
				r.Get("/pagination-options/", s.adminCommon.getPaginationOptions)
				r.Get("/account-statuses/", s.adminCommon.getAccountStatus)
				r.Get("/moderation-options/", s.adminCommon.getModerationOptions)
				r.Get("/currencies/", s.adminCommon.getCurrencies)

				// Staff Management: https://docs.coton.amifactory.network/#/Staff%20Management%20API
				r.Route("/staff", func(r chi.Router) {
					r.Get("/", s.staffManagement.getStaffList)
					r.With(AllowContentTypeJson()).Post("/", s.staffManagement.createStaff)

					r.Group(func(r chi.Router) {
						// Fetch and put staff by ID to context
						r.Use(s.staffManagement.staffById)

						r.Route("/{staffID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
							r.Get("/", s.staffManagement.fetchStaff)
							r.With(AllowContentTypeJson()).Patch("/", s.staffManagement.updateStaff)
							r.Post("/photo", s.staffManagement.updateStaffPhoto)
							r.Delete("/photo", s.staffManagement.deleteStaffPhoto)
							r.Post("/activate", s.staffManagement.activateStaff)
							r.Post("/deactivate", s.staffManagement.deactivateStaff)
							r.Post("/send-email-confirmation", s.staffManagement.sendEmailConfirmationStaff)
							r.Post("/reset-password", s.staffManagement.resetPasswordStaff)
						})
					})
				})

				// Staff Account Management API: https://docs.coton.amifactory.network/#/Staff%20Account%20Management%20API
				r.Route("/accounts", func(r chi.Router) {
					r.Get("/", s.accounts.list)
					r.Get("/all", s.accounts.listAll)

					r.Group(func(r chi.Router) {
						// Fetch and put staff by ID to context
						r.Use(s.accounts.populateAccount)

						r.Route("/{accountID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
							r.Get("/", s.accounts.fetchAccount)
							r.Put("/", s.accounts.updateAccount)
							r.Post("/activate", s.accounts.activateAccount)
							r.Post("/block", s.accounts.blockAccount)
							r.Get("/address", s.accounts.accountAddress)
							r.Get("/members/", s.accounts.accountMembers)
							r.Group(func(r chi.Router) {
								r.Use(s.accounts.populateAccountMember)
								r.Route("/members/{memberID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
									r.Get("/", s.accounts.accountMember)
									r.Put("/", s.accounts.updateAccountMembers)
									r.Delete("/", s.accounts.deleteAccountMembers)
								})
							})
						})
					})
				})

				// Staff Users Management API: https://docs.coton.amifactory.network/#/Staff%20User%20Management%20API
				r.Route("/users", func(r chi.Router) {
					r.Get("/", s.usersManagement.list)
					r.Get("/{userID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", s.usersManagement.fetch)
				})

				// Prices management
				r.Route("/prices", func(r chi.Router) {
					r.Get("/", s.staffPrices.list)
					r.Post("/", s.staffPrices.add)

					r.Group(func(r chi.Router) {
						r.Use(s.staffPrices.populatePriceItem)
						r.Route("/{priceID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
							r.With(AllowContentTypeJson()).Patch("/", s.staffPrices.update)
							r.Delete("/", s.staffPrices.delete)
						})
					})
				})
			})
		})
	})

	// Back-office routes
	rootRouter.Route("/customer", func(adminRouter chi.Router) {
		adminRouter.Route("/api/v1", func(r chi.Router) {
			// public routes
			r.With(AllowContentTypeJson()).Group(func(r chi.Router) {
				r.Post("/sign-in/email", s.userAuth.signInEmail)
				r.Post("/sign-up/email", s.userAuth.signUpEmail)
				r.Post("/sign-out", s.userAuth.signOut)
				r.Post("/confirm-email", s.userProfile.confirmEmail)
				r.Post("/account/confirm-email", s.userAccounts.confirmEmail)
				r.Post("/request-password-reset", s.userAuth.requestPasswordReset)
				r.Post("/password-reset", s.userAuth.passwordResetWithToken)
				r.Post("/refresh", s.userAuth.refreshToken)
				r.Get("/invitation", s.userAccounts.invitationDetails)
			})

			// routes that require authentication
			r.Group(func(r chi.Router) {
				// Check user access token
				r.Use(s.userAuth.authenticator)

				// Fetch user current account (if account-id cookie provided)
				r.Use(s.userAccounts.populateAccount)

				r.Group(func(r chi.Router) {
					// user profile manipulation
					r.Route("/profile", func(r chi.Router) {
						r.Get("/", s.userProfile.get)
						r.With(AllowContentTypeJson()).Put("/", s.userProfile.update)
						r.Post("/photo", s.userProfile.updatePhoto)
						r.Delete("/photo", s.userProfile.deletePhoto)
						r.Post("/send-email-confirmation", s.userProfile.sendEmailConfirmationUser)
					})

					r.Route("/accounts", func(r chi.Router) {
						r.Get("/", s.userAccounts.list)
						r.Post("/{accountID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}/make-current", s.userAccounts.makeCurrent)
					})
				})

				r.Group(func(r chi.Router) {
					// check if account populated
					r.Use(s.userAccounts.ensureAccountPopulated)

					r.Get("/roles/", s.userAccountRoles.availableRoles)
					r.Get("/pagination-options/", s.userCommon.getPaginationOptions)
					r.Get("/campaign-filters/", s.userCommon.campaignFilters)
				})

				r.Route("/invitations/{invitationID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
					r.Post("/accept", s.userAccounts.invitationAccept)
					r.Post("/reject", s.userAccounts.invitationReject)
				})

				// Current Account management
				r.Group(func(r chi.Router) {
					// check if user email confirmed
					r.Use(s.userAuth.ensureConfirmed)

					// check if account populated
					r.Use(s.userAccounts.ensureAccountPopulated)

					r.Route("/account", func(r chi.Router) {
						r.Get("/", s.userAccounts.get)

						// user accounts
						r.Group(func(r chi.Router) {
							r.Use(s.userAccounts.restrictAccess(model.AccountMemberRoleOwner, model.AccountMemberRoleAdmin))

							r.Put("/", s.userAccounts.update)
							r.Post("/send-email-confirmation", s.userAccounts.resendEmailConfirmation)

							r.Route("/members", func(r chi.Router) {
								r.Get("/", s.userAccounts.membersList)
								r.Post("/invite", s.userAccounts.inviteMember)

								r.Group(func(r chi.Router) {
									// Fetch and put Member by ID to context
									r.Use(s.userAccounts.populateAccountMember)
									r.Route("/{memberID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
										r.Get("/", s.userAccounts.getMember)
										r.Put("/", s.userAccounts.updateMember)
										r.Delete("/", s.userAccounts.deleteMember)
										r.Post("/resend-invitation", s.userAccounts.resendMemberInvitation)
									})
								})
							})

							r.Route("/address", func(r chi.Router) {
								r.Get("/", s.userAccounts.serviceAddress)
								r.Put("/", s.userAccounts.serviceAddressUpdate)
							})
						})
					})
				})

				r.Group(func(r chi.Router) {
					// check if user email confirmed
					r.Use(s.userAuth.ensureConfirmed)

					// check if account populated
					r.Use(s.userAccounts.ensureAccountPopulated)

					// Account Contact Groups management
					r.Route("/groups", func(r chi.Router) {
						r.Get("/", s.userContactGroups.list)
						r.Get("/all", s.userContactGroups.listAll)
						r.Post("/", s.userContactGroups.add)

						r.Group(func(r chi.Router) {
							r.Use(s.userContactGroups.populateGroup)
							r.Route("/{groupID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
								r.Get("/", s.userContactGroups.get)
								r.Put("/", s.userContactGroups.update)
								r.Delete("/", s.userContactGroups.delete)

								// Account Group Variables management
								r.Route("/vars", func(r chi.Router) {
									r.Get("/", s.userContactGroupVars.list)
									r.Post("/", s.userContactGroupVars.add)

									r.Group(func(r chi.Router) {
										r.Use(s.userContactGroupVars.populateVariable)
										r.Route("/{varID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
											r.Get("/", s.userContactGroupVars.get)
											r.Put("/", s.userContactGroupVars.update)
											r.Delete("/", s.userContactGroupVars.delete)
										})
									})
								})

								// Account Group Contacts management
								r.Route("/contacts", func(r chi.Router) {
									//r.Get("/group-filter-options/", s.userContacts.groupFilterOptions)
									r.Get("/", s.userContacts.list)
									r.Post("/", s.userContacts.add)
									r.Post("/import-preview", s.userContacts.importPreviewCSV)
									r.Post("/import", s.userContacts.importCSV)
									r.Get("/import-progress", s.userContacts.importProgress)

									r.Group(func(r chi.Router) {
										r.Use(s.userContacts.populateContact)
										r.Route("/{contactID:^[a-f\\d]{8}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{4}-[a-f\\d]{12}$}", func(r chi.Router) {
											r.Put("/", s.userContacts.update)
											r.Delete("/", s.userContacts.delete)
										})
									})
								})
							})
						})
					})

					// Campaigns management
					r.Route("/campaigns", func(r chi.Router) {
						r.Post("/", s.userCampaigns.add)
						r.Get("/constraints", s.userCampaigns.constraints)
					})

					// Prices
					r.Route("/prices", func(r chi.Router) {
						r.Get("/", s.userPrices.list)
					})
				})
			})
		})
	})

	// TODO disable directory indexing
	rootRouter.Handle("/media/*", http.StripPrefix("/media/", http.FileServer(http.Dir(s.MediaDir))))
	rootRouter.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	return rootRouter
}

func (s *Rest) adminControllers() *adminCommon {
	aCommon := adminCommon{}

	return &aCommon
}

func (s *Rest) staffControllers() (*staffAuth, *staffRole, *staffProfile, *staffManagement, *staffAccounts, *usersManagement, *staffPrices) {
	sAuth := staffAuth{
		tokenAuthority:        s.TokenAuthority,
		store:                 s.StaffStore,
		adminPanelHost:        s.AdminPanelHost,
		secureCookiesDisabled: s.SecureCookiesDisabled,
	}

	sRole := staffRole{
		store: s.StaffStore,
	}

	sProfile := staffProfile{
		store:     s.StaffStore,
		mediaDir:  s.MediaDir,
		mediaHost: s.MediaHost,
	}

	sManagement := staffManagement{
		tokenAuthority:             s.TokenAuthority,
		emailTemplates:             s.EmailTemplates,
		linkBuilder:                s.LinkBuilder,
		store:                      s.StaffStore,
		mailSender:                 s.MailSender,
		mediaDir:                   s.MediaDir,
		mediaHost:                  s.MediaHost,
		staffProfilePhotoDir:       path.Join("staff", "photo"),
		maxEmailConfirmationCount:  s.MaxEmailConfirmationCount,
		maxEmailConfirmationPeriod: s.MaxEmailConfirmationPeriod,
	}

	acc := staffAccounts{
		accountStore:          s.AccountsStore,
		secureCookiesDisabled: s.SecureCookiesDisabled,
	}

	prices := staffPrices{
		pricesStore:  s.PricesStore,
		accountStore: s.AccountsStore,
		phoneService: s.PhoneService,
	}

	usersManagement := usersManagement{
		store:     s.UserStore,
		mediaHost: s.MediaHost,
	}

	return &sAuth, &sRole, &sProfile, &sManagement, &acc, &usersManagement, &prices
}

func (s *Rest) userControllers() (*userAuth, *userCommon, *userProfile,
	*userAccounts, *userAccountRoles, *userAccountContactGroups,
	*userAccountContactGroupVars, *userAccountContacts, *userCampaigns, *userPrices) {
	auth := userAuth{
		tokenAuthority:        s.TokenAuthority,
		emailTemplates:        s.EmailTemplates,
		userStore:             s.UserStore,
		accountStore:          s.AccountsStore,
		linkBuilder:           s.LinkBuilder,
		mailSender:            s.MailSender,
		backOfficeHost:        s.BackOfficeHost,
		secureCookiesDisabled: s.SecureCookiesDisabled,
	}

	common := userCommon{}

	profile := userProfile{
		tokenAuthority: s.TokenAuthority,
		emailTemplates: s.EmailTemplates,
		linkBuilder:    s.LinkBuilder,
		mailSender:     s.MailSender,
		userStore:      s.UserStore,
		accountStore:   s.AccountsStore,
		mediaDir:       s.MediaDir,
		mediaHost:      s.MediaHost,
	}

	accounts := userAccounts{
		tokenAuthority:        s.TokenAuthority,
		emailTemplates:        s.EmailTemplates,
		linkBuilder:           s.LinkBuilder,
		mailSender:            s.MailSender,
		accountStore:          s.AccountsStore,
		userStore:             s.UserStore,
		backOfficeHost:        s.BackOfficeHost,
		secureCookiesDisabled: s.SecureCookiesDisabled,
	}

	accountRoles := userAccountRoles{}

	accountContactGroups := userAccountContactGroups{
		contactsStore:         s.ContactsStore,
		secureCookiesDisabled: s.SecureCookiesDisabled,
	}

	accountContactGroupVars := userAccountContactGroupVars{
		contactsStore:         s.ContactsStore,
		secureCookiesDisabled: s.SecureCookiesDisabled,
	}

	accountContacts := userAccountContacts{
		contactsStore:         s.ContactsStore,
		secureCookiesDisabled: s.SecureCookiesDisabled,
		maxImportBodySize:     s.ContactImportFileMaxSize,
		importFilesDir:        s.MediaDir, // TODO: temp solution
		tasks:                 make(map[string]*contacts_import.ContactImportTask),
	}

	campaigns := userCampaigns{
		contactsStore:         s.ContactsStore,
		pricesStore:           s.PricesStore,
		phoneService:          s.PhoneService,
		accountStore:          s.AccountsStore,
		campaignStore:         s.CampaignStore,
		secureCookiesDisabled: s.SecureCookiesDisabled,
	}
	prices := userPrices{
		pricesStore:  s.PricesStore,
		accountStore: s.AccountsStore,
		phoneService: s.PhoneService,
	}

	return &auth, &common, &profile,
		&accounts, &accountRoles, &accountContactGroups,
		&accountContactGroupVars, &accountContacts, &campaigns, &prices
}
