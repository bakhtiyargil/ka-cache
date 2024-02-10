package server

import (
	"context"
	"ka-cache/cache"
	err "ka-cache/http/error"
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
	log.Print("item: " + item.Key + " - successfully set")
	return &pb.Response{
		Message: "success",
		Code:    1,
		Data:    "",
	}, nil
}

func (s *CacheServer) Get(ctx context.Context, obj *pb.Object) (*pb.Response, error) {
	var item = c.Get(obj.Key)
	if item == "" {
		return nil, err.ResourceNotFoundError
	}
	log.Print("item: " + item + " - successfully get")
	return &pb.Response{
		Message: "success",
		Code:    1,
		Data:    item,
	}, nil
}
