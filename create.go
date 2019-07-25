package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type createOrderRequestBody struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type googleMapsAPIResponse struct {
	Rows []row `json:"rows"`
}

type row struct {
	Elements []element `json:"elements"`
}

type element struct {
	Distance distance `json:"distance"`
	Status   string   `json:"status"`
}

type distance struct {
	Value int `json:"value"`
}

func createOrder(c *gin.Context) {
	var createOrderRequestBody createOrderRequestBody
	c.BindJSON(&createOrderRequestBody)

	if error := createOrderRequestBody.validate(); error != "" {
		c.JSON(400, gin.H{"error": error})
	}

	distance, error := getDistance(createOrderRequestBody.Origin, createOrderRequestBody.Destination)

	if error != nil {
		c.JSON(400, gin.H{"error": error.Error()})
	}

	insert, err := db.Prepare("INSERT INTO `order` (distance, status, origin_latitude, origin_longitude, destination_latitude, destination_longitude) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Panic(err.Error())
	}

	result, err := insert.Exec(distance, "UNASSIGNED", createOrderRequestBody.Origin[0], createOrderRequestBody.Origin[1], createOrderRequestBody.Destination[0], createOrderRequestBody.Destination[0])

	defer insert.Close()

	if err != nil {
		log.Panic(err.Error())
	}

	id, _ := result.LastInsertId()

	c.JSON(200, gin.H{
		"id":       id,
		"distance": distance,
		"status":   "UNASSIGNED",
	})
}

func (c *createOrderRequestBody) validate() string {
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

func getDistance(origin []string, destination []string) (int, error) {
	requestURL, _ := url.Parse("https://maps.googleapis.com/maps/api/distancematrix/json")
	q := requestURL.Query()
	q.Add("origins", strings.Join(origin, ","))
	q.Add("destinations", strings.Join(destination, ","))
	q.Add("key", os.Getenv("GOOGLE_MAPS_API_KEY"))
	requestURL.RawQuery = q.Encode()

	var responseData googleMapsAPIResponse
	response, _ := http.Get(requestURL.String())

	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return 0, err
	}

	if responseData.Rows[0].Elements[0].Status != "OK" {
		return 0, fmt.Errorf("Cannot calculate distance")
	}

	return responseData.Rows[0].Elements[0].Distance.Value, nil
}
