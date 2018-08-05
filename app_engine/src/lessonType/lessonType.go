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

type VoiceText struct {
	FileID      string `json:"fileID"`
	LessonID    string `json:"lessonId"`
	Text        string `json:"text"`
	IsTexted    bool   `json:"isTexted"`
	IsConverted bool   `json:"isConverted"`
}

/* The following structs is for json.Unmarshall */

type LessonMaterial struct {
	DurationSec float64            `json:"durationSec"`
	Timelines   []LessonTimeline   `json:"timelines"`
	Poses       []LessonAvatarPose `json:"poses"`
	Published   time.Time          `json:"published"`
	Updated     time.Time          `json:"updated"`
}

type LessonTimeline struct {
	TimeSec  float64                   `json:"timeSec"`
	Text     LessonText                `json:"text"`
	Voice    LessonVoice               `json:"voice"`
	Graphic  []LessonGraphic           `json:"graphics"`
	SPAction LessonAvatarSpecialAction `json:"spAction"`
}

type LessonText struct {
	DurationSec float64 `json:"durationSec"`
	Text        string  `json:"text"`
	Position    string  `json:"position"`
	Style       string  `json:"style"`
	Size        uint8   `json:"size"`
	Color       string  `json:"color"`
}

type LessonVoice struct {
	FileID      string  `json:"fileID"`
	DurationSec float64 `json:"durationSec"`
}

type LessonGraphic struct {
	GraphicID string `json:"graphicID"`
	Action    string `json:"action"`
	WidthPx   uint16 `json:"widthPx"`
	HeightPx  uint16 `json:"heightPx"`
	Position  string `json:"position"`
}

type LessonAvatarSpecialAction struct {
	Action           string `json:"action"`
	FacialExpression string `json:"facialExpression"`
}

type LessonAvatarPose struct {
	TimeSec            float64      `json:"timeSec"`
	LeftElbowAngle     float32      `json:"leftElbowAngle"`
	RightElbowAngle    float32      `json:"rightElbowAngle"`
	LeftShoulderAngle  float32      `json:"leftShoulderAngle"`
	RightShoulderAngle float32      `json:"rightShoulderAngle"`
	LookAt             LessonVector `json:"lookAt"`
	CoreBody           LessonVector `json:"coreBody"`
}

type LessonVector struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}
