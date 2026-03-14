package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func register(c *gin.Context) {
	var input User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp"})
		return
	}

	for _, u := range users {
		if u.Username == input.Username {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User đã tồn tại"})
			return
		}
	}
	users = append(users, input)

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng ký thành công",
	})
}

func login(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	for _, u := range users {
		if u.Username == input.Username && u.Password == input.Password {
			c.JSON(http.StatusOK, gin.H{
				"message": "Đăng nhập thành công",
			})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Sai tài khoản hoặc mật khẩu",
	})
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy thông tin thành công",
	})
}

func deleteUser(c *gin.Context) {
	name := c.Param("name")
	for i, u := range users {
		if u.Username == name {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "đã xóa user",
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
	r.DELETE("/users/:name", deleteUser)
	r.Run(":8080")
}

//