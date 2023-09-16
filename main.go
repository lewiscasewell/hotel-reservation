package main

import (
	"context"
	"flag"

	"github.com/lewiscasewell/hotel-reservation/db"
	"github.com/lewiscasewell/hotel-reservation/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/lewiscasewell/hotel-reservation/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var fiberConfig = fiber.Config{
	// Override default error handler
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", "localhost:5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}

	var (
		store = &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   db.NewMongoHotelStore(client),
			Room:    db.NewMongoRoomStore(client, db.NewMongoHotelStore(client)),
			Booking: db.NewMongoBookingStore(client),
		}
		authHandler    = api.NewAuthHandler(store.User)
		userHandler    = api.NewUserHandler(store.User)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(fiberConfig)
		authApi        = app.Group("/api/auth")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(store.User))
		adminApi       = apiv1.Group("/admin", middleware.AdminAuth)
	)

	// Auth routes
	authApi.Post("/", authHandler.HandleAuthenticate)

	// User routes
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// Hotel routes
	apiv1.Get("/hotel", hotelHandler.HandletGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetHotelRooms)

	// room routes
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/room", roomHandler.HandleGetRooms)

	// booking routes
	adminApi.Get("/booking", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetUserBooking)
	apiv1.Put("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	app.Listen(*listenAddr)
}
