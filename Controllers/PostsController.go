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

type passedData struct {
	Description string `json:"description"`
	Image       string `json:"image"`
}

func Post(c *gin.Context) {

	var body passedData
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
}

func DeletePost(c *gin.Context) {
	deletedPost := Models.Post{}
	err := initializers.DB.First(&deletedPost, c.Param("post_id"))
	currentUserId := c.GetInt("id")
	if err != nil {
		if errors.As(err.Error, &logger.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		} else if deletedPost.UserId != currentUserId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You can not delete posts that are not yours",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post already deleted",
			})
			return
		}
	}
	now := time.Now()
	initializers.DB.Model(&deletedPost).Updates(Models.Post{
		Deleted:   true,
		DeletedAt: &now,
	})
}

func GetPost(c *gin.Context) {
	wantedPost := Models.Post{}
	err := initializers.DB.First(&wantedPost, c.Param("post_id"))
	if err != nil {
		if errors.As(err.Error, &logger.ErrRecordNotFound) {
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
	currentUserId := c.GetInt("id")
	posts := []Models.Post{}
	initializers.DB.Where("user_id=?", currentUserId).Preload("User").Find(&posts)
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
	return
}

//getMyFriendsPostsList

//likePost
//unlikePost
//getPostsLikes
