package structs

import (
	"github.com/google/uuid"
)

type FlightStatus string

const (
	Scheduled FlightStatus = "scheduled"
	Active    FlightStatus = "active"
	Landed    FlightStatus = "landed"
	Cancelled FlightStatus = "cancelled"
	Incident  FlightStatus = "incident"
	Diverted  FlightStatus = "diverted"
)

type LiveFlights struct {
	ID           uuid.UUID    `db:"id"`
	FlightDate   string       `json:"flight_date"`
	FlightStatus FlightStatus `json:"flight_status"`
	Departure    struct {
		Airport         string      `json:"airport"`
		Timezone        string      `json:"timezone"`
		Iata            string      `json:"iata"`
		Icao            string      `json:"icao"`
		Terminal        string      `json:"terminal"`
		Gate            interface{} `json:"gate"`
		Delay           *int        `json:"delay"`
		Scheduled       string      `json:"scheduled"`
		Estimated       string      `json:"estimated"`
		Actual          interface{} `json:"actual"`
		EstimatedRunway interface{} `json:"estimated_runway"`
		ActualRunway    interface{} `json:"actual_runway"`
	} `json:"departure"`
	Arrival struct {
		Airport         string      `json:"airport"`
		Timezone        string      `json:"timezone"`
		Iata            string      `json:"iata"`
		Icao            string      `json:"icao"`
		Terminal        interface{} `json:"terminal"`
		Gate            interface{} `json:"gate"`
		Baggage         interface{} `json:"baggage"`
		Delay           *int        `json:"delay"`
		Scheduled       string      `json:"scheduled"`
		Estimated       string      `json:"estimated"`
		Actual          interface{} `json:"actual"`
		EstimatedRunway interface{} `json:"estimated_runway"`
		ActualRunway    interface{} `json:"actual_runway"`
	} `json:"arrival"`
	Airline struct {
		Name string `json:"name"`
		Iata string `json:"iata"`
		Icao string `json:"icao"`
	} `json:"airline"`
	Flight struct {
		Number     string `json:"number"`
		Iata       string `json:"iata"`
		Icao       string `json:"icao"`
		Codeshared struct {
			AirlineName  string `json:"airline_name"`
			AirlineIata  string `json:"airline_iata"`
			AirlineIcao  string `json:"airline_icao"`
			FlightNumber string `json:"flight_number"`
			FlightIata   string `json:"flight_iata"`
			FlightIcao   string `json:"flight_icao"`
		} `json:"codeshared"`
	} `json:"flight"`
	Aircraft interface{} `json:"aircraft"`
	Live     interface{} `json:"live"`
}

type FlightApiData struct {
	Pagination Pagination    `json:"pagination"`
	Data       []LiveFlights `json:"data"`
}
