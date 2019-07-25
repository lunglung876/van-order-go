package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

type CreateOrderRequestBody struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type GoogleMapsAPIResponse struct {
	Rows []Row `json:"rows"`
}

type Row struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Distance Distance `json:"distance"`
	Status   string   `json:"status"`
}

type Distance struct {
	Value int `json:"value"`
}

func main() {
	router = gin.Default()
	initializeRoutes()
	router.Run()
}

func initializeRoutes() {
	router.POST("/orders", createOrder)
}

func createOrder(c *gin.Context) {
	var createOrderRequestBody CreateOrderRequestBody
	c.BindJSON(&createOrderRequestBody)

	if error := createOrderRequestBody.validate(); error != "" {
		err := map[string]string{"error": error}
		c.JSON(400, err)
	}

	distance, error := getDistance(createOrderRequestBody.Origin, createOrderRequestBody.Destination)

	if error != nil {
		err := map[string]string{"error": error.Error()}
		c.JSON(400, err)
	}

	db := getDatabaseConnection()

	db, err := sql.Open("mysql", "root:root@tcp(mysql)/van")

	if err != nil {
		log.Panic(err.Error())
	}

	defer db.Close()

	insert, err := db.Prepare("INSERT INTO `order` (distance, status, origin_latitude, origin_longitude, destination_latitude, destination_longitude) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Panic(err.Error())
	}

	defer insert.Close()

	_, err = insert.Exec(distance, "UNASSIGNED", createOrderRequestBody.Origin[0], createOrderRequestBody.Origin[1], createOrderRequestBody.Destination[0], createOrderRequestBody.Destination[0])

	if err != nil {
		log.Panic(err.Error())
	}
}

func (c *CreateOrderRequestBody) validate() string {
	if c.Origin[0] == "" || c.Origin[1] == "" {
		return "The origin field is required."
	}

	if c.Destination[0] == "" || c.Destination[1] == "" {
		return "The destination field is required."
	}

	longtitudeRegex, _ := regexp.Compile(`^(\+|-)?(?:180(?:(?:\.0{1,6})?)|(?:[0-9]|[1-9][0-9]|1[0-7][0-9])(?:(?:\.[0-9]{1,6})?))$`)
	latitudeRegex, _ := regexp.Compile(`^(\+|-)?(?:90(?:(?:\.0{1,6})?)|(?:[0-9]|[1-8][0-9])(?:(?:\.[0-9]{1,6})?))$`)

	if !latitudeRegex.MatchString(c.Origin[0]) {
		return "The origin latitude is invalid."
	}

	if !latitudeRegex.MatchString(c.Destination[0]) {
		return "The destination latitude is invalid."
	}

	if !longtitudeRegex.MatchString(c.Origin[1]) {
		return "The origin longtitude is invalid."
	}

	if !longtitudeRegex.MatchString(c.Destination[1]) {
		return "The destination longtitude is invalid."
	}

	return ""
}

func getDatabaseConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:root@mysql/van")

	if err != nil {
		log.Panic(err.Error())
	}

	return db
}

func getDistance(origin []string, destination []string) (int, error) {
	requestURL, _ := url.Parse("https://maps.googleapis.com/maps/api/distancematrix/json")
	q := requestURL.Query()
	q.Add("origins", strings.Join(origin, ","))
	q.Add("destinations", strings.Join(destination, ","))
	q.Add("key", os.Getenv("GOOGLE_MAPS_API_KEY"))
	requestURL.RawQuery = q.Encode()

	var responseData GoogleMapsAPIResponse
	response, _ := http.Get(requestURL.String())

	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return 0, err
	}

	if responseData.Rows[0].Elements[0].Status != "OK" {
		return 0, fmt.Errorf("Cannot calculate distance")
	}

	return responseData.Rows[0].Elements[0].Distance.Value, nil
}
