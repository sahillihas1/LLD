package main

import "fmt"

type ParkingLotService struct {
	ticketingService ITicketingService
	parkingSerivce   IParkingService
}

type Vehicle struct {
	NumberPlate string
	Type        string
}

func NewParkingLotService(ticketingService ITicketingService, parkingSerivce IParkingService) *ParkingLotService {
	return &ParkingLotService{
		ticketingService: ticketingService,
		parkingSerivce:   parkingSerivce,
	}
}

type ITicketingService interface {
	GenerateTicket() string
}

type TicketingService struct {
}

func (t *TicketingService) GenerateTicket() string {
	return "ticket"
}

type IParkingService interface {
	ParkVehicle(vehicle Vehicle) error
}

type ParkingService struct {
}

func (p *ParkingService) ParkVehicle(vehicle Vehicle) error {
	fmt.Println("move vehicle to parking lot")
	return nil
}

func main() {
	ticketingService := &TicketingService{}
	parkingService := &ParkingService{}

	parkingLotService := NewParkingLotService(ticketingService, parkingService)

	vehicle := Vehicle{
		NumberPlate: "123",
		Type:        "car",
	}

	parkingLotService.parkingSerivce.ParkVehicle(vehicle)
}
