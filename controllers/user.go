package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nsallis/gobot/user"
)

func NewUser(context *gin.Context) {
	newUser := user.New()
	newUser.Save()
	context.JSON(200, newUser)
}
