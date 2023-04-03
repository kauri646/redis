package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kauri646/redis.git/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/go-redis/redis"
)
func main() {
	db.RedisInit()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
  
	// Routes
	e.GET("/insert", Insert)
	e.GET("/get", Get)
  
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
  }

type RespJson struct {
	Data interface{}
	Status string
}

type RequestRedis struct {
	Name string
	Age string
}

var key = "app_kauri"

func Insert(c echo.Context) error {
	id := c.QueryParam("id")
	name := c.QueryParam("name")
	age := c.QueryParam("age")

	rdb := db.RedisConnect()

	reqRedis := RequestRedis{
		Name: name,
        Age: age,
	}
	req, _ := json.Marshal(reqRedis)

	err := rdb.HSet(key, id, req).Err()
	if err!= nil {
        return fmt.Errorf("error set redis: %v", err)
    }

	resp := RespJson{
		Data: id,
        Status: "success",
	}

	return c.JSON(http.StatusOK, resp)
  
}

func Get(c echo.Context) error {
	id := c.QueryParam("id")

	rdb := db.RedisConnect()

	val, err := rdb.HGet(key, id).Result()
	if err == redis.Nil{
		return c.JSON(http.StatusNotFound, "data tidak ditemukan")
	}else if err != nil{
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("error get redis %s", err.Error()))
	}

	var requestRedis RequestRedis
	err = json.Unmarshal([]byte(val), &requestRedis)
	if err!= nil {
        return c.JSON(http.StatusBadRequest, fmt.Sprintf("error unmarshal redis %s", err.Error()))
    }

	resp := RespJson{
		Data: requestRedis,
        Status: "success",
	}

	return c.JSON(http.StatusOK, resp)
}