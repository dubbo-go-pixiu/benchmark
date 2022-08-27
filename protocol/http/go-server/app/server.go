/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

import "github.com/apache/dubbo-go-pixiu/pkg/common/constant"

type User struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Age  int32     `json:"age"`
	Time time.Time `json:"time"`
}

func init() {
	cache = &UserDB{
		cacheMap: make(map[string]*User, 16),
		lock:     sync.Mutex{},
	}

	t1, _ := time.Parse(
		time.RFC3339, "%v")

	cache.Add(&User{ID: "0001", Name: "tc", Age: 18, Time: t1})
	cache.Add(&User{ID: "0002", Name: "ic", Age: 88, Time: t1})
}

var cache *UserDB

// UserDB cache user.
type UserDB struct {
	// key is name, value is user obj
	cacheMap map[string]*User
	lock     sync.Mutex
}

// Add adds the user to cache.
func (db *UserDB) Add(u *User) bool {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.cacheMap[u.Name] = u
	return true
}

// Get returns the user.
func (db *UserDB) Get(n string) (*User, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.cacheMap[n]
	return r, ok
}

func main() {
	http.HandleFunc("/user/", user)
	log.Println("Starting http server ...")
	log.Fatal(http.ListenAndServe(":1314", nil))
}

func user(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case constant.FileDateFormat:
		byts, err := ioutil.ReadAll(r.Body)
		if err != nil {
			write, err := w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(write)
			}
			log.Println(byts)
		}
		var user User
		err = json.Unmarshal(byts, &user)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		_, ok := cache.Get(user.Name)
		if ok {
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write([]byte("{\"message\":\"data is exist\"}"))
			return
		}
		user.ID = randSeq(5)
		if cache.Add(&user) {
			b, _ := json.Marshal(&user)
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write(b)
			return
		}
		w.Write(nil)

	case constant.DefaultHTTPType:
		subPath := strings.TrimPrefix(r.URL.Path, "/user/")
		userName := strings.Split(subPath, "/")[0]
		var u *User
		var b bool
		if len(userName) != 0 {
			log.Printf("paths: %v", userName)
			u, b = cache.Get(userName)
		} else {
			q := r.URL.Query()
			u, b = cache.Get(q.Get("name"))
		}
		// w.WriteHeader(200)
		if b {
			b, _ := json.Marshal(u)
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write(b)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
