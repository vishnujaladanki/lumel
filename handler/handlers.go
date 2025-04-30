package api

import (
	"encoding/json"
	"fmt"
	"log"
	"lumel/service"
	"net/http"
	"time"
)

type API struct {
	service service.Service
}

func NewAPI(s service.Service) *API {
	return &API{
		service: s,
	}
}

func (a *API) GetTotalCustomersHandler(w http.ResponseWriter, r *http.Request) {

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

	totalCustomers, err := a.service.GetTotalCustomers(startTime, endTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"total_customers": totalCustomers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (a *API) GetTotalOrdersHandler(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

	totalOrders, err := a.service.GetTotalOrders(startTime, endTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"total_orders": totalOrders,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (a *API) GetAverageOrderValueHandler(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

	averageOrderValue, err := a.service.GetAverageOrderValue(startTime, endTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]float64{
		"average_order_value": averageOrderValue,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (a *API) RefreshDataFromCSVHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "lumel assesment/lumel.csv"

	err := a.service.ReadDataFromCSV(filePath)
	if err != nil {
		log.Printf("Error refreshing data from CSV: %v", err)
		http.Error(w, fmt.Sprintf("Error refreshing data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Data refreshed successfully from CSV."}`))
}
