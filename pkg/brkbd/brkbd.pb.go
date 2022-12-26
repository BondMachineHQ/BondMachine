// Code generated by protoc-gen-go. DO NOT EDIT.
// source: brkbd.proto

package brkbd

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type KeySeq struct {
	Seq                  []*KeySeq_KeyCode `protobuf:"bytes,1,rep,name=seq,proto3" json:"seq,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *KeySeq) Reset()         { *m = KeySeq{} }
func (m *KeySeq) String() string { return proto.CompactTextString(m) }
func (*KeySeq) ProtoMessage()    {}
func (*KeySeq) Descriptor() ([]byte, []int) {
	return fileDescriptor_eb3f5219cbd11c7c, []int{0}
}

func (m *KeySeq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeySeq.Unmarshal(m, b)
}
func (m *KeySeq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeySeq.Marshal(b, m, deterministic)
}
func (m *KeySeq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeySeq.Merge(m, src)
}
func (m *KeySeq) XXX_Size() int {
	return xxx_messageInfo_KeySeq.Size(m)
}
func (m *KeySeq) XXX_DiscardUnknown() {
	xxx_messageInfo_KeySeq.DiscardUnknown(m)
}

var xxx_messageInfo_KeySeq proto.InternalMessageInfo

func (m *KeySeq) GetSeq() []*KeySeq_KeyCode {
	if m != nil {
		return m.Seq
	}
	return nil
}

type KeySeq_KeyCode struct {
	Code                 []byte   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeySeq_KeyCode) Reset()         { *m = KeySeq_KeyCode{} }
func (m *KeySeq_KeyCode) String() string { return proto.CompactTextString(m) }
func (*KeySeq_KeyCode) ProtoMessage()    {}
func (*KeySeq_KeyCode) Descriptor() ([]byte, []int) {
	return fileDescriptor_eb3f5219cbd11c7c, []int{0, 0}
}

func (m *KeySeq_KeyCode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeySeq_KeyCode.Unmarshal(m, b)
}
func (m *KeySeq_KeyCode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeySeq_KeyCode.Marshal(b, m, deterministic)
}
func (m *KeySeq_KeyCode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeySeq_KeyCode.Merge(m, src)
}
func (m *KeySeq_KeyCode) XXX_Size() int {
	return xxx_messageInfo_KeySeq_KeyCode.Size(m)
}
func (m *KeySeq_KeyCode) XXX_DiscardUnknown() {
	xxx_messageInfo_KeySeq_KeyCode.DiscardUnknown(m)
}

var xxx_messageInfo_KeySeq_KeyCode proto.InternalMessageInfo

func (m *KeySeq_KeyCode) GetCode() []byte {
	if m != nil {
		return m.Code
	}
	return nil
}

func init() {
	proto.RegisterType((*KeySeq)(nil), "brkbd.KeySeq")
	proto.RegisterType((*KeySeq_KeyCode)(nil), "brkbd.KeySeq.KeyCode")
}

func init() { proto.RegisterFile("brkbd.proto", fileDescriptor_eb3f5219cbd11c7c) }

var fileDescriptor_eb3f5219cbd11c7c = []byte{
	// 115 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0x2a, 0xca, 0x4e,
	0x4a, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x02, 0xb8, 0xd8, 0xbc,
	0x53, 0x2b, 0x83, 0x53, 0x0b, 0x85, 0xd4, 0xb9, 0x98, 0x8b, 0x53, 0x0b, 0x25, 0x18, 0x15, 0x98,
	0x35, 0xb8, 0x8d, 0x44, 0xf5, 0x20, 0x6a, 0x21, 0x72, 0x20, 0xca, 0x39, 0x3f, 0x25, 0x35, 0x08,
	0xa4, 0x42, 0x4a, 0x96, 0x8b, 0x1d, 0xca, 0x17, 0x12, 0xe2, 0x62, 0x49, 0xce, 0x4f, 0x49, 0x95,
	0x60, 0x54, 0x60, 0xd4, 0xe0, 0x09, 0x02, 0xb3, 0x9d, 0x38, 0xa3, 0xd8, 0xf5, 0xf4, 0xc1, 0xba,
	0x93, 0xd8, 0xc0, 0x56, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x49, 0x12, 0xb5, 0x41, 0x79,
	0x00, 0x00, 0x00,
}
