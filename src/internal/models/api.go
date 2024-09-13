package models

import (
	"encoding/json"
)

type TenderBuilder struct {
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	ServiceType    EnumService `json:"serviceType,omitempty"`
	Status         TenderStatus  `json:"status,omitempty"`
	OrganizationId string      `json:"organizationId,omitempty"`
	UserName       string      `json:"creatorUsername,omitempty"`
}

func (t *TenderBuilder) UnmarshalJSON(data []byte) error {
	input := struct {
		Name           string `json:"name,omitempty"`
		Description    string `json:"description,omitempty"`
		ServiceType    string `json:"serviceType,omitempty"`
		Status         string `json:"status,omitempty"`
		OrganizationId string `json:"organizationId,omitempty"`
		UserName       string `json:"creatorUsername,omitempty"`
	}{}
	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}
	var (
		service EnumService
		status  TenderStatus
	)
	switch input.ServiceType {
	case "Construction":
		service = Construction
	case "Delivery":
		service = Delivery
	case "Manufacture":
		service = Manufacture
	default:
		service = nil
	}

	switch input.Status {
	case "Created":
		status = TenderCreated
	case "Published":
		status = TenderPublished
	case "Closed":
		status = TenderClosed
	default:
		status = nil
	}
	t.Name = input.Name
	t.Description = input.Description
	t.ServiceType = service
	t.Status = status
	t.OrganizationId = input.OrganizationId
	t.UserName = input.UserName
	return nil
}
