package main

import (
	"log"
	"time"

	"github.com/mkusaka/lax/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var conn = db.NewClient(1000 * time.Millisecond)

func main() {
	result, err := conn.NewCustomer("iam")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result)

	res := conn.GetCustomer("iam")
	log.Print(res.ID, res.PrimaryCustomerID)

	cacheKeyConfig := db.CacheKeyConfig{
		HeaderKeys: []string{"Vary"},
		UseURL:     true,
	}

	rules := db.Rules{
		db.Rule{
			Macher:   "localhost/*",
			Matched:  "yo/*",
			Priority: 100,
		},
		db.Rule{
			Macher:   "foo/bar/baz/*",
			Matched:  "yo/*",
			Priority: 200,
		},
		db.Rule{
			Macher:   "/posts/*",
			Matched:  "/posts",
			Priority: 350,
		},
	}

	createdCacheConfig := conn.NewConfig(&res, "localhost:300", "localhost:3000", &cacheKeyConfig, &rules)
	savedCacheConfig, err := conn.SaveConfig(createdCacheConfig)
	if err != nil {
		log.Fatal(err)
	}

	id := (*savedCacheConfig).InsertedID.(primitive.ObjectID)
	stringid := id.Hex()
	convertedID, _ := primitive.ObjectIDFromHex(stringid)

	log.Print("converted")
	log.Print(convertedID)

	config := conn.GetConfig(convertedID)

	proxyed, err := config.ProxyPath("foo/bar/baz/123")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(proxyed)

	inrule, _ := primitive.ObjectIDFromHex("5de26ad00dd8829f348a59b5")
	log.Print(conn.GetConfig(inrule).ID)
}
