package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/pedroagribeiro/bbwf/rest/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

type SubscriberStats struct {
	Count	CounterStats	`json:"count"`	
}

type CounterStats struct {
	Active	int32	`json:"active"`
	Total	int32	`json:"total"`
}

var pb_client pb.SubscriberStatsServiceClient
var contxt context.Context

func getSubscriberStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	stats, err := pb_client.GetSubscriberStats(contxt, &pb.Empty{})
	if err != nil {
		log.Fatalf("gRPC call was unsuccessful: %v", err)
	}
	c := CounterStats{}
	c.Active = stats.Active
	c.Total = stats.Total
	s := SubscriberStats{}
	s.Count = c
	if err := json.NewEncoder(w).Encode(s); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

func registerRouter() {
	r := mux.NewRouter()
	r.Path("/api/v1/subscriberstats").Methods(http.MethodGet).HandlerFunc(getSubscriberStats)
	fmt.Println("Started listenning on localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", r))
}

func main() {

	// grpc related actions
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}	
	defer conn.Close()
	c := pb.NewSubscriberStatsServiceClient(conn)
	pb_client = c
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	contxt = ctx
	defer cancel()

	// rest related
	registerRouter()
}
