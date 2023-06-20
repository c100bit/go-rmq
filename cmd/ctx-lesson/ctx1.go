package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main1() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		err := cancelRequest1(ctx)
		if err != nil {
			cancel()
		}
	}()

	doRequest1(ctx, "https://ya.ru")
}

func cancelRequest1(ctx context.Context) error {
	time.Sleep(100 * time.Millisecond)
	return fmt.Errorf("fail request")
}

func doRequest1(ctx context.Context, requestStr string) {
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
	case <-ctx.Done():
		fmt.Println("req too long")
	}

}
