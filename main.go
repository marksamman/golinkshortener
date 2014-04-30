/*
 * Copyright (c) 2014 Mark Samman <https://github.com/marksamman/golinkshortener>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
)

var db *sql.DB
var redisPool *redis.Pool

func main() {
	// Postgres pool
	var err error
	if db, err = sql.Open("postgres", POSTGRES_CONNECTION_PARAMS); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Redis pool
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", REDIS_HOST)
	}, MAX_IDLE_REDIS_CONNECTIONS)
	defer redisPool.Close()

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/shortened/", shortenedHandler)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", APP_PORT), nil))
}
