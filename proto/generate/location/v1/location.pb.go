// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: location/location.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetCourierLatestPositionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CourierId string `protobuf:"bytes,1,opt,name=courier_id,json=courierId,proto3" json:"courier_id,omitempty"`
}

func (x *GetCourierLatestPositionRequest) Reset() {
	*x = GetCourierLatestPositionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_location_location_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCourierLatestPositionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCourierLatestPositionRequest) ProtoMessage() {}

func (x *GetCourierLatestPositionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_location_location_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCourierLatestPositionRequest.ProtoReflect.Descriptor instead.
func (*GetCourierLatestPositionRequest) Descriptor() ([]byte, []int) {
	return file_location_location_proto_rawDescGZIP(), []int{0}
}

func (x *GetCourierLatestPositionRequest) GetCourierId() string {
	if x != nil {
		return x.CourierId
	}
	return ""
}

type GetCourierLatestPositionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float64 `protobuf:"fixed64,2,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,3,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *GetCourierLatestPositionResponse) Reset() {
	*x = GetCourierLatestPositionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_location_location_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCourierLatestPositionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCourierLatestPositionResponse) ProtoMessage() {}

func (x *GetCourierLatestPositionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_location_location_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCourierLatestPositionResponse.ProtoReflect.Descriptor instead.
func (*GetCourierLatestPositionResponse) Descriptor() ([]byte, []int) {
	return file_location_location_proto_rawDescGZIP(), []int{1}
}

func (x *GetCourierLatestPositionResponse) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *GetCourierLatestPositionResponse) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

var File_location_location_proto protoreflect.FileDescriptor

var file_location_location_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40, 0x0a, 0x1f, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x22, 0x5c, 0x0a, 0x20, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c,
	0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09,
	0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x32, 0x7c, 0x0a, 0x17, 0x43, 0x6f, 0x75,
	0x72, 0x69, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x61, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69,
	0x65, 0x72, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x4c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x4c,
	0x61, 0x74, 0x65, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x16, 0x5a, 0x14, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_location_location_proto_rawDescOnce sync.Once
	file_location_location_proto_rawDescData = file_location_location_proto_rawDesc
)

func file_location_location_proto_rawDescGZIP() []byte {
	file_location_location_proto_rawDescOnce.Do(func() {
		file_location_location_proto_rawDescData = protoimpl.X.CompressGZIP(file_location_location_proto_rawDescData)
	})
	return file_location_location_proto_rawDescData
}

var file_location_location_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_location_location_proto_goTypes = []interface{}{
	(*GetCourierLatestPositionRequest)(nil),  // 0: GetCourierLatestPositionRequest
	(*GetCourierLatestPositionResponse)(nil), // 1: GetCourierLatestPositionResponse
}
var file_location_location_proto_depIdxs = []int32{
	0, // 0: CourierLocationPosition.GetCourierLatestPosition:input_type -> GetCourierLatestPositionRequest
	1, // 1: CourierLocationPosition.GetCourierLatestPosition:output_type -> GetCourierLatestPositionResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_location_location_proto_init() }
func file_location_location_proto_init() {
	if File_location_location_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_location_location_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCourierLatestPositionRequest); i {
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
		file_location_location_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCourierLatestPositionResponse); i {
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
			RawDescriptor: file_location_location_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_location_location_proto_goTypes,
		DependencyIndexes: file_location_location_proto_depIdxs,
		MessageInfos:      file_location_location_proto_msgTypes,
	}.Build()
	File_location_location_proto = out.File
	file_location_location_proto_rawDesc = nil
	file_location_location_proto_goTypes = nil
	file_location_location_proto_depIdxs = nil
}
