package streamqapi

import (
  "fmt"
  "os"
  "io/ioutil"
)

// start of Sound object

type Sound struct {
  Code string `json:"code"`
  Id int `json:"id"`
  Metadata map[string]string `json:"metadata"`
}

/*
  Sound constructors
  Retrieves Sound info from server
*/
func NewSound(code string) (Sound, error) {
  return api.getSound(code)
}

func CreateSound() (Sound, error) {
  return api.newSound()
}

/*
  Refresh metadata from API server
*/
func (sound *Sound) Refresh() (error) {
  nSound, refreshErr := api.refreshSound(sound.Code)
  if refreshErr != nil {
    return refreshErr
  }
  sound.Metadata = nSound.Metadata
  return nil
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
  Code string `json:"code"`
  Id int `json:"id"`
  Index int `json:"index"`
  Items []Sound `json:"items"`
  RepeatAll bool `json:"repeat-all"`
  RepeatOne bool `json:"repeat-one"`
  IsShuffle bool `json:"shuffle"`
}

/*
  SoundQueue constructors
  Retrieves SoundQueue info from server
*/

func NewSoundQueue(code string) (SoundQueue, error) {
  return api.getQueue(code)
}

func CreateSoundQueue() (SoundQueue, error) {
  return api.newQueue()
}

/*
  Update values which an API call may have changed
*/
func (sq *SoundQueue) update(updatedSq SoundQueue) {
  sq.Index = updatedSq.Index
  sq.Items = updatedSq.Items
  sq.RepeatAll = updatedSq.RepeatAll
  sq.RepeatOne = updatedSq.RepeatOne
  sq.IsShuffle = updatedSq.IsShuffle
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
    sq.update(updatedQueue)
  }
  return sq.Now(), nil
}

/*
  Next song, without API sync/calls
*/
func (sq *SoundQueue) Next() (Sound) {
  sq.Index++
  return sq.Now()
}

/*
  Add song to end of queue
*/
func (sq *SoundQueue) Add(s Sound) (error) {
  updatedSq, apiErr := api.queueAdd(sq.Code, s.Code)
  if apiErr != nil {
    return apiErr
  }
  sq.update(updatedSq)
  return nil
}

/*
  Create new sound, upload file and add to queue compound method
*/
func (sq *SoundQueue) AddFile(file *os.File) (Sound, error) {
  // read data for caching
  file.Seek(0,0)
  data, readErr := ioutil.ReadAll(file)
  if readErr != nil {
    fmt.Println(readErr)
  }
  file.Seek(0,0)
  // API-side creation
  nSound, nSoundErr := CreateSound()
  if nSoundErr != nil {
    return nSound, nSoundErr
  }
  uploadErr := api.uploadFile(nSound.Code, file)
  if uploadErr != nil {
    return nSound, uploadErr
  }
  rSoundErr := nSound.Refresh()
  if rSoundErr != nil {
    return nSound, rSoundErr
  }
  // cache audio file
  soundPath := cachedir+"/"+nSound.Code
  soundFile, _ := os.Create(soundPath)
  soundFile.Write(data)
  soundFile.Sync()
  soundFile.Close()
  return nSound, sq.Add(nSound)
}

/*
  Add file to queue from filepath
*/
func (sq *SoundQueue) AddFilePath(path string) (Sound, error) {
  file, fileErr := os.Open(path)
  nSound, aSoundErr := sq.AddFile(file)
  if fileErr != nil {
    return nSound, fileErr
  }
  file.Close()
  if aSoundErr != nil {
    return nSound, aSoundErr
  }
  return nSound, nil
}

/*
  Set shuffle queue mode to -1, 0 or 1 (off, toggle or on, respectively)
*/
func (sq *SoundQueue) Shuffle(mode int) (error) {
  var nQ SoundQueue
  var apiErr error
  switch mode {
  case -1:
    nQ, apiErr = api.queueShuffle(sq.Code, false)
  case 1:
    nQ, apiErr = api.queueShuffle(sq.Code, true)
  case 0:
    nQ, apiErr = api.queueShuffle(sq.Code, !sq.IsShuffle)
  }
  if apiErr != nil {
    return apiErr
  }
  sq.update(nQ)
  return nil
}

// end of Queue object
