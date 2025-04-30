package store

import (
	"fmt"
	"parking_lot/errors"
	"parking_lot/schema"
	"strconv"
)

type leaveStore struct {
	*store
}

// NewLeaveStore returns a new leaveStore object
func NewLeaveStore(st *store) *leaveStore {
	return &leaveStore{st}
}

func (pl *leaveStore) IsHelp(arg string) (string, bool) {
	if arg == string(schema.CMDHelp) {
		return schema.CMDParkHint, true
	}
	return "", false
}

// Execute - `leave` command takes a slot number as an argument,
// and makes it available for future parking.
func (ls *leaveStore) Execute(cmd *schema.Command) (string, error) {
	if res, isHelp := ls.IsHelp(cmd.Arguments[0]); isHelp {
		return res, nil
	}
	if ParkingLot == nil {
		return "", errors.ErrNoParkingLot
	}
	// Validate slot ID
	slotID, err := strconv.Atoi(cmd.Arguments[0])
	if err != nil || slotID <= 0 {
		return "", errors.ErrInvalidSlotID
	}

	slot := ParkingLot.GetSlotByID(slotID)
	if slot == nil {
		return "", errors.ErrInvalidSlotID
	}
	if slot.IsSlotAvailable() {
		return "", errors.ErrInvalidSlotID
	}

	// Remove the vehicle
	vehicle := slot.GetParkedVehicle()
	err = slot.RemoveVehicle()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("done %s", vehicle.RegistrationNumber), nil
}
