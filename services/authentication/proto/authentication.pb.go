// Code generated by protoc-gen-go. DO NOT EDIT.
// source: authentication.proto

package basicblog_authentication

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LoginCommand struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginCommand) Reset()         { *m = LoginCommand{} }
func (m *LoginCommand) String() string { return proto.CompactTextString(m) }
func (*LoginCommand) ProtoMessage()    {}
func (*LoginCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0dbc99083440df2, []int{0}
}

func (m *LoginCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginCommand.Unmarshal(m, b)
}
func (m *LoginCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginCommand.Marshal(b, m, deterministic)
}
func (m *LoginCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginCommand.Merge(m, src)
}
func (m *LoginCommand) XXX_Size() int {
	return xxx_messageInfo_LoginCommand.Size(m)
}
func (m *LoginCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginCommand.DiscardUnknown(m)
}

var xxx_messageInfo_LoginCommand proto.InternalMessageInfo

func (m *LoginCommand) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginCommand) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LoginResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Token                string   `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()    {}
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0dbc99083440df2, []int{1}
}

func (m *LoginResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResponse.Unmarshal(m, b)
}
func (m *LoginResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResponse.Marshal(b, m, deterministic)
}
func (m *LoginResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResponse.Merge(m, src)
}
func (m *LoginResponse) XXX_Size() int {
	return xxx_messageInfo_LoginResponse.Size(m)
}
func (m *LoginResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResponse proto.InternalMessageInfo

func (m *LoginResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *LoginResponse) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type AuthenticateQuery struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthenticateQuery) Reset()         { *m = AuthenticateQuery{} }
func (m *AuthenticateQuery) String() string { return proto.CompactTextString(m) }
func (*AuthenticateQuery) ProtoMessage()    {}
func (*AuthenticateQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0dbc99083440df2, []int{2}
}

func (m *AuthenticateQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthenticateQuery.Unmarshal(m, b)
}
func (m *AuthenticateQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthenticateQuery.Marshal(b, m, deterministic)
}
func (m *AuthenticateQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthenticateQuery.Merge(m, src)
}
func (m *AuthenticateQuery) XXX_Size() int {
	return xxx_messageInfo_AuthenticateQuery.Size(m)
}
func (m *AuthenticateQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthenticateQuery.DiscardUnknown(m)
}

var xxx_messageInfo_AuthenticateQuery proto.InternalMessageInfo

func (m *AuthenticateQuery) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type AuthenticateResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthenticateResponse) Reset()         { *m = AuthenticateResponse{} }
func (m *AuthenticateResponse) String() string { return proto.CompactTextString(m) }
func (*AuthenticateResponse) ProtoMessage()    {}
func (*AuthenticateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0dbc99083440df2, []int{3}
}

func (m *AuthenticateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthenticateResponse.Unmarshal(m, b)
}
func (m *AuthenticateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthenticateResponse.Marshal(b, m, deterministic)
}
func (m *AuthenticateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthenticateResponse.Merge(m, src)
}
func (m *AuthenticateResponse) XXX_Size() int {
	return xxx_messageInfo_AuthenticateResponse.Size(m)
}
func (m *AuthenticateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthenticateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AuthenticateResponse proto.InternalMessageInfo

func (m *AuthenticateResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterType((*LoginCommand)(nil), "basicblog.authentication.LoginCommand")
	proto.RegisterType((*LoginResponse)(nil), "basicblog.authentication.LoginResponse")
	proto.RegisterType((*AuthenticateQuery)(nil), "basicblog.authentication.AuthenticateQuery")
	proto.RegisterType((*AuthenticateResponse)(nil), "basicblog.authentication.AuthenticateResponse")
}

func init() { proto.RegisterFile("authentication.proto", fileDescriptor_d0dbc99083440df2) }

var fileDescriptor_d0dbc99083440df2 = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x49, 0x2c, 0x2d, 0xc9,
	0x48, 0xcd, 0x2b, 0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x92, 0x48, 0x4a, 0x2c, 0xce, 0x4c, 0x4e, 0xca, 0xc9, 0x4f, 0xd7, 0x43, 0x95, 0x57, 0x72,
	0xe3, 0xe2, 0xf1, 0xc9, 0x4f, 0xcf, 0xcc, 0x73, 0xce, 0xcf, 0xcd, 0x4d, 0xcc, 0x4b, 0x11, 0x92,
	0xe2, 0xe2, 0x28, 0x2d, 0x4e, 0x2d, 0xca, 0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0,
	0x0c, 0x82, 0xf3, 0x41, 0x72, 0x05, 0x89, 0xc5, 0xc5, 0xe5, 0xf9, 0x45, 0x29, 0x12, 0x4c, 0x10,
	0x39, 0x18, 0x5f, 0xc9, 0x9e, 0x8b, 0x17, 0x6c, 0x4e, 0x50, 0x6a, 0x71, 0x41, 0x7e, 0x5e, 0x71,
	0xaa, 0x90, 0x04, 0x17, 0x7b, 0x71, 0x69, 0x72, 0x72, 0x6a, 0x71, 0x31, 0xd8, 0x1c, 0x8e, 0x20,
	0x18, 0x57, 0x48, 0x84, 0x8b, 0xb5, 0x24, 0x3f, 0x3b, 0x35, 0x0f, 0x6a, 0x06, 0x84, 0xa3, 0xa4,
	0xc9, 0x25, 0xe8, 0x88, 0x70, 0x5a, 0x6a, 0x60, 0x69, 0x6a, 0x51, 0x25, 0x42, 0x29, 0x23, 0xb2,
	0x52, 0x03, 0x2e, 0x11, 0x64, 0xa5, 0x84, 0xad, 0x34, 0xba, 0xcd, 0xc8, 0xc5, 0xe7, 0x88, 0xe2,
	0x71, 0xa1, 0x28, 0x2e, 0x56, 0xb0, 0x83, 0x85, 0xd4, 0xf4, 0x70, 0x05, 0x8e, 0x1e, 0x72, 0xc8,
	0x48, 0xa9, 0x13, 0x50, 0x07, 0x73, 0x86, 0x12, 0x83, 0x50, 0x2e, 0x17, 0x0f, 0xb2, 0x03, 0x85,
	0xb4, 0x71, 0x6b, 0xc5, 0xf0, 0xb3, 0x94, 0x1e, 0x71, 0x8a, 0x11, 0xd6, 0x25, 0xb1, 0x81, 0x23,
	0xd9, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xda, 0x43, 0x2b, 0x10, 0xfc, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AuthenticationClient is the client API for Authentication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthenticationClient interface {
	// Logs user in
	Login(ctx context.Context, in *LoginCommand, opts ...grpc.CallOption) (*LoginResponse, error)
	// Authenticates user token
	Authenticate(ctx context.Context, in *AuthenticateQuery, opts ...grpc.CallOption) (*AuthenticateResponse, error)
}

type authenticationClient struct {
	cc *grpc.ClientConn
}

func NewAuthenticationClient(cc *grpc.ClientConn) AuthenticationClient {
	return &authenticationClient{cc}
}

func (c *authenticationClient) Login(ctx context.Context, in *LoginCommand, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/basicblog.authentication.Authentication/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationClient) Authenticate(ctx context.Context, in *AuthenticateQuery, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := c.cc.Invoke(ctx, "/basicblog.authentication.Authentication/Authenticate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticationServer is the server API for Authentication service.
type AuthenticationServer interface {
	// Logs user in
	Login(context.Context, *LoginCommand) (*LoginResponse, error)
	// Authenticates user token
	Authenticate(context.Context, *AuthenticateQuery) (*AuthenticateResponse, error)
}

// UnimplementedAuthenticationServer can be embedded to have forward compatible implementations.
type UnimplementedAuthenticationServer struct {
}

func (*UnimplementedAuthenticationServer) Login(ctx context.Context, req *LoginCommand) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (*UnimplementedAuthenticationServer) Authenticate(ctx context.Context, req *AuthenticateQuery) (*AuthenticateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}

func RegisterAuthenticationServer(s *grpc.Server, srv AuthenticationServer) {
	s.RegisterService(&_Authentication_serviceDesc, srv)
}

func _Authentication_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/basicblog.authentication.Authentication/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).Login(ctx, req.(*LoginCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authentication_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/basicblog.authentication.Authentication/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).Authenticate(ctx, req.(*AuthenticateQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _Authentication_serviceDesc = grpc.ServiceDesc{
	ServiceName: "basicblog.authentication.Authentication",
	HandlerType: (*AuthenticationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Authentication_Login_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _Authentication_Authenticate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authentication.proto",
}
