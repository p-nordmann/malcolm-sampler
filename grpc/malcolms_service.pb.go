// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: grpc/malcolms_service.proto

// gRPC service definition.

package grpc

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

type BoundariesUUID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BoundariesUUID) Reset() {
	*x = BoundariesUUID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_malcolms_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoundariesUUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoundariesUUID) ProtoMessage() {}

func (x *BoundariesUUID) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_malcolms_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoundariesUUID.ProtoReflect.Descriptor instead.
func (*BoundariesUUID) Descriptor() ([]byte, []int) {
	return file_grpc_malcolms_service_proto_rawDescGZIP(), []int{0}
}

func (x *BoundariesUUID) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type PosteriorUUID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *PosteriorUUID) Reset() {
	*x = PosteriorUUID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_malcolms_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PosteriorUUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PosteriorUUID) ProtoMessage() {}

func (x *PosteriorUUID) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_malcolms_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PosteriorUUID.ProtoReflect.Descriptor instead.
func (*PosteriorUUID) Descriptor() ([]byte, []int) {
	return file_grpc_malcolms_service_proto_rawDescGZIP(), []int{1}
}

func (x *PosteriorUUID) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Boundaries describes the bounding box of an inversion problem.
//
// Infima are the lower bounds of the parameter space, while suprema are the higher bounds.
type Boundaries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dimension int64     `protobuf:"varint,1,opt,name=dimension,proto3" json:"dimension,omitempty"`
	Infima    []float64 `protobuf:"fixed64,2,rep,packed,name=infima,proto3" json:"infima,omitempty"`
	Suprema   []float64 `protobuf:"fixed64,3,rep,packed,name=suprema,proto3" json:"suprema,omitempty"`
}

func (x *Boundaries) Reset() {
	*x = Boundaries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_malcolms_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Boundaries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Boundaries) ProtoMessage() {}

func (x *Boundaries) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_malcolms_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Boundaries.ProtoReflect.Descriptor instead.
func (*Boundaries) Descriptor() ([]byte, []int) {
	return file_grpc_malcolms_service_proto_rawDescGZIP(), []int{2}
}

func (x *Boundaries) GetDimension() int64 {
	if x != nil {
		return x.Dimension
	}
	return 0
}

func (x *Boundaries) GetInfima() []float64 {
	if x != nil {
		return x.Infima
	}
	return nil
}

func (x *Boundaries) GetSuprema() []float64 {
	if x != nil {
		return x.Suprema
	}
	return nil
}

// PosteriorValuesBatch represents a batch of samples with posterior values.
//
// It expects column-major order for coordinates.
type PosteriorValuesBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid            *BoundariesUUID `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Coordinates     []float64       `protobuf:"fixed64,2,rep,packed,name=coordinates,proto3" json:"coordinates,omitempty"`
	PosteriorValues []float64       `protobuf:"fixed64,3,rep,packed,name=posterior_values,json=posteriorValues,proto3" json:"posterior_values,omitempty"`
}

func (x *PosteriorValuesBatch) Reset() {
	*x = PosteriorValuesBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_malcolms_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PosteriorValuesBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PosteriorValuesBatch) ProtoMessage() {}

func (x *PosteriorValuesBatch) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_malcolms_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PosteriorValuesBatch.ProtoReflect.Descriptor instead.
func (*PosteriorValuesBatch) Descriptor() ([]byte, []int) {
	return file_grpc_malcolms_service_proto_rawDescGZIP(), []int{3}
}

func (x *PosteriorValuesBatch) GetUuid() *BoundariesUUID {
	if x != nil {
		return x.Uuid
	}
	return nil
}

func (x *PosteriorValuesBatch) GetCoordinates() []float64 {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

func (x *PosteriorValuesBatch) GetPosteriorValues() []float64 {
	if x != nil {
		return x.PosteriorValues
	}
	return nil
}

// MakeSamplesRequest represents a request to generate samples.
//
// The underlying algorithm is a random walk (MCMC).
// `origin` is the point where the walk starts from and `amount` is the number of samples to generate.
type MakeSamplesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid   *PosteriorUUID `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Origin []float64      `protobuf:"fixed64,2,rep,packed,name=origin,proto3" json:"origin,omitempty"`
	Amount int64          `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *MakeSamplesRequest) Reset() {
	*x = MakeSamplesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_malcolms_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MakeSamplesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MakeSamplesRequest) ProtoMessage() {}

func (x *MakeSamplesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_malcolms_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MakeSamplesRequest.ProtoReflect.Descriptor instead.
func (*MakeSamplesRequest) Descriptor() ([]byte, []int) {
	return file_grpc_malcolms_service_proto_rawDescGZIP(), []int{4}
}

func (x *MakeSamplesRequest) GetUuid() *PosteriorUUID {
	if x != nil {
		return x.Uuid
	}
	return nil
}

func (x *MakeSamplesRequest) GetOrigin() []float64 {
	if x != nil {
		return x.Origin
	}
	return nil
}

func (x *MakeSamplesRequest) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

// SamplesBatch represents a batch of generated samples.
//
// It stores coordinates in row-major order.
//
// Whether it holds more than 1 point is implementation-specific.
type SamplesBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Coordinates []float64 `protobuf:"fixed64,1,rep,packed,name=coordinates,proto3" json:"coordinates,omitempty"`
}

func (x *SamplesBatch) Reset() {
	*x = SamplesBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_malcolms_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SamplesBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SamplesBatch) ProtoMessage() {}

func (x *SamplesBatch) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_malcolms_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SamplesBatch.ProtoReflect.Descriptor instead.
func (*SamplesBatch) Descriptor() ([]byte, []int) {
	return file_grpc_malcolms_service_proto_rawDescGZIP(), []int{5}
}

func (x *SamplesBatch) GetCoordinates() []float64 {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

var File_grpc_malcolms_service_proto protoreflect.FileDescriptor

var file_grpc_malcolms_service_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x61, 0x6c, 0x63, 0x6f, 0x6c, 0x6d, 0x73, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67,
	0x72, 0x70, 0x63, 0x22, 0x26, 0x0a, 0x0e, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x72, 0x69, 0x65,
	0x73, 0x55, 0x55, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x25, 0x0a, 0x0d, 0x50,
	0x6f, 0x73, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x72, 0x55, 0x55, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x5c, 0x0a, 0x0a, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x72, 0x69, 0x65, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x6d, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x64, 0x69, 0x6d, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x6e, 0x66, 0x69, 0x6d, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06,
	0x69, 0x6e, 0x66, 0x69, 0x6d, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x70, 0x72, 0x65, 0x6d,
	0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x01, 0x52, 0x07, 0x73, 0x75, 0x70, 0x72, 0x65, 0x6d, 0x61,
	0x22, 0x8d, 0x01, 0x0a, 0x14, 0x50, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x72, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x28, 0x0a, 0x04, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x42,
	0x6f, 0x75, 0x6e, 0x64, 0x61, 0x72, 0x69, 0x65, 0x73, 0x55, 0x55, 0x49, 0x44, 0x52, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x01, 0x52, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69,
	0x6e, 0x61, 0x74, 0x65, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x69,
	0x6f, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x01, 0x52,
	0x0f, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x22, 0x6d, 0x0a, 0x12, 0x4d, 0x61, 0x6b, 0x65, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x6f, 0x73, 0x74,
	0x65, 0x72, 0x69, 0x6f, 0x72, 0x55, 0x55, 0x49, 0x44, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x01, 0x52,
	0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x30, 0x0a, 0x0c, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x01, 0x52, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65,
	0x73, 0x32, 0xd1, 0x01, 0x0a, 0x0e, 0x4d, 0x61, 0x6c, 0x63, 0x6f, 0x6c, 0x6d, 0x53, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x0d, 0x41, 0x64, 0x64, 0x42, 0x6f, 0x75, 0x6e, 0x64,
	0x61, 0x72, 0x69, 0x65, 0x73, 0x12, 0x10, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x42, 0x6f, 0x75,
	0x6e, 0x64, 0x61, 0x72, 0x69, 0x65, 0x73, 0x1a, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x42,
	0x6f, 0x75, 0x6e, 0x64, 0x61, 0x72, 0x69, 0x65, 0x73, 0x55, 0x55, 0x49, 0x44, 0x22, 0x00, 0x12,
	0x43, 0x0a, 0x0c, 0x41, 0x64, 0x64, 0x50, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x72, 0x12,
	0x1a, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x72,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x1a, 0x13, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x72, 0x55, 0x55, 0x49, 0x44,
	0x22, 0x00, 0x28, 0x01, 0x12, 0x3f, 0x0a, 0x0b, 0x4d, 0x61, 0x6b, 0x65, 0x53, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x73, 0x12, 0x18, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x61, 0x6b, 0x65, 0x53,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x22, 0x00, 0x30, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x2d, 0x6e, 0x6f, 0x72, 0x64, 0x6d, 0x61, 0x6e, 0x6e, 0x2f, 0x6d,
	0x61, 0x6c, 0x63, 0x6f, 0x6c, 0x6d, 0x2d, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_malcolms_service_proto_rawDescOnce sync.Once
	file_grpc_malcolms_service_proto_rawDescData = file_grpc_malcolms_service_proto_rawDesc
)

func file_grpc_malcolms_service_proto_rawDescGZIP() []byte {
	file_grpc_malcolms_service_proto_rawDescOnce.Do(func() {
		file_grpc_malcolms_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_malcolms_service_proto_rawDescData)
	})
	return file_grpc_malcolms_service_proto_rawDescData
}

var file_grpc_malcolms_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_grpc_malcolms_service_proto_goTypes = []interface{}{
	(*BoundariesUUID)(nil),       // 0: grpc.BoundariesUUID
	(*PosteriorUUID)(nil),        // 1: grpc.PosteriorUUID
	(*Boundaries)(nil),           // 2: grpc.Boundaries
	(*PosteriorValuesBatch)(nil), // 3: grpc.PosteriorValuesBatch
	(*MakeSamplesRequest)(nil),   // 4: grpc.MakeSamplesRequest
	(*SamplesBatch)(nil),         // 5: grpc.SamplesBatch
}
var file_grpc_malcolms_service_proto_depIdxs = []int32{
	0, // 0: grpc.PosteriorValuesBatch.uuid:type_name -> grpc.BoundariesUUID
	1, // 1: grpc.MakeSamplesRequest.uuid:type_name -> grpc.PosteriorUUID
	2, // 2: grpc.MalcolmSampler.AddBoundaries:input_type -> grpc.Boundaries
	3, // 3: grpc.MalcolmSampler.AddPosterior:input_type -> grpc.PosteriorValuesBatch
	4, // 4: grpc.MalcolmSampler.MakeSamples:input_type -> grpc.MakeSamplesRequest
	0, // 5: grpc.MalcolmSampler.AddBoundaries:output_type -> grpc.BoundariesUUID
	1, // 6: grpc.MalcolmSampler.AddPosterior:output_type -> grpc.PosteriorUUID
	5, // 7: grpc.MalcolmSampler.MakeSamples:output_type -> grpc.SamplesBatch
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpc_malcolms_service_proto_init() }
func file_grpc_malcolms_service_proto_init() {
	if File_grpc_malcolms_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_malcolms_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoundariesUUID); i {
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
		file_grpc_malcolms_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PosteriorUUID); i {
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
		file_grpc_malcolms_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Boundaries); i {
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
		file_grpc_malcolms_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PosteriorValuesBatch); i {
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
		file_grpc_malcolms_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MakeSamplesRequest); i {
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
		file_grpc_malcolms_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SamplesBatch); i {
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
			RawDescriptor: file_grpc_malcolms_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_malcolms_service_proto_goTypes,
		DependencyIndexes: file_grpc_malcolms_service_proto_depIdxs,
		MessageInfos:      file_grpc_malcolms_service_proto_msgTypes,
	}.Build()
	File_grpc_malcolms_service_proto = out.File
	file_grpc_malcolms_service_proto_rawDesc = nil
	file_grpc_malcolms_service_proto_goTypes = nil
	file_grpc_malcolms_service_proto_depIdxs = nil
}
