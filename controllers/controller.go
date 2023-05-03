package controllers

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	url := "https://jsonplaceholder.typicode.com/photos"

	// opt, err := redis.ParseURL("redis://default:vPkPKqXYev4rWCKuhPa68FVhWLZFnnu6@redis-17174.c114.us-east-1-4.ec2.cloud.redislabs.com:17174/products-db")
	// if err != nil {
	// 	panic(err)
	// }
	// client := redis.NewClient(opt)

	client := redis.NewClient(&redis.Options{
		Addr:     "redis-17174.c114.us-east-1-4.ec2.cloud.redislabs.com:17174",
		Password: "vPkPKqXYev4rWCKuhPa68FVhWLZFnnu6",
		Username: "default",
		DB:       0,
	})
	ctx := context.Background()

	// err := client.Set(ctx, "foo", "bar", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }
	// val, err := client.Get(ctx, "foo").Result()
	// if err != nil {
	// 	panic(err)
	// }
	photos, err := client.Get(ctx, "photos").Result()
	if err == redis.Nil {
		//key does not exist, perform request and set a key
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		byteData, err := ioutil.ReadAll(response.Body)
		err = client.Set(ctx, "photos", string(byteData), 0).Err()
		if err != nil {
			panic(err)
		}
		w.Write(byteData)
	} else if err != nil {
		panic(err)
	} else {
		w.Write([]byte(photos))
	}

}

func GetProductUsingID(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	url := "https://jsonplaceholder.typicode.com/photos/" + id

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	byteData, err := ioutil.ReadAll(res.Body)
	w.Write(byteData)

}
