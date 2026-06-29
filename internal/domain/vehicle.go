package domain

// VehicleType is the type of car model
type VehicleType string

const (
	Sedan   VehicleType = "sedan"
	SUV     VehicleType = "suv"
	Offroad VehicleType = "offroad"
	Sport   VehicleType = "sport"
)

// Prices is the price for each model
var Prices = map[VehicleType]float64{
	Sedan:   8000,
	SUV:     9500,
	Offroad: 12500,
	Sport:   18200,
}

// SportTaxRate is the tax rate for the sport model
const SportTaxRate = 0.07

// FinalPrice calculates the final price applying taxes if applicable.
func (m VehicleType) FinalPrice() float64 {
	price := Prices[m]
	if m == Sport {
		price += price * SportTaxRate
	}
	return price
}
