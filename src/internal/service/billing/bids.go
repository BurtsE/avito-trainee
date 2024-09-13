package billing

import (
	"avito-test/internal/models"
	"avito-test/internal/service_errors"

	"github.com/pkg/errors"
)

func (s *service) CreateBid(builder models.BidsBuilder) (models.Bid, error) {
	switch builder.AuthorType {
	case models.UserAuthor:
		_, err := s.billingStorage.GetUserById(builder.AuthorId)
		if err != nil {
			return models.Bid{}, service_errors.UserNotFound{Err: err}
		}
	case models.OrgAuthor:
		_, err := s.billingStorage.GetOrganizationById(builder.AuthorId)
		if err != nil {
			return models.Bid{}, service_errors.UserNotFound{Err: err}
		}
	default:
		return models.Bid{}, errors.New("uniplemented")
	}
	_, err := s.billingStorage.GetTender(builder.TenderId)
	if err != nil {
		return models.Bid{}, service_errors.TenderError{Err: err}
	}
	bid, err := s.billingStorage.CreateBid(builder)
	if err != nil {
		return models.Bid{}, err
	}
	return bid, nil
}

func (s *service) GetUserBids(limit, offset int, username string) ([]models.Bid, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return nil, service_errors.UserNotFound{Err: err}
	}
	bids, err := s.billingStorage.GetUserBids(limit, offset, user.Id)
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (s *service) UpdateBid(bidId, username string, builder models.BidsBuilder) (models.Bid, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Bid{}, service_errors.UserNotFound{Err: err}
	}
	responsibleId, err := s.billingStorage.ValidBidRights(user.Id, bidId)
	if err != nil {
		return models.Bid{}, err
	}
	if responsibleId == "" {
		return models.Bid{}, service_errors.AuthError{}
	}
	bid, err := s.billingStorage.UpdateBid(bidId, builder)
	if err != nil {
		return models.Bid{}, errors.Wrap(service_errors.TenderError{}, err.Error())
	}
	return bid, nil
}

func (s *service) SetBidStatus(bidId, username string, status models.BidStatus) (models.Bid, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Bid{}, service_errors.UserNotFound{Err: err}
	}
	responsibleId, err := s.billingStorage.ValidBidRights(user.Id, bidId)
	if err != nil {
		return models.Bid{}, err
	}
	if responsibleId == "" {
		return models.Bid{}, service_errors.AuthError{}
	}
	bid, err := s.billingStorage.SetBidStatus(bidId, status)
	if err != nil {
		return models.Bid{}, errors.Wrap(service_errors.TenderError{}, err.Error())
	}
	return bid, nil
}

func (s *service) GetBid(bidId, username string) (models.Bid, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Bid{}, service_errors.UserNotFound{Err: err}
	}
	responsibleId, err := s.billingStorage.ValidBidRights(user.Id, bidId)
	if err != nil {
		return models.Bid{}, err
	}
	if responsibleId == "" {
		return models.Bid{}, service_errors.AuthError{}
	}
	bid, err := s.billingStorage.GetBid(bidId)
	if err != nil {
		return models.Bid{}, errors.Wrap(service_errors.TenderError{}, err.Error())
	}
	return bid, nil
}

func (s *service) RollbackBid(bidId, username string, version int) (models.Bid, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Bid{}, service_errors.UserNotFound{Err: err}
	}
	responsibleId, err := s.billingStorage.ValidBidRights(user.Id, bidId)
	if err != nil {
		return models.Bid{}, err
	}
	if responsibleId == "" {
		return models.Bid{}, service_errors.AuthError{}
	}
	bid, err := s.billingStorage.RollbackBid(bidId, version)
	if err != nil {
		return models.Bid{}, errors.Wrap(service_errors.TenderError{}, err.Error())
	}
	return bid, nil
}

func (s *service) SubmitDecision(bidId, username string, decision models.Decision) (models.Bid, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Bid{}, service_errors.UserNotFound{Err: err}
	}
	responsibleId, err := s.billingStorage.ValidBidRights(user.Id, bidId)
	if err != nil {
		return models.Bid{}, err
	}
	if responsibleId == "" {
		return models.Bid{}, service_errors.AuthError{}
	}
	err = s.billingStorage.SubmitDecision(bidId, decision)
	if err != nil {
		return models.Bid{}, err
	}
	switch decision {
	case models.BidApproved:
		return s.billingStorage.GetBid(bidId)
	case models.BidRejected:
		return s.billingStorage.SetBidStatus(bidId, models.BidCanceled)
	default:
		return models.Bid{}, errors.New("unimplemented")
	}
}
func (s *service) GetBidsForTender(tenderId, username string, limit, offset int) ([]models.Bid, error) {
	_, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return nil, service_errors.UserNotFound{Err: err}
	}
	bids, err := s.billingStorage.GetBidsForTender(tenderId, limit, offset)
	if err != nil {
		return nil, err
	}
	return bids, nil
}
func (s *service) AddFeedback(bidId, username, bidFeedback string) (models.Bid, error) {
	user, err := s.billingStorage.GetUserByUsername(username)
	if err != nil {
		return models.Bid{}, service_errors.UserNotFound{Err: err}
	}

	tender, err := s.billingStorage.GetTenderForBid(bidId)
	if err != nil {
		return models.Bid{}, service_errors.TenderError{Err: err}
	}

	responsibleId, err := s.billingStorage.ValidTenderRights(user.Id, tender.Id)
	if err != nil {
		return models.Bid{}, err
	}

	if responsibleId == "" {
		return models.Bid{}, service_errors.AuthError{}
	}
	bid, err := s.billingStorage.GetBid(bidId)
	if err != nil {
		return models.Bid{}, service_errors.TenderError{Err: err}
	}

	err = s.billingStorage.AddBidFeedback(bid.Id, bidFeedback)
	if err != nil {
		return models.Bid{}, err
	}
	return bid, nil

}
