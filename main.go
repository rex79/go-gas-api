package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Fornitore struct {
	ID        int64
	Nome      string
	Indirizzo string
}

func getFornitoriSQL() ([]Fornitore, error) {
	var fornitori []Fornitore

	rows, err := db.Query("select * from fornitori")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var forn Fornitore
		if err := rows.Scan(&forn.ID, &forn.Nome, &forn.Indirizzo); err != nil {
			return nil, fmt.Errorf("getFornitoriSQL %v", err)
		}
		fornitori = append(fornitori, forn)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getFornitoriSQL %v", err)
	}
	return fornitori, nil
}

func getFornitori(c *gin.Context) {
	/*fornitori, err := getFornitoriSQL()
	if err != nil {
		log.Fatal("Errore connessioner SQL: %s", err)
	}
	c.IndentedJSON(http.StatusOK, fornitori)*/

	fornitori, err := getFornitoriSQL()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fornitori found: %v\n", fornitori)
	c.IndentedJSON(http.StatusOK, fornitori)
}

func postFornitore(c *gin.Context) {
	var newFornitore = createFornitore()
	c.IndentedJSON(http.StatusCreated, newFornitore)
}

func main() {
	// DB Configuration
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1",
		DBName: "gas_ordini",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	fornitori, err := getFornitoriSQL()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fornitori found: %v\n", fornitori)

	router := gin.Default()
	router.GET("/fornitori", getFornitori)
	router.POST("/fornitori", postFornitore)
	router.Run("localhost:8080")
}
