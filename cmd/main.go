package main

import (
	"fmt"
	"lumel/database"
	api "lumel/handler"

	"lumel/service"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("hello this is lumel assesment test")
	err := start()
	if err != nil {
		log.Panic().Err(err).Send()
	}
}
func start() error {
	log.Info().Msg("main : Started : Initializing db support")
	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}
	svc, err := service.NewService(db)
	svc.ReadDataFromCSV("C:/Users/JaVi204/OC/lumel assesment/lumel.csv")
	api := api.NewAPI(svc)

	http.HandleFunc("/total_customers", api.GetTotalCustomersHandler)
	http.HandleFunc("/total_orders", api.GetTotalOrdersHandler)
	http.HandleFunc("/average_order_value", api.GetAverageOrderValueHandler)
	http.HandleFunc("/refresh_csv", api.RefreshDataFromCSVHandler)

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		return fmt.Errorf("could not stop server gracefully %w", err)
	}
	return nil
}
