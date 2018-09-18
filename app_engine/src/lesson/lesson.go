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
	lessonKey := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if err = datastore.Get(ctx, lessonKey, lesson); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Warningf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusNotFound, err.Error())
		}
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	lesson.ID = id // for json field

	avatar := new(lessonType.Avatar)
	avatarKey := datastore.NewKey(ctx, "Avatar", lesson.AvatarID, 0, nil)
	if err = datastore.Get(ctx, avatarKey, avatar); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Warningf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusNotFound, err.Error())
		}
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	avatar.ID = lesson.AvatarID
	lesson.Avatar = *avatar

	var graphicKeys []*datastore.Key
	for _, id := range lesson.GraphicIDs {
		graphicKeys = append(graphicKeys, datastore.NewKey(ctx, "Graphic", id, 0, nil))
	}
	graphics := make([]lessonType.Graphic, len(lesson.GraphicIDs))
	if err = datastore.GetMulti(ctx, graphicKeys, graphics); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	for i, id := range lesson.GraphicIDs {
		graphics[i].ID = id
	}

	lesson.Graphics = graphics

	return c.JSON(http.StatusOK, lesson)
}

// Create is create lesson function.
func Create(c echo.Context) error {
	id := xid.New().String()
	lesson := new(lessonType.Lesson)
	lesson.Created = time.Now()

	var err error
	ctx := appengine.NewContext(c.Request())
	if err = c.Bind(lesson); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if _, err = datastore.Put(ctx, key, lesson); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	lesson.ID = id // for json response

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

		updateLesson := f.(map[string]interface{})
		mutable := reflect.ValueOf(lesson).Elem()
		for key, lessonField := range updateLesson {
			structKey := strings.Title(key)
			switch v := lessonField.(type) {
			case []interface{}:
				array := make([]string, len(v)) // TODO support another types. reflect.TypeOf(v[0])
				mutable.FieldByName(structKey).Set(reflect.ValueOf(array))
				for i := range v {
					mutable.FieldByName(structKey).Index(i).Set(reflect.ValueOf(v[i]))
				}
			default:
				if structKey == "ViewCount" || structKey == "Version" {
					intValue := int64(v.(float64))
					mutable.FieldByName(structKey).SetInt(intValue)
				} else {
					mutable.FieldByName(structKey).Set(reflect.ValueOf(v))
				}
			}
		}

		_, err := datastore.Put(ctx, lessonKey, lesson)
		return err
	}, nil)

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Warningf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusNotFound, err.Error())
		}
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if !utility.IsValidXIDs(lesson.GraphicIDs) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	return c.JSON(http.StatusOK, lesson)
}
