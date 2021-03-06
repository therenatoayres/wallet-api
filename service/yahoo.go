package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"unicode/utf8"

	"strings"

	"strconv"

	"time"

	"github.com/therenatoayres/wallet-api/dto"
)

const yahoofinanceURL string = "https://download.finance.yahoo.com/d/quotes?f=sl1d1t1&s="

//ExchangeRate used to get the currency rate between one currency and one or more different currencies
func ExchangeRate(from string, to []string) ([]dto.Rate, []dto.Currency, error) {

	var rates []dto.Rate
	var fails []dto.Currency

	if len(to) == 0 {
		err := errors.New("Need to send at least one to currency to be able to convert")
		log.Println(err)
		return rates, fails, err
	}

	rates = make([]dto.Rate, 0)
	fails = make([]dto.Currency, 0)

	for _, code := range to {
		rate, fail := request(from, code)
		if rate == nil {
			fails = append(fails, *fail)
		} else {
			rates = append(rates, *rate)
		}
	}

	return rates, fails, nil
}

func request(from, to string) (*dto.Rate, *dto.Currency) {

	var rate *dto.Rate
	currency := &dto.Currency{
		CodeFrom: from,
		CodeTo:   to,
	}

	url := yahoofinanceURL + from + to + "=X"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return rate, currency
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return rate, currency
	}
	defer response.Body.Close()

	fmt.Println("Yahoo: ", response)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return rate, currency
	}

	parse := string(body)
	value, err := parseYahooResponse(parse)
	if err != nil {
		log.Println(err)
		return rate, currency
	}

	rate = &dto.Rate{
		Conversion: *currency,
		Value:      value,
	}

	return rate, currency
}

func parseYahooResponse(yahoo string) (float64, error) {

	split := strings.Split(yahoo, ",")

	value, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return value, nil
}

//COMMENT
func GetTax(currency *dto.Currency) (*dto.YahooResponse, error) {

	url := yahoofinanceURL + currency.CodeFrom + currency.CodeTo + "=X"
	var yResponse *dto.YahooResponse

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return yResponse, fmt.Errorf("Error creating request")
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return yResponse, fmt.Errorf("Error creating request")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Do: ", err)
		return yResponse, fmt.Errorf("Error creating request")
	}

	response := string(body)

	splittedResponse := strings.Split(response, ",")

	//dateString := strings.Replace(splittedResponse[2], "\"", "", -1) + " " + strings.TrimSpace(strings.Replace(splittedResponse[3], "\"", "", -1))
	date, err := treatDate(splittedResponse[2] + " " + splittedResponse[3])
	if err != nil {
		log.Fatal("Do: ", err)
		return yResponse, fmt.Errorf("Error creating request")
	}

	date = date.UTC()

	rate, err := strconv.ParseFloat(splittedResponse[1], 64)
	if err != nil {
		log.Fatal("Do: ", err)
		return yResponse, fmt.Errorf("Error creating request")
	}

	fmt.Println("Response From YAHOO: ", response)

	yResponse = &dto.YahooResponse{

		Rate: rate,
		Date: date,
	}

	return yResponse, nil
}

func treatDate(dateString string) (time.Time, error) {

	var date time.Time

	dateString = strings.TrimSpace(strings.Replace(dateString, "\"", "", -1))

	dateAndTime := strings.Split(dateString, " ")
	dayMonthandYear := strings.Split(dateAndTime[0], "/")

	day := dayMonthandYear[1]
	month := dayMonthandYear[0]
	year := dayMonthandYear[2]

	if utf8.RuneCountInString(month) == 1 {
		month = "0" + month
	}

	date, err := time.Parse("02/01/2006 15:04pm", day+"/"+month+"/"+year+" "+dateAndTime[1])
	if err != nil {
		log.Fatal("Do: ", err)
		return date, fmt.Errorf("Error creating request")
	}

	return date.UTC(), nil
}
