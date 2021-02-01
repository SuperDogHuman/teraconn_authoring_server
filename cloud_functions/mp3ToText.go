package mp3ToText

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

// Mp3SpeechToText are triggered creating audio file in cloud function, the file used for speech to text, update voice entity of cloud datastore.
func Mp3SpeechToText(ctx context.Context, e GCSEvent) error {
	if !strings.HasPrefix(e.Name, "voice/") {
		return nil
	}

	uri := fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name)

	if e.TimeCreated != e.Updated {
		return nil // ファイルの更新時は何もしない
	}

	// 最初にファイル名でDatastoreからVoiceをとってきて、IsTextedがtrueなら終了する

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return err
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_MP3,
			SampleRateHertz: 44100,
			LanguageCode:    "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: uri},
			//			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
		return err
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	// 失敗してもisTexedはtrueで更新する

	return nil
}
