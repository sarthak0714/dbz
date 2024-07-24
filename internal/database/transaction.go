package database

import (
	"context"
	"sync"

	pb "github.com/sarthak0714/dbz/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Transaction struct {
	db    *Database
	locks []*sync.Mutex
}

func (d *Database) PrepareTransaction(ctx context.Context, req *pb.PrepareRequest) (*pb.PrepareResponse, error) {
	t := &Transaction{
		db:    d,
		locks: make([]*sync.Mutex, 0, len(req.Keys)),
	}

	for _, key := range req.Keys {
		lock := d.getLock(key)
		if !lock.TryLock() {
			// If we can't acquire all locks, release the ones we've acquired
			for _, acquiredLock := range t.locks {
				acquiredLock.Unlock()
			}
			return &pb.PrepareResponse{Ready: false}, status.Errorf(codes.Aborted, "key %s is locked", key)
		}
		t.locks = append(t.locks, lock)
	}

	return &pb.PrepareResponse{Ready: true}, nil
}

func (d *Database) CommitTransaction(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Apply updates
	for key, value := range req.Updates {
		d.data[key] = value
	}

	// Release locks
	for _, lock := range d.locks {
		lock.Unlock()
	}

	return &pb.CommitResponse{Success: true}, nil
}

func (d *Database) AbortTransaction(ctx context.Context, req *pb.AbortRequest) (*pb.AbortResponse, error) {
	// Release locks
	for _, key := range req.Keys {
		lock := d.getLock(key)
		lock.Unlock()
	}

	return &pb.AbortResponse{Success: true}, nil
}
