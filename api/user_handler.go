package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/lewiscasewell/hotel-reservation/db"
	"github.com/lewiscasewell/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {

	id := c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "user not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	insertedUser, err := types.NewUserFromParams(&params)

	if err != nil {
		return err
	}
	if _, err := h.userStore.InsertUser(c.Context(), insertedUser); err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		// values bson.M
		params types.UpdateUserParams
		id     = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": id})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(map[string]string{"deleted": id})
}
