package utils

import (
	"fmt"
	"os"
)

func buildDSNMigration(user string, password string, host string, port string, dbname string, sslmode string) string {
	return fmt.Sprintf(
		"cockroachdb://%s:%s@%s:%s/%s?sslmode=%s&x-migrations-table=schema_migrations",
		user, password, host, port, dbname, sslmode,
	)
}
func buildDSN(user string, password string, host string, port string, dbname string, sslmode string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}
func GetDSN(migrate bool) string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if !migrate {
		return buildDSN(user, password, host, port, dbname, sslmode)
	}
	return buildDSNMigration(user, password, host, port, dbname, sslmode)
}
