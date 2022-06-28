package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type UserProvider1 struct {
	GetUser func(ctx context.Context, req *User) (*User, error)
}

func (u *User) JavaClassName() string {
	return ""
}

type User struct {
	ID   string    `json:"id,omitempty"`
	Code int64     `json:"code,omitempty"`
	Name string    `json:"name,omitempty"`
	Age  int32     `json:"age,omitempty"`
	Time time.Time `json:"time,omitempty"`
}

type userDB struct {
	//key is name ,value is user obj
	nameIndex map[string]*User
	//key is code, value is user obj
	codeIndex map[int64]*User
	lock      sync.Mutex
}

func (db *userDB) GetByName(n string) (*User, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.nameIndex[n]
	return r, ok

}

func outLn(format string, args ...interface{}) {
	fmt.Println("\033[32;40m"+format+"\033[0m\n", args)

}
