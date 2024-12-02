package controllers

import (
	"employee-asset-system/db"
	"employee-asset-system/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateAsset godoc
// @Summary Create a new asset
// @Description Adds a new asset to the database
// @Tags Assets
// @Accept json
// @Produce json
// @Param asset body models.Asset true "Asset data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets [post]
func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	asset.AssetID = uuid.New().String()
	asset.CreatedAt = time.Now()
	asset.UpdatedAt = time.Now()

	_, err := db.Database.Collection("asset").InsertOne(r.Context(), asset)
	if err != nil {
		http.Error(w, "Failed to create asset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Asset created successfully"})
}

// EditAsset godoc
// @Summary Edit an asset's details
// @Description Updates an asset's information by ID
// @Tags Assets
// @Accept json
// @Produce json
// @Param assetId path string true "Asset ID"
// @Param data body map[string]interface{} true "Updated data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets/{assetId} [put]
func EditAsset(w http.ResponseWriter, r *http.Request) {
	assetID := mux.Vars(r)["assetId"]

	var updatedData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedData["updated_at"] = time.Now()
	filter := bson.M{"asset_id": assetID}
	update := bson.M{"$set": updatedData}

	_, err := db.Database.Collection("asset").UpdateOne(r.Context(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update asset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Asset updated successfully"})
}

// GetAllAssets godoc
// @Summary Get all assets
// @Description Fetches all assets from the database
// @Tags Assets
// @Produce json
// @Success 200 {array} models.Asset
// @Failure 500 {object} map[string]string
// @Router /assets [get]
func GetAllAssets(w http.ResponseWriter, r *http.Request) {
	cursor, err := db.Database.Collection("asset").Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch assets", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var assets []models.Asset
	for cursor.Next(r.Context()) {
		var asset models.Asset
		if err := cursor.Decode(&asset); err == nil {
			assets = append(assets, asset)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assets)
}

// GetAssetById godoc
// @Summary Get an asset by ID
// @Description Fetches details of a single asset
// @Tags Assets
// @Produce json
// @Param assetId path string true "Asset ID"
// @Success 200 {object} models.Asset
// @Failure 404 {object} map[string]string
// @Router /assets/{assetId} [get]
func GetAssetById(w http.ResponseWriter, r *http.Request) {
	assetID := mux.Vars(r)["assetId"]

	var asset models.Asset
	err := db.Database.Collection("asset").FindOne(r.Context(), bson.M{"asset_id": assetID}).Decode(&asset)
	if err != nil {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(asset)
}

// DeleteAsset godoc
// @Summary Delete an asset
// @Description Deletes an asset from the database by ID
// @Tags Assets
// @Produce json
// @Param assetId path string true "Asset ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets/{assetId} [delete]
func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	assetID := mux.Vars(r)["assetId"]

	filter := bson.M{"asset_id": assetID}
	_, err := db.Database.Collection("asset").DeleteOne(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to delete asset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Asset deleted successfully"})
}
