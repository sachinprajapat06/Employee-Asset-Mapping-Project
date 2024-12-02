package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Asset struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AssetID   string             `bson:"asset_id" json:"asset_id"`
	AssetName string             `bson:"asset_name" json:"asset_name"`
	AssetType string             `bson:"asset_type" json:"asset_type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
