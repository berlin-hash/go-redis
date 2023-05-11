package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var requrl string = "https://jsonplaceholder.typicode.com/photos"

func createClient() (*redis.Client, context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return client, ctx
}

func getOrSetCache(ctx context.Context, client *redis.Client, key string, callback func() string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		//does not exist
		data := callback()
		client.Set(ctx, key, data, 10*time.Second).Err()
		return data, nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func GetProducts(w http.ResponseWriter, r *http.Request) {

	// opt, err := redis.ParseURL("redis://default:vPkPKqXYev4rWCKuhPa68FVhWLZFnnu6@redis-17174.c114.us-east-1-4.ec2.cloud.redislabs.com:17174/products-db")
	// if err != nil {
	// 	panic(err)
	// }
	// client := redis.NewClient(opt)

	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "redis-17174.c114.us-east-1-4.ec2.cloud.redislabs.com:17174",
	// 	Password: "vPkPKqXYev4rWCKuhPa68FVhWLZFnnu6",
	// 	DB:       0,
	// })

	client, ctx := createClient()

	queryParams := r.URL.Query()

	photos, err := getOrSetCache(ctx, client, fmt.Sprintf("photos?albumId=%s", queryParams.Get("albumId")), func() string {
		res, err := http.Get(fmt.Sprintf("%s?%s", requrl, queryParams.Encode()))
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		byteData, err := ioutil.ReadAll(res.Body)
		return string(byteData)
	})
	if err != nil {
		panic(err)
	}
	w.Write([]byte(photos))

}

func GetProductUsingID(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	client, ctx := createClient()

	photos, err := getOrSetCache(ctx, client, fmt.Sprintf("photos/%s", id), func() string {
		res, err := http.Get(fmt.Sprintf("%s/%s", requrl, id))
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		byteData, err := ioutil.ReadAll(res.Body)
		return string(byteData)
	})
	if err != nil {
		panic(err)
	}
	w.Write([]byte(photos))

}
