package pkg

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

type Gender hessian.JavaEnum

const (
	MAN hessian.JavaEnum = iota
	WOMAN
)

var genderName = map[hessian.JavaEnum]string{
	MAN:   "MAN",
	WOMAN: "WOMAN",
}

var genderValue = map[string]hessian.JavaEnum{
	"MAN":   MAN,
	"WOMAN": WOMAN,
}

func (g Gender) JavaClassName() string {
	return "org.apache.dubbo.sample.Gender"
}

func (g Gender) String() string {
	s, ok := genderName[hessian.JavaEnum(g)]
	if ok {
		return s
	}

	return strconv.Itoa(int(g))
}

func (g Gender) EnumValue(s string) hessian.JavaEnum {
	v, ok := genderValue[s]
	if ok {
		return v
	}

	return hessian.InvalidJavaEnum
}

type User struct {
	// !!! Cannot define lowercase names of variable
	ID   string `hessian:"id"`
	Name string
	Age  int32
	Time time.Time
	Sex  Gender // notice: java enum Object <--> go string
}

func (u User) String() string {
	return fmt.Sprintf(
		"User{ID:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.ID, u.Name, u.Age, u.Time, u.Sex,
	)
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.sample.User"
}

type UserProvider struct {
	GetUsers func(req []string) ([]*User, error)
	GetErr   func(ctx context.Context, req *User) (*User, error)

	GetUser func(ctx context.Context, req *User) (*User, error)

	GetUserNew func(ctx context.Context, req1, req2 *User) (*User, error)

	GetUser0  func(id string, name string) (User, error)
	GetUser2  func(ctx context.Context, req int32) (*User, error) `dubbo:"getUser"`
	GetUser3  func() error
	GetGender func(ctx context.Context, i int32) (Gender, error)
	Echo      func(ctx context.Context, req interface{}) (interface{}, error) // Echo represent EchoFilter will be used
}
