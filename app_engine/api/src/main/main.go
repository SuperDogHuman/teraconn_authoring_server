package main

import (
	"avatar"
	"graphic"
	"lesson"
	"lessonMaterial"
	"lessonPack"
	"lessonVoiceText"
	"net/http"
	"rawVoice"
	"storageObject"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
      "https://authoring.teraconnect.org",
//      "https://teraconnect-authoring-development-dot-teraconnect-209509.appspot.com",
//      "http://localhost:1234",
    },
  }))

	e.GET("/avatars", avatar.Gets)

	e.GET("/graphics", graphic.Gets)

	e.GET("/lessons", lesson.Gets)
	e.GET("/lessons/:id", lesson.Get)
	e.POST("/lessons", lesson.Create)
	e.PATCH("/lessons/:id", lesson.Update)
	e.DELETE("/lessons/:id", lesson.Destroy)

	e.GET("/lessons/:id/materials", lessonMaterial.Gets)
	e.POST("/lessons/:id/materials", lessonMaterial.Put)
	e.PUT("/lessons/:id/materials", lessonMaterial.Put) // same function as POST

	e.GET("/lessons/:id/voice_texts", lessonVoiceText.Gets)

	e.PUT("/lessons/:id/packs", lessonPack.Update)

	e.GET("/storage_objects", storageObject.Gets)
	e.POST("/storage_objects", storageObject.Posts)

	e.POST("/raw_voices", rawVoice.Post)

	http.Handle("/", e)
}
