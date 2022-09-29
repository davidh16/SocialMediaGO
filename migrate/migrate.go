package migrate

import (
	"SocialMediaGO/Models"
	"SocialMediaGO/initializers"
	"fmt"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func Migrate() {
	err := initializers.DB.AutoMigrate(&Models.Post{}, &Models.User{})
	if err != nil {
		fmt.Println(err)
	}

}
