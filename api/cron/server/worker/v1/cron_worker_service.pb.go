// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: api/cron/server/worker/cron_worker_service.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_api_cron_server_worker_cron_worker_service_proto protoreflect.FileDescriptor

var file_api_cron_server_worker_cron_worker_service_proto_rawDesc = []byte{
	0x0a, 0x30, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x5f, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x5f,
	0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x28, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x5f, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xef, 0x0d, 0x0a, 0x06, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x9f, 0x01, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x35, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e,
	0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x33, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x12, 0x19, 0x2f, 0x63,
	0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0xa3, 0x01, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x36, 0x2e, 0x63, 0x72,
	0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72,
	0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x22, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1c, 0x12, 0x1a, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0xab, 0x01,
	0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x12, 0x38, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e,
	0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x3a, 0x01, 0x2a,
	0x22, 0x19, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0xab, 0x01, 0x0a, 0x11,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x12, 0x38, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x63, 0x72,
	0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72,
	0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x3a, 0x01, 0x2a, 0x1a, 0x19,
	0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0xa8, 0x01, 0x0a, 0x11, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x38, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x63, 0x72, 0x6f, 0x6e,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e,
	0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x2a, 0x19, 0x2f, 0x63, 0x72, 0x6f, 0x6e,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x12, 0x8a, 0x01, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x12, 0x30, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12, 0x13, 0x2f, 0x63,
	0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x12, 0x8e, 0x01, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x12, 0x31, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x63,
	0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x73, 0x12, 0x96, 0x01, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x12, 0x33, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e,
	0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x1e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x22, 0x13, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x96, 0x01, 0x0a, 0x0c,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x33, 0x2e, 0x63,
	0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63,
	0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x31, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x1a,
	0x13, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x12, 0xaf, 0x01, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x2e, 0x63, 0x72,
	0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72,
	0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x37, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x25, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x3a, 0x01, 0x2a, 0x1a, 0x1a, 0x2f, 0x63, 0x72, 0x6f,
	0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2f,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x93, 0x01, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x33, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x63,
	0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63,
	0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x2a, 0x13, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x42, 0x35, 0x0a, 0x1e,
	0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x08,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x56, 0x31, 0x50, 0x01, 0x5a, 0x07, 0x2e, 0x2f, 0x76, 0x31,
	0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_cron_server_worker_cron_worker_service_proto_goTypes = []interface{}{
	(*GetWorkerGroupRequest)(nil),     // 0: cron.api.server.cron.worker.v1.GetWorkerGroupRequest
	(*ListWorkerGroupRequest)(nil),    // 1: cron.api.server.cron.worker.v1.ListWorkerGroupRequest
	(*CreateWorkerGroupRequest)(nil),  // 2: cron.api.server.cron.worker.v1.CreateWorkerGroupRequest
	(*UpdateWorkerGroupRequest)(nil),  // 3: cron.api.server.cron.worker.v1.UpdateWorkerGroupRequest
	(*DeleteWorkerGroupRequest)(nil),  // 4: cron.api.server.cron.worker.v1.DeleteWorkerGroupRequest
	(*GetWorkerRequest)(nil),          // 5: cron.api.server.cron.worker.v1.GetWorkerRequest
	(*ListWorkerRequest)(nil),         // 6: cron.api.server.cron.worker.v1.ListWorkerRequest
	(*CreateWorkerRequest)(nil),       // 7: cron.api.server.cron.worker.v1.CreateWorkerRequest
	(*UpdateWorkerRequest)(nil),       // 8: cron.api.server.cron.worker.v1.UpdateWorkerRequest
	(*UpdateWorkerStatusRequest)(nil), // 9: cron.api.server.cron.worker.v1.UpdateWorkerStatusRequest
	(*DeleteWorkerRequest)(nil),       // 10: cron.api.server.cron.worker.v1.DeleteWorkerRequest
	(*GetWorkerGroupReply)(nil),       // 11: cron.api.server.cron.worker.v1.GetWorkerGroupReply
	(*ListWorkerGroupReply)(nil),      // 12: cron.api.server.cron.worker.v1.ListWorkerGroupReply
	(*CreateWorkerGroupReply)(nil),    // 13: cron.api.server.cron.worker.v1.CreateWorkerGroupReply
	(*UpdateWorkerGroupReply)(nil),    // 14: cron.api.server.cron.worker.v1.UpdateWorkerGroupReply
	(*DeleteWorkerGroupReply)(nil),    // 15: cron.api.server.cron.worker.v1.DeleteWorkerGroupReply
	(*GetWorkerReply)(nil),            // 16: cron.api.server.cron.worker.v1.GetWorkerReply
	(*ListWorkerReply)(nil),           // 17: cron.api.server.cron.worker.v1.ListWorkerReply
	(*CreateWorkerReply)(nil),         // 18: cron.api.server.cron.worker.v1.CreateWorkerReply
	(*UpdateWorkerReply)(nil),         // 19: cron.api.server.cron.worker.v1.UpdateWorkerReply
	(*UpdateWorkerStatusReply)(nil),   // 20: cron.api.server.cron.worker.v1.UpdateWorkerStatusReply
	(*DeleteWorkerReply)(nil),         // 21: cron.api.server.cron.worker.v1.DeleteWorkerReply
}
var file_api_cron_server_worker_cron_worker_service_proto_depIdxs = []int32{
	0,  // 0: cron.api.server.cron.worker.v1.Worker.GetWorkerGroup:input_type -> cron.api.server.cron.worker.v1.GetWorkerGroupRequest
	1,  // 1: cron.api.server.cron.worker.v1.Worker.ListWorkerGroup:input_type -> cron.api.server.cron.worker.v1.ListWorkerGroupRequest
	2,  // 2: cron.api.server.cron.worker.v1.Worker.CreateWorkerGroup:input_type -> cron.api.server.cron.worker.v1.CreateWorkerGroupRequest
	3,  // 3: cron.api.server.cron.worker.v1.Worker.UpdateWorkerGroup:input_type -> cron.api.server.cron.worker.v1.UpdateWorkerGroupRequest
	4,  // 4: cron.api.server.cron.worker.v1.Worker.DeleteWorkerGroup:input_type -> cron.api.server.cron.worker.v1.DeleteWorkerGroupRequest
	5,  // 5: cron.api.server.cron.worker.v1.Worker.GetWorker:input_type -> cron.api.server.cron.worker.v1.GetWorkerRequest
	6,  // 6: cron.api.server.cron.worker.v1.Worker.ListWorker:input_type -> cron.api.server.cron.worker.v1.ListWorkerRequest
	7,  // 7: cron.api.server.cron.worker.v1.Worker.CreateWorker:input_type -> cron.api.server.cron.worker.v1.CreateWorkerRequest
	8,  // 8: cron.api.server.cron.worker.v1.Worker.UpdateWorker:input_type -> cron.api.server.cron.worker.v1.UpdateWorkerRequest
	9,  // 9: cron.api.server.cron.worker.v1.Worker.UpdateWorkerStatus:input_type -> cron.api.server.cron.worker.v1.UpdateWorkerStatusRequest
	10, // 10: cron.api.server.cron.worker.v1.Worker.DeleteWorker:input_type -> cron.api.server.cron.worker.v1.DeleteWorkerRequest
	11, // 11: cron.api.server.cron.worker.v1.Worker.GetWorkerGroup:output_type -> cron.api.server.cron.worker.v1.GetWorkerGroupReply
	12, // 12: cron.api.server.cron.worker.v1.Worker.ListWorkerGroup:output_type -> cron.api.server.cron.worker.v1.ListWorkerGroupReply
	13, // 13: cron.api.server.cron.worker.v1.Worker.CreateWorkerGroup:output_type -> cron.api.server.cron.worker.v1.CreateWorkerGroupReply
	14, // 14: cron.api.server.cron.worker.v1.Worker.UpdateWorkerGroup:output_type -> cron.api.server.cron.worker.v1.UpdateWorkerGroupReply
	15, // 15: cron.api.server.cron.worker.v1.Worker.DeleteWorkerGroup:output_type -> cron.api.server.cron.worker.v1.DeleteWorkerGroupReply
	16, // 16: cron.api.server.cron.worker.v1.Worker.GetWorker:output_type -> cron.api.server.cron.worker.v1.GetWorkerReply
	17, // 17: cron.api.server.cron.worker.v1.Worker.ListWorker:output_type -> cron.api.server.cron.worker.v1.ListWorkerReply
	18, // 18: cron.api.server.cron.worker.v1.Worker.CreateWorker:output_type -> cron.api.server.cron.worker.v1.CreateWorkerReply
	19, // 19: cron.api.server.cron.worker.v1.Worker.UpdateWorker:output_type -> cron.api.server.cron.worker.v1.UpdateWorkerReply
	20, // 20: cron.api.server.cron.worker.v1.Worker.UpdateWorkerStatus:output_type -> cron.api.server.cron.worker.v1.UpdateWorkerStatusReply
	21, // 21: cron.api.server.cron.worker.v1.Worker.DeleteWorker:output_type -> cron.api.server.cron.worker.v1.DeleteWorkerReply
	11, // [11:22] is the sub-list for method output_type
	0,  // [0:11] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_cron_server_worker_cron_worker_service_proto_init() }
func file_api_cron_server_worker_cron_worker_service_proto_init() {
	if File_api_cron_server_worker_cron_worker_service_proto != nil {
		return
	}
	file_api_cron_server_worker_cron_worker_group_proto_init()
	file_api_cron_server_worker_cron_worker_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_cron_server_worker_cron_worker_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_cron_server_worker_cron_worker_service_proto_goTypes,
		DependencyIndexes: file_api_cron_server_worker_cron_worker_service_proto_depIdxs,
	}.Build()
	File_api_cron_server_worker_cron_worker_service_proto = out.File
	file_api_cron_server_worker_cron_worker_service_proto_rawDesc = nil
	file_api_cron_server_worker_cron_worker_service_proto_goTypes = nil
	file_api_cron_server_worker_cron_worker_service_proto_depIdxs = nil
}