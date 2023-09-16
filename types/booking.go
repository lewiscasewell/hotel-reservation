package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
	RoomId    primitive.ObjectID `bson:"roomId,omitempty" json:"roomId"`
	StartDate time.Time          `bson:"startDate" json:"startDate"`
	EndDate   time.Time          `bson:"endDate" json:"endDate"`
	NumGuests int                `bson:"numGuests,omitempty" json:"numGuests"`
	Status    string             `bson:"status,omitempty" json:"status"`
}
