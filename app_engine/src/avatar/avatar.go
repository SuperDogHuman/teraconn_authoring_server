package avatar

import (
	"lessonType"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const bucketName = "teraconn_material"
const thumbnailURL = "https://storage.googleapis.com/teraconn_thumbnail/avatar/{id}.png"

// Gets is get lesson avatar.
func Gets(c echo.Context) error {
	// TODO pagination.
	ctx := appengine.NewContext(c.Request())

	var avatars []lessonType.Avatar
	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(avatars) == 0 {
		errMessage := "avatars not found"
		log.Warningf(ctx, "%v\n", errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	for i, key := range keys {
		id := key.StringID()
		avatars[i].ID = id
		avatars[i].ThumbnailURL = strings.Replace(thumbnailURL, "{id}", id, 1)
	}

	return c.JSON(http.StatusOK, avatars)
}
