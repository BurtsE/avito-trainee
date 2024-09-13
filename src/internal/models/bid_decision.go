package models

var (
	BidApproved Decision = &approved{}
	BidRejected Decision = &rejected{}
)

type Decision interface {
	isDesition()
	String() string
}

type approved struct{ Decision }

func (a *approved) String() string {
	return "Approved"
}

type rejected struct{ Decision }

func (r *rejected) String() string {
	return "Rejected"
}
