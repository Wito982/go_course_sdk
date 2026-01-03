package main

import (
	"errors"
	"fmt"
	"os"

	userSDK "github.com/Wito982/go_course_sdk/course"
)

func main() {
	userTrans := userSDK.NewHttpClient("http://localhost:8082", "")

	user, err := userTrans.Get("ffb66817-ea06-45bc-b8b1-60a22f2af")
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
