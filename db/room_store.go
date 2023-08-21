package db

import (
	"context"

	"github.com/lewiscasewell/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	Dropper
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(c *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     c,
		coll:       c.Database(DBNAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.Id = res.InsertedID.(primitive.ObjectID)

	// update hotel with Room ID
	filter := bson.M{"_id": room.HotelId}
	update := bson.M{"$push": bson.M{"rooms": room.Id}}

	err = s.Update(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return room, nil
}
