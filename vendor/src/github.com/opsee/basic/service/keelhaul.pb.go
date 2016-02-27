// Code generated by protoc-gen-gogo.
// source: keelhaul.proto
// DO NOT EDIT!

package service

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/opsee/protobuf/opseeproto"
import _ "github.com/opsee/protobuf/opseeproto/types"
import opsee1 "github.com/opsee/basic/schema"

import github_com_graphql_go_graphql "github.com/graphql-go/graphql"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ListBastionsRequest struct {
	CustomerIds []string `protobuf:"bytes,1,rep,name=customer_ids" json:"customer_ids,omitempty"`
}

func (m *ListBastionsRequest) Reset()         { *m = ListBastionsRequest{} }
func (m *ListBastionsRequest) String() string { return proto.CompactTextString(m) }
func (*ListBastionsRequest) ProtoMessage()    {}

type ListBastionsResponse struct {
	Bastions []*opsee1.Bastion `protobuf:"bytes,1,rep,name=bastions" json:"bastions,omitempty"`
}

func (m *ListBastionsResponse) Reset()         { *m = ListBastionsResponse{} }
func (m *ListBastionsResponse) String() string { return proto.CompactTextString(m) }
func (*ListBastionsResponse) ProtoMessage()    {}

func (m *ListBastionsResponse) GetBastions() []*opsee1.Bastion {
	if m != nil {
		return m.Bastions
	}
	return nil
}

func init() {
	proto.RegisterType((*ListBastionsRequest)(nil), "opsee.ListBastionsRequest")
	proto.RegisterType((*ListBastionsResponse)(nil), "opsee.ListBastionsResponse")
}
func (this *ListBastionsRequest) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*ListBastionsRequest)
	if !ok {
		that2, ok := that.(ListBastionsRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if len(this.CustomerIds) != len(that1.CustomerIds) {
		return false
	}
	for i := range this.CustomerIds {
		if this.CustomerIds[i] != that1.CustomerIds[i] {
			return false
		}
	}
	return true
}
func (this *ListBastionsResponse) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*ListBastionsResponse)
	if !ok {
		that2, ok := that.(ListBastionsResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if len(this.Bastions) != len(that1.Bastions) {
		return false
	}
	for i := range this.Bastions {
		if !this.Bastions[i].Equal(that1.Bastions[i]) {
			return false
		}
	}
	return true
}

type ListBastionsRequestGetter interface {
	GetListBastionsRequest() *ListBastionsRequest
}

var GraphQLListBastionsRequestType *github_com_graphql_go_graphql.Object

type ListBastionsResponseGetter interface {
	GetListBastionsResponse() *ListBastionsResponse
}

var GraphQLListBastionsResponseType *github_com_graphql_go_graphql.Object

func init() {
	GraphQLListBastionsRequestType = github_com_graphql_go_graphql.NewObject(github_com_graphql_go_graphql.ObjectConfig{
		Name:        "serviceListBastionsRequest",
		Description: "",
		Fields: (github_com_graphql_go_graphql.FieldsThunk)(func() github_com_graphql_go_graphql.Fields {
			return github_com_graphql_go_graphql.Fields{
				"customer_ids": &github_com_graphql_go_graphql.Field{
					Type:        github_com_graphql_go_graphql.NewList(github_com_graphql_go_graphql.String),
					Description: "",
					Resolve: func(p github_com_graphql_go_graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*ListBastionsRequest)
						if ok {
							return obj.CustomerIds, nil
						}
						inter, ok := p.Source.(ListBastionsRequestGetter)
						if ok {
							face := inter.GetListBastionsRequest()
							if face == nil {
								return nil, nil
							}
							return face.CustomerIds, nil
						}
						return nil, fmt.Errorf("field customer_ids not resolved")
					},
				},
			}
		}),
	})
	GraphQLListBastionsResponseType = github_com_graphql_go_graphql.NewObject(github_com_graphql_go_graphql.ObjectConfig{
		Name:        "serviceListBastionsResponse",
		Description: "",
		Fields: (github_com_graphql_go_graphql.FieldsThunk)(func() github_com_graphql_go_graphql.Fields {
			return github_com_graphql_go_graphql.Fields{
				"bastions": &github_com_graphql_go_graphql.Field{
					Type:        github_com_graphql_go_graphql.NewList(opsee1.GraphQLBastionType),
					Description: "",
					Resolve: func(p github_com_graphql_go_graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*ListBastionsResponse)
						if ok {
							return obj.Bastions, nil
						}
						inter, ok := p.Source.(ListBastionsResponseGetter)
						if ok {
							face := inter.GetListBastionsResponse()
							if face == nil {
								return nil, nil
							}
							return face.Bastions, nil
						}
						return nil, fmt.Errorf("field bastions not resolved")
					},
				},
			}
		}),
	})
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Keelhaul service

type KeelhaulClient interface {
	ListBastions(ctx context.Context, in *ListBastionsRequest, opts ...grpc.CallOption) (*ListBastionsResponse, error)
}

type keelhaulClient struct {
	cc *grpc.ClientConn
}

func NewKeelhaulClient(cc *grpc.ClientConn) KeelhaulClient {
	return &keelhaulClient{cc}
}

func (c *keelhaulClient) ListBastions(ctx context.Context, in *ListBastionsRequest, opts ...grpc.CallOption) (*ListBastionsResponse, error) {
	out := new(ListBastionsResponse)
	err := grpc.Invoke(ctx, "/opsee.Keelhaul/ListBastions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Keelhaul service

type KeelhaulServer interface {
	ListBastions(context.Context, *ListBastionsRequest) (*ListBastionsResponse, error)
}

func RegisterKeelhaulServer(s *grpc.Server, srv KeelhaulServer) {
	s.RegisterService(&_Keelhaul_serviceDesc, srv)
}

func _Keelhaul_ListBastions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(ListBastionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(KeelhaulServer).ListBastions(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Keelhaul_serviceDesc = grpc.ServiceDesc{
	ServiceName: "opsee.Keelhaul",
	HandlerType: (*KeelhaulServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListBastions",
			Handler:    _Keelhaul_ListBastions_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func NewPopulatedListBastionsRequest(r randyKeelhaul, easy bool) *ListBastionsRequest {
	this := &ListBastionsRequest{}
	v1 := r.Intn(10)
	this.CustomerIds = make([]string, v1)
	for i := 0; i < v1; i++ {
		this.CustomerIds[i] = randStringKeelhaul(r)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedListBastionsResponse(r randyKeelhaul, easy bool) *ListBastionsResponse {
	this := &ListBastionsResponse{}
	if r.Intn(10) != 0 {
		v2 := r.Intn(5)
		this.Bastions = make([]*opsee1.Bastion, v2)
		for i := 0; i < v2; i++ {
			this.Bastions[i] = opsee1.NewPopulatedBastion(r, easy)
		}
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyKeelhaul interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneKeelhaul(r randyKeelhaul) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringKeelhaul(r randyKeelhaul) string {
	v3 := r.Intn(100)
	tmps := make([]rune, v3)
	for i := 0; i < v3; i++ {
		tmps[i] = randUTF8RuneKeelhaul(r)
	}
	return string(tmps)
}
func randUnrecognizedKeelhaul(r randyKeelhaul, maxFieldNumber int) (data []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		data = randFieldKeelhaul(data, r, fieldNumber, wire)
	}
	return data
}
func randFieldKeelhaul(data []byte, r randyKeelhaul, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		data = encodeVarintPopulateKeelhaul(data, uint64(key))
		v4 := r.Int63()
		if r.Intn(2) == 0 {
			v4 *= -1
		}
		data = encodeVarintPopulateKeelhaul(data, uint64(v4))
	case 1:
		data = encodeVarintPopulateKeelhaul(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		data = encodeVarintPopulateKeelhaul(data, uint64(key))
		ll := r.Intn(100)
		data = encodeVarintPopulateKeelhaul(data, uint64(ll))
		for j := 0; j < ll; j++ {
			data = append(data, byte(r.Intn(256)))
		}
	default:
		data = encodeVarintPopulateKeelhaul(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return data
}
func encodeVarintPopulateKeelhaul(data []byte, v uint64) []byte {
	for v >= 1<<7 {
		data = append(data, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	data = append(data, uint8(v))
	return data
}
