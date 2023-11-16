package db

import (
	"context"
	"crypto/md5"
	"embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Init
func Init(connectionURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connectionURL)
	if err != nil {
		return nil, err
	}

	cfg.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		uuid.Register(conn.TypeMap())
		return nil
	}

	return pgxpool.NewWithConfig(context.Background(), cfg)
}

func InitRedis(host, password string, db int) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	}), nil
}

func Migrate(conn *pgxpool.Pool) error {
	//migrate db
	slog.Info("Running migrations")
	ctx := context.Background()
	files, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	slog.Info("Creating migrations table")
	_, err = conn.Exec(ctx, `
		create table if not exists _migrations (
			name text primary key,
			hash text not null,
			created_at timestamp default now()
		);
	`)
	if err != nil {
		return err
	}

	slog.Info("Checking applied migrations")
	rows, _ := conn.Query(ctx, `select name, hash from _migrations order by created_at desc`)
	var name, hash string
	appliedMigrations := make(map[string]string)
	pgx.ForEachRow(rows, []any{&name, &hash}, func() error {
		appliedMigrations[name] = hash
		return nil
	})

	for _, file := range files {
		contents, err := migrationFS.ReadFile("migrations/" + file.Name())
		if err != nil {
			return err
		}

		contentHash := fmt.Sprintf("%x", md5.Sum(contents))

		if prevHash, ok := appliedMigrations[file.Name()]; ok {
			if prevHash != contentHash {
				return fmt.Errorf("hash mismatch for %s", file.Name())
			}

			slog.Info(file.Name() + " already applied")
			continue
		}

		err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, string(contents)); err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, `insert into _migrations (name, hash) values ($1, $2)`, file.Name(), contentHash); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
		slog.Info(file.Name() + " applied")
	}

	slog.Info("Migrations finished")
	return nil
}

//// MigrateRedis performs migrations for Redis
//func MigrateRedis(redisClient *redis.Client) error {
//	ctx := context.Background()
//
//	// Create a set to store applied migrations
//	appliedMigrationsKey := "applied_migrations"
//
//	// Check if the set exists, if not, create it
//	exists, err := redisClient.Exists(ctx, appliedMigrationsKey).Result()
//	if err != nil {
//		return err
//	}
//
//	if exists == 0 {
//		if err := redisClient.SAdd(ctx, appliedMigrationsKey, "initial_migration").Err(); err != nil {
//			return err
//		}
//	}
//
//	// Get the list of applied migrations
//	appliedMigrations, err := redisClient.SMembers(ctx, appliedMigrationsKey).Result()
//	if err != nil {
//		return err
//	}
//
//	// Read migration files
//	migrationFiles, err := ioutil.ReadDir("migrations")
//	if err != nil {
//		return err
//	}
//
//	// Apply pending migrations
//	for _, migrationFile := range migrationFiles {
//		migrationName := migrationFile.Name()
//
//		// Check if the migration has already been applied
//		if contains(appliedMigrations, migrationName) {
//			log.Printf("Migration %s already applied", migrationName)
//			continue
//		}
//
//		// Read the migration file
//		contents, err := ioutil.ReadFile("migrations/" + migrationName)
//		if err != nil {
//			return err
//		}
//
//		// Apply the migration (in this example, we're just hashing the content for simplicity)
//		//contentHash := fmt.Sprintf("%x", md5.Sum(contents))
//
//		// Apply the migration logic, e.g., set a key-value pair in Redis
//		err = applyMigrationLogic(ctx, redisClient, contents)
//		if err != nil {
//			return err
//		}
//
//		// Add the migration to the set of applied migrations
//		if err := redisClient.SAdd(ctx, appliedMigrationsKey, migrationName).Err(); err != nil {
//			return err
//		}
//
//		log.Printf("Migration %s applied", migrationName)
//	}
//
//	log.Println("Migrations finished")
//	return nil
//}
//
//func applyMigrationLogic(ctx context.Context, client *redis.Client, contents []byte) error {
//	// Your migration logic here, e.g., set a key-value pair in Redis
//	key := "example_key"
//	value := string(contents)
//	return client.Set(ctx, key, value, 0).Err()
//}
//
//// Helper function to check if a string is in a slice of strings
//func contains(slice []string, str string) bool {
//	for _, s := range slice {
//		if s == str {
//			return true
//		}
//	}
//	return false
//}

// Small hack to wait for database to start inside docker
func WaitForDB(pgpool *pgxpool.Pool) {
	ctx := context.Background()

	for attempts := 1; ; attempts++ {
		if attempts > 25 {
			break
		}

		if err := pgpool.Ping(ctx); err == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}
}

func WaitForRedis(redis *redis.Client) {
	ctx := context.Background()

	for attempts := 1; ; attempts++ {
		if attempts > 25 {
			break
		}

		if err := redis.Ping(ctx); err == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}
}
