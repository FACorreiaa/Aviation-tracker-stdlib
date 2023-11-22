package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/FACorreiaa/go-ollama/api/structs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"sync"
	"time"
)

type JobRepository struct {
	Conn *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{Conn: db}
}

func (db *JobRepository) StartWorker(wg *sync.WaitGroup, exitChan chan struct{}) error {
	defer wg.Done()

	// Run the worker every week
	//ticker := time.NewTicker(7 * 24 * time.Hour)
	ticker := time.NewTicker(time.Minute)

	defer func() {
		ticker.Stop()
		fmt.Println("Worker ticker stopped.")
	}()

	for {
		select {
		case <-exitChan:
			fmt.Println("Worker received exit signal. Exiting.")
			return nil
		case <-ticker.C:
			fmt.Println("Running weekly migration...")
			err := db.fetchAndInsertNewCities()
			if err != nil {
				fmt.Println("Error during city table weekly migration:", err)
			}
		}
	}
}

func (db *JobRepository) fetchAndInsertNewCities() error {
	// Fetch data from the API
	//apiData, err := fetchAviationStackData("cities", 1000000)
	//	hi, it's me from yesterday
	//	hope you slept well. if the data format is not unknown at least, you can use Time package.
	//	the parser accepts customizable format string, with examples you can look at here:
	//https://pkg.go.dev/time#pkg-constants
	apiData, err := os.ReadFile("./api/cities.json")
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
	existingData, err := db.getExistingCityIDs()

	for _, city := range existingData {
		fmt.Printf("%+v\n", city)
	}

	if err != nil {
		handleError(err, "error getting existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := findNewCityData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {

		if _, err := db.Conn.CopyFrom(
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

// getExistingCityData retrieves existing city data from the database
func (db *JobRepository) getExistingCityIDs() (map[int]struct{}, error) {
	rows, err := db.Conn.Query(context.Background(), "SELECT city_id FROM city")
	if err != nil {
		handleError(err, "Error querying DB")
		return nil, err
	}
	defer rows.Close()

	tableData := make(map[int]struct{})
	for rows.Next() {
		var cityID int
		if err := rows.Scan(&cityID); err != nil {
			handleError(err, "Error scanning IDs")
			return nil, err
		}
		tableData[cityID] = struct{}{}
	}

	return tableData, nil
}

// findNewCityData identifies new city data by comparing the API data with existing data
func findNewCityData(apiData []structs.City, tableData map[int]struct{}) []structs.City {
	var newData []structs.City

	for _, city := range apiData {
		if _, exists := tableData[city.CityID]; !exists {
			newData = append(newData, city)
		}
	}

	return newData
}
