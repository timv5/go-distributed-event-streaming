package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type RestClientInterface interface {
	CallService(responseChannel chan<- string)
}

type RestClient struct {
}

func NewRestClient() RestClient {
	return RestClient{}
}

func (rs RestClient) CallService(responseChannel chan<- string) {
	response, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		responseChannel <- fmt.Sprintf("ERROR: message not received from a client: %s", err.Error())
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Cannot marshal response from a rest client: ", err.Error())
		return
	}

	responseChannel <- string(body)
}
