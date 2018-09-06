package lessonMaterial

import (
	"cloud.google.com/go/storage"
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	//"google.golang.org/appengine/memcache"
	"cloudHelper"
	"encoding/json"
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

	lessonMaterial := new(LessonMaterial)
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
	lessonMaterial := new(LessonMaterial)

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

type LessonMaterial struct {
	DurationSec float64          `json:"durationSec"`
	Timelines   []LessonTimeline `json:"timelines"`
	Pose        LessonAvatarPose `json:"poseKey"`
}

type LessonTimeline struct {
	TimeSec  float64                   `json:"timeSec"`
	Text     LessonMaterialText        `json:"text"`
	Voice    LessonMaterialVoice       `json:"voice"`
	Graphic  []LessonMaterialGraphic   `json:"graphics"`
	SPAction LessonAvatarSpecialAction `json:"spAction"`
}

type LessonMaterialText struct {
	DurationSec     float64 `json:"durationSec"`
	Body            string  `json:"body"`
	HorizontalAlign string  `json:"horizontalAlign"`
	VerticalAlign   string  `json:"verticalAlign"`
	SizeVW          uint8   `json:"sizeVW"`
	BodyColor       string  `json:"bodyColor"`
	BorderColor     string  `json:"borderColor"`
}

type LessonMaterialVoice struct {
	ID          string  `json:"id"`
	DurationSec float64 `json:"durationSec"`
}

type LessonMaterialGraphic struct {
	ID              string `json:"id"`
	Action          string `json:"action"`
	SizePct         uint8  `json:"sizePct"`
	HorizontalAlign string `json:"horizontalAlign"`
	VerticalAlign   string `json:"verticalAlign"`
}

type LessonAvatarSpecialAction struct {
	Action         string `json:"action"`
	FaceExpression string `json:"faceExpression"`
}

type LessonAvatarPose struct {
	LeftHands      []LessonRotation `json:"leftHands"`
	RightHands     []LessonRotation `json:"rightHands"`
	LeftElbows     []LessonRotation `json:"leftElbows"`
	RightElbows    []LessonRotation `json:"rightElbows"`
	LeftShoulders  []LessonRotation `json:"leftShoulders"`
	RightShoulders []LessonRotation `json:"rightShoulders"`
	Necks          []LessonRotation `json:"necks"`
	CoreBodies     []LessonPosition `json:"coreBodies"`
}

type LessonRotation struct {
	Rot  []float32 `json:"rot"`
	Time float32   `json:"time"`
}

type LessonPosition struct {
	Rot  []float32 `json:"pos"`
	Time float32   `json:"time"`
}
