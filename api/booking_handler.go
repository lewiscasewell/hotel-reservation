package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lewiscasewell/hotel-reservation/db"
	"github.com/lewiscasewell/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	bookingStore *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		bookingStore: store,
	}
}

// TODO: this needs to be admin only
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.bookingStore.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// TODO: this needs to be user authenticated
func (h *BookingHandler) HandleGetUserBooking(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg:  "User not found",
		})
	}

	booking, err := h.bookingStore.Booking.GetBookingById(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}

	if booking.UserId != user.Id {
		return c.Status(http.StatusForbidden).JSON(genericResponse{
			Type: "error",
			Msg:  "You are not authorized to view this booking",
		})
	}

	return c.JSON(booking)

}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	booking, err := h.bookingStore.Booking.GetBookingById(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Locals("user").(*types.User)

	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg:  "User not found",
		})
	}

	if booking.UserId != user.Id && user.Role != "admin" {
		return c.Status(http.StatusForbidden).JSON(genericResponse{
			Type: "error",
			Msg:  "You are not authorized to cancel this booking",
		})
	}

	if booking.Status == "cancelled" {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Type: "error",
			Msg:  "Booking is already cancelled",
		})
	}

	booking.Status = "cancelled"

	err = h.bookingStore.Booking.CancelBooking(c.Context(), booking.Id.Hex())
	if err != nil {
		return err
	}

	return c.JSON(genericResponse{
		Type: "success",
		Msg:  "Booking cancelled successfully",
	})
}
