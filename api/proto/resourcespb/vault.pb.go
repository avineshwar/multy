// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: api/proto/resourcespb/vault.proto

package resourcespb

import (
	commonpb "github.com/multycloud/multy/api/proto/commonpb"
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

type CreateVaultRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource *VaultArgs `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *CreateVaultRequest) Reset() {
	*x = CreateVaultRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_vault_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateVaultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateVaultRequest) ProtoMessage() {}

func (x *CreateVaultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_vault_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateVaultRequest.ProtoReflect.Descriptor instead.
func (*CreateVaultRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_vault_proto_rawDescGZIP(), []int{0}
}

func (x *CreateVaultRequest) GetResource() *VaultArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type ReadVaultRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *ReadVaultRequest) Reset() {
	*x = ReadVaultRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_vault_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadVaultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadVaultRequest) ProtoMessage() {}

func (x *ReadVaultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_vault_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadVaultRequest.ProtoReflect.Descriptor instead.
func (*ReadVaultRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_vault_proto_rawDescGZIP(), []int{1}
}

func (x *ReadVaultRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type UpdateVaultRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string     `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Resource   *VaultArgs `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *UpdateVaultRequest) Reset() {
	*x = UpdateVaultRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_vault_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateVaultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateVaultRequest) ProtoMessage() {}

func (x *UpdateVaultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_vault_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateVaultRequest.ProtoReflect.Descriptor instead.
func (*UpdateVaultRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_vault_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateVaultRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *UpdateVaultRequest) GetResource() *VaultArgs {
	if x != nil {
		return x.Resource
	}
	return nil
}

type DeleteVaultRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *DeleteVaultRequest) Reset() {
	*x = DeleteVaultRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_vault_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteVaultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteVaultRequest) ProtoMessage() {}

func (x *DeleteVaultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_vault_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteVaultRequest.ProtoReflect.Descriptor instead.
func (*DeleteVaultRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_vault_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteVaultRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

type VaultArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *commonpb.ResourceCommonArgs `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *VaultArgs) Reset() {
	*x = VaultArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_vault_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VaultArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VaultArgs) ProtoMessage() {}

func (x *VaultArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_vault_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VaultArgs.ProtoReflect.Descriptor instead.
func (*VaultArgs) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_vault_proto_rawDescGZIP(), []int{4}
}

func (x *VaultArgs) GetCommonParameters() *commonpb.ResourceCommonArgs {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *VaultArgs) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type VaultResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonParameters *commonpb.CommonResourceParameters `protobuf:"bytes,1,opt,name=common_parameters,json=commonParameters,proto3" json:"common_parameters,omitempty"`
	Name             string                             `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *VaultResource) Reset() {
	*x = VaultResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_resourcespb_vault_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VaultResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VaultResource) ProtoMessage() {}

func (x *VaultResource) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_resourcespb_vault_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VaultResource.ProtoReflect.Descriptor instead.
func (*VaultResource) Descriptor() ([]byte, []int) {
	return file_api_proto_resourcespb_vault_proto_rawDescGZIP(), []int{5}
}

func (x *VaultResource) GetCommonParameters() *commonpb.CommonResourceParameters {
	if x != nil {
		return x.CommonParameters
	}
	return nil
}

func (x *VaultResource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_api_proto_resourcespb_vault_proto protoreflect.FileDescriptor

var file_api_proto_resourcespb_vault_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x70, 0x62, 0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x13, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x50, 0x0a, 0x12, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x3a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x41, 0x72, 0x67,
	0x73, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x33, 0x0a, 0x10, 0x52,
	0x65, 0x61, 0x64, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64,
	0x22, 0x71, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x64, 0x65, 0x76, 0x2e,
	0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x56, 0x61, 0x75, 0x6c, 0x74, 0x41, 0x72, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x22, 0x35, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x56, 0x61, 0x75,
	0x6c, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0x72, 0x0a, 0x09, 0x56, 0x61,
	0x75, 0x6c, 0x74, 0x41, 0x72, 0x67, 0x73, 0x12, 0x51, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x72, 0x67, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x7c,
	0x0a, 0x0d, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x57, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x64, 0x65, 0x76,
	0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x5e, 0x0a, 0x17,
	0x64, 0x65, 0x76, 0x2e, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x42, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_resourcespb_vault_proto_rawDescOnce sync.Once
	file_api_proto_resourcespb_vault_proto_rawDescData = file_api_proto_resourcespb_vault_proto_rawDesc
)

func file_api_proto_resourcespb_vault_proto_rawDescGZIP() []byte {
	file_api_proto_resourcespb_vault_proto_rawDescOnce.Do(func() {
		file_api_proto_resourcespb_vault_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_resourcespb_vault_proto_rawDescData)
	})
	return file_api_proto_resourcespb_vault_proto_rawDescData
}

var file_api_proto_resourcespb_vault_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_proto_resourcespb_vault_proto_goTypes = []interface{}{
	(*CreateVaultRequest)(nil),                // 0: dev.multy.resources.CreateVaultRequest
	(*ReadVaultRequest)(nil),                  // 1: dev.multy.resources.ReadVaultRequest
	(*UpdateVaultRequest)(nil),                // 2: dev.multy.resources.UpdateVaultRequest
	(*DeleteVaultRequest)(nil),                // 3: dev.multy.resources.DeleteVaultRequest
	(*VaultArgs)(nil),                         // 4: dev.multy.resources.VaultArgs
	(*VaultResource)(nil),                     // 5: dev.multy.resources.VaultResource
	(*commonpb.ResourceCommonArgs)(nil),       // 6: dev.multy.common.ResourceCommonArgs
	(*commonpb.CommonResourceParameters)(nil), // 7: dev.multy.common.CommonResourceParameters
}
var file_api_proto_resourcespb_vault_proto_depIdxs = []int32{
	4, // 0: dev.multy.resources.CreateVaultRequest.resource:type_name -> dev.multy.resources.VaultArgs
	4, // 1: dev.multy.resources.UpdateVaultRequest.resource:type_name -> dev.multy.resources.VaultArgs
	6, // 2: dev.multy.resources.VaultArgs.common_parameters:type_name -> dev.multy.common.ResourceCommonArgs
	7, // 3: dev.multy.resources.VaultResource.common_parameters:type_name -> dev.multy.common.CommonResourceParameters
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_proto_resourcespb_vault_proto_init() }
func file_api_proto_resourcespb_vault_proto_init() {
	if File_api_proto_resourcespb_vault_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_resourcespb_vault_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateVaultRequest); i {
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
		file_api_proto_resourcespb_vault_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadVaultRequest); i {
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
		file_api_proto_resourcespb_vault_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateVaultRequest); i {
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
		file_api_proto_resourcespb_vault_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteVaultRequest); i {
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
		file_api_proto_resourcespb_vault_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VaultArgs); i {
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
		file_api_proto_resourcespb_vault_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VaultResource); i {
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
			RawDescriptor: file_api_proto_resourcespb_vault_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_resourcespb_vault_proto_goTypes,
		DependencyIndexes: file_api_proto_resourcespb_vault_proto_depIdxs,
		MessageInfos:      file_api_proto_resourcespb_vault_proto_msgTypes,
	}.Build()
	File_api_proto_resourcespb_vault_proto = out.File
	file_api_proto_resourcespb_vault_proto_rawDesc = nil
	file_api_proto_resourcespb_vault_proto_goTypes = nil
	file_api_proto_resourcespb_vault_proto_depIdxs = nil
}
