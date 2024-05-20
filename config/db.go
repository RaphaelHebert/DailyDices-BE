package config

import (
	"context"
	"fmt"
	"log"

	"github.com/gofor-little/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Mg MongoInstance

var Mc *mongo.Client

func init(){
	fmt.Println("loading env...")
	err := env.Load(".env"); 
	if err != nil {
		panic(err)
	}
	Mc, err = Connect(); 
	if err != nil {
		log.Fatal(err)
	}
}

func Connect() (*mongo.Client, error) {
	// Database settings
	var dbName = "dices"
	var user = env.Get("MONGO_USER", "")
	var password = env.Get("MONGO_PASSWORD", "")
	var mongoURI = fmt.Sprintf("mongodb://%s:%s@localhost:27017/", user, password) + dbName + "?authSource=admin"

	fmt.Println("connecting to db ",dbName)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		msg := fmt.Errorf("could not ping db: %s", err)
		fmt.Println(msg)
	}


	db := client.Database(dbName)

	Mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	fmt.Println(Mg.Db)
	fmt.Println(Mg.Client)


	return client, nil
}

func Disconnect (mc *mongo.Client){
	if err := mc.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}