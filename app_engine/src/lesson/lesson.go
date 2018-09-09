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
	"utility"
)

// Gets is get multiple lesson function.
func Gets(c echo.Context) error {
	// TODO add pagination
	return c.JSON(http.StatusOK, "")
}

// Get is get lesson function.
func Get(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	id := c.Param("id")

	ids := []string{id}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson := new(lessonType.Lesson)
	lesson.ID = id
	if err := cloudHelper.FetchEntityFromGCD(ctx, lesson, "Lesson"); err != nil {
		log.Errorf(ctx, err.Error())
		if err == datastore.ErrNoSuchEntity {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	avatar := new(lessonType.Avatar)
	avatar.ID = lesson.AvatarID
	if err := cloudHelper.FetchEntityFromGCD(ctx, avatar, "Avatar"); err != nil {
		log.Errorf(ctx, err.Error())
		if err == datastore.ErrNoSuchEntity {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	lesson.Avatar = *avatar

	return c.JSON(http.StatusOK, lesson)
}

// Create is create lesson function.
func Create(c echo.Context) error {
	id := xid.New().String()
	lesson := new(lessonType.Lesson)
	lesson.ID = id
	lesson.Created = time.Now()

	var err error
	ctx := appengine.NewContext(c.Request())
	if err = c.Bind(lesson); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if _, err = datastore.Put(ctx, key, lesson); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lesson)
}

// Update is update lesson function.
func Update(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())
	id := c.Param("id")

	ids := []string{id}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson := new(lessonType.Lesson)
	lesson.ID = id
	lesson.Updated = time.Now()

	var err error

	if err = cloudHelper.FetchEntityFromGCD(ctx, lesson, "Lesson"); err != nil {
		log.Errorf(ctx, err.Error())
		if err == datastore.ErrNoSuchEntity {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err = c.Bind(lesson); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if _, err = datastore.Put(ctx, key, lesson); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}
