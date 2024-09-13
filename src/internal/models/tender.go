package models

import "encoding/json"

type Tender struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      TenderStatus  `json:"status"`
	ServiceType EnumService `json:"serviceType"`
	Version     int         `json:"version"`
	CreatedAt   string      `json:"createdAt"`
}

func (t *Tender) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		ServiceType string `json:"serviceType"`
		Version     int    `json:"version"`
		CreatedAt   string `json:"createdAt"`
	}{
		Id:          t.Id,
		Name:        t.Name,
		Description: t.Description,
		Status:      t.Status.String(),
		ServiceType: t.ServiceType.String(),
		Version:     t.Version,
		CreatedAt:   t.CreatedAt,
	})
}
