// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: service_e_commerce.proto

package pb

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
	Ecommerce_CreateUser_FullMethodName    = "/pb.Ecommerce/CreateUser"
	Ecommerce_LoginUser_FullMethodName     = "/pb.Ecommerce/LoginUser"
	Ecommerce_CreateProduct_FullMethodName = "/pb.Ecommerce/CreateProduct"
)

// EcommerceClient is the client API for Ecommerce service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EcommerceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
	CreateProduct(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[CreateProductRequest, CreateProductResponse], error)
}

type ecommerceClient struct {
	cc grpc.ClientConnInterface
}

func NewEcommerceClient(cc grpc.ClientConnInterface) EcommerceClient {
	return &ecommerceClient{cc}
}

func (c *ecommerceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, Ecommerce_CreateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ecommerceClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, Ecommerce_LoginUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ecommerceClient) CreateProduct(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[CreateProductRequest, CreateProductResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Ecommerce_ServiceDesc.Streams[0], Ecommerce_CreateProduct_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[CreateProductRequest, CreateProductResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Ecommerce_CreateProductClient = grpc.ClientStreamingClient[CreateProductRequest, CreateProductResponse]

// EcommerceServer is the server API for Ecommerce service.
// All implementations must embed UnimplementedEcommerceServer
// for forward compatibility.
type EcommerceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	CreateProduct(grpc.ClientStreamingServer[CreateProductRequest, CreateProductResponse]) error
	mustEmbedUnimplementedEcommerceServer()
}

// UnimplementedEcommerceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEcommerceServer struct{}

func (UnimplementedEcommerceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedEcommerceServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedEcommerceServer) CreateProduct(grpc.ClientStreamingServer[CreateProductRequest, CreateProductResponse]) error {
	return status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}
func (UnimplementedEcommerceServer) mustEmbedUnimplementedEcommerceServer() {}
func (UnimplementedEcommerceServer) testEmbeddedByValue()                   {}

// UnsafeEcommerceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EcommerceServer will
// result in compilation errors.
type UnsafeEcommerceServer interface {
	mustEmbedUnimplementedEcommerceServer()
}

func RegisterEcommerceServer(s grpc.ServiceRegistrar, srv EcommerceServer) {
	// If the following call pancis, it indicates UnimplementedEcommerceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Ecommerce_ServiceDesc, srv)
}

func _Ecommerce_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EcommerceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ecommerce_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EcommerceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ecommerce_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EcommerceServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ecommerce_LoginUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EcommerceServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ecommerce_CreateProduct_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EcommerceServer).CreateProduct(&grpc.GenericServerStream[CreateProductRequest, CreateProductResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Ecommerce_CreateProductServer = grpc.ClientStreamingServer[CreateProductRequest, CreateProductResponse]

// Ecommerce_ServiceDesc is the grpc.ServiceDesc for Ecommerce service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ecommerce_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Ecommerce",
	HandlerType: (*EcommerceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _Ecommerce_CreateUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _Ecommerce_LoginUser_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateProduct",
			Handler:       _Ecommerce_CreateProduct_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "service_e_commerce.proto",
}
