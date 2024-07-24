package database

import (
	"context"
	"sync"

	pb "github.com/sarthak0714/dbz/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Database struct {
	pb.UnimplementedDatabaseServer
	data  map[string]string
	locks map[string]*sync.Mutex
	mutex sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		data:  make(map[string]string),
		locks: make(map[string]*sync.Mutex),
	}
}

func (d *Database) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	value, ok := d.data[req.Key]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "key not found")
	}

	return &pb.GetResponse{Value: value}, nil
}

func (d *Database) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.data[req.Key] = req.Value
	return &pb.PutResponse{Success: true}, nil
}

func (d *Database) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.data, req.Key)
	return &pb.DeleteResponse{Success: true}, nil
}

func (d *Database) getLock(key string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, ok := d.locks[key]; !ok {
		d.locks[key] = &sync.Mutex{}
	}
	return d.locks[key]
}
