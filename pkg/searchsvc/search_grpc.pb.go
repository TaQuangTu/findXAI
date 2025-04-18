// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: api/search.proto

package searchsvc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SearchService_Search_FullMethodName         = "/google.search.v1.SearchService/Search"
	SearchService_DeactivateKeys_FullMethodName = "/google.search.v1.SearchService/DeactivateKeys"
	SearchService_ActivateKeys_FullMethodName   = "/google.search.v1.SearchService/ActivateKeys"
	SearchService_AddKeys_FullMethodName        = "/google.search.v1.SearchService/AddKeys"
	SearchService_GetKeys_FullMethodName        = "/google.search.v1.SearchService/GetKeys"
)

// SearchServiceClient is the client API for SearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchServiceClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	DeactivateKeys(ctx context.Context, in *DeactivateKeysRequest, opts ...grpc.CallOption) (*DeactivateKeysResponse, error)
	ActivateKeys(ctx context.Context, in *ActivateKeysRequest, opts ...grpc.CallOption) (*ActivateKeysResponse, error)
	AddKeys(ctx context.Context, in *AddKeysRequest, opts ...grpc.CallOption) (*AddKeysResponse, error)
	GetKeys(ctx context.Context, in *GetKeysRequest, opts ...grpc.CallOption) (*GetKeysResponse, error)
}

type searchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchServiceClient(cc grpc.ClientConnInterface) SearchServiceClient {
	return &searchServiceClient{cc}
}

func (c *searchServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, SearchService_Search_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) DeactivateKeys(ctx context.Context, in *DeactivateKeysRequest, opts ...grpc.CallOption) (*DeactivateKeysResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeactivateKeysResponse)
	err := c.cc.Invoke(ctx, SearchService_DeactivateKeys_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ActivateKeys(ctx context.Context, in *ActivateKeysRequest, opts ...grpc.CallOption) (*ActivateKeysResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ActivateKeysResponse)
	err := c.cc.Invoke(ctx, SearchService_ActivateKeys_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) AddKeys(ctx context.Context, in *AddKeysRequest, opts ...grpc.CallOption) (*AddKeysResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddKeysResponse)
	err := c.cc.Invoke(ctx, SearchService_AddKeys_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetKeys(ctx context.Context, in *GetKeysRequest, opts ...grpc.CallOption) (*GetKeysResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKeysResponse)
	err := c.cc.Invoke(ctx, SearchService_GetKeys_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServiceServer is the server API for SearchService service.
// All implementations must embed UnimplementedSearchServiceServer
// for forward compatibility.
type SearchServiceServer interface {
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	DeactivateKeys(context.Context, *DeactivateKeysRequest) (*DeactivateKeysResponse, error)
	ActivateKeys(context.Context, *ActivateKeysRequest) (*ActivateKeysResponse, error)
	AddKeys(context.Context, *AddKeysRequest) (*AddKeysResponse, error)
	GetKeys(context.Context, *GetKeysRequest) (*GetKeysResponse, error)
	mustEmbedUnimplementedSearchServiceServer()
}

// UnimplementedSearchServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSearchServiceServer struct{}

func (UnimplementedSearchServiceServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedSearchServiceServer) DeactivateKeys(context.Context, *DeactivateKeysRequest) (*DeactivateKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateKeys not implemented")
}
func (UnimplementedSearchServiceServer) ActivateKeys(context.Context, *ActivateKeysRequest) (*ActivateKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateKeys not implemented")
}
func (UnimplementedSearchServiceServer) AddKeys(context.Context, *AddKeysRequest) (*AddKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddKeys not implemented")
}
func (UnimplementedSearchServiceServer) GetKeys(context.Context, *GetKeysRequest) (*GetKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeys not implemented")
}
func (UnimplementedSearchServiceServer) mustEmbedUnimplementedSearchServiceServer() {}
func (UnimplementedSearchServiceServer) testEmbeddedByValue()                       {}

// UnsafeSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServiceServer will
// result in compilation errors.
type UnsafeSearchServiceServer interface {
	mustEmbedUnimplementedSearchServiceServer()
}

func RegisterSearchServiceServer(s grpc.ServiceRegistrar, srv SearchServiceServer) {
	// If the following call pancis, it indicates UnimplementedSearchServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SearchService_ServiceDesc, srv)
}

func _SearchService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_DeactivateKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeactivateKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).DeactivateKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_DeactivateKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).DeactivateKeys(ctx, req.(*DeactivateKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_ActivateKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).ActivateKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_ActivateKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ActivateKeys(ctx, req.(*ActivateKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_AddKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).AddKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_AddKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).AddKeys(ctx, req.(*AddKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_GetKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetKeys(ctx, req.(*GetKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SearchService_ServiceDesc is the grpc.ServiceDesc for SearchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "google.search.v1.SearchService",
	HandlerType: (*SearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _SearchService_Search_Handler,
		},
		{
			MethodName: "DeactivateKeys",
			Handler:    _SearchService_DeactivateKeys_Handler,
		},
		{
			MethodName: "ActivateKeys",
			Handler:    _SearchService_ActivateKeys_Handler,
		},
		{
			MethodName: "AddKeys",
			Handler:    _SearchService_AddKeys_Handler,
		},
		{
			MethodName: "GetKeys",
			Handler:    _SearchService_GetKeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/search.proto",
}
