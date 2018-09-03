package util

//
//import (
//	"context"
//	"fmt"
//
//	"cloud.google.com/go/texttospeech/apiv1"
//	"google.golang.org/api/option"
//	texttospeech2 "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
//
//	"io/ioutil"
//)
//
//const (
//	GOOGLE_CREDS_FILE = "/home/aman/Downloads/aman-high-c0ffb96700a4.json"
//	OUTPUT_FILE_MP3   = "/home/aman/Downloads/output.mp3"
//)
//
//type TtsClient interface {
//	Speak(text string) (err error)
//}
//
//type GoogleTtsClient struct {
//	lang            string
//	voice           string
//	credentialsFile string
//	client          *texttospeech.Client
//}
//
///**
//util.NewGoogleTtsClient(util.GOOGLE_CREDS_FILE,"en-US-Wavenet-B")
//*/
//func NewGoogleTtsClient(credFilePath, voice string) (ttsClient TtsClient, err error) {
//	var client *texttospeech.Client
//	if client, err = texttospeech.NewClient(context.Background(), option.WithCredentialsFile(credFilePath)); err == nil {
//		ttsClient = &GoogleTtsClient{
//			lang:   "en-US",
//			voice:  voice,
//			client: client,
//		}
//	}
//	return
//}
//
//func (self *GoogleTtsClient) Speak(text string) (err error) {
//	request := self.newSynthesisRequest(text)
//	var resp *texttospeech2.SynthesizeSpeechResponse
//	if resp, err = self.client.SynthesizeSpeech(context.Background(), &request); err == nil {
//		// The resp's AudioContent is binary.
//		if err = ioutil.WriteFile(OUTPUT_FILE_MP3, resp.AudioContent, DEFAULT_PERM); err == nil {
//			fmt.Printf("Audio content written to file: %v\n", OUTPUT_FILE_MP3)
//		}
//	}
//	return
//}
//
//func (self *GoogleTtsClient) newSynthesisRequest(text string) texttospeech2.SynthesizeSpeechRequest {
//	return texttospeech2.SynthesizeSpeechRequest{
//		// Set the text input to be synthesized.
//		Input: &texttospeech2.SynthesisInput{
//			InputSource: &texttospeech2.SynthesisInput_Text{Text: text},
//		},
//		// Build the gender request, select the language code ("en-US") and the SSML
//		// gender gender ("neutral").
//		Voice: &texttospeech2.VoiceSelectionParams{
//			LanguageCode: self.lang,
//			Name:         self.voice,
//		},
//		// Select the type of audio file you want returned.
//		AudioConfig: &texttospeech2.AudioConfig{
//			AudioEncoding: texttospeech2.AudioEncoding_MP3,
//		},
//	}
//}
