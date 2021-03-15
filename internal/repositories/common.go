package repositories

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
	"time"
)

var dbFile = "./data.json"

type CommonRepositoryRecord struct {
	ID        string      `json:"id"`
	Version   string      `json:"version"`
	Content   interface{} `json:"content"`
	CreatedAt time.Time   `json:"createdAt"`
}

var dbFileLocker sync.RWMutex

func writeFile(filename string, data []byte, perm fs.FileMode) error {
	dbFileLocker.Lock()
	defer dbFileLocker.Unlock()

	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			err := ioutil.WriteFile(dbFile, []byte("{}"), 0600)
			if err != nil {
				return err
			}
		}
		return err
	}

	fileData, err := ioutil.ReadFile(dbFile)
	if err != nil {
		return err
	}

	var savedData map[string]string
	err = json.Unmarshal(fileData, &savedData)
	if err != nil {
		return err
	}

	savedData[filename] = string(data)

	savedDataBytes, err := json.Marshal(savedData)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dbFile, savedDataBytes, 0600)
}

func readFile(filename string) ([]byte, error) {
	dbFileLocker.RLock()
	defer dbFileLocker.RUnlock()

	fileData, err := ioutil.ReadFile(dbFile)
	if err != nil {
		return nil, err
	}

	var savedData map[string]string
	err = json.Unmarshal(fileData, &savedData)
	if err != nil {
		return nil, err
	}

	if fileData, ok := savedData[filename]; ok {
		return []byte(fileData), nil
	}

	return nil, errors.New("data not found in db")
}

func truncateData(pattern *regexp.Regexp) error {
	dbFileLocker.Lock()
	defer dbFileLocker.Unlock()

	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	fileData, err := ioutil.ReadFile(dbFile)
	if err != nil {
		return err
	}

	var savedData map[string]string
	err = json.Unmarshal(fileData, &savedData)
	if err != nil {
		return err
	}

	for k := range savedData {
		if pattern.MatchString(k) {
			delete(savedData, k)
		}
	}

	savedDataBytes, err := json.Marshal(savedData)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dbFile, savedDataBytes, 0600)
}
