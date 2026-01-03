package main

import (
	"errors"
	"fmt"
	"os"

	userSDK "github.com/Wito982/go_course_sdk/user"
)

func main() {
	userTrans := userSDK.NewHttpClient("http://localhost:8081", "")

	user, err := userTrans.Get("1c85cf05-3968-40e4-a7b3-f6d5efc5")
	if err != nil {
		if errors.As(err, &userSDK.ErrNotFound{}) {
			fmt.Println("Not found:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Internal Server Error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("User:", user)
}
