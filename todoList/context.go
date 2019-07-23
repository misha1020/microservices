package main

import (
	"context"
	"errors"
	"log"
	"time"
)

func TooLong(ctx context.Context) (int, error){
	select {
	case <- time.After(5 *time.Second):
		return 15, nil
	case <- ctx.Done():
		return 0, errors.New("was waititng for too long")
	}
}

func main() {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(2 * time.Second))
	//time.AfterFunc(3 * time.Second, cancel)
	_, err := TooLong(ctx)

	if err != nil {
		log.Print(err)
	}
}
