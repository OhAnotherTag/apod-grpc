package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"

	pb "github.com/OhAnotherTag/apod-grpc/pkg/apod"
)

const (
	port       = ":9000"
	apiKey     = "DEMO_KEY"
	ApodAPIURL = "https://api.nasa.gov/planetary/apod?api_key=" + apiKey
)

type Server struct {
	pb.UnimplementedApodServiceServer
}

type Apod struct {
	Date           string `json:"date"`
	Explanation    string    `json:"explanation"`
	Hdurl          string    `json:"hdurl"`
	MediaType      string    `json:"media_type"`
	ServiceVersion string    `json:"service_version"`
	Title          string    `json:"title"`
	Url            string    `json:"url"`
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterApodServiceServer(grpcServer, &Server{})

	log.Printf("GRPC server listen on %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func validateDate(date string) bool {
	layoutISO := "2006-01-02"
	t, _ := time.Parse(layoutISO, date)
	return t.Format("2006-01-02") == date
}

func (s *Server) GetApod(ctx context.Context, req *pb.ApodRequest) (*pb.ApodReply, error) {
	res := &pb.ApodReply{}

	if req == nil {
		fmt.Println("request must not be nil")
		return res, xerrors.Errorf("request must not be nil")
	}

	if !validateDate(req.Date) {
		fmt.Println("date must met the YYYY-MM-DD format in the request")
		return res, xerrors.Errorf("date must met the YYYY-MM-DD format in the request")
	}

	log.Printf("Receive: %v", req.GetDate())

	response, err := http.Get(ApodAPIURL + "&date=" + req.GetDate())
	if err != nil {
		log.Fatalf("failed to call Apod API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("Couldn't reach Apod API.")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	var data *pb.Apod
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	res.Data = data

	return res, nil
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the schema gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		startGRPCServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
