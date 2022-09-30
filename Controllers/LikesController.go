package Controllers

import (
	"SocialMediaGO/Models"
	"SocialMediaGO/initializers"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"net/http"
)

func LikePost(c *gin.Context) {
	likedPost := Models.Post{}
	result := initializers.DB.First(&likedPost, c.Param("post_id"))
	if result.Error != nil || likedPost.Deleted == true {
		if errors.As(result.Error, &logger.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Post not found",
		})
		return
	}

	if likedPost.UserId != c.GetInt("id") {
		friend := Models.User{}
		err := initializers.DB.Where("id=?", c.GetInt("id")).Preload("Friends", "id=?", likedPost.UserId).First(&friend)
		if err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error,
			})
			return
		}

		if len(friend.Friends) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You can not access to posts that are not yours or of users that are your friends",
			})
			return
		}
	}

	currentUser := Models.User{}
	result2 := initializers.DB.First(&currentUser, c.GetInt("id"))
	if result2.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result2.Error,
		})
		return
	}

	result3 := initializers.DB.Model(&likedPost).Association("Likes").Append(&currentUser)
	if result3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result3.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post liked",
	})
	return
}

func UnlikePost(c *gin.Context) {
	unlikedPost := Models.Post{}
	result := initializers.DB.First(&unlikedPost, c.Param("post_id"))
	if result.Error != nil || unlikedPost.Deleted == true {
		if errors.As(result.Error, &logger.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Post not found",
		})
		return
	}

	likes := Models.Post{}
	result2 := initializers.DB.Where("id=?", unlikedPost.ID).Preload("Likes", "id=?", c.GetInt("id")).First(&likes)
	if result2.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result2.Error,
		})
		return
	}

	if len(likes.Likes) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Can not unlike posts that have not been liked",
		})
		return
	}

	currentUser := Models.User{}
	result3 := initializers.DB.First(&currentUser, c.GetInt("id"))
	if result3.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result3.Error,
		})
		return
	}

	result4 := initializers.DB.Model(&unlikedPost).Association("Likes").Delete(currentUser)
	if result4 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result4.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked",
	})
	return
}

func GetPostsLikes(c *gin.Context) {
	post := Models.Post{}
	result := initializers.DB.First(&post, c.Param("post_id"))
	if result.Error != nil || post.Deleted == true {
		if errors.As(result.Error, &logger.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Post not found",
		})
		return
	}

	if post.UserId != c.GetInt("id") {
		friend := Models.User{}
		err := initializers.DB.Where("id=?", c.GetInt("id")).Preload("Friends", "id=?", post.UserId).First(&friend)
		if err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error,
			})
			return
		}

		if len(friend.Friends) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You can not access to posts that are not yours or of users that are your friends",
			})
			return
		}
	}

	var likes Models.Post
	result2 := initializers.DB.Where("id=?", post.ID).Preload("Likes").First(&likes)
	if result2.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result2.Error,
		})
		return
	}
	c.JSON(http.StatusForbidden, gin.H{
		"likes": likes.Likes,
	})
	return
} //treba napraviti paginaciju
