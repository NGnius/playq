package interfaces

import (
  "os"
  "time"
)

type PlayqComponent interface {
  Run (chan int)
  Start (chan int)

}

type PlayqComponentStruct struct {
  ControlChannel chan string
  EventChannel chan string
}

type Monitor struct {
  PlayqComponent
  PlayqComponentStruct
  PlaybackControlChannel chan string
  PlaybackEventChannel chan string
  PlaybackFileChannel chan *os.File
  IsAutostart bool
}

type Controller struct {
  PlayqComponent
  PlayqComponentStruct
  MonitorControlChannel chan string
  MonitorEventChannel chan string
}

type Playback struct {
  PlayqComponent
  PlayqComponentStruct
  FileChannel chan *os.File
  BufferedTime time.Duration
}
