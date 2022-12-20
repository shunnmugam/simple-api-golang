package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	Id          string
	FirstName   string `form:"firstName" json:"firstName" binding:"required"`
	LastName    string `form:"lastName" json:"lastName" binding:"required"`
	Email       string `form:"email" json:"email" binding:"required,email"`
	CreatedDate time.Time
}

// local data
var usersData = []User{}

func main() {
	r := gin.Default()

	/**
	 * Create new user
	 */
	r.POST("user", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindJSON(&user); err == nil {

			// create user
			var newUser = User{
				Id:          "user-" + uuid.New().String(),
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Email:       user.Email,
				CreatedDate: time.Now(),
			}

			usersData = append(usersData, newUser)

			ctx.JSON(http.StatusOK, newUser)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

	})

	/**
	 * list all users
	 */
	r.GET("users", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, usersData)
	})

	/**
	 * Get user by id
	 */
	r.GET("user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		user, _ := findUserById(id)
		if user != nil {
			ctx.JSON(http.StatusOK, user)
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{})

	})

	/**
	 * Update user by id
	 */
	r.PATCH("user/:id", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindJSON(&user); err == nil {

			id := ctx.Param("id")

			oldUser, _ := findUserById(id)

			if oldUser != nil {
				(*oldUser).FirstName = user.FirstName
				(*oldUser).LastName = user.LastName
				(*oldUser).Email = user.Email

				ctx.JSON(http.StatusOK, oldUser)
			} else {
				ctx.JSON(http.StatusNotFound, gin.H{})
			}

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

	})

	/**
	 * delete user by id
	 */
	r.DELETE("user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		user, index := findUserById(id)
		fmt.Println(user)
		if user != nil {
			usersData = append(usersData[:index], usersData[index+1:]...)
			ctx.JSON(http.StatusOK, gin.H{})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{})
		}

	})

	r.Run(":8080")
}

func findUserById(id string) (user *User, index int) {
	for i, value := range usersData {
		if value.Id == id {
			index = i
			// user = &value // value is copy
			user = &usersData[i]
			break
		}
	}

	return
}
