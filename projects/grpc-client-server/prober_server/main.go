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
	numOfErrors := 0
	start := time.Now()

	timeOut := 4 * time.Second
	if in.TimeOutMsecs != nil {
		timeOut = time.Duration(in.GetTimeOutMsecs()) * time.Millisecond
	}
	client := http.Client{
		Timeout: timeOut,
	}
	fmt.Printf("Time out : %v\n", timeOut)

	result := pb.ProbeReply{StatusCodes: make(map[int64]int64)}

	for i := 0; i < int(numOFRequests); i++ {

		resp, err := client.Get(in.GetEndpoint())
		if err != nil {
			numOfErrors++
			fmt.Printf("could not probe: %+v", err)
		} else {
			defer resp.Body.Close()
			code := int64(resp.StatusCode)
			result.StatusCodes[code]++
		}
	}

	elapsed := time.Since(start)
	elapsedMsecs := float32(elapsed / time.Millisecond)
	fmt.Printf("Elapsed time: %f\n", elapsedMsecs)

	result.AvgLatencyMsecs = elapsedMsecs / float32(numOFRequests)
	result.PercentageErrors = (float32(numOfErrors) / float32(numOFRequests)) * float32(100)
	return &result, nil
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
