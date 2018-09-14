package graphic

import (
	"cloudHelper"
	"lessonType"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const bucketName = "teraconn_material"

// Gets is get lesson graphic.
func Gets(c echo.Context) error {
	// TODO pagination.
	ctx := appengine.NewContext(c.Request())

	var graphics []lessonType.Graphic
	query := datastore.NewQuery("Graphic").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &graphics)
	if err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(graphics) == 0 {
		errMessage := "graphics not found"
		log.Warningf(ctx, "%v\n", errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	for i, graphic := range graphics {
		id := keys[i].StringID()
		filePath := "graphic/" + id + "." + graphic.FileType
		fileType := "" // this is unnecessary when GET request
		url, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)

		if err != nil {
			log.Errorf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		log.Infof(ctx, "%v\n", url)
		graphics[i].ID = id
		graphics[i].URL = url
	}

	return c.JSON(http.StatusOK, graphics)
}
