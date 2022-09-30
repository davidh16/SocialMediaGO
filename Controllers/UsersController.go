package Controllers

import (
	"SocialMediaGO/Models"
	"SocialMediaGO/initializers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Register(c *gin.Context) {

	var passedData struct {
		Name     string
		Surname  string
		Email    string
		Password string
	}
	if c.Bind(&passedData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read passed data",
		})
		return
	}

	registeredUser := Models.User{}
	registeredUser = Models.User{Name: passedData.Name, Surname: passedData.Surname, Email: passedData.Email, Password: passedData.Password}
	result := initializers.DB.Create(&registeredUser)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "User already exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": registeredUser,
	})
} //treba implementirati validator

func Login(c *gin.Context) {
	var loginData struct {
		Email    string
		Password string
	}

	if c.Bind(&loginData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read passed data",
		})
		return
	}

	currentUser := Models.User{}
	initializers.DB.Where("email = ?", loginData.Email).First(&currentUser)

	if currentUser.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid credentials",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": currentUser.ID,
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*8, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
	return
} //treba implementirati validator

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "test", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func UpdateProfile(c *gin.Context) {
	userId, _ := c.Get("id")
	updateData := Models.User{}
	if c.Bind(&updateData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read passed data",
		})
		return
	}
	updatedUser := Models.User{}
	initializers.DB.First(&updatedUser, userId)
	initializers.DB.Model(&updatedUser).Updates(Models.User{
		Name:     updateData.Name,
		Surname:  updateData.Surname,
		Password: updateData.Password,
	})
	c.JSON(http.StatusOK, gin.H{
		"user": updatedUser,
	})
} //treba implementirati validator

func Deactivate(c *gin.Context) {
	userId, _ := c.Get("id")
	deactivatedUser := Models.User{}
	result := initializers.DB.First(&deactivatedUser, userId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error,
		})
	}

	now := time.Now()

	result2 := initializers.DB.Model(&deactivatedUser).Updates(Models.User{
		Deactivated: true,
		DeletedAt:   &now,
	})
	if result2.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result2.Error,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been deactivated",
	})
	return
}

func Me(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
