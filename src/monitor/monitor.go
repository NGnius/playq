package monitor

import (
    "../streamqapi"
    "fmt"
    "os"
    "strings"
)

type Monitor struct {
  EventChannel chan string
  ControlChannel chan string
  PlaybackControlChannel chan string
  PlaybackEventChannel chan string
  PlaybackFileChannel chan *os.File
  ActiveQueue streamqapi.SoundQueue
  discardNextNext bool
  discardNextBadfile bool
}

func New(qcode string) Monitor{
  streamqapi.InitAPI("http://localhost:5000")
  q, apiErr := streamqapi.NewSoundQueue(qcode)
  if apiErr != nil {
    fmt.Println("Monitor may not have started properly due to following API error")
    fmt.Println(apiErr)
  }
  return Monitor{ActiveQueue:q, EventChannel:make(chan string, 1), ControlChannel:make(chan string), discardNextNext:false, discardNextBadfile:false}
}

func (m Monitor) Start(end chan int) {
  go m.Run(end)
}

func (m Monitor) Run(end chan int) {
  m.discardNextNext = true
  switch m.ActiveQueue.Index {
  case -1:
    m.PlaybackControlChannel <- "next"
    m.playNext()
  case len(m.ActiveQueue.Items):
    fmt.Println("Queue is already complete, Monitor did not start")
    m.PlaybackControlChannel <- "end"
    m.EventChannel <- "end"
    end <- 1
    return
  default:
    m.PlaybackControlChannel <- "next"
    m.playCurrent()
  }
  go m.preloadNext()
  monitorLoop: for {
    select {
    case msg := <- m.ControlChannel:
      args := strings.Split(msg, " ")
      switch args[0] {
      case "end", "shutdown":
        m.PlaybackControlChannel <- "end"
      case "next":
        m.discardNextNext = true
        m.PlaybackControlChannel <- "next"
        m.playNext()
        go m.preloadNext()
        m.EventChannel <- ""
      case "pause", "play", "toggle":
        m.PlaybackControlChannel <- args[0]
        m.EventChannel <- ""
      default:
        m.EventChannel <- "Monitor: Bad input"
      }
    case event := <- m.PlaybackEventChannel:
      switch event {
      case "next":
        if m.discardNextNext {
          m.discardNextNext = false
        } else {
          m.playNext()
          go m.preloadNext()
        }
      case "badfile":
        if m.discardNextBadfile {
          m.discardNextBadfile = false
        } else {
          m.PlaybackControlChannel <- "next"
          m.discardNextNext = true
          m.playNext()
          go m.preloadNext()
        }
      case "end":
        m.EventChannel <- "end"
        break monitorLoop
      }

    }
  }
  fmt.Println("Monitor end")
  end <- 0
}

func (m Monitor) play(s streamqapi.Sound) {
  f := s.GetFile()
  m.PlaybackFileChannel <- f
}

func (m *Monitor) playNext() {
  nextSound, nextErr := m.ActiveQueue.GetNext()
  if nextErr != nil {
    dummyFile, _ := os.Open("dummy")
    m.discardNextBadfile = true
    m.PlaybackFileChannel <- dummyFile
    m.PlaybackControlChannel <- "end"
    //fmt.Println("Queue completed")
    return
  }
  m.play(nextSound)
}

func (m Monitor) preloadNext() {
  if m.ActiveQueue.Index+1 >= len(m.ActiveQueue.Items) {
    return
  }
  probableNextSound := m.ActiveQueue.Items[m.ActiveQueue.Index+1]
  probableNextSound.GetFile().Close()
}

func (m Monitor) playCurrent() {
  m.play(m.ActiveQueue.Now())
}
