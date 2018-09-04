package lessonType

import "time"

type Lesson struct {
	ID           string    `json:"id"`
	AvatarID     string    `json:"-"`
	Avatar       Avatar    `json:"avatar" datastore:"-"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DurationSec  float64   `json:"durationSec"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	GraphicIDs   []string  `json:"graphicIDs"`
	ViewCount    int64     `json:"viewCount"`
	Version      int64     `json:"version"`
	IsPublic     bool      `json:"isPublic"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type Avatar struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userID"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	Name         string    `json:"name"`
	Version      int64     `json:"version"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type LessonAuthor struct {
	ID       string `json:"id"`
	LessonID string `json:"lessonID"`
	UserID   string `json:"userID"`
	Role     string `json:"role"`
}

type Graphic struct {
	ID     string   `json:"id"`
	UserID string   `json:"userID"`
	TeamID []string `json:"ownerIDs"`
}

type LessonGraphic struct {
	ID         string    `json:"id"`
	LessonID   string    `json:"lessonID"`
	GraphicIDs []string  `json:"graphicIDs"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type LessonVoiceText struct {
	FileID      string `json:"fileID"`
	LessonID    string `json:"lessonId"`
	Text        string `json:"text"`
	IsTexted    bool   `json:"isTexted"`
	IsConverted bool   `json:"isConverted"`
}

/* The following structs is for json.Unmarshall */

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
