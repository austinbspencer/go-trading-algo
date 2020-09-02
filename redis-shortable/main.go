package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/my/repo/go/src/github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

type alpacaClientContainer struct {
	api          *alpaca.Client
	nasdaq       []alpaca.Asset
	shortable    []string
	nonShortable []string
}

var alpacaClient alpacaClientContainer

func init() {
	os.Setenv(common.EnvApiKeyID, goDotEnvVariable("KEY"))
	os.Setenv(common.EnvApiSecretKey, goDotEnvVariable("SECRET"))

	// fmt.Printf("Running w/ credentials [%v %v]\n", common.Credentials().ID, common.Credentials().Secret)

	alpaca.SetBaseUrl(goDotEnvVariable("BASE_URL"))
	alpacaClient = alpacaClientContainer{
		alpaca.NewClient(common.Credentials()),
		[]alpaca.Asset{},
		[]string{},
		[]string{},
	}
}

func runningtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}

func main() {
	defer track(runningtime("main"))
	alpacaClient.setAssets()
	alpacaClient.getShortable()
}

func (alp *alpacaClientContainer) setAssets() {
	// Get a list of all active assets.
	status := "active"
	assets, err := alp.api.ListAssets(&status)
	if err != nil {
		panic(err)
	}
	// Filter the assets down to just those on NASDAQ.
	for _, asset := range assets {
		if asset.Exchange == "NASDAQ" {
			alp.nasdaq = append(alp.nasdaq, asset)
		}
	}
}

func (alp *alpacaClientContainer) getShortable() {
	fmt.Println("Nasdaq size:", len(alp.nasdaq))
	for i := 0; i < len(alp.nasdaq); i++ {
		if alp.nasdaq[i].Tradable && alp.nasdaq[i].Shortable {
			alp.shortable = append(alp.shortable, alp.nasdaq[i].Symbol)
		} else {
			alp.nonShortable = append(alp.nonShortable, alp.nasdaq[i].Symbol)
		}
	}

	fmt.Println("Shortable Assets:", len(alp.shortable))
	fmt.Println("non-Shortable Assets:", len(alp.nonShortable))
}
