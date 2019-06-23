package monitor

import (
    "../streamqapi"
    //"http"
    //"json"
)

type Monitor struct {
  ControlChannel chan string
  API streamqapi.API
}

func New(qcode string) Monitor{
  return Monitor{API: streamqapi.NewAPI(qcode)}
}
