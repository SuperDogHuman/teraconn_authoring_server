package lessonMaterial

import (
	"cloud.google.com/go/storage"
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	//"google.golang.org/appengine/memcache"
	"cloudHelper"
	"encoding/json"
	"lessonType"
	"net/http"
)

const bucketName = "teraconn_material"

// GetLessonMaterials is get material of the lesson function.
func Gets(c echo.Context) error {
	// increment view cont in memorycache
	// https://cloud.google.com/appengine/docs/standard/go/memcache/reference
	// https://cloud.google.com/appengine/docs/standard/go/memcache/using?hl=ja

	lessonID := c.Param("id")
	ctx := appengine.NewContext(c.Request())
	filePath := "lesson/" + lessonID + ".json"

	bytes, err := cloudHelper.GetObjectFromGCS(ctx, bucketName, filePath)
	if err != nil {
		log.Errorf(ctx, err.Error())
		if err == storage.ErrObjectNotExist {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	lessonMaterial := new(lessonType.LessonMaterial)
	if err := json.Unmarshal(bytes, lessonMaterial); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lessonMaterial)
}

// PutLessonMaterials is put material of the lesson function.
func Put(c echo.Context) error {
	lessonID := c.Param("id")
	ctx := appengine.NewContext(c.Request())
	lessonMaterial := new(lessonType.LessonMaterial)

	if err := c.Bind(lessonMaterial); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	contents, err := json.Marshal(lessonMaterial)
	if err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	filePath := "lesson/" + lessonID + ".json"
	contentType := "application/json"

	if err := cloudHelper.CreateObjectToGCS(ctx, bucketName, filePath, contentType, contents); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeed")
}
