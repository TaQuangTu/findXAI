// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: api/content.proto

package contentsvc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ExtractContentFromLinksRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Links         []string               `protobuf:"bytes,1,rep,name=links,proto3" json:"links,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractContentFromLinksRequest) Reset() {
	*x = ExtractContentFromLinksRequest{}
	mi := &file_api_content_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractContentFromLinksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractContentFromLinksRequest) ProtoMessage() {}

func (x *ExtractContentFromLinksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_content_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractContentFromLinksRequest.ProtoReflect.Descriptor instead.
func (*ExtractContentFromLinksRequest) Descriptor() ([]byte, []int) {
	return file_api_content_proto_rawDescGZIP(), []int{0}
}

func (x *ExtractContentFromLinksRequest) GetLinks() []string {
	if x != nil {
		return x.Links
	}
	return nil
}

type ExtractContentFromLinksReponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Contents      []*ExtractedContent    `protobuf:"bytes,1,rep,name=contents,proto3" json:"contents,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractContentFromLinksReponse) Reset() {
	*x = ExtractContentFromLinksReponse{}
	mi := &file_api_content_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractContentFromLinksReponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractContentFromLinksReponse) ProtoMessage() {}

func (x *ExtractContentFromLinksReponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_content_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractContentFromLinksReponse.ProtoReflect.Descriptor instead.
func (*ExtractContentFromLinksReponse) Descriptor() ([]byte, []int) {
	return file_api_content_proto_rawDescGZIP(), []int{1}
}

func (x *ExtractContentFromLinksReponse) GetContents() []*ExtractedContent {
	if x != nil {
		return x.Contents
	}
	return nil
}

type ExtractedContent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Link          string                 `protobuf:"bytes,1,opt,name=link,proto3" json:"link,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Content       string                 `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractedContent) Reset() {
	*x = ExtractedContent{}
	mi := &file_api_content_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractedContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractedContent) ProtoMessage() {}

func (x *ExtractedContent) ProtoReflect() protoreflect.Message {
	mi := &file_api_content_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractedContent.ProtoReflect.Descriptor instead.
func (*ExtractedContent) Descriptor() ([]byte, []int) {
	return file_api_content_proto_rawDescGZIP(), []int{2}
}

func (x *ExtractedContent) GetLink() string {
	if x != nil {
		return x.Link
	}
	return ""
}

func (x *ExtractedContent) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ExtractedContent) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_api_content_proto protoreflect.FileDescriptor

const file_api_content_proto_rawDesc = "" +
	"\n" +
	"\x11api/content.proto\x12\x10google.search.v1\"6\n" +
	"\x1eExtractContentFromLinksRequest\x12\x14\n" +
	"\x05links\x18\x01 \x03(\tR\x05links\"`\n" +
	"\x1eExtractContentFromLinksReponse\x12>\n" +
	"\bcontents\x18\x01 \x03(\v2\".google.search.v1.ExtractedContentR\bcontents\"V\n" +
	"\x10ExtractedContent\x12\x12\n" +
	"\x04link\x18\x01 \x01(\tR\x04link\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12\x18\n" +
	"\acontent\x18\x03 \x01(\tR\acontent2\x8f\x01\n" +
	"\x0eContentService\x12}\n" +
	"\x17ExtractContentFromLinks\x120.google.search.v1.ExtractContentFromLinksRequest\x1a0.google.search.v1.ExtractContentFromLinksReponseB\x10Z\x0epkg/contentsvcb\x06proto3"

var (
	file_api_content_proto_rawDescOnce sync.Once
	file_api_content_proto_rawDescData []byte
)

func file_api_content_proto_rawDescGZIP() []byte {
	file_api_content_proto_rawDescOnce.Do(func() {
		file_api_content_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_content_proto_rawDesc), len(file_api_content_proto_rawDesc)))
	})
	return file_api_content_proto_rawDescData
}

var file_api_content_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_content_proto_goTypes = []any{
	(*ExtractContentFromLinksRequest)(nil), // 0: google.search.v1.ExtractContentFromLinksRequest
	(*ExtractContentFromLinksReponse)(nil), // 1: google.search.v1.ExtractContentFromLinksReponse
	(*ExtractedContent)(nil),               // 2: google.search.v1.ExtractedContent
}
var file_api_content_proto_depIdxs = []int32{
	2, // 0: google.search.v1.ExtractContentFromLinksReponse.contents:type_name -> google.search.v1.ExtractedContent
	0, // 1: google.search.v1.ContentService.ExtractContentFromLinks:input_type -> google.search.v1.ExtractContentFromLinksRequest
	1, // 2: google.search.v1.ContentService.ExtractContentFromLinks:output_type -> google.search.v1.ExtractContentFromLinksReponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_content_proto_init() }
func file_api_content_proto_init() {
	if File_api_content_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_content_proto_rawDesc), len(file_api_content_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_content_proto_goTypes,
		DependencyIndexes: file_api_content_proto_depIdxs,
		MessageInfos:      file_api_content_proto_msgTypes,
	}.Build()
	File_api_content_proto = out.File
	file_api_content_proto_goTypes = nil
	file_api_content_proto_depIdxs = nil
}
