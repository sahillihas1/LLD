package store

import (
	"fmt"
	"parking_lot/schema"
	"strconv"
	"strings"
)

type Query struct {
	*store
}

// NewQueryStore returns the store object
func NewQueryStore(st *store) *Query {
	q := &Query{st}
	return q
}

func (qc *Query) Execute(cmd *schema.Command) (string, error) {
	queryStrategy := qc.getQueryStrategy(cmd)
	resp, err := queryStrategy.ExecuteQuery(cmd.Arguments[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resp), nil
}

func (qc *Query) getQueryStrategy(cmd *schema.Command) QueryStrategy {
	switch cmd.Command {
	case string("registration_numbers_for_cars_with_colour"):
		return &RegNumbersByColourHandler{qc.store}
	case string("slot_numbers_for_cars_with_colour"):
		return &SlotNumbersByColourHandler{qc.store}
	case string("slot_number_for_registration_number"):
		return &SlotNumberByRegHandler{qc.store}
	default:
		return nil
	}
}

type QueryStrategy interface {
	ExecuteQuery(key string) (interface{}, error)
}

type RegNumbersByColourHandler struct {
	*store
}

func (h *RegNumbersByColourHandler) ExecuteQuery(key string) (interface{}, error) {
	var results []string
	for _, slot := range ParkingLot.Slots {
		if slot.Vehicle != nil && strings.EqualFold(slot.Vehicle.Colour, key) {
			results = append(results, slot.Vehicle.RegistrationNumber)
		}
	}
	if len(results) == 0 {
		return "Not found", nil
	}
	return strings.Join(results, ", "), nil
}

type SlotNumbersByColourHandler struct {
	*store
}

func (h *SlotNumbersByColourHandler) ExecuteQuery(key string) (interface{}, error) {
	var results []string
	for _, slot := range ParkingLot.Slots {
		if slot.Vehicle != nil && strings.EqualFold(slot.Vehicle.Colour, key) {
			results = append(results, strconv.Itoa(int(slot.ID)))
		}
	}
	if len(results) == 0 {
		return "Not found", nil
	}
	return strings.Join(results, ", "), nil
}

type SlotNumberByRegHandler struct {
	*store
}

func (h *SlotNumberByRegHandler) ExecuteQuery(key string) (interface{}, error) {
	for _, slot := range ParkingLot.Slots {
		if slot.Vehicle != nil && strings.EqualFold(slot.Vehicle.RegistrationNumber, key) {
			return strconv.Itoa(int(slot.ID)), nil
		}
	}
	return "Not found", nil
}
