package models

var (
	OrgAuthor  Author = &organization{}
	UserAuthor Author = &user{}
)

type Author interface {
	isTenderStatus()
	String() string
}

type organization struct{ TenderStatus }

func (o *organization) String() string {
	return "Organization"
}

type user struct{ TenderStatus }

func (u *user) String() string {
	return "User"
}
