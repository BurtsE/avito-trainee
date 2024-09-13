package billing

import (
	"avito-test/internal/models"
	"avito-test/internal/service_errors"

	"github.com/pkg/errors"
)

func (s *service) GetTender(tenderId, username string) (models.Tender, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.UserNotFound{}, err.Error())
	}
	authorId, err := s.billingStorage.ValidTenderRights(user.Id, tenderId)
	if err != nil {
		return models.Tender{}, err
	}
	if authorId == "" {
		return models.Tender{}, service_errors.AuthError{}
	}
	tender, err := s.billingStorage.GetTender(tenderId)
	if err != nil {
		return models.Tender{}, err
	}
	return tender, nil
}

func (s *service) CreateTender(builder models.TenderBuilder) (models.Tender, error) {
	result := models.Tender{}
	user, err := s.billingStorage.GetUserByUsername(builder.UserName)
	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.UserNotFound{}, err.Error())
	}
	responsibleOrgId, err := s.billingStorage.GetResponsibleOrganization(user.Id, builder.OrganizationId)
	if err != nil {
		return models.Tender{}, service_errors.AuthError{}
	}

	id, creationDate, err := s.billingStorage.CreateTender(builder.Name, builder.Description, builder.ServiceType.String(),
		builder.Status.String(), responsibleOrgId)
	if err != nil {
		return models.Tender{}, err
	}
	result.Id = id
	result.Name = builder.Name
	result.ServiceType = builder.ServiceType
	result.Status = builder.Status
	result.Version = 1
	result.Description = builder.Description
	result.CreatedAt = creationDate
	return result, nil
}

func (s *service) UpdateTender(tenderId, userName string, builder models.TenderBuilder) (models.Tender, error) {
	user, err := s.billingStorage.GetUserByUsername(userName)
	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.UserNotFound{}, err.Error())
	}
	responsibleId, err := s.billingStorage.ValidTenderRights(user.Id, tenderId)
	if err != nil {
		return models.Tender{}, err
	}
	if responsibleId == "" {
		return models.Tender{}, service_errors.AuthError{}
	}
	tender, err := s.billingStorage.UpdateTender(tenderId, builder)
	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.TenderError{}, err.Error())
	}
	return tender, nil
}

func (s *service) SetTenderStatus(tenderId, username string, status models.TenderStatus) (models.Tender, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.UserNotFound{}, err.Error())
	}
	responsibleId, err := s.billingStorage.ValidTenderRights(user.Id, tenderId)
	if err != nil {
		return models.Tender{}, err
	}
	if responsibleId == "" {
		return models.Tender{}, service_errors.AuthError{}
	}
	tender, err := s.billingStorage.SetTenderStatus(tenderId, status)
	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.TenderError{}, err.Error())
	}
	return tender, nil
}

func (s *service) RollbackTender(tenderId, username string, version int) (models.Tender, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.UserNotFound{}, err.Error())
	}

	responsibleId, err := s.billingStorage.ValidTenderRights(user.Id, tenderId)
	if err != nil {
		return models.Tender{}, err
	}
	if responsibleId == "" {
		return models.Tender{}, service_errors.AuthError{}
	}

	tender, err := s.billingStorage.RollbackTender(tenderId, version)

	if err != nil {
		return models.Tender{}, errors.Wrap(service_errors.TenderError{}, err.Error())
	}
	return tender, nil
}

func (s *service) GetTenders(limit, offset int, serviceTypes []models.EnumService) ([]models.Tender, error) {
	tenders, err := s.billingStorage.GetTenders(limit, offset, serviceTypes)
	if err != nil {
		return nil, err
	}
	return tenders, nil
}

func (s *service) GetUserTenders(limit, offset int, username string) ([]models.Tender, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return nil, service_errors.UserNotFound{Err: err}
	}
	tenders, err := s.billingStorage.GetUserTenders(limit, offset, user.Id)
	if err != nil {
		return nil, err
	}
	return tenders, nil
}
