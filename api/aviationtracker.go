package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FACorreiaa/go-ollama/api/structs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*Airline Migration function */

func MigrateAirlineAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM airline").Scan(&count); err != nil {
		fmt.Println("Error querying the table", err)
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertAirlineData(conn); err != nil {
			fmt.Println("Error inserting data", err)
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAirlineData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("airlines")
	if err != nil {
		fmt.Printf("error fetching airline data %v", err)
		return err
	}

	var response structs.AirlineApiData
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("error unmarshaling API response: %v", err)
		return err
	}

	responseData := response.Data

	//Insert data from the json
	if _, err = conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"airline"},
		[]string{"fleet_average_age", "airline_id", "callsign", "hub_code", "iata_code", "icao_code", "country_iso2",
			"date_founded", "iata_prefix_accounting", "airline_name", "country_name", "fleet_size", "status", "type",
			"created_at",
		},
		pgx.CopyFromSlice(len(responseData), func(i int) ([]interface{}, error) {
			createdAt := responseData[i].CreatedAt.Time.Format("2006-01-02 15:04:05-07:00")

			return []interface{}{
				responseData[i].FleetAverageAge,
				responseData[i].AirlineId,
				responseData[i].Callsign,
				responseData[i].HubCode,
				responseData[i].IataCode,
				responseData[i].IcaoCode,
				responseData[i].CountryIso2,
				responseData[i].DateFounded,
				responseData[i].IataPrefixAccounting,
				responseData[i].AirlineName,
				responseData[i].CountryName,
				responseData[i].FleetSize,
				responseData[i].Status,
				responseData[i].Type,
				createdAt,
			}, nil
		}),
	); err != nil {
		log.Printf("error inserting data into airline table: %v", err)
		return err
	}

	slog.Info("Data inserted into the airline table")

	return nil
}

/*Aircraft Migration function */

func MigrateAircraftAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM aircraft").Scan(&count); err != nil {
		fmt.Println("Error querying the table", err)
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertAircraftData(conn); err != nil {
			fmt.Println("Error inserting data", err)
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAircraftData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("aircraft_types")
	if err != nil {
		fmt.Printf("error fetching aircraft data %v", err)
		return err
	}

	var response structs.AircraftApiData
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("error unmarshaling API response: %v", err)
		return err
	}

	responseData := response.Data

	//Insert data from the json
	if _, err = conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"aircraft"},
		[]string{"iata_code", "aircraft_name", "plane_type_id", "created_at"},
		pgx.CopyFromSlice(len(responseData), func(i int) ([]interface{}, error) {
			createdAt := responseData[i].CreatedAt.Time.Format("2006-01-02 15:04:05-07:00")

			return []interface{}{
				responseData[i].IataCode,
				responseData[i].AircraftName,
				responseData[i].PlaneTypeId,
				createdAt,
			}, nil
		}),
	); err != nil {
		log.Printf("error inserting data into aircraft table: %v", err)
		return err
	}

	slog.Info("Data inserted into the aircraft table")

	return nil
}

/*Tax Migration function */

func MigrateTaxAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM tax").Scan(&count); err != nil {
		fmt.Println("Error querying the table", err)
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertTaxData(conn); err != nil {
			fmt.Println("Error inserting data", err)
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertTaxData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("taxes")
	if err != nil {
		fmt.Printf("error fetching tax data %v", err)
		return err
	}

	var response structs.TaxApiData
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("error unmarshaling API response: %v", err)
		return err
	}

	responseData := response.Data

	//Insert data from the json
	if _, err = conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"tax"},
		[]string{"tax_id", "tax_name", "iata_code", "created_at"},
		pgx.CopyFromSlice(len(responseData), func(i int) ([]interface{}, error) {
			createdAt := responseData[i].CreatedAt.Time.Format("2006-01-02 15:04:05-07:00")

			return []interface{}{
				responseData[i].TaxId,
				responseData[i].TaxName,
				responseData[i].IataCode,
				createdAt,
			}, nil
		}),
	); err != nil {
		log.Printf("error inserting data into aircraft table: %v", err)
		return err
	}

	slog.Info("Data inserted into the aircraft table")

	return nil
}

/* Airplane */

func MigrateAirplaneAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM airplane").Scan(&count); err != nil {
		fmt.Println("Error querying the table", err)
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertAirplaneData(conn); err != nil {
			fmt.Println("Error inserting data", err)
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAirplaneData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("airplanes")
	if err != nil {
		fmt.Printf("error fetching airplanes data %v", err)
		return err
	}

	var response structs.AirplaneApiData
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("error unmarshaling API response: %v", err)
		return err
	}

	responseData := response.Data

	//Insert data from the json
	if _, err = conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"airplane"},
		[]string{"iata_type", "airplane_id", "airline_iata_code", "iata_code_long", "iata_code_short",
			"airline_icao_code", "construction_number", "delivery_date", "engines_count", "engines_type",
			"first_flight_date", "icao_code_hex", "line_number", "model_code", "registration_number",
			"test_registration_number", "plane_age", "plane_class", "model_name", "plane_owner", "plane_series",
			"plane_status", "production_line", "registration_date", "rollout_date", "created_at",
		},
		pgx.CopyFromSlice(len(responseData), func(i int) ([]interface{}, error) {
			createdAt := responseData[i].CreatedAt.Time.Format("2006-01-02 15:04:05-07")
			deliveryDate := responseData[i].DeliveryDate.Time
			firstFlightDate := responseData[i].FirstFlightDate.Time
			registrationDate := responseData[i].RegistrationDate.Time
			rolloutDate := responseData[i].RolloutDate.Time

			return []interface{}{
				responseData[i].IataType,
				responseData[i].AirplaneId,
				responseData[i].AirlineIataCode,
				responseData[i].IataCodeLong,
				responseData[i].IataCodeShort,
				responseData[i].AirlineIcaoCode,
				responseData[i].ConstructionNumber,
				deliveryDate,
				responseData[i].EnginesCount,
				responseData[i].EnginesType,
				firstFlightDate,
				responseData[i].IcaoCodeHex,
				responseData[i].LineNumber,
				responseData[i].ModelCode,
				responseData[i].RegistrationNumber,
				responseData[i].TestRegistrationNumber,
				responseData[i].PlaneAge,
				responseData[i].PlaneClass,
				responseData[i].ModelName,
				responseData[i].PlaneOwner,
				responseData[i].PlaneSeries,
				responseData[i].PlaneStatus,
				responseData[i].ProductionLine,
				registrationDate,
				rolloutDate,
				createdAt,
			}, nil
		}),
	); err != nil {
		log.Printf("error inserting data into airplane table: %v", err)
		return err
	}

	slog.Info("Data inserted into the airplane table")

	return nil
}

/* Airports */

func MigrateAirportAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM airport").Scan(&count); err != nil {
		fmt.Println("Error querying the table", err)
		return err
	}

	if count == 0 {
		// No data in the airport table, fetch from the external API
		if err := fetchDataAndInsertAirportData(conn); err != nil {
			fmt.Println("Error inserting data", err)
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAirportData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("airports")
	if err != nil {
		fmt.Printf("error fetching airport data %v", err)
		return err
	}

	var response structs.AirportApiData
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("error unmarshaling API response: %v", err)
		return err
	}

	responseData := response.Data

	//Insert data from the json
	if _, err = conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"airport"},
		[]string{"gmt", "airport_id", "iata_code", "city_iata_code", "icao_code",
			"country_iso2", "geoname_id", "latitude", "longitude", "airport_name",
			"country_name", "phone_number", "timezone", "created_at",
		},
		pgx.CopyFromSlice(len(responseData), func(i int) ([]interface{}, error) {
			createdAt := responseData[i].CreatedAt.Time.Format("2006-01-02 15:04:05-07")

			return []interface{}{
				responseData[i].GMT,
				responseData[i].AirportId,
				responseData[i].IataCode,
				responseData[i].CityIataCode,
				responseData[i].IcaoCode,
				responseData[i].CountryIso2,
				responseData[i].GeonameId,
				responseData[i].Latitude,
				responseData[i].Longitude,
				responseData[i].AirportName,
				responseData[i].CountryName,
				responseData[i].PhoneNumber,
				responseData[i].Timezone,
				createdAt,
			}, nil
		}),
	); err != nil {
		log.Printf("error inserting data into airport table: %v", err)
		return err
	}

	slog.Info("Data inserted into the airport table")

	return nil
}

/* Countries */

func MigrateCountryAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM country").Scan(&count); err != nil {
		fmt.Println("Error querying the table", err)
		return err
	}

	if count == 0 {
		// No data in the country table, fetch from the external API
		if err := fetchDataAndInsertCountryData(conn); err != nil {
			fmt.Println("Error inserting data", err)
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertCountryData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("countries")
	if err != nil {
		fmt.Printf("error fetching country data %v", err)
		return err
	}

	var response structs.CountryApiData
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("error unmarshaling API response: %v", err)
		return err
	}

	responseData := response.Data

	//Insert data from the json
	if _, err = conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"country"},
		[]string{"country_name", "country_iso2", "country_iso3", "country_iso_numeric", "population",
			"capital", "continent", "currency_name", "currency_code", "fips_code",
			"phone_prefix", "created_at",
		},
		pgx.CopyFromSlice(len(responseData), func(i int) ([]interface{}, error) {
			createdAt := responseData[i].CreatedAt.Time.Format("2006-01-02 15:04:05-07")

			return []interface{}{
				responseData[i].CountryName,
				responseData[i].CountryIso2,
				responseData[i].CountryIso3,
				responseData[i].CountryIsoNumeric,
				responseData[i].Population,
				responseData[i].Capital,
				responseData[i].Continent,
				responseData[i].CurrencyName,
				responseData[i].CurrencyCode,
				responseData[i].FipsCode,
				responseData[i].PhonePrefix,
				createdAt,
			}, nil
		}),
	); err != nil {
		log.Printf("error inserting data into country table: %v", err)
		return err
	}

	slog.Info("Data inserted into the country table")

	return nil
}

/* Cities */

func MigrateCityAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM city").Scan(&count); err != nil {
		fmt.Println("Error querying the table", err)
		return err
	}

	if count == 0 {
		// No data in the airport table, fetch from the external API
		if err := fetchDataAndInsertCityData(conn); err != nil {
			fmt.Println("Error inserting data", err)
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertCityData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("cities")
	if err != nil {
		fmt.Printf("error fetching city data %v", err)
		return err
	}

	var response structs.CityApiData
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("error unmarshaling API response: %v", err)
		return err
	}

	responseData := response.Data

	//Insert data from the json
	if _, err = conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"city"},
		[]string{"gmt", "city_id", "iata_code", "country_iso2", "geoname_id",
			"latitude", "longitude", "city_name", "timezone", "created_at",
		},
		pgx.CopyFromSlice(len(responseData), func(i int) ([]interface{}, error) {
			createdAt := responseData[i].CreatedAt.Time.Format("2006-01-02 15:04:05-07")

			return []interface{}{
				responseData[i].GMT,
				responseData[i].CityId,
				responseData[i].IataCode,
				responseData[i].CountryIso2,
				responseData[i].GeonameId,
				responseData[i].Latitude,
				responseData[i].Longitude,
				responseData[i].CityName,
				responseData[i].Timezone,
				createdAt,
			}, nil
		}),
	); err != nil {
		log.Printf("error inserting data into city table: %v", err)
		return err
	}

	slog.Info("Data inserted into the city table")

	return nil
}

/* Fetch data from endpoint */

func fetchAviationStackData(endpoint string, queryParams ...string) ([]byte, error) {
	accessKey := os.Getenv("AVIATION_STACK_API_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("missing API access key")
	}

	baseURL := "http://api.aviationstack.com/v1/"

	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Set the endpoint path
	parsedURL.Path += endpoint

	// Create a new query parameters object
	query := parsedURL.Query()

	// Add the access key parameter
	query.Set("access_key", accessKey)

	// Add additional query parameters
	if len(queryParams) > 0 {
		for _, param := range queryParams {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) == 2 {
				query.Set(parts[0], parts[1])
			}
		}
	}

	parsedURL.RawQuery = query.Encode()

	finalURL := parsedURL.String()

	response, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("something is not ok")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}
