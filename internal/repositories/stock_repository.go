package repositories

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"stonk-watcher/internal/entities"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetStockInfo(ticker string) (*entities.StockInfoDTO, error) {
	ticker = strings.ToLower(ticker)
	key := fmt.Sprintf("stock-info-%s.json", ticker)
	bytes, err := readFile(key)
	if err != nil {
		return nil, err
	}

	var res CommonRepositoryRecord
	res.Content = &entities.StockInfoDTO{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}

	content, ok := res.Content.(*entities.StockInfoDTO)
	if !ok {
		return nil, errors.New("invalid stock data")
	}

	version, err := getStockInfoVersion()
	if err != nil {
		return nil, err
	}

	if res.Version != version {
		return nil, nil
	}

	return content, nil
}

func PersistStockInfo(ticker string, dto *entities.StockInfoDTO) error {
	version, err := getStockInfoVersion()
	if err != nil {
		return err
	}

	ticker = strings.ToLower(ticker)
	key := fmt.Sprintf("stock-info-%s.json", ticker)
	bytes, err := json.Marshal(CommonRepositoryRecord{
		ID:        uuid.NewString(),
		Version:   version,
		Content:   dto,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return writeFile(key, bytes, 0600)
}

func getStockInfoVersion() (string, error) {
	defaultStockInfoDTO := entities.StockInfoDTO{
		FinvizStockInfoDTO:        &entities.FinvizStockInfoDTO{},
		MarketWatchInfoDTO:        &entities.MarketWatchInfoDTO{},
		MorningStarPerformanceDTO: &entities.MorningStarPerformanceDTO{},
	}
	bytes, err := json.Marshal(defaultStockInfoDTO)
	if err != nil {
		return "", err
	}

	h := md5.New()
	h.Write(bytes)

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func DeleteStockInfo(ticker string) error {
	ticker = strings.ToLower(ticker)
	return truncateData(regexp.MustCompile(fmt.Sprintf("^stock-info-%s\\.json$", ticker)))
}

func TruncateStockInfo() error {
	return truncateData(regexp.MustCompile("^stock-info-.*\\.json$"))
}
