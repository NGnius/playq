package main

import (
    "fmt"
    "flag"
    "./controller"
    "./monitor"
    "./playback"
    //TODO: "./interfaces"
    "time"
)

var ABOUT_STRING string = `playq v0.0.1 for streamq
Developed by NGnius (Graham)`

// command line flag arguments
var autostart bool
var queueCode string
var queueCreate bool
var bufferTime time.Duration
var apiBase string

func init() {
  const (
    usageAutostart = "Automatically start audio playback. Otherwise, a 'play' command must be sent to start playback."
    usageQueueCode = "Queue code to use for playback. This is used for adding and getting songs to play."
    usageQueueCodeShort = "queue (shorthand)"
    usageQueueCreate = "Create a new queue instead of using a pre-existing one. The queue code will be chosen by the server, not '--queue'."
    usageBufferTime = "Audio time to buffer. Shorter times reduce errors but increase lag."
    usageApiBase = "Base url address to API endpoint. Connection method (http:// or https://) must be included."
  )
  flag.BoolVar(&autostart, "auto", false, usageAutostart)
  flag.StringVar(&queueCode, "queue", "A", usageQueueCode)
  flag.StringVar(&queueCode, "q", "A", usageQueueCodeShort)
  flag.BoolVar(&queueCreate, "create", false, usageQueueCreate)
  flag.DurationVar(&bufferTime, "buffer", time.Second/100, usageBufferTime)
  flag.StringVar(&apiBase, "url", "http://localhost:5000", usageApiBase)
}

func main() {
  var cli controller.CommandLineInterface
  var mon monitor.Monitor
  var player playback.Playback

  flag.Parse()
  fmt.Println(ABOUT_STRING)
  // create CLI
  cli = controller.NewCLI()
  cli_complete := make(chan int)

  // create monitor
  if queueCreate {
    mon = monitor.NewAndCreateQueue(apiBase)
  } else {
    mon = monitor.NewAndRetrieveQueue(queueCode, apiBase)
  }
  cli.MonitorControlChannel = mon.ControlChannel
  cli.MonitorEventChannel = mon.EventChannel
  monitor_complete := make(chan int)

  // create player
  player = playback.New()
  mon.PlaybackControlChannel = player.ControlChannel
  mon.PlaybackFileChannel = player.FileChannel
  mon.PlaybackEventChannel = player.EventChannel
  playback_complete := make(chan int)

  // command line actions
  mon.IsAutostart = autostart
  player.BufferedTime = bufferTime

  // start components
  cli.Start(cli_complete)
  player.Start(playback_complete)
  mon.Start(monitor_complete)

  // wait until all components complete
  //fmt.Println(<- cli_complete)
  <- cli_complete
  //fmt.Println(<- monitor_complete)
  <- monitor_complete
  //fmt.Println(<- playback_complete)
  <- playback_complete
}
