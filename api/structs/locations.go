package structs

// cities

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

type City struct {
	ID          string     `json:"id"`
	GMT         string     `json:"gmt,omitempty"`
	CityID      int        `json:"city_id,string,omitempty"`
	IataCode    string     `json:"iata_code"`
	CountryISO2 string     `json:"country_iso2"`
	GeonameID   string     `json:"geoname_id,omitempty"`
	Latitude    float64    `json:"latitude,string,omitempty"`
	Longitude   float64    `json:"longitude,string,omitempty"`
	CityName    string     `json:"city_name"`
	Timezone    string     `json:"timezone"`
	CreatedAt   CustomTime `json:"created_at"`
}

type CityApiData struct {
	Pagination Pagination `json:"pagination"`
	Data       []City     `json:"data"`
}

type CityListResponse []City

//countries

type Country struct {
	ID                string     `json:"id"`
	CountryName       string     `json:"country_name"`
	CountryISO2       string     `json:"country_iso2"`
	CountryIso3       string     `json:"country_iso3"`
	CountryIsoNumeric int        `json:"country_iso_numeric,string"`
	Population        int        `json:"population,string"`
	Capital           string     `json:"capital"`
	Continent         string     `json:"continent"`
	CurrencyName      string     `json:"currency_name"`
	CurrencyCode      string     `json:"currency_code"`
	FipsCode          string     `json:"fips_code"`
	PhonePrefix       string     `json:"phone_prefix"`
	CreatedAt         CustomTime `db:"created_at" json:"created_at"`
}

type CountryApiData struct {
	Data []Country `json:"data"`
}

type CountryListResponse []Country

//func (c *CustomGMT) UnmarshalJSON(data []byte) error {
//	var gmtString string
//	if err := json.Unmarshal(data, &gmtString); err != nil {
//		return err
//	}
//	println(gmtString)
//	// Parse the time string into hours (float64)
//	offset, err := parseTimeZoneOffset(gmtString)
//	if err != nil {
//		return err
//	}
//	println(offset)
//
//	c.GMT = &offset
//
//	return nil
//}
//
//func parseTimeZoneOffset(offsetString string) (float64, error) {
//	parts := strings.Split(offsetString, ":")
//	if len(parts) != 2 {
//		return 0, fmt.Errorf("invalid time zone offset format: %s", offsetString)
//	}
//
//	// Parse hours and minutes
//	hours, err := strconv.ParseFloat(parts[0], 64)
//	if err != nil {
//		return 0, err
//	}
//
//	// Convert minutes to fractional hours
//	minutes, err := strconv.ParseFloat(parts[1], 64)
//	if err != nil {
//		return 0, err
//	}
//
//	offset := hours + minutes/60.0
//	return offset, nil
//}

type CustomFloat struct {
	float64
}

//type CustomGMT float64
//
//func (g *CustomGMT) UnmarshalJSON(data []byte) error {
//	var gmtStr string
//	if err := json.Unmarshal(data, &gmtStr); err != nil {
//		return err
//	}
//
//	// Handle the special case for "-9:30" and convert it to a float64 value
//	if gmtStr == "-9:30" {
//		*g = CustomGMT(-9.5)
//	} else {
//		// If it's not the special case, parse the string to float64
//		gmt, err := strconv.ParseFloat(gmtStr, 64)
//		if err != nil {
//			return err
//		}
//		*g = CustomGMT(gmt)
//	}
//	return nil
//}
