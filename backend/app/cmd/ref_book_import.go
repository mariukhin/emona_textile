package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"amifactory.team/sequel/coton-app-backend/app/model"
	"github.com/sirupsen/logrus"
)

type RefBookImport struct {
	RefBookJson    string `long:"ref-book-json" required:"true" description:"Json file with refbook records"`
	CountryCodeCSV string `long:"country-code-csv" required:"true" description:"CSV file with country codes"`
	Verbose        bool   `long:"verbose" description:"Print verbose messages"`
	MongoDBCommand
}

func (c *RefBookImport) Execute(args []string) error {
	appLogger := logrus.New()
	appLogger.Out = os.Stdout
	if c.Verbose {
		appLogger.Level = logrus.DebugLevel
	} else {
		appLogger.Level = logrus.WarnLevel
	}

	appLogger.SetFormatter(&logrus.TextFormatter{})

	appLogger.Infof("CotonSMS refbook import started")

	data, err := ioutil.ReadFile(c.RefBookJson)
	if err != nil {
		panic(err)
	}

	type ParsedCode struct {
		ID          string `json:"_id"`
		PhonePrefix string `json:"dialcode_src"`
		Mcc         int    `json:"mcc"`
		Mnc         *int   `json:"mnc"`
		MinLen      int    `json:"min_length"`
		MaxLen      int    `json:"max_length"`
	}

	parsedCodes := make([]ParsedCode, 0)
	err = json.Unmarshal(data, &parsedCodes)
	if err != nil {
		panic(err)
	}

	filteredParsedCodes := make([]ParsedCode, 0)

	for _, codeItem := range parsedCodes {
		if codeItem.Mnc != nil {
			filteredParsedCodes = append(filteredParsedCodes, codeItem)
		}
	}

	csvFile, err := os.Open(c.CountryCodeCSV)
	if err != nil {
		panic(err)
	}

	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	_, err = csvReader.Read()
	if err != nil {
		panic(err)
	}

	type CountryInfo struct {
		MCC              int
		MNC              int
		CountryCode2     string
		CountryPhoneCode string
		Country          string
		Operator         string
	}

	records := make([]CountryInfo, 0)

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		mcc, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}
		mnc, err := strconv.Atoi(record[1])
		if err != nil {
			continue
		}
		records = append(records, CountryInfo{
			MCC:              mcc,
			MNC:              mnc,
			CountryCode2:     record[2],
			Country:          record[3],
			CountryPhoneCode: record[4],
			Operator:         record[5],
		})
	}

	refBookRecords := make([]model.RefBookRecord, 0)
	for _, code := range filteredParsedCodes {
		refBookRecord := model.RefBookRecord{
			ID:          code.ID,
			PhonePrefix: code.PhonePrefix,
			Mcc:         code.Mcc,
			Mnc:         *code.Mnc,
			MinLen:      code.MinLen,
			MaxLen:      code.MaxLen,
		}

		for _, record := range records {
			if record.MNC == refBookRecord.Mnc && record.MCC == refBookRecord.Mcc {
				refBookRecord.CountryPhoneCode = record.CountryPhoneCode
				refBookRecord.CountryCode2 = strings.ToUpper(record.CountryCode2)
				refBookRecord.Country = record.Country
				refBookRecord.Operator = record.Operator
				break
			}
		}

		refBookRecords = append(refBookRecords, refBookRecord)
	}

	storeOpts := model.NewStorageOptions()
	storeOpts.MongoURI = c.MongoURI
	storeOpts.MongoDatabase = c.MongoDatabase
	store, err := storeOpts.Storage()
	if err != nil {
		panic(err)
	}

	refBookStore, err := model.NewRefBookStore(store)
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	for _, rec := range refBookRecords {
		err = refBookStore.Save(ctx, &rec)
		if err != nil {
			panic(err)
		}
	}

	appLogger.Infof("CotonSMS refbook import finished")
	return nil
}
