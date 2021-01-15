package cache

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func Connect() *redis.Client {

	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT") )

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return client
}

/*
	this funcntion for ser value to redis
	like setex
	@param		key		string
	@param		value	string
	@param		expire	int
*/
func Set(key, value string, expire time.Duration) error {
	client := Connect()
	err := client.Set(key, value, expire).Err()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func Exists(key string) bool {
	client := Connect()
	value := client.Exists(key)
	v := value.Val()

	if v == 0 {
		return false
	} else if v == 1 {
		return true
	}

	return false
}

/*
	this function for get velue from redis
	@param		key		string
	@param		interface{}
	@param		error
*/
func Get(key string) (string, error) {


	client := Connect()
	value, err := client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("no value found")
		return "", err
	} else if err != nil {
		// handle error
		fmt.Println(err.Error())
		return "", err
	}

	return value, nil
}

/*
	delete key in redis
	@param		key		string
	@return		error
*/
func Del(key string) error {
	client := Connect()
	err := client.Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}
