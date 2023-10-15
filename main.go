package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
  "io"
	"encoding/json"
)

type user struct {
	Name string `json:"name"`
}

func main() {
	router := gin.Default()
	router.GET("/hello/go", getGreeting)

	router.Run("127.0.0.1:8090")
}

func getGreeting(c *gin.Context) {
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
	data := user{}
	unmarshalErr := json.Unmarshal([]byte(body), &data)

	if unmarshalErr != nil {
		fmt.Print(unmarshalErr.Error())
	}
	c.IndentedJSON(http.StatusOK, data)
}
