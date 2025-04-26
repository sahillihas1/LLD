package main

type IParkingService interface {
	ParkVehicle(vehicle Vehicle) error
	UnparkVehicle(vehicle Vehicle) error
}

type Vehicle struct {
	NumberPlate string
	Color       string
	Type        VehicleType
}

type VehicleType int

const (
	Car VehicleType = iota
	Bike
)

type IparkingRepo interface {
	UpdateSpot() error
	GetSpot() error
}

type ParkingRepo struct {
	parkingSpots map[int]*ParkingSpot
}

type ParkingSpot struct {
	ID     int
	Level  int
	status bool
}

type IPaymentService interface {
	MakePayment() error
}

type CardService struct {
}

func (c *CardService) MakePayment() error {
	return nil
}

type CashService struct {
}

func (c *CashService) MakePayment() error {
	return nil
}

type ParkingService struct {
	parkingRepo IparkingRepo
	paymentServ IPaymentService
}

func (p *ParkingService) ParkVehicle(vehicle Vehicle) error {
	getParkingFactory
}
