package controller

import (
  "os"
  "bufio"
  "fmt"
  "strings"
  "path/filepath"
)

var cli_prompt string = "cli@playq $ "

var Commands []string = []string{"pause", "play", "next", "toggle", "shuffle"}

type CommandLineInterface struct {
  MonitorControlChannel chan string
  MonitorEventChannel chan string
}

func NewCLI() CommandLineInterface{
  return CommandLineInterface{}
}

func (c CommandLineInterface) Start(end chan int) {
  go c.Run(end)
}

func (c CommandLineInterface) Run(end chan int) {
  cliReader := bufio.NewReader(os.Stdin)
  commandLoop: for {
    fmt.Print(cli_prompt)
    text, _ := cliReader.ReadString('\n')
    text = text[:len(text)-1]
    args := strings.Split(text, " ")
    select {
    case msg := <- c.MonitorEventChannel:
      switch msg {
      case "end":
        break commandLoop
      }
    default:
      switch args[0] {
      // commands with cli functionality
      case "end","shutdown":
          c.MonitorControlChannel <- "end"
          for {
            resp := <- c.MonitorEventChannel
            if resp == "end"{ break }
          }
          break commandLoop
      case "prompt":
        if len(args) == 2 {
          cli_prompt = args[1]+" "
        } else {
          fmt.Println("Usage: prompt NEW_PROMPT")
        }
      case "add":
        msg, ok := buildAddMessage(text)
        if !ok {
          fmt.Println("Invalid directory/file path")
        } else {
          c.sendAndPrintResp(msg)
        }
      case "help":
        fmt.Println("Commands: "+strings.Join(Commands, ", "))
        fmt.Println("Extra commands: add, end, prompt, help")
      // commands with monitor-only functionality
      case func(candidate string) string {
        // if args[0] in Commands, it's a match
        for i,_ := range Commands {
          if Commands[i] == candidate {return candidate}
        }
        return ""
        }(args[0]):
        if args[0] != "" {
          c.sendAndPrintResp(text)
        }
      default:
        fmt.Println("Invalid command '"+args[0]+"'")
        fmt.Println("Use help for a list of commands.")
      }
    }
  }
  // fmt.Println("Controller end")
  end <- 0
}

func (c CommandLineInterface) sendAndPrintResp(msg string){
  c.MonitorControlChannel <- msg
  resp := <- c.MonitorEventChannel
  if resp != "" {
    fmt.Println(resp)
  }
}

func buildAddMessage(command string) (string, bool) {
  msg := "add \n"
  path := strings.Trim(command[4:], " ")
  pathInfo, statErr := os.Stat(path)
  if statErr != nil {
    // assume path is actually sound code
    return msg+path, true
  }
  if pathInfo.IsDir() {
    filepaths := getFilesInDirRecurse(path)
    for _, elem := range filepaths {
      msg += elem+"\n"
    }
    msg = msg[:len(msg)-1]
    return msg, true
  } else {
    msg += path
    return msg, true
  }
}

func getFilesInDirRecurse(dirpath string) []string {
  // recursively find files
  var files []string
  dir, dirErr := os.Open(dirpath)
  if dirErr != nil {
    return files
  }
  paths, readErr := dir.Readdirnames(0)
  if readErr != nil {
    return files
  }
  for _, path := range paths {
    fullpath := filepath.Join(dirpath, path)
    info, statErr := os.Stat(fullpath)
    if statErr == nil { // ignore errored pathes
      if info.IsDir() {
        files = append(files, getFilesInDirRecurse(fullpath)...)
      } else {
        files = append(files, fullpath)
      }
    } else {
      fmt.Println(statErr)
    }
  }
  return files
}
