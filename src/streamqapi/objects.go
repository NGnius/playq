package streamqapi

import (
  "os"
)

// start of Sound object

type Sound struct {
  Code string `json:code`
  Id int `json:id`
  Metadata map[string]string `json:metadata`
}

/*
  Sound constructor
  Retrieves Sound info from server
*/
func NewSound(code string) (Sound, error) {
  return api.getSound(code)
}

/*
  Retrieve file from cache or API server
*/
func (sound Sound) GetFile() *os.File {
  soundFile := cachedir+"/"+sound.Code
  _, fileErr := os.Stat(soundFile)
  if fileErr == nil || os.IsExist(fileErr) {
    // load from cache
    file, _ := os.Open(soundFile)
    file.Seek(0,0)
    return file
  } else {
    // download and cache
    file, _ := os.Create(soundFile)
    fileBytes, apiErr := api.getFile(sound.Code)
    if apiErr != nil {
      return file
    }
    file.Write(fileBytes)
    file.Sync()
    file.Seek(0,0)
    return file
  }
}

// end of Sound object

// start of SoundQueue object

type SoundQueue struct {
  Code string `json:code`
  Id int `json:id`
  Items []Sound `json:items`
  Index int `json:index`
}

/*
  SoundQueue constructor
  Retrieves SoundQueue info from server
*/

func NewSoundQueue(code string) (SoundQueue, error){
  return api.getQueue(code)
}

/*
  Get current song from API server
  Equivalent to Now() under ideal circumstances
*/
func (sq SoundQueue) GetNow() (Sound, error) {
  return api.queueNow(sq.Code)
}

/*
  Currently playing song
*/
func (sq SoundQueue) Now() (Sound) {
  return sq.Items[sq.Index]
}

/*
  Get next song from API server
  If local next is not the same as the API's next, the queue is resynced
*/
func (sq *SoundQueue) GetNext() (Sound, error) {
  nextSound, apiErr := api.queueNext(sq.Code)
  if apiErr != nil {
    return nextSound, apiErr
  }
  sq.Index++
  if sq.Now().Code != nextSound.Code {
    // resync queue
    updatedQueue, apiErr := api.getQueue(sq.Code)
    if apiErr != nil {
      return nextSound, apiErr
    }
    sq.Index = updatedQueue.Index
    sq.Items = updatedQueue.Items
  }
  return sq.Now(), nil
}

/*
  Next song, without API sync/calls
*/
func (sq SoundQueue) Next() (Sound) {
  sq.Index++
  return sq.Now()
}

// end of Queue object
