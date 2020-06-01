package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// UserModel ... mongodb mapping
type UserModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	Email        string        `bson:"email" json:"email"`
	CreatedAt    time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt    time.Time     `bson:"updated_at",omitempty json:"updated_at,omitempty"`
}

// Token ... mongodb mapping
type TokenModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	UserId        bson.ObjectId        `bson:"user_id" json:"user_id"`
	Token     string        `bson:"token,omitempty" json:"token,omitempty" sql:"type:text"`
	CreatedAt    time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt    time.Time     `bson:"updated_at",omitempty json:"updated_at,omitempty"`
}