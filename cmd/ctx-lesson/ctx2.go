package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	duration := 1 * time.Second
	ctx := context.Background()
	ctx = context.WithValue(ctx, "test", "hello")
	ctx, cancel := context.WithTimeout(ctx, duration)

	defer cancel()

	doRequest(ctx, "https://ya.ru")
}

func doRequest(ctx context.Context, requestStr string) {
	req, _ := http.NewRequest(http.MethodGet, requestStr, nil)
	req = req.WithContext(ctx)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Printf("status code =%v", res.StatusCode)
		fmt.Println(ctx.Value("test"))
	case <-ctx.Done():
		fmt.Println("req too long")
	}

}
