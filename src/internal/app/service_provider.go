package app

import (
	"avito-test/internal/config"
	router "avito-test/internal/router"
	"avito-test/internal/service"

	billingService "avito-test/internal/service/billing"
	storages "avito-test/internal/storage"
	billingStorage "avito-test/internal/storage/billing"
	"log"

	"github.com/sirupsen/logrus"
)

// TODO auth service, repo init

type serviceProvider struct {
	cfg            *config.Config
	bilingStorage  storages.Storage
	billingService service.Service
	router         *router.Router
}

func NewSericeProvider() *serviceProvider {
	s := &serviceProvider{}
	s.Router()
	return s
}

func (s *serviceProvider) Config() *config.Config {
	if s.cfg == nil {
		cfg, err := config.InitConfig()
		if err != nil {
			log.Fatal(err)
		}
		s.cfg = cfg
	}
	return s.cfg
}

func (s *serviceProvider) Storage() storages.Storage {
	if s.bilingStorage == nil {
		storage, err := billingStorage.NewRepository(s.Config())
		if err != nil {
			log.Fatalf("could not init storage: %s", err.Error())
		}
		s.bilingStorage = storage
	}
	return s.bilingStorage
}

func (s *serviceProvider) Service() service.Service {
	if s.billingService == nil {
		s.billingService = billingService.NewService(s.Storage(), s.Config())
	}
	return s.billingService
}

func (s *serviceProvider) Router() *router.Router {
	if s.router == nil {
		s.router = router.NewRouter(logrus.New(), s.Config(), s.Service())
	}
	return s.router
}
