package main

import (
	"fmt"
	"net/http"
	"weight-loss-challenge/internal/database"

	"github.com/gin-gonic/gin"
)

func (app *application) getMyInfo(c *gin.Context) {
	currentUser := app.GetUserFromContext(c)
	if currentUser == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

func (app *application) logWeight(c *gin.Context) {
	var weight database.Weight
	if err := c.ShouldBindJSON(&weight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	weight.UserId = user.Id

	err := app.models.Weight.Insert(&weight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log weight"})
		return
	}

	user.CurrentWeightGr = weight.Weight

	err = app.models.Users.UpdateUserWeights(user)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user weights"})
		return
	}

	c.JSON(http.StatusCreated, user)

}

func (app *application) getAllWeight(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	weights, err := app.models.Weight.GetWeightLogByUser(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get user's weight logs"})
		return
	}
	c.JSON(http.StatusCreated, weights)
}

func (app *application) getAllUser(c *gin.Context) {
	currentUser := app.GetUserFromContext(c)
	if currentUser == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	usersInfoList, err := app.models.Users.GetAllUserInfo()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get all user info"})
		return
	}

	c.JSON(http.StatusOK, usersInfoList)
}
