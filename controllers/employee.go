package controllers

import (
	"employee-asset-system/db"
	"employee-asset-system/models"
	"employee-asset-system/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Adds a new employee to the database
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees [post]
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	employee.EmpID = uuid.New().String()
	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	employee.Password = utils.HashPassword(employee.Password)

	_, err := db.Database.Collection("employee").InsertOne(r.Context(), employee)
	if err != nil {
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee created successfully"})
}

// EditEmployee godoc
// @Summary Edit an employee's details
// @Description Updates an employee's information by ID
// @Tags Employees
// @Accept json
// @Produce json
// @Param employeeId path string true "Employee ID"
// @Param data body map[string]interface{} true "Updated data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{employeeId} [put]
func EditEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID := mux.Vars(r)["employeeId"]

	var updatedData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedData["updated_at"] = time.Now()
	filter := bson.M{"emp_id": employeeID}
	update := bson.M{"$set": updatedData}

	_, err := db.Database.Collection("employee").UpdateOne(r.Context(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee updated successfully"})
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Deletes an employee from the database by ID
// @Tags Employees
// @Produce json
// @Param employeeId path string true "Employee ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{employeeId} [delete]
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID := mux.Vars(r)["employeeId"]

	filter := bson.M{"emp_id": employeeID}
	_, err := db.Database.Collection("employee").DeleteOne(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee deleted successfully"})
}

// GetEmployeeById godoc
// @Summary Get an employee by ID
// @Description Fetches details of a single employee
// @Tags Employees
// @Produce json
// @Param employeeId path string true "Employee ID"
// @Success 200 {object} models.Employee
// @Failure 404 {object} map[string]string
// @Router /employees/{employeeId} [get]
func GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	employeeID := mux.Vars(r)["employeeId"]

	var employee models.Employee
	err := db.Database.Collection("employee").FindOne(r.Context(), bson.M{"emp_id": employeeID}).Decode(&employee)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employee)
}

// GetAllEmployees godoc
// @Summary Get all employees with asset count
// @Description Fetches all employees and their asset counts
// @Tags Employees
// @Produce json
// @Success 200 {object} models.EmployeeList
// @Failure 500 {object} map[string]string
// @Router /employees [get]
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {

	coll := db.Database.Collection("employee")
	// Define aggregation pipeline
	pipeline := mongo.Pipeline{
		{
			{"$lookup", bson.D{
				{"from", "mapping"},
				{"localField", "emp_id"},
				{"foreignField", "employee_id"},
				{"as", "assets"},
			}},
		},
		{
			{"$addFields", bson.D{
				{"asset_count", bson.D{{"$size", "$assets"}}},
			}},
		},
		{
			{"$project", bson.D{
				{"emp_id", 1},
				{"first_name", 1},
				{"last_name", 1},
				{"gender", 1},
				{"phone_number", 1},
				{"employee_email", 1},
				{"address", 1},
				{"blood_group", 1},
				{"emergency_contact_number", 1},
				{"asset_count", 1},
			}},
		},
	}

	cursor, err := coll.Find(r.Context(), pipeline)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var employees []models.DashboardEmployee
	if err := cursor.All(r.Context(), &employees); err != nil {
		http.Error(w, "Failed to parse aggregation results", http.StatusInternalServerError)
		return
	}
	var data models.EmployeeList
	data.Employees = employees

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
