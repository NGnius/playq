package playback

import (
    //"fmt"
    "github.com/faiface/beep"
    "github.com/faiface/beep/speaker"
    "github.com/faiface/beep/mp3"
    "github.com/faiface/beep/wav"
    "github.com/faiface/beep/flac"
    "github.com/faiface/beep/vorbis"
    "time"
    "os"
)

type Playback struct {
  currentStream beep.Streamer
  currentFormat beep.SampleRate
  FileChannel chan *os.File
  EventChannel chan string
  ControlChannel chan string
  BufferedTime time.Duration
}

func New() Playback{
  return Playback{FileChannel:make(chan *os.File, 1), EventChannel:make(chan string, 2), ControlChannel:make(chan string), BufferedTime:time.Second/100}
}

func (p Playback) Start(end chan int) {
  go p.Run(end)
}

func (p Playback) Run(end chan int) {
  isSpeakerLocked := false
  isSpeakerInited := false
  songDone := make(chan bool, 1)
  playLoop: for {
    select {
    case <- songDone:
      p.EventChannel <- "next"
      nextStream, nextFormat, nextErr := decodeAudioFile( <- p.FileChannel)
      if nextErr == nil {
        p.currentStream = nextStream
        if !isSpeakerInited /*|| nextFormat.SampleRate != p.currentFormat*/ {
          speaker.Init(nextFormat.SampleRate, nextFormat.SampleRate.N(p.BufferedTime))
          p.currentFormat = nextFormat.SampleRate
          isSpeakerInited = true
        } else {
          if isSpeakerLocked { speaker.Unlock() }
          speaker.Clear()
        }
        speaker.Play(beep.Seq(nextStream, beep.Callback(func(){songDone <- true})))
        if isSpeakerLocked { speaker.Lock() }
      } else {
        p.EventChannel <- "badfile"
      }
    case msg := <- p.ControlChannel:
      switch msg {
      case "next":
        songDone <- false
      case "toggle":
        if isSpeakerInited {
          if isSpeakerLocked {
            speaker.Unlock()
          } else {
            speaker.Lock()
          }
        }
        isSpeakerLocked = !isSpeakerLocked
      case "play":
        if isSpeakerLocked {
          if isSpeakerInited {
            speaker.Unlock()
          }
          isSpeakerLocked = false
        }
      case "pause":
        if !isSpeakerLocked {
          if isSpeakerInited {
            speaker.Lock()
          }
          isSpeakerLocked = true
        }
      case "end", "stop", "close":
        if isSpeakerLocked && isSpeakerInited { speaker.Unlock() }
        speaker.Close()
        p.EventChannel <- "end"
        break playLoop
      case "clear":
        if isSpeakerLocked && isSpeakerInited { speaker.Unlock() }
        speaker.Clear()
        if isSpeakerLocked && isSpeakerInited { speaker.Lock() }
      }
    }
  }
  // fmt.Println("Playback end")
  end <- 0
}

func decodeAudioFile(f *os.File) (beep.Streamer, beep.Format, error) {
  // try all known audio formats, return when no error occurs
  /*var b []byte = []byte{0,}
  _, readErr := f.Read(b)
  if readErr != nil {
    fmt.Println("Decoding: File may be corrupted?")
    fmt.Println(b)
  }
  f.Seek(0,0)*/
  // mp3
  streamer, format, decodeErr := mp3.Decode(f)
  if decodeErr == nil {
    //fmt.Println("Decoded as mp3")
    return streamer, format, nil
  }
  // wav
  streamer, format, decodeErr = wav.Decode(f)
  if decodeErr == nil {
    //fmt.Println("Decoded as wav")
    return streamer, format, nil
  }
  // flac
  streamer, format, decodeErr = flac.Decode(f)
  if decodeErr == nil {
    //fmt.Println("Decoded as flac")
    return streamer, format, nil
  }
  // vorbis
  streamer, format, decodeErr = vorbis.Decode(f)
  if decodeErr == nil {
    //fmt.Println("Decoded as vorbis")
    return streamer, format, nil
  }
  return streamer, format, decodeErr
}
