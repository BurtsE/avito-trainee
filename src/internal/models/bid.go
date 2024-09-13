package models

import "encoding/json"

const (
	CreatorOrganization = "Organization"
	CreatorUser         = "User"
)

type BidsBuilder struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	TenderId    string `json:"tenderId,omitempty"`
	AuthorType  Author `json:"authorType,omitempty"`
	AuthorId    string `json:"authorId,omitempty"`
}

func (b *BidsBuilder) UnmarshalJSON(data []byte) error {
	input := struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		TenderId    string `json:"tenderId,omitempty"`
		AuthorType  string `json:"authorType,omitempty"`
		AuthorId    string `json:"authorId,omitempty"`
	}{}
	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}
	var (
		authorType Author
	)
	switch input.AuthorType {
	case "Organization":
		authorType = OrgAuthor
	case "User":
		authorType = UserAuthor
	default:
		authorType = nil
	}

	b.Name = input.Name
	b.Description = input.Description
	b.TenderId = input.TenderId
	b.AuthorType = authorType
	b.AuthorId = input.AuthorId
	return nil
}

type Bid struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Status      BidStatus `json:"status"`
	AuthorType  Author    `json:"authorType"`
	Description string    `json:"description"`
	AuthorId    string    `json:"authorId"`
	Version     int       `json:"version"`
	CreatedAt   string    `json:"createdAt"`
}

func (b *Bid) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Status      string `json:"status"`
		AuthorType  string `json:"authorType"`
		Description string `json:"description"`
		AuthorId    string `json:"authorId"`
		Version     int    `json:"version"`
		CreatedAt   string `json:"createdAt"`
	}{
		Id:          b.Id,
		Name:        b.Name,
		Status:      b.Status.String(),
		AuthorType:  b.AuthorType.String(),
		Description: b.Description,
		Version:     b.Version,
		AuthorId:    b.AuthorId,
		CreatedAt:   b.CreatedAt,
	})
}
