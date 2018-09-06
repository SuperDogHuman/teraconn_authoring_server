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
)

// Gets is get lesson graphic.
func Gets(c echo.Context) error {
	lessonGraphic := new(lessonType.LessonGraphic)
	lessonGraphic.ID = c.Param("id") // LessonGraphicID is the same as lessonID

	var err error
	ctx := appengine.NewContext(c.Request())
	if err = cloudHelper.FetchEntityFromGCD(ctx, lessonGraphic, "LessonGraphic"); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		if err == datastore.ErrNoSuchEntity {
			return c.JSON(http.StatusNotFound, err.Error())
		}
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
	lessonGraphic := new(lessonType.LessonGraphic)
	lessonGraphic.ID = c.Param("id") // LessonGraphicID is the same as lessonID
	lessonGraphic.Created = time.Now()

	// TODO check exist entity.

	ctx := appengine.NewContext(c.Request())
	if err := cloudHelper.CreateEntityToGCD(ctx, c, lessonGraphic, "LessonGraphic"); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lessonGraphic)
}
