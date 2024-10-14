// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.1
// source: services/echo/echo.proto

package echo

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EchoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *EchoMessage) Reset() {
	*x = EchoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_echo_echo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoMessage) ProtoMessage() {}

func (x *EchoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_services_echo_echo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoMessage.ProtoReflect.Descriptor instead.
func (*EchoMessage) Descriptor() ([]byte, []int) {
	return file_services_echo_echo_proto_rawDescGZIP(), []int{0}
}

func (x *EchoMessage) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type EchoNMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	N     int32  `protobuf:"varint,2,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *EchoNMessage) Reset() {
	*x = EchoNMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_echo_echo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoNMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoNMessage) ProtoMessage() {}

func (x *EchoNMessage) ProtoReflect() protoreflect.Message {
	mi := &file_services_echo_echo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoNMessage.ProtoReflect.Descriptor instead.
func (*EchoNMessage) Descriptor() ([]byte, []int) {
	return file_services_echo_echo_proto_rawDescGZIP(), []int{1}
}

func (x *EchoNMessage) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *EchoNMessage) GetN() int32 {
	if x != nil {
		return x.N
	}
	return 0
}

var File_services_echo_echo_proto protoreflect.FileDescriptor

var file_services_echo_echo_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x65, 0x63, 0x68, 0x6f, 0x2f,
	0x65, 0x63, 0x68, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x65, 0x63, 0x68, 0x6f,
	0x22, 0x23, 0x0a, 0x0b, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x32, 0x0a, 0x0c, 0x45, 0x63, 0x68, 0x6f, 0x4e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x6e, 0x32, 0x69, 0x0a, 0x04, 0x45, 0x63, 0x68,
	0x6f, 0x12, 0x2e, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x11, 0x2e, 0x65, 0x63, 0x68, 0x6f,
	0x2e, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x11, 0x2e, 0x65,
	0x63, 0x68, 0x6f, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x00, 0x12, 0x31, 0x0a, 0x05, 0x45, 0x63, 0x68, 0x6f, 0x4e, 0x12, 0x12, 0x2e, 0x65, 0x63, 0x68,
	0x6f, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x4e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x12,
	0x2e, 0x65, 0x63, 0x68, 0x6f, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x4e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x30, 0x01, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6a, 0x62, 0x61, 0x72, 0x72, 0x69, 0x65, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2d, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x65, 0x63, 0x68, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_echo_echo_proto_rawDescOnce sync.Once
	file_services_echo_echo_proto_rawDescData = file_services_echo_echo_proto_rawDesc
)

func file_services_echo_echo_proto_rawDescGZIP() []byte {
	file_services_echo_echo_proto_rawDescOnce.Do(func() {
		file_services_echo_echo_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_echo_echo_proto_rawDescData)
	})
	return file_services_echo_echo_proto_rawDescData
}

var file_services_echo_echo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_services_echo_echo_proto_goTypes = []any{
	(*EchoMessage)(nil),  // 0: echo.EchoMessage
	(*EchoNMessage)(nil), // 1: echo.EchoNMessage
}
var file_services_echo_echo_proto_depIdxs = []int32{
	0, // 0: echo.Echo.Echo:input_type -> echo.EchoMessage
	1, // 1: echo.Echo.EchoN:input_type -> echo.EchoNMessage
	0, // 2: echo.Echo.Echo:output_type -> echo.EchoMessage
	1, // 3: echo.Echo.EchoN:output_type -> echo.EchoNMessage
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_services_echo_echo_proto_init() }
func file_services_echo_echo_proto_init() {
	if File_services_echo_echo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_echo_echo_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*EchoMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_echo_echo_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*EchoNMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_echo_echo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_echo_echo_proto_goTypes,
		DependencyIndexes: file_services_echo_echo_proto_depIdxs,
		MessageInfos:      file_services_echo_echo_proto_msgTypes,
	}.Build()
	File_services_echo_echo_proto = out.File
	file_services_echo_echo_proto_rawDesc = nil
	file_services_echo_echo_proto_goTypes = nil
	file_services_echo_echo_proto_depIdxs = nil
}
