// Copyright (c) ArenaML Labs Pvt Ltd.

package helper

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kr/pretty"
)

// Example usage with your structs
type JSONEngineConfig struct {
	Id       *string        `json:"id,omitempty"`
	Inactive *bool          `json:"inactive,omitempty"`
	Name     string         `json:"name,omitempty"`
	Index    *float64       `json:"index"`
	PStr     *TFEngineModel `json:"pStr"`
}

type TFEngineModel struct {
	Id       types.String      `tfsdk:"id"`
	Inactive types.Bool        `tfsdk:"inactive"`
	Kind     types.String      `tfsdk:"kind"`
	Name     types.String      `tfsdk:"name"`
	Index    types.Float64     `tfsdk:"index"`
	PStr     *JSONEngineConfig `tfsdk:"pStr"`
}

func forTestJSONObjWithEmb() JSONEngineConfig {
	return JSONEngineConfig{
		Id:       PointerOf("some-id"),
		Inactive: PointerOf(false),
		Name:     "name-yp",
		Index:    PointerOf(1024.0),
	}
}

func forEmbedTestTFObj() *TFEngineModel {
	return &TFEngineModel{
		Id:       types.StringValue("some-id"),
		Inactive: types.BoolValue(false),
		Name:     types.StringValue("name-yp"),
		Index:    types.Float64Value(1024),
	}
}

func TestConvertStructToTerraformModel(t *testing.T) {
	type args[T any] struct {
		source T
	}
	type testCase[T any, U any] struct {
		name string
		args args[T]
		want *U
	}
	tests := []testCase[JSONEngineConfig, TFEngineModel]{
		{
			name: "json-tf",
			args: args[JSONEngineConfig]{
				source: forTestJSONObjWithEmb(),
			},
			want: forEmbedTestTFObj(),
		},
	}
	resp := TFEngineModel{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ConvertJSONStructToSimpleTF(context.Background(), tt.args.source, &resp); !reflect.DeepEqual(&resp, tt.want) {
				t.Errorf("ConvertJSONStructToSimpleTF() = %v, want %v", &resp, tt.want)
			}
		})
	}
}

// PointerOf returns a pointer to a .
func PointerOf[A any](a A) *A {
	return &a
}

type JsonTFTagTest struct {
	Abc string `json:"abc"`
	XYZ int64  `json:"xyz"`
}

type JsonTFTagSrc struct {
	Abc types.String `tfsdk:"abc"`
	XYZ types.Int64  `tfsdk:"xyz"`
}

func forTestTfSrc() JsonTFTagSrc {
	return JsonTFTagSrc{
		Abc: types.StringValue("abc"),
		XYZ: types.Int64Value(12),
	}
}

func forTestGoJSONDest() JsonTFTagTest {
	return JsonTFTagTest{
		Abc: "abc",
		XYZ: 12,
	}
}

func TestTFtoJSONObj(t *testing.T) {
	type args struct {
		src interface{}
		obj interface{}
	}

	s1 := forTestTfSrc()
	t1 := &JsonTFTagTest{}
	w1 := forTestGoJSONDest()

	s2 := forEmbedTestTFObj()
	t2 := &JSONEngineConfig{}
	w2 := forTestJSONObjWithEmb()

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				src: s1,
				obj: t1,
			},
			want: &w1,
		}, {
			name: "embedded-struct",
			args: args{
				src: s2,
				obj: t2,
			},
			want: &w2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := tt.args.src
			dest := tt.args.obj
			err := ConvertTfModelToApiJSON(context.Background(), src, dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertTfModelToApiJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, dest) {
				t.Errorf("ConvertTfModelToApiJSON() got = %v, want %v", dest, tt.want)
			}
			pretty.Println(dest)

		})
	}
}
