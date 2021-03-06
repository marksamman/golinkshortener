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
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func registerClick(id int, req *http.Request) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return
	}

	db.Exec("INSERT INTO clicks (inserted, ip, user_agent, link_id) VALUES ($1, $2, $3, $4)",
		time.Now().Unix(), ip, req.UserAgent(), id)
}

func shortenHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		http.Error(w, "failed to split host/port", http.StatusInternalServerError)
		return
	}

	randValue := rand.Intn(4096)
	randomString := []byte{urlSafe[randValue>>6], urlSafe[randValue&63]}

	var linkId int
	if err := db.QueryRow("INSERT INTO links (url, creator_ip, created, random) VALUES ($1, $2, $3, $4) RETURNING id", req.FormValue("url"), ip, time.Now().Unix(), randomString).Scan(&linkId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/shortened/%s%s", encodeInt(linkId), randomString), http.StatusFound)
}

func shortenedHandler(w http.ResponseWriter, req *http.Request) {
	linkId := req.URL.Path[11:]
	if len(linkId) < 3 {
		http.NotFound(w, req)
		return
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		http.Error(w, "failed to split host/port", http.StatusInternalServerError)
		return
	}

	id := decodeInt(linkId[:len(linkId)-2])
	if id == 0 {
		http.NotFound(w, req)
		return
	}

	var url, creator_ip, random string
	var createdInt int64
	if err := db.QueryRow("SELECT url, host(creator_ip), created, random FROM links WHERE id = $1", id).Scan(&url, &creator_ip, &createdInt, &random); err != nil {
		http.NotFound(w, req)
		return
	}

	created := time.Unix(createdInt, 0).String()

	if random != linkId[len(linkId)-2:] {
		http.NotFound(w, req)
		return
	}

	if creator_ip != ip {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query("SELECT inserted, user_agent FROM clicks WHERE link_id = $1 ORDER BY id DESC", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clicks := []Click{}
	for rows.Next() {
		click := Click{}

		var inserted int64
		if err := rows.Scan(&inserted, &click.UserAgent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		click.Inserted = time.Unix(inserted, 0).String()
		clicks = append(clicks, click)
	}

	if err := shortenedTemplate.Execute(w, ShortenedTemplateData{req.Host, linkId, url, created, clicks}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	linkId := req.URL.Path[1:]
	if len(linkId) == 0 {
		http.ServeFile(w, req, "static/index.html")
		return
	}

	if len(linkId) < 3 {
		http.NotFound(w, req)
		return
	}

	conn := redisPool.Get()
	defer conn.Close()
	if res, _ := conn.Do("GET", linkId); res != nil {
		http.Redirect(w, req, string(res.([]uint8)), http.StatusMovedPermanently)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		conn.Do("EXPIRE", linkId, REDIS_CACHE_TTL)
		registerClick(decodeInt(linkId[:len(linkId)-2]), req)
		return
	}

	id := decodeInt(linkId[:len(linkId)-2])
	if id == 0 {
		http.NotFound(w, req)
		return
	}

	var url, random string
	if err := db.QueryRow("SELECT url, random FROM links WHERE id = $1", id).Scan(&url, &random); err != nil {
		http.NotFound(w, req)
		return
	}

	if random != linkId[len(linkId)-2:] {
		http.NotFound(w, req)
		return
	}

	http.Redirect(w, req, url, http.StatusMovedPermanently)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	if _, err := conn.Do("SET", linkId, url); err != nil {
		conn.Do("EXPIRE", linkId, REDIS_CACHE_TTL)
	}
	registerClick(id, req)
}
