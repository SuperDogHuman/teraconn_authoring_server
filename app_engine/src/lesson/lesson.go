package lesson

import (
	"context"
	"reflect"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"bytes"
	"encoding/json"
	"io"
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

	var err error

	lesson := new(lessonType.Lesson)
	lesson.ID = id
	lessonKey := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if err = datastore.Get(ctx, lessonKey, lesson); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		if err == datastore.ErrNoSuchEntity {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	avatar := new(lessonType.Avatar)
	avatar.ID = lesson.AvatarID
	avatarKey := datastore.NewKey(ctx, "Avatar", avatar.ID, 0, nil)
	if err = datastore.Get(ctx, avatarKey, avatar); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		if err == datastore.ErrNoSuchEntity {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	lesson.Avatar = *avatar

	var graphicKeys []*datastore.Key
	for _, id := range lesson.GraphicIDs {
		graphicKeys = append(graphicKeys, datastore.NewKey(ctx, "Graphic", id, 0, nil))
	}
	graphicCount := len(lesson.GraphicIDs)
	graphics := make([]lessonType.Graphic, graphicCount)
	if err = datastore.GetMulti(ctx, graphicKeys, graphics); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	lesson.Graphics = graphics

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
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
	if _, err = datastore.Put(ctx, key, lesson); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
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

	buf := new(bytes.Buffer)
	io.Copy(buf, c.Request().Body)
	var f interface{}
	if err := json.Unmarshal(buf.Bytes(), &f); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	lesson := new(lessonType.Lesson)
	lesson.Updated = time.Now()
	lessonKey := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := datastore.Get(ctx, lessonKey, lesson); err != nil {
			return err
		}

		newLesson := f.(map[string]interface{})
		mutable := reflect.ValueOf(lesson).Elem()
		for k, v := range newLesson {
			structKey := strings.Title(k)
			mutable.FieldByName(structKey).Set(reflect.ValueOf(v))
		}

		_, err := datastore.Put(ctx, lessonKey, lesson)
		return err
	}, nil)

	if err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}
