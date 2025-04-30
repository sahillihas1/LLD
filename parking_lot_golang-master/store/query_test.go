package store

import (
	"fmt"

	"parking_lot/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parking lot store tests", func() {
	var (
		connection Store
	)
	connection = NewStore()
	It("Tear Down Store Data", func() {
		TearDown()
	})

	Context("parking_lot store execute", func() {
		TearDown()
		It("Create a parking lot with 2 slots", func() {
			cmd := &schema.Command{
				Command:   "create_parking_lot",
				Arguments: []string{"2"},
			}
			res, err := connection.CreateParkingLot().Execute(cmd)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(res).To(Equal(fmt.Sprintf(ParkinglotCreatedInfo, 2)))
		})

		It("park a vehicle", func() {
			cmd := &schema.Command{
				Command:   schema.CMDregistration_numbers_for_cars_with_colour,
				Arguments: []string{"Red"},
			}
			res, err := connection.Park().Execute(cmd)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(res).To(Equal("Allocated slot number: 1"))
		})

		It("Get Status", func() {
			cmd := &schema.Command{
				Command:   "status",
				Arguments: []string{},
			}
			res, err := connection.Status().Execute(cmd)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(res).To(Equal(res))
		})
	})
})
