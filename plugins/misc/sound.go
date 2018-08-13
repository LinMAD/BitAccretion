package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"log"
	"os"
	"time"
)

const (
	resourcePath = "/resources/sound/"
	sTag         = "SOUND"
)

// SendSoundAlert executes sound stream in speakers where runs program
func SendSoundAlert() {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf(sTag + ": Could not retrieve working directory")

		return
	}

	playAlert(wd)
	playAttention(wd)

	return
}

// playAttention sound
func playAttention(basePath string) {
	attentionFile, _ := os.Open(basePath + resourcePath + "f_confirmcivilstatus_1_cut_spkr.wav")
	alarmStream, format, _ := wav.Decode(attentionFile)

	play(alarmStream, format.SampleRate)
}

// playAlert sound 3 times
func playAlert(basePath string) {
	alarmFile, _ := os.Open(basePath + resourcePath + "manhack_alert_pass1_cut.wav")
	alarmStream, format, _ := wav.Decode(alarmFile)

	play(alarmStream, format.SampleRate)
}

// play given streamer and rate of it
func play(s beep.Streamer, rate beep.SampleRate) {
	speaker.Init(rate, rate.N(time.Second/10))

	playing := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(playing)
	})))

	<-playing
}
