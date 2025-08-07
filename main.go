package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpc2 "ka-cache/server/grpc"
	"log"
	"time"
)

func main() {
	var optsDial []grpc.DialOption

	optsDial = append(optsDial, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:3000", optsDial...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := grpc2.NewCacheClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.Put(ctx, &grpc2.Item{
		Key:   "Isaq",
		Value: "A",
	})

	fmt.Println(response)
}
