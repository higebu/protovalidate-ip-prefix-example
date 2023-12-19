package server

import (
	"context"
	"log"
	"net"
	"testing"

	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/bufbuild/protovalidate-go"
	examplev1 "github.com/higebu/protovalidate-ip-prefix-example/gen/go/proto/example/v1"
)

var l *bufconn.Listener

func init() {
	l = bufconn.Listen(1024 * 1024)
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal(err)
	}
	s := New()
	g := grpc.NewServer(grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)))
	examplev1.RegisterExampleServiceServer(g, s)
	go func() {
		if err := g.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return l.Dial()
}

func TestAddIPv4Address(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := examplev1.NewExampleServiceClient(conn)
	resp, err := client.AddIPv4Address(ctx, &examplev1.AddIPv4AddressRequest{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	if resp.GetName() != "test" {
		t.Fatalf("got %s, want %s", resp.GetName(), "test")
	}
}
