// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: artist.proto

package musicalpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type Artist struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Author        string                 `protobuf:"bytes,3,opt,name=author,proto3" json:"author,omitempty"`
	Producer      string                 `protobuf:"bytes,4,opt,name=producer,proto3" json:"producer,omitempty"`
	Country       string                 `protobuf:"bytes,5,opt,name=country,proto3" json:"country,omitempty"`
	Description   string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Artist) Reset() {
	*x = Artist{}
	mi := &file_artist_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Artist) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Artist) ProtoMessage() {}

func (x *Artist) ProtoReflect() protoreflect.Message {
	mi := &file_artist_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Artist.ProtoReflect.Descriptor instead.
func (*Artist) Descriptor() ([]byte, []int) {
	return file_artist_proto_rawDescGZIP(), []int{0}
}

func (x *Artist) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Artist) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Artist) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *Artist) GetProducer() string {
	if x != nil {
		return x.Producer
	}
	return ""
}

func (x *Artist) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *Artist) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Artist) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Artist) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type ListArtistsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	SearchQuery   *string                `protobuf:"bytes,3,opt,name=search_query,json=searchQuery,proto3,oneof" json:"search_query,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListArtistsRequest) Reset() {
	*x = ListArtistsRequest{}
	mi := &file_artist_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListArtistsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListArtistsRequest) ProtoMessage() {}

func (x *ListArtistsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_artist_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListArtistsRequest.ProtoReflect.Descriptor instead.
func (*ListArtistsRequest) Descriptor() ([]byte, []int) {
	return file_artist_proto_rawDescGZIP(), []int{1}
}

func (x *ListArtistsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListArtistsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListArtistsRequest) GetSearchQuery() string {
	if x != nil && x.SearchQuery != nil {
		return *x.SearchQuery
	}
	return ""
}

type ListArtistsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Artists       []*Artist              `protobuf:"bytes,1,rep,name=artists,proto3" json:"artists,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListArtistsResponse) Reset() {
	*x = ListArtistsResponse{}
	mi := &file_artist_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListArtistsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListArtistsResponse) ProtoMessage() {}

func (x *ListArtistsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_artist_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListArtistsResponse.ProtoReflect.Descriptor instead.
func (*ListArtistsResponse) Descriptor() ([]byte, []int) {
	return file_artist_proto_rawDescGZIP(), []int{2}
}

func (x *ListArtistsResponse) GetArtists() []*Artist {
	if x != nil {
		return x.Artists
	}
	return nil
}

var File_artist_proto protoreflect.FileDescriptor

const file_artist_proto_rawDesc = "" +
	"\n" +
	"\fartist.proto\x12\amusical\x1a\x1fgoogle/protobuf/timestamp.proto\"\x92\x02\n" +
	"\x06Artist\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n" +
	"\x06author\x18\x03 \x01(\tR\x06author\x12\x1a\n" +
	"\bproducer\x18\x04 \x01(\tR\bproducer\x12\x18\n" +
	"\acountry\x18\x05 \x01(\tR\acountry\x12 \n" +
	"\vdescription\x18\x06 \x01(\tR\vdescription\x129\n" +
	"\n" +
	"created_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"~\n" +
	"\x12ListArtistsRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\x05R\bpageSize\x12&\n" +
	"\fsearch_query\x18\x03 \x01(\tH\x00R\vsearchQuery\x88\x01\x01B\x0f\n" +
	"\r_search_query\"@\n" +
	"\x13ListArtistsResponse\x12)\n" +
	"\aartists\x18\x01 \x03(\v2\x0f.musical.ArtistR\aartistsBCZAgithub.com/TeoPlow/online-music-service/src/musical/pkg/musicalpbb\x06proto3"

var (
	file_artist_proto_rawDescOnce sync.Once
	file_artist_proto_rawDescData []byte
)

func file_artist_proto_rawDescGZIP() []byte {
	file_artist_proto_rawDescOnce.Do(func() {
		file_artist_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_artist_proto_rawDesc), len(file_artist_proto_rawDesc)))
	})
	return file_artist_proto_rawDescData
}

var file_artist_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_artist_proto_goTypes = []any{
	(*Artist)(nil),                // 0: musical.Artist
	(*ListArtistsRequest)(nil),    // 1: musical.ListArtistsRequest
	(*ListArtistsResponse)(nil),   // 2: musical.ListArtistsResponse
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_artist_proto_depIdxs = []int32{
	3, // 0: musical.Artist.created_at:type_name -> google.protobuf.Timestamp
	3, // 1: musical.Artist.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: musical.ListArtistsResponse.artists:type_name -> musical.Artist
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_artist_proto_init() }
func file_artist_proto_init() {
	if File_artist_proto != nil {
		return
	}
	file_artist_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_artist_proto_rawDesc), len(file_artist_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_artist_proto_goTypes,
		DependencyIndexes: file_artist_proto_depIdxs,
		MessageInfos:      file_artist_proto_msgTypes,
	}.Build()
	File_artist_proto = out.File
	file_artist_proto_goTypes = nil
	file_artist_proto_depIdxs = nil
}
