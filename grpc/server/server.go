package server

import (
	"context"
	"ka-cache/cache"
	"log"
)
import pb "ka-cache/grpc"

type CacheServer struct {
	pb.UnimplementedCacheServer
}

func NewServer() *CacheServer {
	s := &CacheServer{}
	return s
}

var c = cache.NewCache(5)

func (s *CacheServer) Put(ctx context.Context, item *pb.Item) (*pb.Response, error) {
	c.Set(item.Key, item.Value)
	log.Print("item: " + item.Key + " - successfully set.")
	return &pb.Response{
		Message: "success",
		Code:    1,
		Data:    "",
	}, nil
}

func (s *CacheServer) Get(ctx context.Context, obj *pb.Object) (*pb.Response, error) {
	var data = c.Get(obj.Key)
	log.Print("item: " + data + " - successfully get.")
	return &pb.Response{
		Message: "success",
		Code:    1,
		Data:    data,
	}, nil
}
