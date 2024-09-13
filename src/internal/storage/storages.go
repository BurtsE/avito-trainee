package storage

import "avito-test/internal/models"

type Storage interface {
	CreateTender(name, description, serviceType, moderationStatus, responsibleId string) (id, time string, err error)
	UpdateTender(tenderId string, builder models.TenderBuilder) (models.Tender, error)
	SetTenderStatus(tenderId string, status models.TenderStatus) (models.Tender, error)
	RollbackTender(tenderId string, version int) (models.Tender, error)
	GetTender(tenderId string) (models.Tender, error)
	GetTenderForBid(bidId string) (models.Tender, error)
	GetTenders(limit, offset int, serviceTypes []models.EnumService) ([]models.Tender, error)
	GetUserTenders(limit, offset int, userId string) ([]models.Tender, error)

	GetUserByUsername(username string) (models.User, error)
	GetUserById(id string) (models.User, error)

	GetOrganizationById(id string) (models.DummyOrganization, error)

	GetResponsibleOrganization(userId, organizationId string) (string, error)
	GetOrganizationId(userId string) (string, error)

	GetBid(bidId string) (models.Bid, error)
	CreateBid(models.BidsBuilder) (models.Bid, error)
	GetUserBids(limit, offset int, userId string) ([]models.Bid, error)
	UpdateBid(bidId string, builder models.BidsBuilder) (models.Bid, error)
	SetBidStatus(bidId string, status models.BidStatus) (models.Bid, error)
	RollbackBid(bidId string, version int) (models.Bid, error)
	GetBidsForTender(tenderId string, limit, offset int) ([]models.Bid, error)

	SubmitDecision(bidId string, decision models.Decision) error

	AddBidFeedback(bidId, feedback string) error

	ValidTenderRights(userId, tenderId string) (string, error)
	ValidBidRights(userId, bidId string) (string, error)
}
