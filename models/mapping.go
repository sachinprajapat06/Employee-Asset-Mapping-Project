package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmployeeAssetMapping struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MappingID    string             `bson:"mapping_id" json:"mapping_id"`
	EmployeeID   string             `bson:"employee_id" json:"employee_id"`
	AssetID      string             `bson:"asset_id" json:"asset_id"`
	AssignedDate time.Time          `bson:"assigned_date" json:"assigned_date"`
	Status       string             `bson:"status" json:"status"`
	Notes        string             `bson:"notes" json:"notes"`
}
