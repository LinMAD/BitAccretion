package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

const (
	resourcePath = "/resources/sound"
	alarmPath    = "alarm"
	voicePath    = "voice"
)

var wd string

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	wd, _ = os.Getwd()
}

// play given streamer and rate of it
func play(s beep.Streamer, r beep.SampleRate) {
	_ = speaker.Init(r, r.N(time.Second/10))

	playing := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(playing)
	})))

	<-playing
}

// execRandomAlarm sound file
func execRandomAlarm(dir string) {
	files, filesErr := ioutil.ReadDir(dir)
	if filesErr != nil {
		return
	}

	name := files[rand.Intn(len(files))].Name()
	af, _ := os.Open(path.Join(dir, name))
	stream, format, _ := wav.Decode(af)

	play(stream, format.SampleRate)
}

// execVoice matching voice name if exist
func execVoice(dir, reqName string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	var voiceFileName string
	for _, f := range files {
		if strings.Contains(f.Name(), reqName) {
			voiceFileName = f.Name()
			break
		}
	}

	if voiceFileName == "" {
		return
	}

	audioFile, _ := os.Open(path.Join(dir, voiceFileName))
	stream, format, _ := wav.Decode(audioFile)
	play(stream, format.SampleRate)
}

// PlayAlert for given name
func PlayAlert(name string) {
	if wd == "" {
		return
	}

	execRandomAlarm(path.Join(resourcePath, alarmPath))
	execVoice(path.Join(resourcePath, voicePath), name)
}
