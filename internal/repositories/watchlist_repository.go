package repositories

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func GetWatchlist() ([]string, error) {
	key := "watchlist.json"
	bytes, err := readFile(dataPath + key)
	if err != nil {
		return nil, err
	}

	var res CommonRepositoryRecord
	res.Content = new([]string)

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}

	return *res.Content.(*[]string), nil
}

func PersistWatchlist(watchlist []string) error {
	key := "watchlist.json"
	bytes, err := json.Marshal(CommonRepositoryRecord{
		ID:        uuid.NewString(),
		Version:   "1.0.0",
		Content:   watchlist,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return writeFile(dataPath+key, bytes, 0600)
}
