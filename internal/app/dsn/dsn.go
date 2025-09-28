package dsn

import "fmt"

func FromEnv(dbHost string, dbPort int, dbUser string, dbPass string, dbName string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
}
