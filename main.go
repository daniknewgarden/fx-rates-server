package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"fx-rates-rest-api/db"
	"fx-rates-rest-api/rates"
	"fx-rates-rest-api/scraper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client


func main() {
	// Connect to mongoDB
	mongoClient = db.GetMongoClient()

	// Fetch the document
    document, err := scraper.FetchDocument(scraper.URL, 3)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	// Scrape the USD FX rates
	usdFXMatrixMap := scraper.ScrapeUSDRate(document)

	// Save parsed USD rates to MongoDB
	usdFxMatrixCollection := mongoClient.Database("test").Collection("fx-matrix") // TODO: thing about data type we save to database (currently map[string]float64)

	filter := bson.M{"baseCurrency": "USD"}
	update := bson.M{"$set": bson.M{"fxRates": usdFXMatrixMap}}

	updateOptions := options.Update().SetUpsert(true)
	_, err = usdFxMatrixCollection.UpdateOne(context.TODO(), filter, update, updateOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("USD FX Matrix saved to database!")

	// Start the server
	http.HandleFunc("/getRates", rates.GetRatesHandler(mongoClient))
	http.ListenAndServe(":8080", nil)

	// Close the mongoDB connection
	defer mongoClient.Disconnect(context.Background())
}