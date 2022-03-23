package domain

type Rider struct {
	ID       uint
	Name     string
	Status   string
	Location Location
}

func NewRider(name string, status string, location Location) Rider {
	return Rider{Name: name, Status: status, Location: location}
}
