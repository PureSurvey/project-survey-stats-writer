// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.0--rc1
// source: trackingevent.proto

package pb

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

type TrackingEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType      int32   `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Timestamp      int64   `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	UnitId         int32   `protobuf:"varint,3,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	Geo            string  `protobuf:"bytes,4,opt,name=geo,proto3" json:"geo,omitempty"`
	Lang           string  `protobuf:"bytes,5,opt,name=lang,proto3" json:"lang,omitempty"`
	Gender         int32   `protobuf:"varint,6,opt,name=gender,proto3" json:"gender,omitempty"`
	SurveyIds      []int32 `protobuf:"varint,7,rep,packed,name=survey_ids,json=surveyIds,proto3" json:"survey_ids,omitempty"`
	MismatchReason []int32 `protobuf:"varint,8,rep,packed,name=mismatch_reason,json=mismatchReason,proto3" json:"mismatch_reason,omitempty"`
}

func (x *TrackingEvent) Reset() {
	*x = TrackingEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trackingevent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrackingEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrackingEvent) ProtoMessage() {}

func (x *TrackingEvent) ProtoReflect() protoreflect.Message {
	mi := &file_trackingevent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrackingEvent.ProtoReflect.Descriptor instead.
func (*TrackingEvent) Descriptor() ([]byte, []int) {
	return file_trackingevent_proto_rawDescGZIP(), []int{0}
}

func (x *TrackingEvent) GetEventType() int32 {
	if x != nil {
		return x.EventType
	}
	return 0
}

func (x *TrackingEvent) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *TrackingEvent) GetUnitId() int32 {
	if x != nil {
		return x.UnitId
	}
	return 0
}

func (x *TrackingEvent) GetGeo() string {
	if x != nil {
		return x.Geo
	}
	return ""
}

func (x *TrackingEvent) GetLang() string {
	if x != nil {
		return x.Lang
	}
	return ""
}

func (x *TrackingEvent) GetGender() int32 {
	if x != nil {
		return x.Gender
	}
	return 0
}

func (x *TrackingEvent) GetSurveyIds() []int32 {
	if x != nil {
		return x.SurveyIds
	}
	return nil
}

func (x *TrackingEvent) GetMismatchReason() []int32 {
	if x != nil {
		return x.MismatchReason
	}
	return nil
}

var File_trackingevent_proto protoreflect.FileDescriptor

var file_trackingevent_proto_rawDesc = []byte{
	0x0a, 0x13, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0xeb, 0x01, 0x0a, 0x0d, 0x54, 0x72,
	0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x65, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x67, 0x65, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x61, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6c, 0x61, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x05, 0x52, 0x09, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x49, 0x64, 0x73, 0x12, 0x27,
	0x0a, 0x0f, 0x6d, 0x69, 0x73, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x08, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0e, 0x6d, 0x69, 0x73, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x42, 0x36, 0x5a, 0x34, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x2d, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x2d, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2d,
	0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trackingevent_proto_rawDescOnce sync.Once
	file_trackingevent_proto_rawDescData = file_trackingevent_proto_rawDesc
)

func file_trackingevent_proto_rawDescGZIP() []byte {
	file_trackingevent_proto_rawDescOnce.Do(func() {
		file_trackingevent_proto_rawDescData = protoimpl.X.CompressGZIP(file_trackingevent_proto_rawDescData)
	})
	return file_trackingevent_proto_rawDescData
}

var file_trackingevent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_trackingevent_proto_goTypes = []interface{}{
	(*TrackingEvent)(nil), // 0: pb.TrackingEvent
}
var file_trackingevent_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_trackingevent_proto_init() }
func file_trackingevent_proto_init() {
	if File_trackingevent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trackingevent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrackingEvent); i {
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
			RawDescriptor: file_trackingevent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trackingevent_proto_goTypes,
		DependencyIndexes: file_trackingevent_proto_depIdxs,
		MessageInfos:      file_trackingevent_proto_msgTypes,
	}.Build()
	File_trackingevent_proto = out.File
	file_trackingevent_proto_rawDesc = nil
	file_trackingevent_proto_goTypes = nil
	file_trackingevent_proto_depIdxs = nil
}
