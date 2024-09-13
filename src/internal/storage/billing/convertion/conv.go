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

func TenderStatusFromString(input string) (models.TenderStatus, error) {

	switch input {
	case "Created":
		return models.TenderCreated, nil
	case "Published":
		return models.TenderPublished, nil
	case "Closed":
		return models.TenderClosed, nil
	default:
		return nil, errors.New("status convertion error")
	}
}

func AuthorTypeFromString(input string) (models.Author, error) {

	switch input {
	case "Organization":
		return models.OrgAuthor, nil
	case "User":
		return models.UserAuthor, nil
	default:
		return nil, errors.New("author type conversion error")
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
		return nil, errors.New("status convertion error")
	}
}
