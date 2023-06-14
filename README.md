# grpchttp1

Offers a gRPC flow that uses HTTP 1 for transporting and dispatching Unary gRPC calls.
This package was created to have support for gRPC when implementing GCP Cloud Functions that, at the time of writing, do not support HTTP 2. 

Note that the client needs to be aware of the fact that the gRPC function requires HTTP 1 transport.

## server(less) side

This handler implements http.Handler (ServeHTTP)

    httpHandler := grpchttp1.NewHTTPHandlerWithRegistrar()

use it to register your gRPC service implementation

    svc := new(YourServiceImpl)
    RegisterYourServiceServer(httpHandler, svc)

and use it as a regular HTTP handler in your server

    err := http.ListenAndServe(":8080", httpHandler)

## client side

This connection implements grpc.ClientConnInterface

    conn := grpchttp1.NewClientConn(http.DefaultClient, "https://your.service.endpoint.net")

use it to create your gRPC service client

    client := NewYourServiceClient(conn)

## grpc error handling

Grpc errors are transported in binary using the `google/rpc/status.proto` definition.
So if a service call returns a gRPC error then the error is marshaled using protobuf and transported to the client with HTTP status 500 (InternalServerError). 
On the client, this error is unmarshaled into a Status message and return to the caller as a gRPC error.

### (un) supported features

- Unary calls
- gRPC error codes
- Streaming calls not implemented
- grpc metadata (transported as http header)
- Client- and Server Interceptors support not implemented