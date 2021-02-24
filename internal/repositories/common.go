package repositories

import (
	"io/ioutil"
	"time"
)

type CommonRepositoryRecord struct {
	ID        string      `json:"id"`
	Version   string      `json:"version"`
	Content   interface{} `json:"content"`
	CreatedAt time.Time   `json:"createdAt"`
}

var writeFile = ioutil.WriteFile
var readFile = ioutil.ReadFile
