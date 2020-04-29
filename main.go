package main

import (
	"fmt"
	"gee/gee"
	"io/ioutil"
	"log"
	"net/http"
)


func main() {
	r := gee.New()

	r.Get("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hell world!")
	})

	r.POST("/hello", func(c *gee.Context) {
		data, err := ioutil.ReadAll(c.Req.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "心累服务器崩了")
		}

		defer c.Req.Body.Close()

		fmt.Println(string(data))

	})

	log.Fatal(r.Run(":8000"))
}