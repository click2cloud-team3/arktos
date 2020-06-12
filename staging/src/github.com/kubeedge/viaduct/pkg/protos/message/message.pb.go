// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
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

type MessageRouter struct {
	// where the message come from
	Source string `protobuf:"bytes,1,opt,name=Source,proto3" json:"Source,omitempty"`
	// where the message will send to
	Group string `protobuf:"bytes,2,opt,name=Group,proto3" json:"Group,omitempty"`
	// what's the operation on resource
	Operaion string `protobuf:"bytes,3,opt,name=Operaion,proto3" json:"Operaion,omitempty"`
	// what's the resource want to operate
	Resouce              string   `protobuf:"bytes,4,opt,name=Resouce,proto3" json:"Resouce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageRouter) Reset()         { *m = MessageRouter{} }
func (m *MessageRouter) String() string { return proto.CompactTextString(m) }
func (*MessageRouter) ProtoMessage()    {}
func (*MessageRouter) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *MessageRouter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageRouter.Unmarshal(m, b)
}
func (m *MessageRouter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageRouter.Marshal(b, m, deterministic)
}
func (m *MessageRouter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageRouter.Merge(m, src)
}
func (m *MessageRouter) XXX_Size() int {
	return xxx_messageInfo_MessageRouter.Size(m)
}
func (m *MessageRouter) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageRouter.DiscardUnknown(m)
}

var xxx_messageInfo_MessageRouter proto.InternalMessageInfo

func (m *MessageRouter) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *MessageRouter) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *MessageRouter) GetOperaion() string {
	if m != nil {
		return m.Operaion
	}
	return ""
}

func (m *MessageRouter) GetResouce() string {
	if m != nil {
		return m.Resouce
	}
	return ""
}

type MessageHeader struct {
	// the message uuid
	ID string `protobuf:"bytes,1,opt,name=MessageID,proto3" json:"MessageID,omitempty"`
	// the response message parent id must be same with the message received
	// please use NewRespByMessage to new response message
	ParentID string `protobuf:"bytes,2,opt,name=ParentID,proto3" json:"ParentID,omitempty"`
	// the time of creating
	Timestamp int64 `protobuf:"zigzag64,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	// the flag will be set in send sync
	Sync bool `protobuf:"varint,4,opt,name=Sync,proto3" json:"Sync,omitempty"`
	// message type
	MessageType          string   `protobuf:"bytes,5,opt,name=MessageType,proto3" json:"MessageType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageHeader) Reset()         { *m = MessageHeader{} }
func (m *MessageHeader) String() string { return proto.CompactTextString(m) }
func (*MessageHeader) ProtoMessage()    {}
func (*MessageHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *MessageHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageHeader.Unmarshal(m, b)
}
func (m *MessageHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageHeader.Marshal(b, m, deterministic)
}
func (m *MessageHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageHeader.Merge(m, src)
}
func (m *MessageHeader) XXX_Size() int {
	return xxx_messageInfo_MessageHeader.Size(m)
}
func (m *MessageHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageHeader.DiscardUnknown(m)
}

var xxx_messageInfo_MessageHeader proto.InternalMessageInfo

func (m *MessageHeader) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *MessageHeader) GetParentID() string {
	if m != nil {
		return m.ParentID
	}
	return ""
}

func (m *MessageHeader) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *MessageHeader) GetSync() bool {
	if m != nil {
		return m.Sync
	}
	return false
}

func (m *MessageHeader) GetMessageType() string {
	if m != nil {
		return m.MessageType
	}
	return ""
}

type Message struct {
	Header               *MessageHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Router               *MessageRouter `protobuf:"bytes,2,opt,name=router,proto3" json:"router,omitempty"`
	Content              []byte         `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetHeader() *MessageHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Message) GetRouter() *MessageRouter {
	if m != nil {
		return m.Router
	}
	return nil
}

func (m *Message) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*MessageRouter)(nil), "message.MessageRouter")
	proto.RegisterType((*MessageHeader)(nil), "message.MessageHeader")
	proto.RegisterType((*Message)(nil), "message.Message")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x95, 0xd0, 0x26, 0xed, 0x95, 0x32, 0x9c, 0x50, 0x65, 0x21, 0x86, 0x2a, 0x13, 0x53,
	0x06, 0x78, 0x04, 0x22, 0x41, 0x06, 0x04, 0x72, 0xfb, 0x02, 0x26, 0x9c, 0xa0, 0x43, 0x6c, 0xcb,
	0x76, 0x86, 0xcc, 0x3c, 0x00, 0xaf, 0x8c, 0x72, 0x71, 0x02, 0x48, 0x6c, 0xfe, 0xee, 0x7e, 0xdd,
	0xff, 0xeb, 0x37, 0x6c, 0x5b, 0xf2, 0x5e, 0xbd, 0x53, 0x69, 0x9d, 0x09, 0x06, 0xf3, 0x88, 0x85,
	0x87, 0xed, 0xd3, 0xf8, 0x94, 0xa6, 0x0b, 0xe4, 0x70, 0x07, 0xd9, 0xc1, 0x74, 0xae, 0x21, 0x91,
	0xec, 0x93, 0x9b, 0xb5, 0x8c, 0x84, 0x97, 0xb0, 0x7c, 0x70, 0xa6, 0xb3, 0x22, 0xe5, 0xf1, 0x08,
	0x78, 0x05, 0xab, 0x67, 0x4b, 0x4e, 0x9d, 0x8c, 0x16, 0x67, 0xbc, 0x98, 0x19, 0x05, 0xe4, 0x92,
	0xbc, 0xe9, 0x1a, 0x12, 0x0b, 0x5e, 0x4d, 0x58, 0x7c, 0x25, 0xb3, 0xeb, 0x23, 0xa9, 0x37, 0x72,
	0x78, 0x01, 0x69, 0x5d, 0x45, 0xc7, 0xb4, 0xae, 0x86, 0xbb, 0x2f, 0xca, 0x91, 0x0e, 0x75, 0x15,
	0x0d, 0x67, 0xc6, 0x6b, 0x58, 0x1f, 0x4f, 0x2d, 0xf9, 0xa0, 0x5a, 0xcb, 0xa6, 0x28, 0x7f, 0x06,
	0x88, 0xb0, 0x38, 0xf4, 0xba, 0x61, 0xcb, 0x95, 0xe4, 0x37, 0xee, 0x61, 0x13, 0xed, 0x8e, 0xbd,
	0x25, 0xb1, 0xe4, 0x83, 0xbf, 0x47, 0xc5, 0x67, 0x02, 0x79, 0x64, 0x2c, 0x21, 0xfb, 0xe0, 0x54,
	0x9c, 0x67, 0x73, 0xbb, 0x2b, 0xa7, 0xee, 0xfe, 0x64, 0x96, 0x51, 0x35, 0xe8, 0x1d, 0x77, 0xc7,
	0x49, 0xff, 0xd1, 0x8f, 0xcd, 0xca, 0xa8, 0x1a, 0x7a, 0xb9, 0x37, 0x3a, 0x90, 0x0e, 0x9c, 0xfe,
	0x5c, 0x4e, 0xf8, 0x9a, 0xf1, 0xe7, 0xdc, 0x7d, 0x07, 0x00, 0x00, 0xff, 0xff, 0xd1, 0xe5, 0x28,
	0xb7, 0xad, 0x01, 0x00, 0x00,
}