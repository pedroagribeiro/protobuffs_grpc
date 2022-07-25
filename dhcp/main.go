package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	pb "github.com/pedroagribeiro/bbwf/dhcp/proto"
	"github.com/pedroagribeiro/bbwf/dhcp/stats"
	"google.golang.org/grpc"
)

const (
	fileName = "./stats/sqlite.db"
	port = ":50051"
)

var repository *stats.SQLiteRepository

type StatsServer struct {
	pb.UnimplementedSubscriberStatsServiceServer	
}

func (s *StatsServer) GetSubscriberStats(ctx context.Context, in *pb.Empty) (*pb.SubscriberStats, error) {
	log.Printf("Received a gRPC request to send SubscriberStats")
	stats, err := repository.GetFirst()
	if err != nil {
		return nil, err
	}
	return &pb.SubscriberStats{Active: stats.Active, Total: stats.Total}, nil
}

func main() {

	// database related actions
	os.Remove(fileName)
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	statsRepository := stats.NewSQLiteRepository(db)
	repository = statsRepository
	if err := statsRepository.Migrate(); err != nil {
		log.Fatal(err)
	}
	statsSample := stats.Stats{
		ID: 0,
		Active: 5,
		Total: 10,
	}
	_, err = statsRepository.Create(statsSample)
	if err != nil {
		log.Fatal(err)
	}

	// grpc related actions	
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSubscriberStatsServiceServer(s, &StatsServer{})
	log.Printf("server listenning at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}