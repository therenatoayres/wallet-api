package controller

import (
	"errors"
	"fmt"
	"net/http"

	"log"

	"encoding/json"

	"github.com/therenatoayres/wallet-api/dto"
	"github.com/therenatoayres/wallet-api/service"
)

//GetCurrencyRate  used to get the currency rate between one currency to one or more diffent currencies
func GetCurrencyRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	parameters := r.URL.Query()

	from := parameters["from"]
	to := parameters["to"]

	fmt.Println("Currencies: " + fmt.Sprint(to))

	if from == nil {
		err := errors.New("No from currency suplied")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(from) > 1 {
		err := errors.New("To many from currencies suplied")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if to == nil {
		err := errors.New("No to currency suplied")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rates, fails, err := service.ExchangeRate(from[0], to)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := dto.Response{
		Rates: rates,
		Fails: fails,
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
