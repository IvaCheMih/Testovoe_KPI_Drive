package services

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

type ConfigurationService struct {
	RabbitURL     string
	KpiPostApiURL string
	KpiGetApiURL  string
	Token         string
}

func CreateConfigurationService() (ConfigurationService, error) {
	RabbitURL, exists := os.LookupEnv("RABBIT_URL")

	if !exists {
		log.Error("LookupEnv error")
		return ConfigurationService{}, errors.New("LookupEnv: RABBIT_URL not found")
	}

	KpiPostApiURL, exists := os.LookupEnv("KPI_POST_API_URL")

	if !exists {
		log.Error("LookupEnv error")
		return ConfigurationService{}, errors.New("LookupEnv: KPI_POST_API_URL not found")
	}

	KpiGetApiURL, exists := os.LookupEnv("KPI_GET_API_URL")

	if !exists {
		log.Error("LookupEnv error")
		return ConfigurationService{}, errors.New("LookupEnv: KPI_GET_API_URL not found")
	}

	Token, exists := os.LookupEnv("KPI_GET_API_URL")

	if !exists {
		log.Error("LookupEnv error")
		return ConfigurationService{}, errors.New("LookupEnv: TOKEN not found")
	}

	return ConfigurationService{
		RabbitURL:     RabbitURL,
		KpiPostApiURL: KpiPostApiURL,
		KpiGetApiURL:  KpiGetApiURL,
		Token:         Token,
	}, nil
}
