package main

import (
	"net/http"
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

var users []User
var nextID = 1
var loginAttempts = make(map[string]int)
var jwtKey = []byte("my_secret_key")

func register(c *gin.Context) {
	var input User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	for _, u := range users {
		if u.Username == input.Username {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User đã tồn tại"})
			return
		}
	}

	input.ID = nextID
	nextID++
	users = append(users, input)
	c.JSON(http.StatusOK, gin.H{
		"message": "Tạo thông tin thành công",
	})
}

func login(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	//kiem tra spam
	if loginAttempts[input.Username] >= 5 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"message": "Spam quá nhiều lần, thử lại sau",
		})
		return
	}

	for _, u := range users {
		if u.Username == input.Username && u.Password == input.Password {

			loginAttempts[input.Username] = 0
			c.JSON(http.StatusOK, gin.H{
				"message": "Đăng nhập thành công",
			})
			return
		}
	}

	loginAttempts[input.Username]++
	c.JSON(http.StatusUnauthorized, gin.H{
		"message":  "Đăng nhập thất bại",
		"attempts": loginAttempts[input.Username],
	})
}

// get users
func getUsers(c *gin.Context) {
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Không có thông tin users",
		})
		return
	}
	c.JSON(http.StatusOK, users)
}


func getUserID(c *gin.Context) {
	idParam := c.Param("id")

	id, _ := strconv.Atoi(idParam)
	for _, u := range users {
		if u.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"message": "Đã tìm thấy id",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Không tìm thấy id",
	})
}

func filterUsersByAge(c *gin.Context) {
	ageParam := c.Query("age")

	if ageParam == "" {
		c.JSON(400, gin.H{"error": "Chưa nhập tuổi"})
		return
	}
	age, err := strconv.Atoi(ageParam)
	if err != nil {
		c.JSON(400, gin.H{"err": "Age không hợp lệ"})
		return
	}

	var result []User
	for _, u := range users {
		if u.Age == age {
			result = append(result, u)
		}
	}
	c.JSON(200, result)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	for i, u := range users {
		if u.ID == idInt {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Đã xóa user",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Không tìm thấy user",
	})
}

func main() {
	r := gin.Default()
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/users", getUsers)
	r.GET("/user/:id", getUserID)
	r.GET("/users/filter", filterUsersByAge)
	r.DELETE("/user/:id", deleteUser)
	r.Run(":8080")
}
