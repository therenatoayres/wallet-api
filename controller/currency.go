package controller

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

func getRate(w http.ResponseWriter, r *http.Request) {

	var currency dto.Currency

	parameters := r.URL.Query()

	currency.CodeFrom = parameters["from"][0]
	currency.CodeTo = parameters["to"][0]

	t, err := service.getTax(&currency)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // bad request
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(t))

}

func getTotalConversion(w http.ResponseWriter, r *http.Request) {

	//1. create a conversion Slice and a con Slice
	var coins dto.WalletItems

	//2. Get the value to convert to, from the query parameter in the URL
	to := r.URL.Query()["to"][0]

	fmt.Println("Gotta convert stuff to ", to)
	fmt.Println("")

	//3. Read the body of the request, which will be a JSON array of Currency values
	// and convert it to our Currency object
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &coins); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//4. For each Currency Value, we have to connect to Yahoo's server and get our conversion rate
	//Then, get the asrate and calculate the values for each coin
	var total float64
	for i := 0; i < len(coins); i++ {

		conversion := dto.Currency{coins[i].coin, to}

		t, err := service.getTax(&conversion)

		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(500) // bad request
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			break
		}

		rate := strings.Split(t, ",")[1]

		c, err := strconv.ParseFloat(rate, 64)
		if err != nil {
			panic(err)
		}

		fmt.Println("Rate for "+coins[i].coin+" is ", c)
		total = total + c*coins[i].value
	}

	fmt.Println("---------------------")
	fmt.Println("Total: ", total)
	fmt.Println("---------------------")

	//5, Set headers for response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//6. Create and send response
	response := dto.WalletItem{to, total}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // server error
		panic(err)
	}
}
