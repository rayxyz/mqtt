package main

import (
	"log"
	"net/http"

	mqttclient "mqtt/client"

	"github.com/gorilla/mux"
)

var (
	client = &mqttclient.Client{}
)

func runMQTTClient() {
	client.Run("localhost:8080")
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	client.Publish("I want to fuck 小姐姐...")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi, welcome to MQTTv3.1.1"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/publish", publishHandler)
	http.Handle("/", r)

	go runMQTTClient()

	log.Fatal(http.ListenAndServe("localhost:8081", r))
}
