package lesson

import (
	"github.com/labstack/echo"
	"github.com/rs/xid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"cloudHelper"
	"lessonType"
	"net/http"
	"time"
)

// Gets is get multiple lesson function.
func Gets(c echo.Context) error {
	// TODO add pagination
	return c.JSON(http.StatusOK, "")
}

// Get is get lesson function.
func Get(c echo.Context) error {
	id := c.Param("id")
	lesson := new(lessonType.Lesson)
	lesson.ID = id

	ctx := appengine.NewContext(c.Request())
	if err := cloudHelper.FetchObjectFromGCD(ctx, lesson); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Errorf(ctx, err.Error())
			return c.JSON(http.StatusNotFound, err.Error())
		}
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}

// Create is create lesson function.
func Create(c echo.Context) error {
	id := xid.New().String()
	lesson := new(lessonType.Lesson)
	lesson.ID = id
	lesson.Created = time.Now()

	ctx := appengine.NewContext(c.Request())
	if err := cloudHelper.PutObjectToGCD(ctx, c, lesson); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lesson)
}

// Update is update lesson function.
func Update(c echo.Context) error {
	id := c.Param("id")
	lesson := new(lessonType.Lesson)
	lesson.ID = id

	ctx := appengine.NewContext(c.Request())
	if err := cloudHelper.FetchObjectFromGCD(ctx, lesson); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Errorf(ctx, err.Error())
			return c.JSON(http.StatusNotFound, err.Error())
		}
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := cloudHelper.PutObjectToGCD(ctx, c, lesson); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}
