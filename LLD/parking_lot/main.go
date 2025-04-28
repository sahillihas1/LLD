package main

import "fmt"

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
	UpdateSpot(spotId *ParkingSpot) error
	GetSpots() []*ParkingSpot
}

type ParkingRepo struct {
	parkingSpots map[int]*ParkingSpot
}

func (p *ParkingRepo) UpdateSpot(spotId *ParkingSpot) error {
	if spotId == nil {
		return fmt.Errorf("invalid parking spot")
	}
	if _, exists := p.parkingSpots[spotId.ID]; !exists {
		return fmt.Errorf("parking spot with ID %d does not exist", spotId.ID)
	}
	p.parkingSpots[spotId.ID] = spotId
	return nil
}

func (p *ParkingRepo) GetSpots() []*ParkingSpot {
	spots := make([]*ParkingSpot, 0, len(p.parkingSpots))
	for _, spot := range p.parkingSpots {
		spots = append(spots, spot)
	}
	return spots
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
	getParkingStategy := p.getParkingStrategy(vehicle)
	if getParkingStategy == nil {
		return fmt.Errorf("no parking strategy found for vehicle type")
	}
	err := getParkingStategy.ParkVehicle(vehicle)
	if err != nil {
		return fmt.Errorf("error parking vehicle: %v", err)
	}
	return nil
}

func (p *ParkingService) UnParkVehicle(vehicle Vehicle) error {
	getParkingStategy := p.getParkingStrategy(vehicle)
	if getParkingStategy == nil {
		return fmt.Errorf("no parking strategy found for vehicle type")
	}
	err := getParkingStategy.ParkVehicle(vehicle)
	if err != nil {
		return fmt.Errorf("error parking vehicle: %v", err)
	}
	return nil
}

func getFeesStrategy(vehicle Vehicle) (int, error) {
	switch vehicle.Type {
	case Car:
		return 10, nil
	case Bike:
		return 5, nil
	default:
		return 0, fmt.Errorf("unknown vehicle type")
	}
}

func (p *ParkingService) getParkingStrategy(vehicle Vehicle) IParkingStrategy {
	switch vehicle.Type {
	case Car:
		return &CarParkingStrategy{}
	case Bike:
		return &BikeParkingStrategy{}
	}
	return nil
}

type IParkingStrategy interface {
	ParkVehicle(vehicle Vehicle) error
}

type CarParkingStrategy struct {
	ParkingRepo IparkingRepo
}

func (c *CarParkingStrategy) ParkVehicle(vehicle Vehicle) error {
	for _, spot := range c.ParkingRepo.GetSpots() {
		if !spot.status && spot.Level == 1 { // Assuming Level 1 is for cars
			spot.status = true
			err := c.ParkingRepo.UpdateSpot(spot)
			if err != nil {
				return fmt.Errorf("failed to update parking spot: %v", err)
			}
			fmt.Printf("Vehicle %s parked in spot %d\n", vehicle.NumberPlate, spot.ID)
			return nil
		}
	}
	return fmt.Errorf("no available parking spots for cars")
}

type BikeParkingStrategy struct {
	ParkingRepo IparkingRepo
}

func (b *BikeParkingStrategy) ParkVehicle(vehicle Vehicle) error {
	for _, spot := range b.ParkingRepo.GetSpots() {
		if !spot.status && spot.Level == 2 { // Assuming Level 2 is for bikes
			spot.status = true
			err := b.ParkingRepo.UpdateSpot(spot)
			if err != nil {
				return fmt.Errorf("failed to update parking spot: %v", err)
			}
			fmt.Printf("Vehicle %s parked in bike spot %d\n", vehicle.NumberPlate, spot.ID)
			return nil
		}
	}
	// If no bike spots are available, try car spots
	for _, spot := range b.ParkingRepo.GetSpots() {
		if !spot.status && spot.Level == 1 { // Assuming Level 1 is for cars
			spot.status = true
			err := b.ParkingRepo.UpdateSpot(spot)
			if err != nil {
				return fmt.Errorf("failed to update parking spot: %v", err)
			}
			fmt.Printf("Vehicle %s parked in car spot %d\n", vehicle.NumberPlate, spot.ID)
			return nil
		}
	}
	return fmt.Errorf("no available parking spots for bikes or cars")
}

func main() {
	// Example usage
	parkingRepo := &ParkingRepo{
		parkingSpots: make(map[int]*ParkingSpot),
	}
	paymentService := &CardService{}

	parkingService := &ParkingService{
		parkingRepo: parkingRepo,
		paymentServ: paymentService,
	}

	vehicle := Vehicle{
		NumberPlate: "ABC123",
		Color:       "Red",
		Type:        Car,
	}

	err := parkingService.ParkVehicle(vehicle)
	if err != nil {
		panic(err)
	}
}
