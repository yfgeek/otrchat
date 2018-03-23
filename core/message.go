package core

import (
	"net"
	"time"
	"fmt"
)

type Message struct{
	Status int
	UserID int
	UserName string
	Content string
	ToUserID int
	FromIP   *net.UDPAddr
	ToIP 	 *net.UDPAddr
	Time string
}


func NowTime() string{
	now := time.Now()
	local, err2 := time.LoadLocation("Local")
	if err2 != nil {
		fmt.Println(err2)
	}
	return now.In(local).Format("15:04:05")

}