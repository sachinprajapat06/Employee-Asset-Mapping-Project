package controllers

import (
	"employee-asset-system/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey = []byte("your_secret_key")

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// @Summary Login for employees
// @Description Employee login using phone number or email and password
// @Tags Login
// @Accept  json
// @Produce  json
// @Param   request body LoginRequest true "Login Request Body"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /login/auth [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// MongoDB lookup (replace with proper client/DB setup)
	client, _ := mongo.NewClient() // Add appropriate options
	collection := client.Database("employee_asset_db").Collection("employees")

	var employee bson.M
	err := collection.FindOne(r.Context(), bson.M{
		"$or": []bson.M{
			{"phone_number": req.Identifier},
			{"employee_email": req.Identifier},
		},
	}).Decode(&employee)

	if err != nil || !utils.CheckPassword(req.Password, employee["password"].(string)) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := GenerateJWT(employee["emp_id"].(string))
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func GenerateJWT(empId string) string {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Subject:   empId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}
