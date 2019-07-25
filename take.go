package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type takeOrderRequestBody struct {
	Status string `json:"status"`
}

func takeOrder(c *gin.Context) {
	var takeOrderRequestBody takeOrderRequestBody
	c.BindJSON(&takeOrderRequestBody)

	if takeOrderRequestBody.Status != "TAKEN" {
		c.JSON(400, gin.H{"error": "Invalid status."})
	}

	rows, err := db.Query("SELECT status FROM `order` WHERE id = ?", c.Param("id"))

	if err != nil {
		log.Panic(err.Error())
	}

	if !rows.Next() {
		c.JSON(404, gin.H{"error": "Order not found."})
	}

	var status string
	rows.Scan(&status)

	if status != "UNASSIGNED" {
		c.JSON(400, gin.H{"error": "Order is already taken."})
	}

	update, err := db.Prepare("UPDATE `order` SET status = 'TAKEN' WHERE id = ?")

	if err != nil {
		log.Panic(err.Error())
	}

	defer update.Close()

	c.JSON(200, gin.H{"status": "SUCCESS"})
}
