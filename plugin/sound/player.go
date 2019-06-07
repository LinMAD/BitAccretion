package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

const (
	resourcePath = "/resources/sound"
	voicePath    = resourcePath + "/voice/"
	alarmPath    = resourcePath + "/alarm/"
	sTag         = "SOUND"
)

// SendSoundAlert executes sound stream in speakers where runs program
func SendSoundAlert() {
	rand.Seed(time.Now().UTC().UnixNano())

	wd, err := os.Getwd()
	if err != nil {
		log.Printf(sTag + ": Errpr -> Could not retrieve working directory")

		return
	}

	playAlert(wd + alarmPath)
	playVoice(wd + voicePath)

	return
}

// playVoice random voice file
func playVoice(basePath string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Printf("%s: Error -> %v", sTag, err.Error())

		return
	}

	randVoice := files[rand.Intn(len(files))].Name()
	voiceFile, _ := os.Open(basePath + randVoice)
	voiceStream, format, _ := wav.Decode(voiceFile)

	play(voiceStream, format.SampleRate)
}

// playAlert random alert file
func playAlert(basePath string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Printf("%s: Error -> %v", sTag, err.Error())

		return
	}

	randAlarm := files[rand.Intn(len(files))].Name()
	alarmFile, _ := os.Open(basePath + randAlarm)
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
