package structs

//cities

type City struct {
	ID          string     `json:"id"`
	GMT         float64    `json:"gmt,string"`
	CityId      int        `json:"city_id,string"`
	IataCode    string     `json:"iata_code"`
	CountryIso2 string     `json:"country_iso2"`
	GeonameId   float64    `json:"geoname_id,string"`
	Latitude    float64    `json:"latitude,string"`
	Longitude   float64    `json:"longitude,string"`
	CityName    string     `json:"city_name"`
	Timezone    string     `json:"timezone"`
	CreatedAt   CustomTime `db:"created_at" json:"created_at"`
}

type CityListResponse []City

//countries

type Country struct {
	ID                string     `json:"id"`
	CountryName       string     `json:"country_name"`
	CountryIso2       string     `json:"country_iso2"`
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

type CountryListResponse []Country

type CountryApiData struct {
	Data []Country `json:"data"`
}

type CityApiData struct {
	Data []City `json:"data"`
}
