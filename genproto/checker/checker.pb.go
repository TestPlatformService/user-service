// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: checker.proto

package checker

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

type InputOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	In  string `protobuf:"bytes,1,opt,name=in,proto3" json:"in,omitempty"`
	Out string `protobuf:"bytes,2,opt,name=out,proto3" json:"out,omitempty"`
}

func (x *InputOutput) Reset() {
	*x = InputOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputOutput) ProtoMessage() {}

func (x *InputOutput) ProtoReflect() protoreflect.Message {
	mi := &file_checker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputOutput.ProtoReflect.Descriptor instead.
func (*InputOutput) Descriptor() ([]byte, []int) {
	return file_checker_proto_rawDescGZIP(), []int{0}
}

func (x *InputOutput) GetIn() string {
	if x != nil {
		return x.In
	}
	return ""
}

func (x *InputOutput) GetOut() string {
	if x != nil {
		return x.Out
	}
	return ""
}

type SubmitReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code        string         `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Lang        string         `protobuf:"bytes,2,opt,name=lang,proto3" json:"lang,omitempty"`
	TimeLimit   int32          `protobuf:"varint,3,opt,name=time_limit,json=timeLimit,proto3" json:"time_limit,omitempty"`
	MemoryLimit int64          `protobuf:"varint,4,opt,name=memory_limit,json=memoryLimit,proto3" json:"memory_limit,omitempty"`
	Io          []*InputOutput `protobuf:"bytes,5,rep,name=io,proto3" json:"io,omitempty"`
}

func (x *SubmitReq) Reset() {
	*x = SubmitReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitReq) ProtoMessage() {}

func (x *SubmitReq) ProtoReflect() protoreflect.Message {
	mi := &file_checker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitReq.ProtoReflect.Descriptor instead.
func (*SubmitReq) Descriptor() ([]byte, []int) {
	return file_checker_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitReq) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *SubmitReq) GetLang() string {
	if x != nil {
		return x.Lang
	}
	return ""
}

func (x *SubmitReq) GetTimeLimit() int32 {
	if x != nil {
		return x.TimeLimit
	}
	return 0
}

func (x *SubmitReq) GetMemoryLimit() int64 {
	if x != nil {
		return x.MemoryLimit
	}
	return 0
}

func (x *SubmitReq) GetIo() []*InputOutput {
	if x != nil {
		return x.Io
	}
	return nil
}

type SubmitResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SubmitResp) Reset() {
	*x = SubmitResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitResp) ProtoMessage() {}

func (x *SubmitResp) ProtoReflect() protoreflect.Message {
	mi := &file_checker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitResp.ProtoReflect.Descriptor instead.
func (*SubmitResp) Descriptor() ([]byte, []int) {
	return file_checker_proto_rawDescGZIP(), []int{2}
}

var File_checker_proto protoreflect.FileDescriptor

var file_checker_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x22, 0x2f, 0x0a, 0x0b, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x75, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6f, 0x75, 0x74, 0x22, 0x9b, 0x01, 0x0a, 0x09, 0x53, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c,
	0x61, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x61, 0x6e, 0x67, 0x12,
	0x1d, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x12, 0x24, 0x0a, 0x02, 0x69, 0x6f, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x52, 0x02, 0x69, 0x6f, 0x22, 0x0c, 0x0a, 0x0a, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x32, 0x43, 0x0a, 0x0e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x12, 0x12, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x75, 0x62, 0x6d,
	0x69, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x2e,
	0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x42, 0x12, 0x5a, 0x10, 0x67, 0x65,
	0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_checker_proto_rawDescOnce sync.Once
	file_checker_proto_rawDescData = file_checker_proto_rawDesc
)

func file_checker_proto_rawDescGZIP() []byte {
	file_checker_proto_rawDescOnce.Do(func() {
		file_checker_proto_rawDescData = protoimpl.X.CompressGZIP(file_checker_proto_rawDescData)
	})
	return file_checker_proto_rawDescData
}

var file_checker_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_checker_proto_goTypes = []any{
	(*InputOutput)(nil), // 0: checker.InputOutput
	(*SubmitReq)(nil),   // 1: checker.SubmitReq
	(*SubmitResp)(nil),  // 2: checker.SubmitResp
}
var file_checker_proto_depIdxs = []int32{
	0, // 0: checker.SubmitReq.io:type_name -> checker.InputOutput
	1, // 1: checker.CheckerService.Submit:input_type -> checker.SubmitReq
	2, // 2: checker.CheckerService.Submit:output_type -> checker.SubmitResp
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_checker_proto_init() }
func file_checker_proto_init() {
	if File_checker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_checker_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*InputOutput); i {
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
		file_checker_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SubmitReq); i {
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
		file_checker_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SubmitResp); i {
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
			RawDescriptor: file_checker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_checker_proto_goTypes,
		DependencyIndexes: file_checker_proto_depIdxs,
		MessageInfos:      file_checker_proto_msgTypes,
	}.Build()
	File_checker_proto = out.File
	file_checker_proto_rawDesc = nil
	file_checker_proto_goTypes = nil
	file_checker_proto_depIdxs = nil
}
