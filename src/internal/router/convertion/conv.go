package convertion

import (
	"avito-test/internal/models"
	"errors"
)

func ServiceFromString(input string) (models.EnumService, error) {
	switch input {
	case "Construction":
		return models.Construction, nil
	case "Delivery":
		return models.Delivery, nil
	case "Manufacture":
		return models.Manufacture, nil
	default:
		return nil, errors.New("service conversion error")
	}

}

func StatusFromString(input string) (models.TenderStatus, error) {
	switch input {
	case "Created":
		return models.TenderCreated, nil
	case "Published":
		return models.TenderPublished, nil
	case "Closed":
		return models.TenderClosed, nil
	default:
		return nil, errors.New("status conversion error")
	}
}

func BidStatusFromString(input string) (models.BidStatus, error) {
	switch input {
	case "Created":
		return models.BidCreated, nil
	case "Published":
		return models.BidPublished, nil
	case "Canceled":
		return models.BidCreated, nil
	default:
		return nil, errors.New("status conversion error")
	}
}

func BidDecisionFromString(input string) (models.Decision, error) {
	switch input {
	case "Approved":
		return models.BidApproved, nil
	case "Rejected":
		return models.BidRejected, nil
	default:
		return nil, errors.New("decision conversion error")
	}
}
