package lessonGraphic

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"cloudHelper"
	"lessonType"
	"net/http"
	"time"
	"utility"
)

// Gets is get lesson graphic.
func Gets(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())
	id := c.Param("id")

	ids := []string{id}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonGraphic := new(lessonType.LessonGraphic)
	lessonGraphic.ID = id // LessonGraphicID is the same as lessonID

	var err error
	if err = cloudHelper.FetchEntityFromGCD(ctx, lessonGraphic, "LessonGraphic"); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Warningf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusNotFound, err.Error())
		}
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var keys []*datastore.Key
	for _, id := range lessonGraphic.GraphicIDs {
		keys = append(keys, datastore.NewKey(ctx, "Graphic", id, 0, nil))
	}

	graphicCount := len(lessonGraphic.GraphicIDs)
	graphics := make([]lessonType.Graphic, graphicCount)
	if err = datastore.GetMulti(ctx, keys, graphics); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	lessonGraphic.Graphics = graphics

	return c.JSON(http.StatusOK, lessonGraphic)
}

// Create is create lesson graphic.
func Create(c echo.Context) error {
	id := c.Param("id")
	ctx := appengine.NewContext(c.Request())

	ids := []string{id}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonGraphic := new(lessonType.LessonGraphic)
	lessonGraphic.ID = id // LessonGraphicID is the same as lessonID
	lessonGraphic.Created = time.Now()

	// TODO check exist entity.

	var err error
	if err = c.Bind(lessonGraphic); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if !utility.IsValidXIDs(lessonGraphic.GraphicIDs) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	key := datastore.NewKey(ctx, "LessonGraphic", id, 0, nil)
	if _, err = datastore.Put(ctx, key, lessonGraphic); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lessonGraphic)
}
