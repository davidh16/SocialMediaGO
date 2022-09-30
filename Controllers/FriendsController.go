package Controllers

import (
	"SocialMediaGO/Models"
	"SocialMediaGO/initializers"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"net/http"
)

func AddFriend(c *gin.Context) {
	currentUser := Models.User{}
	addedFriend := Models.User{}
	result := initializers.DB.Find(&currentUser, c.GetInt("id"))
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": result.Error,
		})
		return
	}

	result2 := initializers.DB.First(&addedFriend, c.Param("friend_id"))
	if result2.Error != nil {
		if errors.As(result2.Error, &logger.ErrRecordNotFound) || addedFriend.Deactivated == true {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User does not exist",
			})
			return
		} else if addedFriend.Deactivated == true {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User does not exist",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unexpected error",
			})
			return
		}
	}

	result3 := initializers.DB.Model(&currentUser).Association("Friends").Append(&addedFriend)
	if result3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unexpected error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You added " + addedFriend.Name + " as your friend",
	})

	return
}

func DeleteFriend(c *gin.Context) {
	currentUser := Models.User{}
	deletedFriend := Models.User{}
	result := initializers.DB.Find(&currentUser, c.GetInt("id"))
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": result.Error,
		})
		return
	}

	result2 := initializers.DB.First(&deletedFriend, c.Param("friend_id"))
	if result2.Error != nil {
		if errors.As(result2.Error, &logger.ErrRecordNotFound) || deletedFriend.Deactivated == true {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User does not exist",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unexpected error",
			})
			return
		}
	}

	result3 := initializers.DB.Model(&currentUser).Association("Friends").Delete(deletedFriend)
	if result3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unexpected error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You deleted " + deletedFriend.Name + " from your friend list",
	})

	return
}

func GetFriend(c *gin.Context) {
	friend := Models.User{}
	result := initializers.DB.Where("id=?", c.GetInt("id")).Preload("Friends", "id=?", c.Param("friend_id")).First(&friend)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	} else if len(friend.Friends) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Friend not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"friend": friend.Friends[0],
	})
	return
}

func GetFriendsList(c *gin.Context) {
	friends := Models.User{}
	result := initializers.DB.Where("id=?", c.GetInt("id")).Preload("Friends").First(&friends)
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
	c.JSON(http.StatusOK, gin.H{
		"friends": friends.Friends,
	})
	return
} //fali paginacija
