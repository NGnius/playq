package streamqapi

type API struct {
  QueueCode string
}

func NewAPI(code string) API {
  return API{QueueCode:code}
}
