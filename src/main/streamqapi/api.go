package streamqapi

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "errors"
  "io"
)

var cachedir string = "./temp"//os.TempDir()
var api API

// api endpoints
var fileGetUrl string = "/api/file/"
var fileUploadPostUrl string = "/api/file/"
var soundGetUrl string = "/api/sound/"
var soundNewUrl string = "/api/sound/new"
var soundMetadataUrl string ="/api/sound/?refresh=true"
var queueGetUrl string = "/api/queue/"
var queueNewUrl string = "/api/queue/new"
var queueNextUrl string = "/api/queue/next"
var queuePreviousUrl string = "/api/queue/previous"
var queueNowUrl string = "/api/queue/now"
var queueAddUrl string = "/api/queue/add"

// start of API object

type API struct {
  Base string
  //client http.Client
}

func (this API) _getBytes(url string) ([]byte, error) {
  resp, respErr := http.Get(url)
  if respErr != nil {
    return nil, respErr
  }
  if resp.StatusCode != 200 {
    return nil, errors.New("Get "+url+" status: "+resp.Status) //"API Response "+resp.Status+" for url "+url)
  }
  defer resp.Body.Close()
  body, ioErr := ioutil.ReadAll(resp.Body)
  if ioErr != nil {
    return nil, ioErr
  }
  return body, nil
}

func (this API) getFile(code string) ([]byte, error) {
  return this._getBytes(this.Base+fileGetUrl+"?code="+code)
}

func (this API) uploadFile(code string, file io.Reader) (error) {
  url := this.Base+fileUploadPostUrl+"?code="+code
  resp, respErr := http.Post(url, "audio/mp3", file)
  if respErr != nil {
    return respErr
  }
  if resp.StatusCode != 200 {
    return errors.New("Post "+url+" status: "+resp.Status)
  }
  return nil
}

func (this API) getSound(code string) (Sound, error) {
  var s Sound
  body, respErr := this._getBytes(this.Base+soundGetUrl+"?code="+code)
  if respErr != nil {
    return s, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &s)
  if unmarsharlErr != nil {
    return s, unmarsharlErr
  }
  return s, nil
}

func (this API) newSound() (Sound, error) {
  var s Sound
  body, respErr := this._getBytes(this.Base+soundNewUrl)
  if respErr != nil {
    return s, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &s)
  if unmarsharlErr != nil {
    return s, unmarsharlErr
  }
  return s, nil
}

func (this API) refreshSound(code string) (Sound, error) {
  var s Sound
  body, respErr := this._getBytes(this.Base+soundMetadataUrl+"&code="+code)
  if respErr != nil {
    return s, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &s)
  if unmarsharlErr != nil {
    return s, unmarsharlErr
  }
  return s, nil
}

func (this API) getQueue(code string) (SoundQueue, error) {
  var q SoundQueue
  body, respErr := this._getBytes(this.Base+queueGetUrl+"?code="+code)
  if respErr != nil {
    return q, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &q)
  if unmarsharlErr != nil {
    return q, unmarsharlErr
  }
  return q, nil
}

func (this API) newQueue() (SoundQueue, error) {
  var q SoundQueue
  body, respErr := this._getBytes(this.Base+queueNewUrl)
  if respErr != nil {
    return q, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &q)
  if unmarsharlErr != nil {
    return q, unmarsharlErr
  }
  return q, nil
}

func (this API) queueNext(code string) (Sound, error) {
  var s Sound
  body, respErr := this._getBytes(this.Base+queueNextUrl+"?code="+code)
  if respErr != nil {
    return s, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &s)
  if unmarsharlErr != nil {
    return s, unmarsharlErr
  }
  return s, nil
}

func (this API) queuePrevious(code string) (Sound, error) {
  var s Sound
  body, respErr := this._getBytes(this.Base+queuePreviousUrl+"?code="+code)
  if respErr != nil {
    return s, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &s)
  if unmarsharlErr != nil {
    return s, unmarsharlErr
  }
  return s, nil
}

func (this API) queueNow(code string) (Sound, error) {
  var s Sound
  body, respErr := this._getBytes(this.Base+queueNowUrl+"?code="+code)
  if respErr != nil {
    return s, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &s)
  if unmarsharlErr != nil {
    return s, unmarsharlErr
  }
  return s, nil
}

func (this API) queueAdd(code string, soundCode string) (SoundQueue, error) {
  var q SoundQueue
  body, respErr := this._getBytes(this.Base+queueAddUrl+"?code="+code+"&sound-code="+soundCode)
  if respErr != nil {
    return q, respErr
  }
  unmarsharlErr := json.Unmarshal(body, &q)
  if unmarsharlErr != nil {
    return q, unmarsharlErr
  }
  return q, nil
}

// end of API object

func InitAPI(baseUrl string){
  api = API{Base:baseUrl/*, client:http.Client()*/}
}
