package graphic

import (
	"cloudHelper"
	"lessonType"
	"net/http"
	"strings"
	"utility"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const thumbnailURL = "https://storage.googleapis.com/teraconn_thumbnail/graphic/{id}.{fileType}"

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
		bucketName := utility.MaterialBucketName(ctx)
		url, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)

		if err != nil {
			log.Errorf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		graphics[i].ID = id
		graphics[i].URL = url

		replacedURL := strings.Replace(thumbnailURL, "{id}", id, 1)
		graphics[i].ThumbnailURL = strings.Replace(replacedURL, "{fileType}", graphic.FileType, 1)
	}

	return c.JSON(http.StatusOK, graphics)
}
