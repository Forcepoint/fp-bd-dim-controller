package api

import (
	"fmt"
	"fp-dynamic-elements-manager-controller/api/auth"
	"fp-dynamic-elements-manager-controller/api/backup"
	"fp-dynamic-elements-manager-controller/api/batch"
	"fp-dynamic-elements-manager-controller/api/docker"
	"fp-dynamic-elements-manager-controller/api/elements"
	"fp-dynamic-elements-manager-controller/api/export"
	"fp-dynamic-elements-manager-controller/api/health"
	"fp-dynamic-elements-manager-controller/api/logging"
	"fp-dynamic-elements-manager-controller/api/modules"
	"fp-dynamic-elements-manager-controller/api/notification"
	"fp-dynamic-elements-manager-controller/api/queue"
	"fp-dynamic-elements-manager-controller/api/registration"
	"fp-dynamic-elements-manager-controller/api/stats"
	"fp-dynamic-elements-manager-controller/api/update"
	"fp-dynamic-elements-manager-controller/api/user"
	"fp-dynamic-elements-manager-controller/api/util"
	authfuncs "fp-dynamic-elements-manager-controller/internal/auth"
	backup2 "fp-dynamic-elements-manager-controller/internal/backup"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	docker2 "fp-dynamic-elements-manager-controller/internal/docker"
	healthfuncs "fp-dynamic-elements-manager-controller/internal/health"
	logstructs "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/modules/structs"
	queuefuncs "fp-dynamic-elements-manager-controller/internal/queue"
	userfuncs "fp-dynamic-elements-manager-controller/internal/user"
	"github.com/gammazero/workerpool"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	authPathPrefix     = "/api"
	internalPathPrefix = "/internal"
	ingressPathPrefix  = "/ingress"
)

type server struct {
	logger         *logstructs.AppLogger
	dbReadyChan    <-chan struct{}
	addRoutesChan  chan structs.ModuleMetadata
	doneChan       chan struct{}
	router         *mux.Router
	authRouter     *mux.Router
	internalRouter *mux.Router
	ingressRouter  *mux.Router
	dao            *persistence.DataAccessObject
	pusher         queuefuncs.Pusher
	wp             *workerpool.WorkerPool
	handler        *docker2.CommandHandler
	provider       backup2.Provider
}

func NewServer(
	logger *logstructs.AppLogger,
	dbReadyChan <-chan struct{},
	dao *persistence.DataAccessObject,
	pusher queuefuncs.Pusher,
	handler *docker2.CommandHandler,
	provider backup2.Provider,
) *server {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(util.AddHeaders)
	authRouter := router.PathPrefix(authPathPrefix).Subrouter()
	authRouter.Use(authfuncs.JwtVerify)
	internalRouter := router.PathPrefix(internalPathPrefix).Subrouter()
	internalRouter.Use(authfuncs.InternalAuthVerify)
	ingressRouter := router.PathPrefix(ingressPathPrefix).Subrouter()
	return &server{
		logger:         logger,
		dbReadyChan:    dbReadyChan,
		addRoutesChan:  make(chan structs.ModuleMetadata),
		doneChan:       make(chan struct{}),
		router:         router,
		authRouter:     authRouter,
		internalRouter: internalRouter,
		ingressRouter:  ingressRouter,
		dao:            dao,
		pusher:         pusher,
		handler:        handler,
		provider:       provider,
	}
}

func (s *server) StartServer() {
	defer func() { s.doneChan <- struct{}{} }()
	upgrader := websocket.Upgrader{
		ReadBufferSize:  8 << 10,
		WriteBufferSize: 8 << 10,
		CheckOrigin: func(r *http.Request) bool {
			hdr := r.Header.Get("Origin")
			if hdr == "" {
				return false
			}
			return os.Getenv("HOST_DOMAIN") == strings.ReplaceAll(hdr, "https://", "")
		},
	}

	s.router.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "x-access-token"}),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodDelete, http.MethodPut}),
		handlers.AllowedOrigins([]string{os.Getenv("HOST_DOMAIN")})))

	s.router.Handle("/login", auth.Login(s.dao.UserRepo)).Methods(http.MethodPost)

	s.authRouter.Handle("/ws", notification.Handler(upgrader, s.logger.NotificationService))

	s.authRouter.Handle("/export", export.Handler(s.dao.ListElementRepo))
	s.authRouter.Handle("/backup", backup.Handler(s.provider, s.logger.NotificationService))
	s.authRouter.Handle("/keys", auth.GetRegistrationKey())
	s.authRouter.Handle("/health", health.Handler(s.dao))
	s.authRouter.Handle("/stats", stats.Handler(s.dao.ListElementRepo))
	s.authRouter.Handle("/logs", logging.Handler(s.dao.LogEntryRepo))
	s.authRouter.Handle("/modules", modules.Handler(s.dao))
	s.authRouter.Handle("/docker", docker.Handler(s.handler, s.logger.NotificationService))
	s.authRouter.Handle("/batch", batch.Handler(s.dao.UpdateStatusRepo))
	s.authRouter.Handle("/user", user.Handler(s.dao.UserRepo, s.logger))
	s.authRouter.Handle("/elements", elements.Handler(s.pusher, s.dao, s.logger))

	s.internalRouter.Handle("/register", registration.Handler(s.addRoutesChan))
	s.internalRouter.Handle("/queue", queue.Handler(s.pusher, s.dao, s.logger))
	s.internalRouter.Handle("/update", update.Handler(s.dao.UpdateStatusRepo))
	s.internalRouter.Handle("/logevent", logging.Handler(s.dao.LogEntryRepo))

	s.startDynamicRouteHandler()

	s.logger.SystemLogger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("CONTROLLER_PORT")), handlers.CompressHandler(s.router)), "error running server")
}

func (s *server) startDynamicRouteHandler() {
	go func() {
		for {
			select {
			case data := <-s.addRoutesChan:
				s.logger.UserLogger.Info(fmt.Sprintf("Adding new module: %s", data.ModuleDisplayName))
				s.dao.ModuleMetadataRepo.UpsertModuleMetadata(data)
				s.addModuleRoutes(data)
			case <-s.dbReadyChan:
				s.logger.UserLogger.Info("Adding module routes from persistence...")
				err := userfuncs.CreateAdminUserIfNotExists(s.dao.UserRepo)
				if err != nil {
					s.logger.SystemLogger.Error(err, "error creating admin user")
					return
				}
				metadata, err := s.dao.ModuleMetadataRepo.GetAllModuleMetadata()

				if err != nil {
					s.logger.SystemLogger.Error(err, "error retrieving module metadata to add routes")
					return
				}

				for _, dbModule := range metadata {
					s.addModuleRoutes(dbModule)
				}
				s.createAndSetRegistrationToken()
			case <-s.doneChan:
				return
			}
		}
	}()
}

func (s *server) addModuleRoutes(data structs.ModuleMetadata) {
	if !healthfuncs.IsUp(data.ModuleServiceName, data.InternalPort) {
		return
	}
	for _, v := range data.ModuleEndpoints {
		inboundRoute := data.InboundRoute + v.Endpoint
		parsedUrl, _ := url.Parse(fmt.Sprintf("http://%s:%s%s", data.ModuleServiceName, data.InternalPort, v.Endpoint))

		if v.Secure {
			s.authRouter.Handle(inboundRoute, util.NewReverseProxy(parsedUrl))
		} else {
			s.ingressRouter.Handle(inboundRoute, util.NewReverseProxy(parsedUrl))
		}
	}
}

func (s *server) createAndSetRegistrationToken() {
	if !viper.IsSet("internaltoken") {
		token := shortuuid.New()
		viper.Set("internaltoken", token)

		if err := viper.WriteConfig(); err != nil {
			s.logger.SystemLogger.Error(err, "error writing last run time to config file")
		}
	}
}
