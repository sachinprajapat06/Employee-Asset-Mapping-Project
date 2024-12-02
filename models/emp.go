package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmpID                  string             `bson:"emp_id" json:"emp_id"`
	FirstName              string             `bson:"first_name" json:"first_name"`
	LastName               string             `bson:"last_name" json:"last_name"`
	Gender                 string             `bson:"gender" json:"gender"`
	PhoneNumber            string             `bson:"phone_number" json:"phone_number"`
	EmployeeEmail          string             `bson:"employee_email" json:"employee_email"`
	Address                string             `bson:"address" json:"address"`
	BloodGroup             string             `bson:"blood_group" json:"blood_group"`
	EmergencyContactNumber string             `bson:"emergency_contact_number" json:"emergency_contact_number"`
	Password               string             `bson:"password" json:"password,omitempty"` // Use `omitempty` to exclude in JSON responses.
	CreatedAt              time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt              time.Time          `bson:"updated_at" json:"updated_at"`
}

type DashboardEmployee struct {
	EmpId                  string `bson:"emp_id" json:"EmpId"`
	FirstName              string `bson:"first_name" json:"FirstName"`
	LastName               string `bson:"last_name" json:"LastName"`
	Gender                 string `bson:"gender" json:"Gender"`
	PhoneNumber            string `bson:"phone_number" json:"PhoneNumber"`
	EmployeeEmail          string `bson:"employee_email" json:"EmployeeEmail"`
	Address                string `bson:"address" json:"Address"`
	BloodGroup             string `bson:"blood_group" json:"BloodGroup"`
	EmergencyContactNumber string `bson:"emergency_contact_number" json:"EmergencyContactNumber"`
	AssetCount             int    `bson:"asset_count" json:"AssetCount"`
}

type EmployeeList struct {
	Employees []DashboardEmployee `bson:"employees" json:"EmployeeList"`
}
