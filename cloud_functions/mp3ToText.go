package mp3ToText

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	speech "cloud.google.com/go/speech/apiv1p1beta1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
)

// GCSEvent are events by Cloud Storage.
type GCSEvent struct {
	Name        string    `json:"name"`
	Bucket      string    `json:"bucket"`
	TimeCreated time.Time `json:"timeCreated"`
	Updated     time.Time `json:"updated"`
}

type Voice struct {
	UserID      int64     `json:"userID"`
	LessonID    int64     `json:"lessonID"`
	Speeched    float64   `json:"speeched"`
	DurationSec float64   `json:"durationSec"`
	Text        string    `json:"text"`
	IsTexted    bool      `json:"isTexted"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// Mp3SpeechToText are triggered creating audio file in cloud function, the file used for speech to text, update voice entity of cloud datastore.
func Mp3SpeechToText(ctx context.Context, e GCSEvent) error {
	if !strings.HasPrefix(e.Name, "voice/") {
		return nil
	}

	if e.TimeCreated != e.Updated {
		return nil // ファイルの更新時は何もしない
	}

	fileName := strings.TrimLeft(e.Name, "voice/")
	voiceID, err := strconv.ParseInt(strings.TrimRight(fileName, ".mp3"), 10, 64)
	if err != nil {
		return err
	}

	datastoreClient, err := datastore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		return err
	}

	var voice Voice
	key := datastore.IDKey("Voice", voiceID, nil)
	err = getVoiceFromCloudStorage(ctx, datastoreClient, key, &voice)
	if err != nil {
		return err
	}

	if voice.IsTexted {
		return nil
	}

	if voice.DurationSec < 1.0 {
		// 音声が短すぎるときは、処理済みのフラグだけ立てて終了する
		updateVoiceToCloudStorage(ctx, datastoreClient, key, &voice)
		if err != nil {
			return err
		}

		return nil
	}

	uri := fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name)
	text, err := getSpeechFromURI(ctx, uri)
	if err != nil {
		return err
	}

	voice.Text = text
	err = updateVoiceToCloudStorage(ctx, datastoreClient, key, &voice)
	if err != nil {
		return err
	}

	return nil
}

func getVoiceFromCloudStorage(ctx context.Context, client *datastore.Client, key *datastore.Key, voice *Voice) error {
	if err := client.Get(ctx, key, voice); err != nil {
		return err
	}
	return nil
}

func updateVoiceToCloudStorage(ctx context.Context, client *datastore.Client, key *datastore.Key, voice *Voice) error {
	datastoreClient, err := datastore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	voice.IsTexted = true
	voice.Updated = time.Now()
	_, err = datastoreClient.Put(ctx, key, voice)
	if err != nil {
		return err
	}

	return nil
}

func getSpeechFromURI(ctx context.Context, uri string) (string, error) {
	client, err := speech.NewClient(ctx)
	if err != nil {
		return "", err
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_MP3,
			SampleRateHertz: 44100,
			LanguageCode:    "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: uri},
		},
	})

	if err != nil {
		return "", err
	}

	text := ""
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			if text != "" {
				text = text + "。"
			}
			text = text + alt.Transcript
		}
	}
	return text, nil
}
