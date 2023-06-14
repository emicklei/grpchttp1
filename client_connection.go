package grpchttp1

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const httpHeaderMetadataKey = "grpchttp1-md-json-base64"

// Assert *HttpClientConn implements ClientConnInterface.
var _ grpc.ClientConnInterface = (*HttpClientConn)(nil)

type HttpClientConn struct {
	client      *http.Client
	endpointURL string
}

func NewClientConn(client *http.Client, endpoint string) *HttpClientConn {
	return &HttpClientConn{client: client, endpointURL: endpoint}
}

// Invoke performs a unary RPC and returns after the response is received
// into reply.
func (c *HttpClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	body := new(bytes.Buffer)
	msg, ok := args.(proto.Message)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "argument not a proto.Message")
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "proto marshal failed")
	}
	body.Write(data)
	req, err := http.NewRequest("POST", c.endpointURL+method, body)
	if err != nil {
		return status.Error(codes.InvalidArgument, "endpoint not valid")
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		data, err := json.Marshal(md)
		if err != nil {
			return status.Errorf(codes.Internal, "metadata could not be serialized:%v", err)
		}
		req.Header.Add(httpHeaderMetadataKey, base64.StdEncoding.EncodeToString(data))
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return status.Error(codes.Unavailable, "HTTP call failed")
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return status.Error(codes.Unknown, "reading HTTP response body failed")
	}
	if resp.StatusCode != http.StatusOK {
		if resp.Header.Get("Content-Type") == "application/protobuf" {
			// unmarshal into proto version of status.Status
			stat := new(spb.Status)
			err = proto.Unmarshal(data, stat)
			if err != nil {
				return status.Error(codes.Unknown, "HTTP response handling failed (+ not a gRPC status)")
			}
			return status.FromProto(stat).Err()
		}
		// TODO map HTTP codes to gRPC codes?
		return status.Error(codes.Unknown, "HTTP response handling failed (+ not a protobuf response)")
	} else {
		err = proto.Unmarshal(data, reply.(protoreflect.ProtoMessage))
		if err != nil {
			return status.Error(codes.Unknown, "HTTP response protobuf unmarshal failed")
		}
	}
	return nil
}

// NewStream begins a streaming RPC.
func (c *HttpClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "streaming not (yet) supported")
}
