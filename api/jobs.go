package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/FACorreiaa/go-ollama/api/structs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
	"time"
)

type JobAPIMethods interface {
	startCityJobVerification()
}

type RepositoryJob struct {
	Conn *pgxpool.Pool
}

func NewRepositoryJob(db *pgxpool.Pool) *RepositoryJob {
	return &RepositoryJob{Conn: db}
}

func NewServiceJob(repo *RepositoryJob) *ServiceJob {
	return &ServiceJob{repo: repo}
}

type ServiceJob struct {
	repo    *RepositoryJob
	cityJob JobAPIMethods
}

// getExistingID retrieves existing table_id from the database
func (s *ServiceJob) getExistingID(query string, id int, tableData map[int]struct{}) (map[int]struct{}, error) {
	rows, err := s.repo.Conn.Query(context.Background(), query)
	if err != nil {
		handleError(err, "Error querying DB")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			handleError(err, "Error scanning IDs")
			return nil, err
		}
		tableData[id] = struct{}{}
	}

	return tableData, nil
}

// findNewCityData identifies new city data by comparing the API data with existing data
func (s *ServiceJob) findNewCityData(apiData []structs.City, tableData map[int]struct{}) []structs.City {
	var newData []structs.City

	for _, a := range apiData {
		if _, exists := tableData[a.CityID]; !exists {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewCountryData(apiData []structs.Country, tableData map[int]struct{}) []structs.Country {
	var newData []structs.Country

	for _, a := range apiData {
		if _, exists := tableData[a.CountryIsoNumeric]; !exists {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAirportData(apiData []structs.Airport, tableData map[int]struct{}) []structs.Airport {
	var newData []structs.Airport

	for _, a := range apiData {
		if _, exists := tableData[a.AirportId]; !exists {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAirplaneData(apiData []structs.Airplane, tableData map[int]struct{}) []structs.Airplane {
	var newData []structs.Airplane

	for _, a := range apiData {
		if _, exists := tableData[a.AirplaneId]; !exists {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewTaxData(apiData []structs.Tax, tableData map[int]struct{}) []structs.Tax {
	var newData []structs.Tax

	for _, a := range apiData {
		if _, exists := tableData[a.TaxId]; !exists {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAirlineData(apiData []structs.Airline, tableData map[int]struct{}) []structs.Airline {
	var newData []structs.Airline

	for _, a := range apiData {
		if _, exists := tableData[a.AirlineId]; !exists {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAircraftData(apiData []structs.Aircraft, tableData map[int]struct{}) []structs.Aircraft {
	var newData []structs.Aircraft

	for _, a := range apiData {
		if _, exists := tableData[a.PlaneTypeId]; !exists {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) insertNewCities() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("cities", 1000000)
	query := `select city_id from city`
	tableData := make(map[int]struct{})
	var cityID int

	apiData, err := os.ReadFile("./api/data/cities.json")
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.CityApiData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, cityID, tableData)

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewCityData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"city"},
			[]string{"gmt", "city_id", "iata_code", "country_iso2", "geoname_id",
				"latitude", "longitude", "city_name", "timezone", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				city := newDataMap[i]
				return []interface{}{
					city.GMT,
					city.CityID,
					city.IataCode,
					city.CountryISO2,
					city.GeonameID,
					city.Latitude,
					city.Longitude,
					city.CityName,
					city.Timezone,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into cities table")
			return err
		}

		slog.Info("New data inserted into the city table")
	} else {
		slog.Info("No new data to insert into the city table")
	}

	return nil
}

func (s *ServiceJob) insertNewCountries() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("countries", 1000000)
	query := `select country_iso_numeric from country`
	tableData := make(map[int]struct{})
	var countryIsoNumeric int

	apiData, err := os.ReadFile("./api/data/countries.json")
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.CountryApiData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, countryIsoNumeric, tableData)

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewCountryData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"country"},
			[]string{"country_name", "country_iso2", "country_iso3", "country_iso_numeric", "population",
				"capital", "continent", "currency_name", "currency_code", "fips_code",
				"phone_prefix", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				country := newDataMap[i]
				return []interface{}{
					country.CountryName,
					country.CountryISO2,
					country.CountryIso3,
					country.CountryIsoNumeric,
					country.Population,
					country.Capital,
					country.Continent,
					country.CurrencyName,
					country.CurrencyCode,
					country.FipsCode,
					country.PhonePrefix,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into country table")
			return err
		}

		slog.Info("New data inserted into the country table")
	} else {
		slog.Info("No new data to insert into the country table")
	}

	return nil
}

func (s *ServiceJob) insertNewAirports() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("countries", 1000000)
	query := `select airport_id from airport`
	tableData := make(map[int]struct{})
	var airport_id int

	apiData, err := os.ReadFile("./api/data/airports.json")
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AirportApiData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, airport_id, tableData)

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAirportData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"airport"},
			[]string{"gmt", "airport_id", "iata_code", "city_iata_code", "icao_code",
				"country_iso2", "geoname_id", "latitude", "longitude", "airport_name",
				"country_name", "phone_number", "timezone", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				airport := newDataMap[i]
				return []interface{}{
					airport.GMT, airport.AirportId, airport.IataCode,
					airport.CityIataCode, airport.IcaoCode, airport.CountryISO2,
					airport.GeonameID, airport.Latitude, airport.Longitude,
					airport.AirportName, airport.CountryName, airport.PhoneNumber,
					airport.Timezone, formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into airport table")
			return err
		}

		slog.Info("New data inserted into the airport table")
	} else {
		slog.Info("No new data to insert into the airport table")
	}

	return nil
}

func (s *ServiceJob) insertNewAirplanes() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("countries", 1000000)
	query := `select airplane_id from airplane`
	tableData := make(map[int]struct{})
	var airplaneID int

	apiData, err := os.ReadFile("./api/data/airplane.json")
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AirplaneApiData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, airplaneID, tableData)

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAirplaneData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"airplane"},
			[]string{"iata_type", "airplane_id", "airline_iata_code", "iata_code_long", "iata_code_short",
				"airline_icao_code", "construction_number", "delivery_date", "engines_count", "engines_type",
				"first_flight_date", "icao_code_hex", "line_number", "model_code", "registration_number",
				"test_registration_number", "plane_age", "plane_class", "model_name", "plane_owner", "plane_series",
				"plane_status", "production_line", "registration_date", "rollout_date", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				airplane := newDataMap[i]
				return []interface{}{
					airplane.IataType,
					airplane.AirplaneId,
					airplane.AirlineIataCode,
					airplane.IataCodeLong,
					airplane.IataCodeShort,
					airplane.AirlineIcaoCode,
					airplane.ConstructionNumber,
					airplane.DeliveryDate.Time,
					airplane.EnginesCount,
					airplane.EnginesType,
					airplane.FirstFlightDate.Time,
					airplane.IcaoCodeHex,
					airplane.LineNumber,
					airplane.ModelCode,
					airplane.RegistrationNumber,
					airplane.TestRegistrationNumber,
					airplane.PlaneAge,
					airplane.PlaneClass,
					airplane.ModelName,
					airplane.PlaneOwner,
					airplane.PlaneSeries,
					airplane.PlaneStatus,
					airplane.ProductionLine,
					airplane.RegistrationDate.Time,
					airplane.RolloutDate.Time,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into airport table")
			return err
		}

		slog.Info("New data inserted into the airport table")
	} else {
		slog.Info("No new data to insert into the airport table")
	}

	return nil
}

func (s *ServiceJob) insertNewTax() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("countries", 1000000)
	query := `select tax_id from tax`
	tableData := make(map[int]struct{})
	var tax_id int

	apiData, err := os.ReadFile("./api/data/tax.json")
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.TaxApiData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, tax_id, tableData)

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewTaxData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"tax"},
			[]string{"tax_id", "tax_name", "iata_code", "created_at"},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				tax := newDataMap[i]
				return []interface{}{
					tax.TaxId, tax.TaxName, tax.IataCode,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into tax table")
			return err
		}

		slog.Info("New data inserted into the tax table")
	} else {
		slog.Info("No new data to insert into the tax table")
	}

	return nil
}

func (s *ServiceJob) insertNewAirline() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("countries", 1000000)
	query := `select airline_id from airline`
	tableData := make(map[int]struct{})
	var airlineID int

	apiData, err := os.ReadFile("./api/data/airline.json")
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AirlineApiData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, airlineID, tableData)

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAirlineData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"airline"},
			[]string{"fleet_average_age", "airline_id", "callsign", "hub_code", "iata_code", "icao_code", "country_iso2",
				"date_founded", "iata_prefix_accounting", "airline_name", "country_name", "fleet_size", "status", "type",
				"created_at",
			}, pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				airline := newDataMap[i]
				return []interface{}{
					airline.FleetAverageAge,
					airline.AirlineId,
					airline.Callsign,
					airline.HubCode,
					airline.IataCode,
					airline.IcaoCode,
					airline.CountryISO2,
					airline.DateFounded,
					airline.IataPrefixAccounting,
					airline.AirlineName,
					airline.CountryName,
					airline.FleetSize,
					airline.Status,
					airline.Type,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into airline table")
			return err
		}

		slog.Info("New data inserted into the airline table")
	} else {
		slog.Info("No new data to insert into the airline table")
	}

	return nil
}

func (s *ServiceJob) insertNewAircraft() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("countries", 1000000)
	query := `select plane_type_id from aircraft`
	tableData := make(map[int]struct{})
	var planeTypeID int

	apiData, err := os.ReadFile("./api/data/aircraft.json")
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AircraftApiData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, planeTypeID, tableData)

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAircraftData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"aircraft"},
			[]string{"iata_code", "aircraft_name", "plane_type_id", "created_at"},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				aircraft := newDataMap[i]
				return []interface{}{
					aircraft.IataCode,
					aircraft.AircraftName,
					aircraft.PlaneTypeId,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into aircraft table")
			return err
		}

		slog.Info("New data inserted into the aircraft table")
	} else {
		slog.Info("No new data to insert into the aircraft table")
	}

	return nil
}

func (s *ServiceJob) StartAPICheckCronJob() {
	c := cron.New()
	slog.Info("Insert city job started")
	c.AddFunc("@every 1m", func() {
		startTime := time.Now()
		err := s.insertNewCities()
		slog.Info("City job finished in: ", time.Since(startTime))
		handleError(err, "Error checking for new cities")
	})
	c.AddFunc("@every 1m", func() {
		startTime := time.Now()
		err := s.insertNewCountries()
		slog.Info("Country job finished in: ", time.Since(startTime))
		handleError(err, "Error checking for new countries")
	})
	c.AddFunc("@every 1m", func() {
		startTime := time.Now()
		err := s.insertNewAirports()
		slog.Info("Airport job finished in: ", time.Since(startTime))
		handleError(err, "Error checking for new countries")
	})
	c.AddFunc("@every 1m", func() {
		startTime := time.Now()
		err := s.insertNewAirplanes()
		slog.Info("Airport job finished in: ", time.Since(startTime))
		handleError(err, "Error checking for new countries")
	})
	c.AddFunc("@every 1m", func() {
		startTime := time.Now()
		err := s.insertNewTax()
		slog.Info("Tax job finished in: ", time.Since(startTime))
		handleError(err, "Error checking for new tax")
	})
	c.AddFunc("@every 1m", func() {
		startTime := time.Now()
		err := s.insertNewAirline()
		slog.Info("Airline job finished in: ", time.Since(startTime))
		handleError(err, "Error checking for new airline")
	})
	c.AddFunc("@every 1m", func() {
		startTime := time.Now()
		err := s.insertNewAircraft()
		slog.Info("Aircraft job finished in: ", time.Since(startTime))
		handleError(err, "Error checking for new aircraft")
	})

	c.Start()
}
