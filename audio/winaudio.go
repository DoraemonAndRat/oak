package audio

import (
	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"errors"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

var (
	loadedWavs = make(map[string]Audio, 0)
)

// We alias the winaudio package's interface here
// so game files don't need to import winaudio
type Audio winaudio.Audio

func InitWinAudio() {
	winaudio.InitWinAudio()
}

func GetActiveWavChannel(frequency, freqRand int, fileNames ...string) (chan int, error) {
	var err error

	sounds := make([]Audio, len(fileNames))
	for i, f := range fileNames {
		sounds[i], err = GetWav(f)
		if err != nil {
			return nil, err
		}
	}

	soundCh := make(chan int)
	go func() {
		for {
			delay := time.Duration(rand.Intn(freqRand) + frequency)
			<-time.After(delay * time.Millisecond)
			// Every once in a while, after some delay,
			// we play an audio that slipped through the
			// above routine.
			a := <-soundCh
			sounds[a].Play()
		}
	}()
	return soundCh, nil
}

func GetWavChannel(frequency, freqRand int, fileNames ...string) (chan int, error) {

	soundCh, err := GetActiveWavChannel(frequency, freqRand, fileNames...)
	if err != nil {
		return soundCh, err
	}
	// This routine serves to steal almost every
	// attempt to play audio
	go func() {
		for {
			<-soundCh
		}
	}()
	return soundCh, nil
}

func GetWav(fileName string) (Audio, error) {
	if _, ok := loadedWavs[fileName]; !ok {
		return nil, errors.New("File not loaded")
	}
	return loadedWavs[fileName], nil
}

func LoadWav(directory, fileName string) (Audio, error) {
	dlog.Verb("Loading", directory, fileName)
	if _, ok := loadedWavs[fileName]; !ok {
		buffer, err := winaudio.LoadWav(filepath.Join(directory, fileName))
		if err != nil {
			return buffer, err
		}
		loadedWavs[fileName] = buffer
	}
	return loadedWavs[fileName], nil
}

func BatchLoad(baseFolder string) error {

	files, _ := ioutil.ReadDir(baseFolder)

	for _, file := range files {
		if !file.IsDir() {
			n := file.Name()
			dlog.Verb(n)
			switch strings.ToLower(n[len(n)-4:]) {
			case ".wav":
				dlog.Verb("loading file ", n)
				LoadWav(baseFolder, n)
			default:
				dlog.Error("Unsupported file ending for batchLoad: ", n)
			}
		}
	}
	dlog.Verb("Loading complete")
	return nil
}
