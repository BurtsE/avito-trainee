package billing

import (
	"avito-test/internal/config"
	def "avito-test/internal/service"
	"avito-test/internal/storage"
)

var _ def.Service = (*service)(nil)

type service struct {
	billingStorage storage.Storage
}

func NewService(billingStorage storage.Storage, cfg *config.Config) *service {
	return &service{
		billingStorage: billingStorage,
	}
}



