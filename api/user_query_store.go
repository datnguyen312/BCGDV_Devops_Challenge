package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "3x4mpl3"
	dbname   = "postgres"
)

// DB is deliberately global as it should live between requests 
var db *sql.DB

func createTable(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS UserQueries (
		Time timestamp PRIMARY KEY,
		City text
	)`

	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatal("Unable to create table: ", err)
		return err
	}
	fmt.Println("Created UserQueries table")

	return nil
}

type CityQueryStorer interface {
	StoreCityQuery(city string) error
}

type CityQueryRetriever interface {
	RetrieveAllQueriedCities() []UserQuery
}

type PostgresCityStore struct {
	db *sql.DB
}

func InitDatabase() (storer CityQueryStorer, retriever CityQueryRetriever, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Unable to connect to database")
		return
	}
	pcs := &PostgresCityStore{db: db}
	storer = pcs
	retriever = pcs

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err != nil {
			fmt.Println("Unable to ping database")
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		fmt.Println("Unable to ping database")
		return
	}

	fmt.Println("Successfully connected to database")

	err = createTable(pcs.db)
	return
}

func (p *PostgresCityStore) StoreCityQuery(city string) (err error) {
	var result sql.Result
	result, err = p.db.Exec("INSERT INTO UserQueries VALUES ($1, $2)", time.Now(), city)
	if err != nil {
		log.Println("Unable to store user query: ", err)
		return
	}

	var lastInsertId int64
	lastInsertId, err = result.LastInsertId()
	fmt.Println("Database insert succeeded:", lastInsertId)
	return
}

type UserQuery struct {
	Time time.Time
	City string
}

func (p *PostgresCityStore) RetrieveAllQueriedCities() []UserQuery {
	rows, err := p.db.Query("SELECT * FROM UserQueries ORDER BY Time DESC")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var queriedCities []UserQuery
	for rows.Next() {
		userQuery := UserQuery{}
		err := rows.Scan(&userQuery.Time, &userQuery.City)
		if err != nil {
			log.Fatal(err)
		}

		queriedCities = append(queriedCities, userQuery)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return queriedCities
}
