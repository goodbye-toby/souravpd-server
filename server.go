package main

import(
	"net/http"
	"log"
	"./models.go"
)

func main(){
	port := ":8080"
	log.Fatal( http.ListenAndServe(port , nil))
}
