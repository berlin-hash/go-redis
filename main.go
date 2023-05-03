package main

import (
	"fmt"
	"go-redis/routers"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Go redis")

	r := routers.Router()
	log.Fatal(http.ListenAndServe(":8080", r))

}
