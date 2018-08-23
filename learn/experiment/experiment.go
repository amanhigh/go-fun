package main

import (
	"github.com/amanhigh/go-fun/util"
	texttospeech "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	// Instantiates a client.
	if ttsClient, err := util.NewGoogleTtsClient(util.GOOGLE_CREDS_FILE, texttospeech.SsmlVoiceGender_FEMALE); err == nil {
		ttsClient.Speak("Hi, My name is Aman. How are you ?")
	}

}
