package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eoria17/AWS-Golang-Music-Sub/controllers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	controllers.Route(router)

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/"))))
	http.Handle("/assets/", router)

	fmt.Println("Currently Listening to port 8080..")

	log.Println(http.ListenAndServe(":8080", router))

}
