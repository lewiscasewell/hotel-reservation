package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/lewiscasewell/hotel-reservation/db"
	"github.com/lewiscasewell/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(lname, fname, email, password string) {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Id:       primitive.NewObjectID(),
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}
	// rooms with random base prices
	rooms := []types.Room{
		{
			Size:      "small",
			BasePrice: float64(100) + (rand.Float64() * 100),
		},
		{
			Size:      "normal",
			BasePrice: float64(200) + (rand.Float64() * 100),
		},
		{
			Size:      "kingsize",
			BasePrice: float64(350) + (rand.Float64() * 100),
		},
	}
	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {

		room.HotelId = insertedHotel.Id
		_, err := roomStore.Insert(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func main() {
	seedHotel("Hilton", "London", 4)
	seedHotel("Marriot", "New York", 4)
	seedHotel("Hilton", "New York", 5)
	seedHotel("Marriot", "London", 3)
	seedUser("Casewell", "Lewis", "lewis@test.com", "password")
}

func init() {
	var err error
	ctx := context.Background()
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		panic(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
