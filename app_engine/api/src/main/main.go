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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/appengine"
)

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"https://authoring.teraconnect.org",
			//			"https://teraconnect-authoring-development-dot-teraconnect-209509.appspot.com",
			//			"http://localhost:1234",
		},
	}))

	e.GET("/lessons", lesson.Gets)
	e.GET("/lessons/:id", lesson.Get)

	auth := e.Group("", middleware.JWT([]byte("secret")))

	auth.GET("/avatars", avatar.Gets)

	auth.GET("/graphics", graphic.Gets)

	auth.POST("/lessons", lesson.Create)
	auth.PATCH("/lessons/:id", lesson.Update)
	auth.DELETE("/lessons/:id", lesson.Destroy)

	auth.GET("/lessons/:id/materials", lessonMaterial.Gets)
	auth.POST("/lessons/:id/materials", lessonMaterial.Put)
	auth.PUT("/lessons/:id/materials", lessonMaterial.Put) // same function as POST

	auth.GET("/lessons/:id/voice_texts", lessonVoiceText.Gets)

	auth.PUT("/lessons/:id/packs", lessonPack.Update)

	auth.GET("/storage_objects", storageObject.Gets)
	auth.POST("/storage_objects", storageObject.Posts)

	auth.POST("/raw_voices", rawVoice.Post)

	http.Handle("/", e)
	appengine.Main()
}
