package main

import (
	"CotoChallenge/internal/domain"
	"CotoChallenge/internal/repository"
	"CotoChallenge/internal/service"
	handler "CotoChallenge/internal/transport/http"
	"CotoChallenge/internal/transport/http/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"time"
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

	// adding records to repository for testing purposes
	seedSales(repository)

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
	mux.HandleFunc("GET /api/v1/volume", saleHandler.GetTotalVolume)
	mux.HandleFunc("GET /api/v1/volume/center", saleHandler.GetVolumeByCenter)
	mux.HandleFunc("GET /api/v1/sales/percentage", saleHandler.GetSalesPercentegeByCenterOverTotalSales)

	handler := middleware.TimingMiddleware(logger)(mux)
	handler = middleware.Recovery(logger)(handler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	logger.Printf("Starting server on port %s", port)
	serverErr := server.ListenAndServe()
	if serverErr != nil {
		logger.Fatalf("Error starting server: %v", serverErr)
	}
}

// seedSales pre-loads the repository with mock sales data.
// Distribution is intentionally uneven to produce varied percentages per center/model.
// Total: 20 sales → each sale = 5% of global units.
//
// Expected percentages:
// C1 → Sedan: 20%, SUV: 10%, Sport: 5%
// C2 → Sport: 15%, Offroad: 5%, Sedan: 5%
// C3 → Offroad: 15%, SUV: 5%
// C4 → Sedan: 5%, SUV: 5%, Sport: 5%, Offroad: 5%
func seedSales(repo repository.SaleRepository) {
	sales := []domain.Sale{
		// Centro 1 — dominant in Sedan (4 units)
		{Id: "1", Vehicle: domain.Sedan, Price: domain.Sedan.FinalPrice(), Date: time.Now().AddDate(0, 0, -15), Center: "C1"},
		{Id: "2", Vehicle: domain.Sedan, Price: domain.Sedan.FinalPrice(), Date: time.Now().AddDate(0, 0, -12), Center: "C1"},
		{Id: "3", Vehicle: domain.Sedan, Price: domain.Sedan.FinalPrice(), Date: time.Now().AddDate(0, 0, -10), Center: "C1"},
		{Id: "4", Vehicle: domain.Sedan, Price: domain.Sedan.FinalPrice(), Date: time.Now().AddDate(0, 0, -8), Center: "C1"},
		{Id: "5", Vehicle: domain.SUV, Price: domain.SUV.FinalPrice(), Date: time.Now().AddDate(0, 0, -7), Center: "C1"},
		{Id: "6", Vehicle: domain.SUV, Price: domain.SUV.FinalPrice(), Date: time.Now().AddDate(0, 0, -5), Center: "C1"},
		{Id: "7", Vehicle: domain.Sport, Price: domain.Sport.FinalPrice(), Date: time.Now().AddDate(0, 0, -3), Center: "C1"},

		// Centro 2 — dominant in Sport (3 units)
		{Id: "8", Vehicle: domain.Sport, Price: domain.Sport.FinalPrice(), Date: time.Now().AddDate(0, 0, -14), Center: "C2"},
		{Id: "9", Vehicle: domain.Sport, Price: domain.Sport.FinalPrice(), Date: time.Now().AddDate(0, 0, -11), Center: "C2"},
		{Id: "10", Vehicle: domain.Sport, Price: domain.Sport.FinalPrice(), Date: time.Now().AddDate(0, 0, -9), Center: "C2"},
		{Id: "11", Vehicle: domain.Offroad, Price: domain.Offroad.FinalPrice(), Date: time.Now().AddDate(0, 0, -6), Center: "C2"},
		{Id: "12", Vehicle: domain.Sedan, Price: domain.Sedan.FinalPrice(), Date: time.Now().AddDate(0, 0, -4), Center: "C2"},

		// Centro 3 — dominant in Offroad (3 units)
		{Id: "13", Vehicle: domain.Offroad, Price: domain.Offroad.FinalPrice(), Date: time.Now().AddDate(0, 0, -13), Center: "C3"},
		{Id: "14", Vehicle: domain.Offroad, Price: domain.Offroad.FinalPrice(), Date: time.Now().AddDate(0, 0, -10), Center: "C3"},
		{Id: "15", Vehicle: domain.Offroad, Price: domain.Offroad.FinalPrice(), Date: time.Now().AddDate(0, 0, -7), Center: "C3"},
		{Id: "16", Vehicle: domain.SUV, Price: domain.SUV.FinalPrice(), Date: time.Now().AddDate(0, 0, -2), Center: "C3"},

		// Centro 4 — balanced, one of each model
		{Id: "17", Vehicle: domain.Sedan, Price: domain.Sedan.FinalPrice(), Date: time.Now().AddDate(0, 0, -9), Center: "C4"},
		{Id: "18", Vehicle: domain.SUV, Price: domain.SUV.FinalPrice(), Date: time.Now().AddDate(0, 0, -6), Center: "C4"},
		{Id: "19", Vehicle: domain.Sport, Price: domain.Sport.FinalPrice(), Date: time.Now().AddDate(0, 0, -3), Center: "C4"},
		{Id: "20", Vehicle: domain.Offroad, Price: domain.Offroad.FinalPrice(), Date: time.Now().AddDate(0, 0, -1), Center: "C4"},
	}

	for _, s := range sales {
		err := repo.CreateSale(context.Background(), s)
		if err != nil {
			continue
		}
	}
}
