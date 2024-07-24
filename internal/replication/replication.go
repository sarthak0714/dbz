package replication

import (
	"io"

	"github.com/sarthak0714/dbz/internal/database"
	pb "github.com/sarthak0714/dbz/pkg/api"
)

type ReplicationServer struct {
	pb.UnimplementedReplicationServer
	db *database.Database
}

func NewReplicationServer(db *database.Database) *ReplicationServer {
	return &ReplicationServer{db: db}
}

func (s *ReplicationServer) Replicate(stream pb.Replication_ReplicateServer) error {
	for {
		update, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Apply the update locally
		_, err = s.db.Put(stream.Context(), &pb.PutRequest{Key: update.Key, Value: update.Value})
		if err != nil {
			return err
		}

		// Acknowledge the update
		if err := stream.Send(&pb.ReplicationAck{Success: true}); err != nil {
			return err
		}
	}
}
