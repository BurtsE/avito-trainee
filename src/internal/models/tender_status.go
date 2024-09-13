package models

var (
	TenderCreated   TenderStatus = &created{}
	TenderPublished TenderStatus = &published{}
	TenderClosed    TenderStatus = &closed{}
)

type TenderStatus interface {
	isTenderStatus()
	String() string
}

type created struct{ TenderStatus }

func (c *created) String() string {
	return "Created"
}

type published struct{ TenderStatus }

func (p *published) String() string {
	return "Published"
}

type closed struct{ TenderStatus }

func (c *closed) String() string {
	return "Closed"
}
