package events

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jz222/loggy/models"
)

type db struct {
	Errs []models.Error
	sync.RWMutex
}

var instance db

func getInstance() *db {
	return &instance
}

func (s *db) Add(err models.Error) {
	s.Lock()
	defer s.Unlock()
	tmp := append(s.Errs, err)
	s.Errs = tmp
}

func (s *db) Get() []models.Error {
	s.RLock()
	defer s.RUnlock()
	return s.Errs
}

func SaveError(errorLog models.Error) {
	hash := md5.Sum([]byte(errorLog.Message + errorLog.Stacktrace))
	errorLog.Fingerprint = hex.EncodeToString(hash[:])

	data, _ := json.MarshalIndent(errorLog, "", "  ")
	fmt.Println(string(data))

	instance.Add(errorLog)
}

func GetErrors() []models.Error {
	return instance.Get()
}
