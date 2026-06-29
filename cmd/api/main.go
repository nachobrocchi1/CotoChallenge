package main

import (
	"CotoChallenge/internal/repository"
	"CotoChallenge/internal/service"
	handler "CotoChallenge/internal/transport/http"
	"CotoChallenge/internal/transport/http/middleware"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default value
	}

	logger := log.New(os.Stdout, "CotoChallenge: ", log.LstdFlags)

	repository, err := repository.NewInMemorySaleRepository()
	if err != nil {
		log.Fatalf("init Repository error: %v", err)
	}

	service, err := service.NewSaleService(logger, repository)
	if err != nil {
		log.Fatalf("init Service error: %v", err)
	}
	saleHandler, err := handler.NewSaleHandler(service)
	if err != nil {
		log.Fatalf("init Handler error: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/sales", saleHandler.CreateSale)
	mux.HandleFunc("GET /api/v1/sales", saleHandler.GetSales)
	mux.HandleFunc("GET /api/v1/sales/center", saleHandler.GetSalesByCenter)
	mux.HandleFunc("GET /api/v1/sales/percentage", saleHandler.GetSalesPercentegeByCenterOverTotalSales)

	handlerWithTiming := middleware.TimingMiddleware(logger)(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handlerWithTiming,
	}

	logger.Printf("Starting server on port %s", port)
	serverErr := server.ListenAndServe()
	if serverErr != nil {
		logger.Fatalf("Error starting server: %v", serverErr)
	}
}
