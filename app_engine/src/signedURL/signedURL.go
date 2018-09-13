package signedURL

import (
	"cloudHelper"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const bucketName = "teraconn_material"

// Gets is get signed URLs of files.
func Gets(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	jsonString := c.Request().Header.Get("X-Get-Params")
	var fileRequests []FileRequest
	if err := json.Unmarshal([]byte(jsonString), &fileRequests); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	urlsLength := len(fileRequests)
	urls := make([]string, urlsLength)

	for i, fileRequest := range fileRequests {
		// TODO check user permission
		// TODO check file exists

		filePath := strings.ToLower(fileRequest.Entity) + "/" + fileRequest.ID + "." + fileRequest.Extension
		signedURL, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
		if err != nil {
			log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		}
		urls[i] = signedURL
	}

	return c.JSON(http.StatusOK, GetsResponses{SignedURLs: urls})
}

type FileRequest struct {
	ID        string `json:"id"`
	Entity    string `json:"entity"`
	Extension string `json:"extension"`
}

type GetsResponses struct {
	SignedURLs []string `json:"signed_urls"`
}
