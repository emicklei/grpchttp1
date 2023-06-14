package grpchttp1

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Assert *HTTPCallHandler implements ServiceRegistrar.
var _ grpc.ServiceRegistrar = (*HTTPCallHandler)(nil)

// Assert *HTTPCallHandler implements http.Handler.
var _ http.Handler = (*HTTPCallHandler)(nil)

type serviceImpl struct {
	desc *grpc.ServiceDesc
	impl any
}

type HTTPCallHandler struct {
	serviceImpls map[string]serviceImpl
}

func (s *HTTPCallHandler) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.serviceImpls[desc.ServiceName] = serviceImpl{desc: desc, impl: impl}
}

func NewHTTPHandlerWithRegistrar() *HTTPCallHandler {
	return &HTTPCallHandler{
		serviceImpls: make(map[string]serviceImpl),
	}
}

func (s *HTTPCallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	svc := parts[1]
	method := parts[2]
	reg, ok := s.serviceImpls[svc]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "no service found")
		return
	}
	for _, each := range reg.desc.Methods {
		if each.MethodName == method {
			ctx := r.Context()
			if mdv := r.Header.Get(httpHeaderMetadataKey); mdv != "" {
				mdb, err := base64.StdEncoding.DecodeString(mdv)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					io.WriteString(w, "unable to base64 decode metadata")
					return
				}
				md := metadata.MD{}
				err = json.Unmarshal(mdb, &md)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					io.WriteString(w, "unable to json decode metadata")
					return
				}
				ctx = metadata.NewIncomingContext(ctx, md)
			}
			// use generated handler
			output, err := each.Handler(reg.impl, ctx, func(input any) error {
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return fmt.Errorf("unable to read request body: %v", err)
				}
				protoInput, ok := input.(proto.Message)
				if !ok {
					return errors.New("expected input proto.Message")
				}
				if err := proto.Unmarshal(data, protoInput); err != nil {
					return fmt.Errorf("unable to proto unmarshal request body: %v", err)
				}
				return nil
			}, nil)
			if err != nil {
				handleError(w, err)
				return
			}
			// write output on the wire
			protoOutput, ok := output.(proto.Message)
			if !ok {
				handleError(w, fmt.Errorf("expected output proto.Message"))
				return
			}
			data, err := proto.Marshal(protoOutput)
			if err != nil {
				handleError(w, fmt.Errorf("unable to proto marshal output for response body"))
				return
			}
			w.Write(data)
			return
		}
	}
	handleError(w, status.Error(codes.NotFound, "no service or method found"))
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	stat := status.Convert(err).Proto()
	data, _ := proto.Marshal(stat)
	w.Write(data)
}
