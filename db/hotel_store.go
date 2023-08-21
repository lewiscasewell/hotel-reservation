package db

import (
	"context"

	"github.com/lewiscasewell/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	Dropper
	InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(c *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: c,
		coll:   c.Database(DBNAME).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func (s *MongoHotelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.Id = res.InsertedID.(primitive.ObjectID)

	return hotel, nil
}
