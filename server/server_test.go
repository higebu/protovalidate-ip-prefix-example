package server

import (
	"context"
	"log"
	"net"
	"testing"

	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/stretchr/testify/require"
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

	valErr := &protovalidate.ValidationError{}
	tests := []struct {
		name string
		addr string
		want error
	}{
		{name: "valid", addr: "192.168.100.5/24", want: nil},
		{name: "invalid_without_prefix_length", addr: "192.168.100.5", want: valErr},
		{name: "invalid_ipv6", addr: "2001::1/64", want: valErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.AddIPv4Address(ctx, &examplev1.AddIPv4AddressRequest{Name: tt.name, Address: tt.addr})
			if tt.want == nil {
				require.NoError(t, err)
			} else {
				t.Logf("err: %v", err)
				require.ErrorAs(t, err, &tt.want)
			}
		})
	}
}

func TestAddIPv6Address(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := examplev1.NewExampleServiceClient(conn)

	valErr := &protovalidate.ValidationError{}
	tests := []struct {
		name string
		addr string
		want error
	}{
		{name: "valid", addr: "2001::1/64", want: nil},
		{name: "invalid_without_prefix_length", addr: "2001::1", want: valErr},
		{name: "invalid_ipv4", addr: "192.168.100.5/24", want: valErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.AddIPv6Address(ctx, &examplev1.AddIPv6AddressRequest{Name: tt.name, Address: tt.addr})
			if tt.want == nil {
				require.NoError(t, err)
			} else {
				t.Logf("err: %v", err)
				require.ErrorAs(t, err, &tt.want)
			}
		})
	}
}

func TestAddIPv4Route(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := examplev1.NewExampleServiceClient(conn)

	valErr := &protovalidate.ValidationError{}
	tests := []struct {
		name    string
		prefix  string
		nexthop string
		want    error
	}{
		{name: "valid", prefix: "192.168.100.0/24", nexthop: "192.168.100.1", want: nil},
		{name: "invalid_prefix", prefix: "192.168.100.5/24", nexthop: "192.168.100.1", want: valErr},
		{name: "invalid_ipv6", prefix: "2001::/64", nexthop: "2001::1", want: valErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.AddIPv4Route(ctx, &examplev1.AddIPv4RouteRequest{Name: tt.name, Prefix: tt.prefix, Nexthop: tt.nexthop})
			if tt.want == nil {
				require.NoError(t, err)
			} else {
				t.Logf("err: %v", err)
				require.ErrorAs(t, err, &tt.want)
			}
		})
	}
}

func TestAddIPv6Route(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := examplev1.NewExampleServiceClient(conn)

	valErr := &protovalidate.ValidationError{}
	tests := []struct {
		name    string
		prefix  string
		nexthop string
		want    error
	}{
		{name: "valid", prefix: "2001::/64", nexthop: "2001::1", want: nil},
		{name: "invalid_prefix", prefix: "2001::1/64", nexthop: "2001::1", want: valErr},
		{name: "invalid_ipv4", prefix: "192.168.100.0/24", nexthop: "192.168.100.1", want: valErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.AddIPv6Route(ctx, &examplev1.AddIPv6RouteRequest{Name: tt.name, Prefix: tt.prefix, Nexthop: tt.nexthop})
			if tt.want == nil {
				require.NoError(t, err)
			} else {
				t.Logf("err: %v", err)
				require.ErrorAs(t, err, &tt.want)
			}
		})
	}
}
