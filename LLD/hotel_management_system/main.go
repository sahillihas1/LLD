package hotel_management_system

type HotelManagementSystem struct {
	roomService    IRoomService
	PaymentService IPaymentService
	BookingService IBookingService
}

func NewHotelManagementSystem(roomService IRoomService, paymentService IPaymentService, bookingService IBookingService) *HotelManagementSystem {
	return &HotelManagementSystem{
		roomService:    roomService,
		PaymentService: paymentService,
		BookingService: bookingService,
	}
}

type IRoomDecorator interface {
	BookRoom(room Room) error
}

type SingleRoomDecorator struct {
}

func (s SingleRoomDecorator) BookRoom(room Room) error {
	//TODO implement me
	panic("implement me")
}

type Room struct {
	roomType string
}

type IRoomService interface {
	GetAvailableRooms() []Room
	BookRoom(room Room) error
}

type RoomService struct {
	rooms               map[string]Room
	singleRoomDecorator IRoomDecorator
}

func NewRoomService() *RoomService {
	return &RoomService{}
}

func (r *RoomService) GetAvailableRooms() []Room {
	return nil
}

func getRoomDecorator(roomType string) IRoomDecorator {
	if roomType == "single" {
		return &SingleRoomDecorator{}
	}
	return nil
}

func (r *RoomService) BookRoom(room Room) error {
	if room.roomType == "single" {
		err := getRoomDecorator(room.roomType).BookRoom(room)
		if err != nil {
			return err
		}
	}
	return nil
}

type IPaymentService interface {
	MakePayment(amount float64) error
}

type CardPayment struct {
}

func (c *CardPayment) MakePayment(amount float64) error {
	return nil
}

type UPIPayment struct {
}

func (u *UPIPayment) MakePayment(amount float64) error {
	return nil
}

type IBookingService interface {
	BookRoom(room Room, paymentService IPaymentService) error
}

type BookingService struct {
	roomService    IRoomService
	PaymentService IPaymentService
}

func NewBookingService(roomService IRoomService, paymentService IPaymentService) *BookingService {
	return &BookingService{
		roomService:    roomService,
		PaymentService: paymentService,
	}
}

func (b *BookingService) BookRoom(room Room, paymentService IPaymentService) error {
	return nil
}

func main() {
	roomService := NewRoomService()
	cardPayment := &CardPayment{}
	bookingService := NewBookingService(roomService, cardPayment)
	hotelManagementSystem := NewHotelManagementSystem(roomService, cardPayment, bookingService)
	room := Room{}
	hotelManagementSystem.BookingService.BookRoom(room, cardPayment)
}
