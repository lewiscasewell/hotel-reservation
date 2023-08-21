package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	Id       primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	TwinRoomType
	SuiteRoomType
)

type Room struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelId   primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}
