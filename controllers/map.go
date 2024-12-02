package controllers

import (
	"employee-asset-system/db"
	"employee-asset-system/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// AssignAssetMapping godoc
// @Summary Assign an asset to an employee
// @Description Assigns a new asset to an employee
// @Tags Asset Mapping
// @Accept json
// @Produce json
// @Param mapping body models.EmployeeAssetMapping true "Asset mapping data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /asset-mapping [post]
func AssignAssetMapping(w http.ResponseWriter, r *http.Request) {
	var mapping models.EmployeeAssetMapping
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	fmt.Println("here    ", mapping.EmployeeID)

	mapping.MappingID = uuid.New().String()
	mapping.AssignedDate = time.Now()
	mapping.Status = "active"

	_, err := db.Database.Collection("mapping").InsertOne(r.Context(), mapping)
	if err != nil {
		http.Error(w, "Failed to assign asset mapping", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Asset mapping assigned successfully"})
}

// GetAllAssetsMappedToEmployee godoc
// @Summary Get all assets mapped to an employee
// @Description Fetches all asset mappings for a specific employee
// @Tags Asset Mapping
// @Produce json
// @Param employeeId path string true "Employee ID"
// @Success 200 {array} models.EmployeeAssetMapping
// @Failure 500 {object} map[string]string
// @Router /asset-mapping/employee/{employeeId} [get]
func GetAllAssetsMappedToEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID := mux.Vars(r)["employeeId"]

	cursor, err := db.Database.Collection("mapping").Find(r.Context(), bson.M{"employee_id": employeeID})
	if err != nil {
		http.Error(w, "Failed to fetch mappings", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var mappings []models.EmployeeAssetMapping
	for cursor.Next(r.Context()) {
		var mapping models.EmployeeAssetMapping
		if err := cursor.Decode(&mapping); err == nil {
			mappings = append(mappings, mapping)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mappings)
}

// RemoveAssetMapping godoc
// @Summary Remove an asset mapping
// @Description Deletes a specific asset mapping by its ID
// @Tags Asset Mapping
// @Produce json
// @Param mappingId path string true "Mapping ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /asset-mapping/{mappingId} [delete]
func RemoveAssetMapping(w http.ResponseWriter, r *http.Request) {
	mappingID := mux.Vars(r)["mappingId"]

	filter := bson.M{"mapping_id": mappingID}
	_, err := db.Database.Collection("mapping").DeleteOne(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to remove asset mapping", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Asset mapping removed successfully"})
}
