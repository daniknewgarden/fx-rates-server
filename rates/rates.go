package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Rate struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type GetRatesResponse struct {
	Rates []Rate `json:"rates"`
}

// TODO: add in memory cache for rates to prevent unnecessary calls to the database
func GetRates(mongoClient *mongo.Client, baseCurrency string) (primitive.M, error) {
	usdFxMatrixCollection := mongoClient.Database("test").Collection("fx-matrix")

	filter := bson.M{"baseCurrency": "USD"}
	var usdRates = bson.M{"baseCurrency": "USD", "fxRates": make(map[string]float64)}

	err := usdFxMatrixCollection.FindOne(context.TODO(), filter).Decode(&usdRates)
	if err != nil {
		return nil, err
	}


    fxRates, ok := usdRates["fxRates"].(primitive.M)
    if !ok {
        return nil, fmt.Errorf("failed to convert usdRates[\"fxRates\"] to primitive.M")
    }

    if baseCurrency == "USD" {
        return fxRates, nil
    }

	baseCurrencyRate, ok := fxRates[baseCurrency].(float64)
	if !ok {
		return nil, fmt.Errorf("failed to convert baseCurrencyRate to float64")
	}

	for currency, usdRate := range fxRates {
		usdRateFloat, ok := usdRate.(float64)
		if !ok {
			return nil, fmt.Errorf("failed to convert usdRate to float64")
		}
		fxRates[currency] = usdRateFloat / baseCurrencyRate
	}

    return fxRates, nil
}

func GetRatesHandler(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseCurrency := r.URL.Query().Get("baseCurrency")
		result, err := GetRates(mongoClient, baseCurrency)
		if err != nil {
			log.Fatal(err)
		}

		var rates []Rate
		for currency, value := range result {
			rates = append(rates, Rate{
				Currency: currency,
				Value: value.(float64),
			})
		}

		response := GetRatesResponse{
			Rates: rates,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}