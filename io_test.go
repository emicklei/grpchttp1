package grpchttp1

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/grpchttp1/testpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestSimple(t *testing.T) {

	// serverless side for HTTP 1.1 transport
	hndl := NewHTTPHandlerWithRegistrar()

	// normal grpc service registration
	svc := new(IOServiceImpl)
	testpb.RegisterIOServiceServer(hndl, svc)

	// run test http server
	httpSrv := httptest.NewServer(hndl)
	defer httpSrv.Close()

	// client connection for HTTP 1.1 transport
	cc := NewClientConn(http.DefaultClient, httpSrv.URL)

	// normal grpc client setup and call
	req := new(testpb.Input)
	req.Name = "test"
	client := testpb.NewIOServiceClient(cc)
	md := metadata.MD{}
	md.Append("key", "value")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	out, err := client.Call(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("  result:", out.Result)
	// simulate error
	req.Fail = true
	out, err = client.Call(ctx, req)
	t.Log(out, err)
}

type IOServiceImpl struct {
	testpb.UnimplementedIOServiceServer
}

func (s *IOServiceImpl) Call(ctx context.Context, in *testpb.Input) (*testpb.Output, error) {
	if in.Fail {
		return nil, status.Error(codes.Internal, "fail request from IOServiceImpl")
	}
	vs := metadata.ValueFromIncomingContext(ctx, "key")
	log.Println("meta data values key=", vs)
	return &testpb.Output{Result: strings.ToUpper(in.Name)}, nil
}

func TestExampleNewHTTPHandlerWithRegistrar(t *testing.T) {
	go func() {
		err := http.ListenAndServe(":8080", NewHTTPHandlerWithRegistrar())
		t.Log(err)
	}()
}
