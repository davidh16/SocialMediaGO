package Controllers

import (
	"SocialMediaGO/Models"
	"SocialMediaGO/initializers"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"net/http"
	"time"
)

type PostInput struct {
	Description string `json:"description"`
	Image       string `json:"image"`
}

func Post(c *gin.Context) {

	body := PostInput{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read passed data",
		})
		return
	}

	currentUserId := c.GetInt("id")
	if &currentUserId == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed",
		})
		return
	}

	newPost := Models.Post{
		Description: body.Description,
		Image:       body.Image,
		UserId:      currentUserId,
	}

	result := initializers.DB.Create(&newPost)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "New post could not be created",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": newPost,
	})
	return
} //treba implementirati validator, treba omogućiti uploadanje slike

func DeletePost(c *gin.Context) {
	deletedPost := Models.Post{}
	result := initializers.DB.First(&deletedPost, c.Param("post_id"))
	if result.Error != nil {
		if errors.As(result.Error, &logger.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		} else if deletedPost.UserId != c.GetInt("id") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You can not delete posts that are not yours",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post already deleted",
			})
			return
		}
	}

	now := time.Now()

	result2 := initializers.DB.Model(&deletedPost).Updates(Models.Post{
		Deleted:   true,
		DeletedAt: &now,
	})
	if result2.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result2.Error,
		})
		return
	}
}

func GetPost(c *gin.Context) {
	wantedPost := Models.Post{}
	result := initializers.DB.First(&wantedPost, c.Param("post_id"))
	if result.Error != nil {
		if errors.As(result.Error, &logger.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		}
	}

	if wantedPost.Deleted == true {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": wantedPost,
	})
	return
}

func GetMyPostsList(c *gin.Context) {
	posts := []Models.Post{}

	result := initializers.DB.Where("user_id=?", c.GetInt("id")).Preload("User").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
	return
} //potrebno je napraviti paginaciju

func GetMyFriendsPostsList(c *gin.Context) {
	friends := Models.User{}

	result := initializers.DB.Where("id=?", c.GetInt("id")).Preload("Friends", "id=?", c.Param("friend_id")).First(&friends)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	} else if len(friends.Friends) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Friend not found",
		})
		return
	}

	posts := []Models.Post{}

	result2 := initializers.DB.Where("user_id=?", c.Param("friend_id")).Preload("User", "id=?", c.Param("friend_id")).Find(&posts)
	if result2.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result2.Error,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"friend posts": posts,
	})
	return
} // treba napraviti paginaciju
