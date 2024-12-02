package routes

import (
	"employee-asset-system/controllers"
	_ "employee-asset-system/docs"
	"employee-asset-system/middleware"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {

	// Public Routes
	r.HandleFunc("/login/auth", controllers.Login).Methods("POST")

	// Swagger endpoint
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Serve Swagger JSON explicitly
	r.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../docs/swagger.json")
	})

	// Protected Routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Employee Routes
	api.HandleFunc("/employee/createemployee", controllers.CreateEmployee).Methods("POST")
	api.HandleFunc("/employee/editemployee/{employeeId}", controllers.EditEmployee).Methods("PUT")
	api.HandleFunc("/employee/deleteemployee/{employeeId}", controllers.DeleteEmployee).Methods("DELETE")
	api.HandleFunc("/employee/employee/{employeeId}", controllers.GetEmployeeById).Methods("GET")

	// Asset Routes
	api.HandleFunc("/asset/createasset", controllers.CreateAsset).Methods("POST")
	api.HandleFunc("/asset/editasset/{assetId}", controllers.EditAsset).Methods("PUT")
	api.HandleFunc("/asset/deleteasset/{assetId}", controllers.DeleteAsset).Methods("DELETE")
	api.HandleFunc("/asset/asset/{assetId}", controllers.GetAssetById).Methods("GET")
	api.HandleFunc("/asset/getallasset", controllers.GetAllAssets).Methods("GET")

	// Mapping Routes
	api.HandleFunc("/mapping/assignassetmapping", controllers.AssignAssetMapping).Methods("POST")
	api.HandleFunc("/mapping/getallassets/{employeeId}", controllers.GetAllAssetsMappedToEmployee).Methods("GET")
	api.HandleFunc("/mapping/removeassetmapping/{mappingId}", controllers.RemoveAssetMapping).Methods("DELETE")

	// Dashboard
	api.HandleFunc("/dashboard", controllers.GetAllEmployees).Methods("GET")

}
