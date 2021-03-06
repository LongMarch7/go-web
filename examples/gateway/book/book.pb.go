// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/github.com/LongMarch7/go-web/examples/gateway/book/book.proto

package book

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

// 请求书详情的参数结构  book_id 32位整形
type BookInfoParams struct {
	BookId               int32    `protobuf:"varint,1,opt,name=book_id,json=bookId,proto3" json:"book_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BookInfoParams) Reset()         { *m = BookInfoParams{} }
func (m *BookInfoParams) String() string { return proto.CompactTextString(m) }
func (*BookInfoParams) ProtoMessage()    {}
func (*BookInfoParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_b052b9cede179df0, []int{0}
}

func (m *BookInfoParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BookInfoParams.Unmarshal(m, b)
}
func (m *BookInfoParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BookInfoParams.Marshal(b, m, deterministic)
}
func (m *BookInfoParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookInfoParams.Merge(m, src)
}
func (m *BookInfoParams) XXX_Size() int {
	return xxx_messageInfo_BookInfoParams.Size(m)
}
func (m *BookInfoParams) XXX_DiscardUnknown() {
	xxx_messageInfo_BookInfoParams.DiscardUnknown(m)
}

var xxx_messageInfo_BookInfoParams proto.InternalMessageInfo

func (m *BookInfoParams) GetBookId() int32 {
	if m != nil {
		return m.BookId
	}
	return 0
}

// 书详情信息的结构   book_name字符串类型
type BookInfo struct {
	BookId               int32    `protobuf:"varint,1,opt,name=book_id,json=bookId,proto3" json:"book_id,omitempty"`
	BookName             string   `protobuf:"bytes,2,opt,name=book_name,json=bookName,proto3" json:"book_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BookInfo) Reset()         { *m = BookInfo{} }
func (m *BookInfo) String() string { return proto.CompactTextString(m) }
func (*BookInfo) ProtoMessage()    {}
func (*BookInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_b052b9cede179df0, []int{1}
}

func (m *BookInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BookInfo.Unmarshal(m, b)
}
func (m *BookInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BookInfo.Marshal(b, m, deterministic)
}
func (m *BookInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookInfo.Merge(m, src)
}
func (m *BookInfo) XXX_Size() int {
	return xxx_messageInfo_BookInfo.Size(m)
}
func (m *BookInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_BookInfo.DiscardUnknown(m)
}

var xxx_messageInfo_BookInfo proto.InternalMessageInfo

func (m *BookInfo) GetBookId() int32 {
	if m != nil {
		return m.BookId
	}
	return 0
}

func (m *BookInfo) GetBookName() string {
	if m != nil {
		return m.BookName
	}
	return ""
}

// 请求书列表的参数结构  page、limit   32位整形
type BookListParams struct {
	Page                 int32    `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit                int32    `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BookListParams) Reset()         { *m = BookListParams{} }
func (m *BookListParams) String() string { return proto.CompactTextString(m) }
func (*BookListParams) ProtoMessage()    {}
func (*BookListParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_b052b9cede179df0, []int{2}
}

func (m *BookListParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BookListParams.Unmarshal(m, b)
}
func (m *BookListParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BookListParams.Marshal(b, m, deterministic)
}
func (m *BookListParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookListParams.Merge(m, src)
}
func (m *BookListParams) XXX_Size() int {
	return xxx_messageInfo_BookListParams.Size(m)
}
func (m *BookListParams) XXX_DiscardUnknown() {
	xxx_messageInfo_BookListParams.DiscardUnknown(m)
}

var xxx_messageInfo_BookListParams proto.InternalMessageInfo

func (m *BookListParams) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *BookListParams) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

// 书列表的结构    BookInfo结构数组
type BookList struct {
	BookList             []*BookInfo `protobuf:"bytes,1,rep,name=book_list,json=bookList,proto3" json:"book_list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *BookList) Reset()         { *m = BookList{} }
func (m *BookList) String() string { return proto.CompactTextString(m) }
func (*BookList) ProtoMessage()    {}
func (*BookList) Descriptor() ([]byte, []int) {
	return fileDescriptor_b052b9cede179df0, []int{3}
}

func (m *BookList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BookList.Unmarshal(m, b)
}
func (m *BookList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BookList.Marshal(b, m, deterministic)
}
func (m *BookList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookList.Merge(m, src)
}
func (m *BookList) XXX_Size() int {
	return xxx_messageInfo_BookList.Size(m)
}
func (m *BookList) XXX_DiscardUnknown() {
	xxx_messageInfo_BookList.DiscardUnknown(m)
}

var xxx_messageInfo_BookList proto.InternalMessageInfo

func (m *BookList) GetBookList() []*BookInfo {
	if m != nil {
		return m.BookList
	}
	return nil
}

func init() {
	proto.RegisterType((*BookInfoParams)(nil), "book.BookInfoParams")
	proto.RegisterType((*BookInfo)(nil), "book.BookInfo")
	proto.RegisterType((*BookListParams)(nil), "book.BookListParams")
	proto.RegisterType((*BookList)(nil), "book.BookList")
}

func init() {
	proto.RegisterFile("src/github.com/LongMarch7/go-web/examples/gateway/book/book.proto", fileDescriptor_b052b9cede179df0)
}

var fileDescriptor_b052b9cede179df0 = []byte{
	// 280 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x5f, 0x4b, 0xf3, 0x30,
	0x14, 0xc6, 0xd7, 0xf7, 0x5d, 0xe7, 0x96, 0xc2, 0x2e, 0xc2, 0xc0, 0xa2, 0x37, 0x25, 0x57, 0x15,
	0xb1, 0x81, 0x89, 0x0c, 0xbc, 0x52, 0x6f, 0x64, 0x30, 0x45, 0xea, 0x07, 0x90, 0xb4, 0x3b, 0x76,
	0x61, 0x4d, 0x4f, 0x69, 0xa3, 0x53, 0xfc, 0xf2, 0x92, 0xa4, 0xd3, 0x2a, 0x78, 0x13, 0xce, 0xbf,
	0xdf, 0x79, 0xce, 0x13, 0x72, 0xdd, 0x36, 0x39, 0x2f, 0xa4, 0xde, 0xbc, 0x64, 0x49, 0x8e, 0x8a,
	0xaf, 0xb0, 0x2a, 0xee, 0x44, 0x93, 0x6f, 0x16, 0xbc, 0xc0, 0xb3, 0x1d, 0x64, 0x1c, 0xde, 0x84,
	0xaa, 0x4b, 0x68, 0x79, 0x21, 0x34, 0xec, 0xc4, 0x3b, 0xcf, 0x10, 0xb7, 0xf6, 0x49, 0xea, 0x06,
	0x35, 0xd2, 0xa1, 0x89, 0xd9, 0x09, 0x99, 0xde, 0x20, 0x6e, 0x97, 0xd5, 0x33, 0x3e, 0x88, 0x46,
	0xa8, 0x96, 0x1e, 0x92, 0x03, 0xd3, 0x79, 0x92, 0xeb, 0xd0, 0x8b, 0xbc, 0xd8, 0x4f, 0x47, 0x26,
	0x5d, 0xae, 0xd9, 0x15, 0x19, 0xef, 0x47, 0xff, 0x1c, 0xa2, 0xc7, 0x64, 0x62, 0x1b, 0x95, 0x50,
	0x10, 0xfe, 0x8b, 0xbc, 0x78, 0x92, 0x8e, 0x4d, 0xe1, 0x5e, 0x28, 0x60, 0x97, 0x4e, 0x6c, 0x25,
	0x5b, 0xdd, 0x89, 0x51, 0x32, 0xac, 0x45, 0x01, 0xdd, 0x12, 0x1b, 0xd3, 0x19, 0xf1, 0x4b, 0xa9,
	0xa4, 0xb6, 0xb8, 0x9f, 0xba, 0x84, 0x2d, 0x9c, 0xba, 0x61, 0xe9, 0x69, 0x27, 0x52, 0xca, 0x56,
	0x87, 0x5e, 0xf4, 0x3f, 0x0e, 0xe6, 0xd3, 0xc4, 0x5a, 0xdb, 0x1f, 0xe8, 0x44, 0xcd, 0xf0, 0xfc,
	0x83, 0x04, 0xa6, 0xfa, 0x08, 0xcd, 0xab, 0xcc, 0x81, 0x5e, 0x90, 0xe0, 0x16, 0xf4, 0x97, 0x91,
	0xd9, 0x4f, 0xce, 0x9d, 0x75, 0xf4, 0x6b, 0x1b, 0x1b, 0xf4, 0x30, 0x7b, 0x41, 0x0f, 0xfb, 0x76,
	0xd3, 0xc7, 0x4c, 0x95, 0x0d, 0xb2, 0x91, 0xfd, 0xeb, 0xf3, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xdc, 0xd9, 0xd8, 0x20, 0xb0, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BookServiceClient is the client API for BookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BookServiceClient interface {
	GetBookInfo(ctx context.Context, in *BookInfoParams, opts ...grpc.CallOption) (*BookInfo, error)
	GetBookList(ctx context.Context, in *BookListParams, opts ...grpc.CallOption) (*BookList, error)
}

type bookServiceClient struct {
	cc *grpc.ClientConn
}

func NewBookServiceClient(cc *grpc.ClientConn) BookServiceClient {
	return &bookServiceClient{cc}
}

func (c *bookServiceClient) GetBookInfo(ctx context.Context, in *BookInfoParams, opts ...grpc.CallOption) (*BookInfo, error) {
	out := new(BookInfo)
	err := c.cc.Invoke(ctx, "/book.BookService/GetBookInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) GetBookList(ctx context.Context, in *BookListParams, opts ...grpc.CallOption) (*BookList, error) {
	out := new(BookList)
	err := c.cc.Invoke(ctx, "/book.BookService/GetBookList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookServiceServer is the server API for BookService service.
type BookServiceServer interface {
	GetBookInfo(context.Context, *BookInfoParams) (*BookInfo, error)
	GetBookList(context.Context, *BookListParams) (*BookList, error)
}

// UnimplementedBookServiceServer can be embedded to have forward compatible implementations.
type UnimplementedBookServiceServer struct {
}

func (*UnimplementedBookServiceServer) GetBookInfo(ctx context.Context, req *BookInfoParams) (*BookInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBookInfo not implemented")
}
func (*UnimplementedBookServiceServer) GetBookList(ctx context.Context, req *BookListParams) (*BookList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBookList not implemented")
}

func RegisterBookServiceServer(s *grpc.Server, srv BookServiceServer) {
	s.RegisterService(&_BookService_serviceDesc, srv)
}

func _BookService_GetBookInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookInfoParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).GetBookInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/GetBookInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).GetBookInfo(ctx, req.(*BookInfoParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_GetBookList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookListParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).GetBookList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/GetBookList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).GetBookList(ctx, req.(*BookListParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _BookService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "book.BookService",
	HandlerType: (*BookServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBookInfo",
			Handler:    _BookService_GetBookInfo_Handler,
		},
		{
			MethodName: "GetBookList",
			Handler:    _BookService_GetBookList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/github.com/LongMarch7/go-web/examples/gateway/book/book.proto",
}
