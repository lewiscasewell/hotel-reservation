package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lewiscasewell/hotel-reservation/db"
	"github.com/lewiscasewell/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)

	hotel := types.Hotel{
		Id:       primitive.NewObjectID(),
		Name:     "Hotel California",
		Location: "California",
		Rooms:    []primitive.ObjectID{},
	}
	room := types.Room{
		Type:      types.DoubleRoomType,
		BasePrice: 100,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	room.HotelId = insertedHotel.Id
	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("seeding database", insertedHotel)
	fmt.Println("seeding database", insertedRoom)
}
