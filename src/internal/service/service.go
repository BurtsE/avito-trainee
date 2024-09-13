package service

import "avito-test/internal/models"

type Service interface {
	GetTenders(limit, offset int, serviceTypes []models.EnumService) ([]models.Tender, error)
	GetTender(tenderId, username string) (models.Tender, error)
	GetUserTenders(limit, offset int, username string) ([]models.Tender, error)
	CreateTender(models.TenderBuilder) (models.Tender, error)
	UpdateTender(string, string, models.TenderBuilder) (models.Tender, error)
	SetTenderStatus(tenderId, username string, status models.TenderStatus) (models.Tender, error)
	RollbackTender(tenderId, username string, version int) (models.Tender, error)

	CreateBid(models.BidsBuilder) (models.Bid, error)
	GetUserBids(limit, offset int, username string) ([]models.Bid, error)
	UpdateBid(string, string, models.BidsBuilder) (models.Bid, error)
	SetBidStatus(bidId, username string, status models.BidStatus) (models.Bid, error)
	GetBid(bidId, username string) (models.Bid, error)
	RollbackBid(bidId, username string, version int) (models.Bid, error)
	GetBidsForTender(tenderId, username string, limit, offset int) ([]models.Bid, error)

	AddFeedback(bidId, username, bidFeedback string)(models.Bid, error)

	SubmitDecision(bidId, username string, decision models.Decision) (models.Bid, error)
}
