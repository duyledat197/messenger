// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.14.0
// source: chat/message.proto

package chat

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Message struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChannelId     int64                  `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	MessageId     int64                  `protobuf:"varint,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	UserId        string                 `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content       string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Reaction      []string               `protobuf:"bytes,5,rep,name=reaction,proto3" json:"reaction,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_chat_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_chat_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_chat_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *Message) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *Message) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetReaction() []string {
	if x != nil {
		return x.Reaction
	}
	return nil
}

type SendMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChannelId     int64                  `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Content       string                 `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	UserId        string                 `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendMessageRequest) Reset() {
	*x = SendMessageRequest{}
	mi := &file_chat_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageRequest) ProtoMessage() {}

func (x *SendMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageRequest.ProtoReflect.Descriptor instead.
func (*SendMessageRequest) Descriptor() ([]byte, []int) {
	return file_chat_message_proto_rawDescGZIP(), []int{1}
}

func (x *SendMessageRequest) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *SendMessageRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *SendMessageRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type SendMessageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MessageId     int64                  `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendMessageResponse) Reset() {
	*x = SendMessageResponse{}
	mi := &file_chat_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageResponse) ProtoMessage() {}

func (x *SendMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chat_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageResponse.ProtoReflect.Descriptor instead.
func (*SendMessageResponse) Descriptor() ([]byte, []int) {
	return file_chat_message_proto_rawDescGZIP(), []int{2}
}

func (x *SendMessageResponse) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

type GetMessageListChannelRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChannelId     int64                  `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Offset        int64                  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit         int64                  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMessageListChannelRequest) Reset() {
	*x = GetMessageListChannelRequest{}
	mi := &file_chat_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMessageListChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageListChannelRequest) ProtoMessage() {}

func (x *GetMessageListChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageListChannelRequest.ProtoReflect.Descriptor instead.
func (*GetMessageListChannelRequest) Descriptor() ([]byte, []int) {
	return file_chat_message_proto_rawDescGZIP(), []int{3}
}

func (x *GetMessageListChannelRequest) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *GetMessageListChannelRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *GetMessageListChannelRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type GetMessageListChannelResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Messages      []*Message             `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMessageListChannelResponse) Reset() {
	*x = GetMessageListChannelResponse{}
	mi := &file_chat_message_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMessageListChannelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageListChannelResponse) ProtoMessage() {}

func (x *GetMessageListChannelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chat_message_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageListChannelResponse.ProtoReflect.Descriptor instead.
func (*GetMessageListChannelResponse) Descriptor() ([]byte, []int) {
	return file_chat_message_proto_rawDescGZIP(), []int{4}
}

func (x *GetMessageListChannelResponse) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

var File_chat_message_proto protoreflect.FileDescriptor

const file_chat_message_proto_rawDesc = "" +
	"\n" +
	"\x12chat/message.proto\x12\x04chat\x1a\x1cgoogle/api/annotations.proto\"\x96\x01\n" +
	"\aMessage\x12\x1d\n" +
	"\n" +
	"channel_id\x18\x01 \x01(\x03R\tchannelId\x12\x1d\n" +
	"\n" +
	"message_id\x18\x02 \x01(\x03R\tmessageId\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\tR\x06userId\x12\x18\n" +
	"\acontent\x18\x04 \x01(\tR\acontent\x12\x1a\n" +
	"\breaction\x18\x05 \x03(\tR\breaction\"f\n" +
	"\x12SendMessageRequest\x12\x1d\n" +
	"\n" +
	"channel_id\x18\x01 \x01(\x03R\tchannelId\x12\x18\n" +
	"\acontent\x18\x02 \x01(\tR\acontent\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\tR\x06userId\"4\n" +
	"\x13SendMessageResponse\x12\x1d\n" +
	"\n" +
	"message_id\x18\x01 \x01(\x03R\tmessageId\"k\n" +
	"\x1cGetMessageListChannelRequest\x12\x1d\n" +
	"\n" +
	"channel_id\x18\x01 \x01(\x03R\tchannelId\x12\x16\n" +
	"\x06offset\x18\x02 \x01(\x03R\x06offset\x12\x14\n" +
	"\x05limit\x18\x03 \x01(\x03R\x05limit\"J\n" +
	"\x1dGetMessageListChannelResponse\x12)\n" +
	"\bmessages\x18\x01 \x03(\v2\r.chat.MessageR\bmessages2\xe6\x01\n" +
	"\x0eMessageService\x12B\n" +
	"\vSendMessage\x12\x18.chat.SendMessageRequest\x1a\x19.chat.SendMessageResponse\x12\x8f\x01\n" +
	"\x15GetMessageListChannel\x12\".chat.GetMessageListChannelRequest\x1a#.chat.GetMessageListChannelResponse\"-\x82\xd3\xe4\x93\x02':\x01*\"\"/v1/channels/{channel_id}/messagesB\x1cZ\x1aopenmyth/messgener/pb/chatb\x06proto3"

var (
	file_chat_message_proto_rawDescOnce sync.Once
	file_chat_message_proto_rawDescData []byte
)

func file_chat_message_proto_rawDescGZIP() []byte {
	file_chat_message_proto_rawDescOnce.Do(func() {
		file_chat_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_chat_message_proto_rawDesc), len(file_chat_message_proto_rawDesc)))
	})
	return file_chat_message_proto_rawDescData
}

var file_chat_message_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_chat_message_proto_goTypes = []any{
	(*Message)(nil),                       // 0: chat.Message
	(*SendMessageRequest)(nil),            // 1: chat.SendMessageRequest
	(*SendMessageResponse)(nil),           // 2: chat.SendMessageResponse
	(*GetMessageListChannelRequest)(nil),  // 3: chat.GetMessageListChannelRequest
	(*GetMessageListChannelResponse)(nil), // 4: chat.GetMessageListChannelResponse
}
var file_chat_message_proto_depIdxs = []int32{
	0, // 0: chat.GetMessageListChannelResponse.messages:type_name -> chat.Message
	1, // 1: chat.MessageService.SendMessage:input_type -> chat.SendMessageRequest
	3, // 2: chat.MessageService.GetMessageListChannel:input_type -> chat.GetMessageListChannelRequest
	2, // 3: chat.MessageService.SendMessage:output_type -> chat.SendMessageResponse
	4, // 4: chat.MessageService.GetMessageListChannel:output_type -> chat.GetMessageListChannelResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_chat_message_proto_init() }
func file_chat_message_proto_init() {
	if File_chat_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_chat_message_proto_rawDesc), len(file_chat_message_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_message_proto_goTypes,
		DependencyIndexes: file_chat_message_proto_depIdxs,
		MessageInfos:      file_chat_message_proto_msgTypes,
	}.Build()
	File_chat_message_proto = out.File
	file_chat_message_proto_goTypes = nil
	file_chat_message_proto_depIdxs = nil
}
