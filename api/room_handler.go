package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lewiscasewell/hotel-reservation/db"
	"github.com/lewiscasewell/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	NumGuests int       `json:"numGuests"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.StartDate) || now.After(p.EndDate) {
		return fmt.Errorf("Start date and end date must be in the future")
	}

	if p.StartDate.After(p.EndDate) {
		return fmt.Errorf("Start date must be before end date")
	}

	if p.NumGuests <= 0 {
		return fmt.Errorf("Number of guests must be greater than 0")
	}

	return nil
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	id := c.Params("id")

	roomId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user, ok := c.Locals("user").(*types.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return fmt.Errorf("Error parsing body: %w", err)
	}

	if err := params.validate(); err != nil {
		return fmt.Errorf("Invalid params: %w", err)
	}

	available, err := h.isRoomAvailableToBook(c.Context(), roomId, params)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Type: "error",
			Msg:  "Error checking room availability",
		})
	}

	if !available {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Type: "error",
			Msg:  "Room is not available for those dates",
		})
	}

	booking := &types.Booking{
		UserId:    user.Id,
		RoomId:    roomId,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		NumGuests: params.NumGuests,
	}

	booking, err = h.store.Booking.Insert(c.Context(), booking)

	if err != nil {
		return fmt.Errorf("Error inserting booking: %w", err)
	}

	return c.JSON(booking)
}

func (h *RoomHandler) isRoomAvailableToBook(ctx context.Context, roomId primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomId": roomId,
		"startDate": bson.M{
			"$gte": params.StartDate,
		},
		"endDate": bson.M{
			"$lte": params.EndDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	return len(bookings) == 0, nil

}
