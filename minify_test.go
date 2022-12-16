package structpbConv

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

type TestCase[T proto.Message] struct {
	name string
	args T
	want T
}

func TestMinifyList(t *testing.T) {
	tests := []TestCase[*structpb.ListValue]{
		{name: "Nil"},
		{name: "Empty list", args: &structpb.ListValue{}},
		{
			name: "Non empty list",
			args: &structpb.ListValue{Values: []*structpb.Value{structpb.NewStringValue("test-value")}},
			want: &structpb.ListValue{Values: []*structpb.Value{structpb.NewStringValue("test-value")}},
		},
		{
			name: "Sparse empty list",
			args: &structpb.ListValue{Values: []*structpb.Value{structpb.NewStringValue("test-value"), nil, structpb.NewStringValue("another-value")}},
			want: &structpb.ListValue{Values: []*structpb.Value{structpb.NewStringValue("test-value"), structpb.NewStringValue("another-value")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := MinifyList(tt.args)
			if nil == tt.want {
				require.Nil(t, output)
			} else {
				require.True(t, proto.Equal(output, tt.want))
			}
		})
	}
}

func TestMinifyStruct(t *testing.T) {
	tests := []TestCase[*structpb.Struct]{
		{name: "Nil"},
		{name: "Empty map", args: &structpb.Struct{}},
		{name: "Map of nils map", args: &structpb.Struct{Fields: map[string]*structpb.Value{
			"nil-key":   nil,
			"null-key":  structpb.NewNullValue(),
			"null-kind": {},
		}}},
		{
			name: "Non empty map",
			args: &structpb.Struct{Fields: map[string]*structpb.Value{
				"string-key":  structpb.NewStringValue(""),
				"bool-key":    structpb.NewBoolValue(true),
				"numeric-key": structpb.NewNumberValue(123),
				"struct-key": structpb.NewStructValue(
					&structpb.Struct{Fields: map[string]*structpb.Value{
						"string-key": structpb.NewStringValue(""),
					}},
				),
				"list-key": structpb.NewListValue(
					&structpb.ListValue{Values: []*structpb.Value{
						structpb.NewStringValue(""),
					}},
				),
			}},
			want: &structpb.Struct{Fields: map[string]*structpb.Value{
				"string-key":  structpb.NewStringValue(""),
				"bool-key":    structpb.NewBoolValue(true),
				"numeric-key": structpb.NewNumberValue(123),
				"struct-key": structpb.NewStructValue(
					&structpb.Struct{Fields: map[string]*structpb.Value{
						"string-key": structpb.NewStringValue(""),
					}},
				),
				"list-key": structpb.NewListValue(
					&structpb.ListValue{Values: []*structpb.Value{
						structpb.NewStringValue(""),
					}},
				),
			}},
		},
		{
			name: "Sparse map",
			args: &structpb.Struct{Fields: map[string]*structpb.Value{
				"nil-key":     nil,
				"null-key":    structpb.NewNullValue(),
				"string-key":  structpb.NewStringValue(""),
				"bool-key":    structpb.NewBoolValue(true),
				"numeric-key": structpb.NewNumberValue(123),
				"struct-key": structpb.NewStructValue(
					&structpb.Struct{Fields: map[string]*structpb.Value{
						"nil-key":    nil,
						"null-key":   structpb.NewNullValue(),
						"string-key": structpb.NewStringValue(""),
					}},
				),
				"list-key": structpb.NewListValue(
					&structpb.ListValue{Values: []*structpb.Value{
						nil,
						structpb.NewNullValue(),
						structpb.NewStringValue(""),
					}},
				),
			}},
			want: &structpb.Struct{Fields: map[string]*structpb.Value{
				"string-key":  structpb.NewStringValue(""),
				"bool-key":    structpb.NewBoolValue(true),
				"numeric-key": structpb.NewNumberValue(123),
				"struct-key": structpb.NewStructValue(
					&structpb.Struct{Fields: map[string]*structpb.Value{
						"string-key": structpb.NewStringValue(""),
					}},
				),
				"list-key": structpb.NewListValue(
					&structpb.ListValue{Values: []*structpb.Value{
						structpb.NewStringValue(""),
					}},
				),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := MinifyStruct(tt.args)
			if nil == tt.want {
				require.Nil(t, output)
			} else {
				require.True(t, proto.Equal(output, tt.want))
			}
		})
	}
}

func TestMinifyValue(t *testing.T) {
	tests := []TestCase[*structpb.Value]{
		{name: "Nil"},
		{name: "Empty value", args: &structpb.Value{}},
		{
			name: "Null value",
			args: &structpb.Value{Kind: &structpb.Value_NullValue{}},
			want: nil,
		},
		{
			name: "String value",
			args: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "test-value"}},
			want: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "test-value"}},
		},
		{
			name: "Bool value",
			args: &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: true}},
			want: &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: true}},
		},
		{
			name: "Numeric value",
			args: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 123}},
			want: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 123}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := MinifyValue(tt.args)
			if nil == tt.want {
				require.Nil(t, output)
			} else {
				require.True(t, proto.Equal(output, tt.want))
			}
		})
	}
}

func TestMinify(t *testing.T) {
	const input = `
	{
		"key": "value",
		"nil": null,
		"empty-list": [],
		"nested-struct": {
			"nested": "value"
		}
	}
	`

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(input), &data)
	require.NoError(t, err)
	require.NotNil(t, data)

	structpbData, _ := structpb.NewStruct(data)
	output := MinifyStruct(structpbData)

	outputJson, _ := json.Marshal(output)
	require.Equal(t, string(outputJson), `{"key":"value","nested-struct":{"nested":"value"}}`)
}
