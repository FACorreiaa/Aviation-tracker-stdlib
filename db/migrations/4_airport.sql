CREATE TABLE airport (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       gmt FLOAT DEFAULT '0.00',
                       airport_id INT DEFAULT 0,
                       iata_code varchar(255),
                       city_iata_code varchar(255),
                       icao_code varchar(255),
                       country_iso2 varchar(255),
                       geoname_id INT DEFAULT 0,
                       latitude float8 DEFAULT '0.00',
                       longitude float8 DEFAULT '0.00',
                       airport_name varchar(255),
                       country_name varchar(255),
                       phone_number varchar(255),
                       timezone varchar(255),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW ()
);

