package postgres

import (
	"fmt"
	"github.com/FACorreiaa/go-ollama/internal/logs"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type Postgres struct {
	db *gorm.DB
}

type QueryExecMode uint

const (
	CacheStatement = iota
)

func (m QueryExecMode) value() string {
	switch m {
	case CacheStatement:
		return "cache_statement"
	default:
		return ""
	}
}

// OptionParams string `env:"DB_OPTION_PARAMS" envDefault:""`
type Config struct {
	scheme               string
	host                 string
	port                 string
	username             string
	password             string
	dbName               string
	sslMode              string
	maxConnWaitingTime   time.Duration
	defaultQueryExecMode QueryExecMode
}

func NewConfig(
	scheme string,
	host string,
	port string,
	username string,
	password string,
	dbName string,
	sslMode string,
	maxConnWaitingTime time.Duration,
	defaultQueryExecMode QueryExecMode,
) Config {
	return Config{
		scheme:               scheme,
		host:                 host,
		port:                 port,
		username:             username,
		password:             password,
		dbName:               dbName,
		sslMode:              sslMode,
		maxConnWaitingTime:   maxConnWaitingTime,
		defaultQueryExecMode: defaultQueryExecMode,
	}
}

func (c Config) formatDSN() string {
	password := os.Getenv("DB_PASSWORD")
	fmt.Printf("Password from env: %s\n", password) // Add this line for debugging

	u := url.URL{
		Scheme: c.scheme,
		User:   url.UserPassword(c.username, password),
		Host:   c.host,
		Path:   c.dbName,
	}
	q := make(url.Values)
	q.Add("sslmode", c.sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func newDB(config Config) (*gorm.DB, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(filepath.Join(dir, ".env"))

	dsn := config.formatDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	//need to fix this
	// Ping the database to ensure the connection is established.
	//sqlDB, err := db.DB()
	//if err != nil {
	//	return nil, err
	//}
	//pingCtx, cancel := context.WithTimeout(context.Background(), config.maxConnWaitingTime)
	//defer cancel()
	//if err := sqlDB.PingContext(pingCtx); err != nil {
	//	return nil, err
	//}

	return db, nil
}

func NewPostgres(config Config) *Postgres {
	db, err := newDB(config)
	if err != nil {
		logs.DefaultLogger.WithError(err).Fatal("Error on postgres init")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
	return &Postgres{db: db}
}

func (p *Postgres) GetDB() *gorm.DB {
	return p.db
}
