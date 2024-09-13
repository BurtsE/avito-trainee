package models

var (
	BidCreated   BidStatus = &bidCreated{}
	BidPublished BidStatus = &bidPublished{}
	BidCanceled  BidStatus = &bidCanceled{}
)

type BidStatus interface {
	isBidStatus()
	String() string
}

type bidCreated struct{ BidStatus }

func (c *bidCreated) String() string {
	return "Created"
}

type bidPublished struct{ BidStatus }

func (p *bidPublished) String() string {
	return "Published"
}

type bidCanceled struct{ BidStatus }

func (c *bidCanceled) String() string {
	return "Canceled"
}
