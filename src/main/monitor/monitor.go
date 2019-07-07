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
  IsAutostart bool
  discardNextNext bool
  discardNextBadfile bool
  isCompleted bool
}

func New(q streamqapi.SoundQueue) Monitor {
  return Monitor{ActiveQueue:q, EventChannel:make(chan string, 1), ControlChannel:make(chan string), IsAutostart:false,discardNextNext:false, discardNextBadfile:false, isCompleted:false}
}

func NewAndRetrieveQueue(qcode string) Monitor {
  streamqapi.InitAPI("http://localhost:5000")
  q, apiErr := streamqapi.NewSoundQueue(qcode)
  if apiErr != nil {
    fmt.Println("Monitor may not have started properly due to API error:")
    fmt.Println(apiErr)
  }
  return New(q)
}

func NewAndCreateQueue() Monitor {
  streamqapi.InitAPI("http://localhost:5000")
  q, apiErr := streamqapi.CreateSoundQueue()
  if apiErr != nil {
    fmt.Println("Monitor may not have started properly due to API error:")
    fmt.Println(apiErr)
  } else {
    fmt.Print("Created new queue with code ")
    fmt.Println(q.Code)
  }
  return New(q)
}

func (m Monitor) Start(end chan int) {
  go m.Run(end)
}

func (m Monitor) Run(end chan int) {
  // start up duties
  m.doStartup()
  monitorLoop: for {
    select {
    case msg := <- m.ControlChannel:
      args := strings.Split(msg, " ")
      switch args[0] {
      case "end", "shutdown":
        // send terminate signal to player
        m.PlaybackControlChannel <- "end"
      case "next", "skip":
        m.controlNext()
        // mostly API-based controls
      case "add":
        m.controlAdd(msg, args)
      case "shuffle":
        m.controlShuffle(msg, args)
        // WIP
      /*case "repeat-all":
        m.controlRepeatAll()
      case "repeat-one":
        m.controlRepeatOne()*/
      case "pause", "play", "toggle":
        // commands which only involve playback
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
        // when player shutdowns, trigger shutdown here too
        m.EventChannel <- "end"
        break monitorLoop
      }

    }
  }
  // fmt.Println("Monitor end")
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
    m.isCompleted = true
    m.PlaybackFileChannel <- dummyFile
    m.PlaybackControlChannel <- "clear"
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

func (m *Monitor) doStartup() {
  // start up duties
  if !m.IsAutostart {
    m.PlaybackControlChannel <- "pause"
  }
  switch m.ActiveQueue.Index {
  case -1:
    m.discardNextNext = true
    m.PlaybackControlChannel <- "next"
    m.playNext()
  case len(m.ActiveQueue.Items):
    fmt.Println("Queue is already complete, audio did not start")
    m.isCompleted = true
  default:
    m.discardNextNext = true
    m.PlaybackControlChannel <- "next"
    m.playCurrent()
  }
  go m.preloadNext()
}

// control methods

func (m *Monitor) controlNext() {
  m.discardNextNext = true // next next event will be caused by this
  m.PlaybackControlChannel <- "next"
  m.playNext()
  go m.preloadNext()
  m.EventChannel <- ""
}

func (m *Monitor) controlAdd(msg string, _ []string) {
  // add file or files
  fileArgs := strings.Split(msg, "\n")
  if len(fileArgs) == 1{
    m.EventChannel <- "Missing add target"
  } else {
    failed := false
    addLoop: for _, elem := range fileArgs[1:] {
      var addErr error
      if strings.Contains(elem, "/") {
        // add by (local) file path
        _, addErr = m.ActiveQueue.AddFilePath(elem)
      } else {
        // add by sound ID
        var s streamqapi.Sound
        s, addErr = streamqapi.NewSound(elem)
        if addErr == nil {
          addErr = m.ActiveQueue.Add(s)
        }
      }
      if addErr != nil {
        fmt.Println(addErr)
        m.EventChannel <- "Failed at "+elem
        failed = true
        break addLoop
      }
    }
    if !failed {
      if m.isCompleted {
        m.isCompleted = false
        m.discardNextNext = true
        m.PlaybackControlChannel <- "next"
        m.playCurrent()
      }
      m.EventChannel <- ""
    }
  }
}

func (m *Monitor) controlShuffle(_ string, args []string) {
  resp := ""
  if len(args) == 1 || args[1] == "toggle" {
    err := m.ActiveQueue.Shuffle(0)
    if err != nil {
      resp = "Shuffle toggle failed"
    }
  } else if args[1] == "off" {
    err := m.ActiveQueue.Shuffle(-1)
    if err != nil {
      resp = "Shuffle disable failed"
    }
  } else if args[1] == "on"{
    err := m.ActiveQueue.Shuffle(1)
    if err != nil {
      resp = "Shuffle enable failed"
    }
  } else {
    resp = "Unrecognised command format"
  }
  m.EventChannel <- resp
}

/* WIP
func (m Monitor) controlRepeatAll(_ string, args []string) {
  resp = ""
  if len(args) == 1 || args[1] == "toggle" {
    err := m.ActiveQueue.RepeatAll()
    if err != nil {
      resp = "Repeat All toggle failed"
    }
  } else if args[1] == "off" {
    err := m.ActiveQueue.RepeatAllDisable()
    if err != nil {
      resp = "Repeat All disable failed"
    }
  } else if args[1] == "on"{
    err := m.ActiveQueue.RepeatAllEnable()
    if err != nil {
      resp = "Repeat All enable failed"
    }
  } else {
    resp = "Unrecognised command format"
  }
  m.EventChannel <- resp
}

func (m Monitor) controlRepeatOne(_ string, args []string) {
  resp = ""
  if len(args) == 1 || args[1] == "toggle" {
    err := m.ActiveQueue.RepeatOne()
    if err != nil {
      resp = "Repeat One toggle failed"
    }
  } else if args[1] == "off" {
    err := m.ActiveQueue.RepeatOneDisable()
    if err != nil {
      resp = "Repeat One disable failed"
    }
  } else if args[1] == "on"{
    err := m.ActiveQueue.RepeatOneEnable()
    if err != nil {
      resp = "Repeat One enable failed"
    }
  } else {
    resp = "Unrecognised command format"
  }
  m.EventChannel <- resp
}*/
