package avatar

import (
	"cloudHelper"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const bucketName = "teraconn_material"

// Get is get signing of raw voice files function.
func Get(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	filePath := "avatar/" + c.Param("id") + ".zip"
	contentType := "application/zip"

	bytes, err := cloudHelper.GetObjectFromGCS(ctx, bucketName, filePath)

	if err != nil {
		log.Errorf(ctx, err.Error())
		if err == storage.ErrObjectNotExist {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, contentType, bytes)
}

// URLGet is get signed URL of avatar file.
func URLGet(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	// TODO check avatar exist
	filePath := "avatar/" + c.Param("id") + ".vrm"

	signedURL, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")

	if err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, avatarSign{SignedURL: signedURL})
}

type avatarSign struct {
	SignedURL string `json:"signed_url"`
}
