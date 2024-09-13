package models

var (
	Construction EnumService = &construction{}
	Delivery     EnumService = &delivery{}
	Manufacture  EnumService = &manufacture{}
)

type EnumService interface {
	isEnumSevice()
	String() string
}

type construction struct{ EnumService }

func (c *construction) String() string {
	return "Construction"
}

type delivery struct{ EnumService }

func (d *delivery) String() string {
	return "Delivery"
}

type manufacture struct{ EnumService }

func (m *manufacture) String() string {
	return "Manufacture"
}
