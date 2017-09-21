// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/lightlog.cmd.proto

/*
Package lightlog is a generated protocol buffer package.

It is generated from these files:
	protos/lightlog.cmd.proto

It has these top-level messages:
	Cmd
*/
package lightlog

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Cmd struct {
	Cmd       string `protobuf:"bytes,1,opt,name=cmd" json:"cmd,omitempty"`
	Seq       int32  `protobuf:"varint,2,opt,name=seq" json:"seq,omitempty"`
	StartTime int64  `protobuf:"varint,3,opt,name=startTime" json:"startTime,omitempty"`
	EndTime   int64  `protobuf:"varint,4,opt,name=endTime" json:"endTime,omitempty"`
	Data      []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Cmd) Reset()                    { *m = Cmd{} }
func (m *Cmd) String() string            { return proto.CompactTextString(m) }
func (*Cmd) ProtoMessage()               {}
func (*Cmd) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Cmd) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *Cmd) GetSeq() int32 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *Cmd) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *Cmd) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *Cmd) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Cmd)(nil), "lightlog.cmd")
}

func init() { proto.RegisterFile("protos/lightlog.cmd.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0xcf, 0xc9, 0x4c, 0xcf, 0x28, 0xc9, 0xc9, 0x4f, 0xd7, 0x4b, 0xce, 0x4d, 0xd1,
	0x03, 0x8b, 0x09, 0x71, 0xc0, 0xc4, 0x94, 0x4a, 0xb9, 0x98, 0x93, 0x73, 0x53, 0x84, 0x04, 0xc0,
	0x94, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x4c, 0xa4, 0x38, 0xb5, 0x50, 0x82, 0x49, 0x81,
	0x51, 0x83, 0x35, 0x08, 0xc4, 0x14, 0x92, 0xe1, 0xe2, 0x2c, 0x2e, 0x49, 0x2c, 0x2a, 0x09, 0xc9,
	0xcc, 0x4d, 0x95, 0x60, 0x56, 0x60, 0xd4, 0x60, 0x0e, 0x42, 0x08, 0x08, 0x49, 0x70, 0xb1, 0xa7,
	0xe6, 0xa5, 0x80, 0xe5, 0x58, 0xc0, 0x72, 0x30, 0xae, 0x90, 0x10, 0x17, 0x4b, 0x4a, 0x62, 0x49,
	0xa2, 0x04, 0xab, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x98, 0x9d, 0xc4, 0x06, 0x76, 0x87, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x6b, 0x3a, 0x64, 0xdf, 0xa4, 0x00, 0x00, 0x00,
}
