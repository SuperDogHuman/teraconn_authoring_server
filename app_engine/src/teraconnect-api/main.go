package main

import (
	"graphic"
	"lesson"
	"lessonGraphic"
	"lessonMaterial"
	"lessonPack"
	"lessonVoiceText"
	"net/http"
	"rawVoiceSigning"
	"signedURL"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	e.GET("/raw_voice_signing", rawVoiceSigning.Get) // TODO change request params to include url

	e.GET("/graphics", graphic.Gets)

	e.GET("/lessons", lesson.Gets)
	e.GET("/lessons/:id", lesson.Get)
	e.POST("/lessons", lesson.Create)
	e.PATCH("/lessons/:id", lesson.Update)

	e.GET("/lessons/:id/materials", lessonMaterial.Gets)
	e.POST("/lessons/:id/materials", lessonMaterial.Put)
	e.PUT("/lessons/:id/materials", lessonMaterial.Put) // same function as POST

	e.GET("/lessons/:id/voice_texts", lessonVoiceText.Gets)

	e.GET("/lessons/:id/graphics", lessonGraphic.Gets)
	e.POST("/lessons/:id/graphics", lessonGraphic.Create)

	e.PUT("/lessons/:id/packs", lessonPack.Update)

	e.GET("/signed_urls", signedURL.Gets)

	http.Handle("/", e)
}
