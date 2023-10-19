package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
  "io"
	"encoding/json"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/google/uuid"
)

type user struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type username struct {
	Name string `json:"name"`
}

const mongoUri = "mongodb://root:example@localhost:27017"

func main() {
	router := gin.Default()
	router.GET("/hello/go", getGreeting)

	router.Run("127.0.0.1:8090")
}

func getGreeting(c *gin.Context) {
	userCh := make(chan username)
	go retrieveUserName(userCh)

	go storeUser(userCh)
	c.IndentedJSON(http.StatusOK, <- userCh)
}

func storeUser(userCh chan username) {
	userName := <- userCh
	user := user {
		Id: uuid.New().String(), 
		Name: userName.Name,
	}
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("admin").Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	fmt.Printf("Inserted user with _id: %v\n", result.InsertedID)
	userCh <- userName
}

func retrieveUserName(userCh chan username) {
	url := "http://localhost:8080/hello/wiremock"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
			fmt.Print(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
			fmt.Print(err.Error())
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
			fmt.Print(readErr.Error())
	}
	userName := username{}
	unmarshalErr := json.Unmarshal([]byte(body), &userName)

	if unmarshalErr != nil {
		fmt.Print(unmarshalErr.Error())
	}
	userCh <- userName
}
