package main

import (
	"avatar"
	"lesson"
	"lessonMaterial"
	"net/http"
	"rawVoiceSigning"
	"voiceText"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	e.GET("/raw_voice_signing", rawVoiceSigning.Get) // TODO change request params to include url
	e.GET("/voice_text/:lesson_id", voiceText.Gets)

	e.GET("/lessons", lesson.Gets)
	e.GET("/lessons/:id", lesson.Get)
	e.POST("/lessons", lesson.Create)
	e.PUT("/lessons/:id", lesson.Update)

	e.GET("/lessons/:id/materials", lessonMaterial.Gets)
	e.POST("/lessons/:id/materials", lessonMaterial.Put)
	e.PUT("/lessons/:id/materials", lessonMaterial.Put) // same function as POST

	e.GET("/avatars/:id", avatar.Get)

	http.Handle("/", e)
}
