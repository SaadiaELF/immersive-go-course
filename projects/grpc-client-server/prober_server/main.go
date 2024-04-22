package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement prober.ProberServer.
type server struct {
	pb.UnimplementedProberServer
}

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	numOFRequests := in.GetNumOfRequests()
	start := time.Now()

	for i := 0; i < int(numOFRequests); i++ {
		resp, err := http.Get(in.GetEndpoint())
		if err != nil {
			log.Fatalf("could not probe: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %v", resp.Status)
		}
	}

	elapsed := time.Since(start)
	elapsedMsecs := float32(elapsed / time.Millisecond)
	fmt.Printf("Elapsed time: %f\n", elapsedMsecs)

	return &pb.ProbeReply{AvgLatencyMsecs: elapsedMsecs / float32(numOFRequests)}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProberServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
