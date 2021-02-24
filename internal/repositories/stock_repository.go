package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"stonk-watcher/internal/config"
	"stonk-watcher/internal/entities"
	"strings"
	"time"

	"github.com/google/uuid"
)

const dataPath = "./data/"

func GetStockInfo(ticker string) (*entities.StockInfoDTO, error) {
	ticker = strings.ToLower(ticker)
	key := fmt.Sprintf("stock-info-%s.json", ticker)
	bytes, err := readFile(dataPath + key)
	if err != nil {
		return nil, err
	}

	var res CommonRepositoryRecord
	res.Content = &entities.StockInfoDTO{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}

	return res.Content.(*entities.StockInfoDTO), nil
}

func PersistStockInfo(ticker string, dto *entities.StockInfoDTO) error {
	ticker = strings.ToLower(ticker)
	key := fmt.Sprintf("stock-info-%s.json", ticker)
	bytes, err := json.Marshal(CommonRepositoryRecord{
		ID:        uuid.NewString(),
		Version:   config.GetVersion(),
		Content:   dto,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return writeFile(dataPath+key, bytes, 0600)
}

func TruncateStockInfo() error {
	files, err := filepath.Glob(dataPath + "stock-info-*.json")
	if err != nil {
		return err
	}

	for _, f := range files {
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}
	return nil
}
