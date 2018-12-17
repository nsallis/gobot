package main

import (
	// "bufio"
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nsallis/gobot/controllers"
	"github.com/nsallis/gobot/user"
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func main() {
	Init()
	routes := gin.Default()

	routes.POST("/api/new_user", controllers.NewUser)

	routes.Run()
	// u := user.New()
	// u.UserHash["email"] = "nsallis@enservio.com"
	// fmt.Println("User sessionID: " + u.SessionID)
	// u.Save()
	// user := user.Load(u.SessionID)
	// fmt.Printf("Loaded user: %v", user)
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter text: ")
	// text, _ := reader.ReadString('\n')
	// fmt.Println(text)
}

func Init() {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log = logrus.New()
	log.Out = os.Stdout
	user.Log = log
	user.DB = db
}
