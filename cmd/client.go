package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/OhAnotherTag/apod-grpc/pkg/apod"
)

const (
	addr = "localhost:9000"
)

var (
	defaultDate = time.Now().Format("2006-01-02")
)
// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Query the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewApodServiceClient(conn)

		fmt.Println(args)

		date := defaultDate
		if len(os.Args) > 2 {
			date = os.Args[2]
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		r, err := client.GetApod(ctx, &pb.ApodRequest{Date: date})
		if err != nil {
			log.Fatalf("server responded with error: %v", err)
		}

		log.Println(r.GetData())
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
