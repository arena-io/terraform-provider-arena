// Copyright (c) ArenaML Labs Pvt Ltd.

package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TFType string

const (
	TfString    TFType = "string"
	TfInt       TFType = "int"
	TfInt64     TFType = "int64"
	TfBoolean   TFType = "boolean"
	TFFloat     TFType = "float"
	TFFloat64   TFType = "float64"
	TFList      TFType = "list"
	TFSet       TFType = "set"
	TfMap       TFType = "map"
	TfObject    TFType = "object"
	TfJSON      TFType = "json"
	TfJSONExact TFType = "json-exact"
)

func (t TFType) ResSchemaAttr() rschema.Attribute {
	switch t {
	case TfString:
		return rschema.StringAttribute{}
	case TfInt:
		return rschema.Int32Attribute{}
	case TfInt64:
		return rschema.Int64Attribute{}
	case TfBoolean:
		return rschema.BoolAttribute{}
	case TFFloat:
		return rschema.Float32Attribute{}
	case TFFloat64:
		return rschema.Float64Attribute{}
	case TFList:
		return rschema.ListAttribute{}
	case TfMap:
		return rschema.MapAttribute{}
	case TfObject:
		return rschema.ObjectAttribute{}
	default:
		return nil
	}
}

// BaseSchema is for simple fields common to both resource and data source schema of an entity
type BaseSchema struct {
	Name      string
	AttrType  TFType
	SubType   TFType // underlying type for list, map, etc
	Required  bool   // resource-only field
	Optional  bool   // resource-only field
	Computed  bool   // resource-only field
	Sensitive bool
	Desc      string
	MdDesc    string
	Default   *attrDefault
}

type attrDefault struct {
	StaticBool   bool
	StaticString string
	StaticInt    int
	StaticInt64  int64
	StaticFloat  float32
	StaticDouble float64
}

func attrBoolDefault(ad *attrDefault) defaults.Bool {
	if ad == nil {
		return nil
	}

	return booldefault.StaticBool(ad.StaticBool)
}

func attrStringDefault(ad *attrDefault) defaults.String {
	if ad == nil {
		return nil
	}
	return stringdefault.StaticString(ad.StaticString)
}

func attrIntDefault(ad *attrDefault) defaults.Int32 {
	if ad == nil {
		return nil
	}
	return int32default.StaticInt32(int32(ad.StaticInt))
}

func attrInt64Default(ad *attrDefault) defaults.Int64 {
	if ad == nil {
		return nil
	}
	return int64default.StaticInt64(ad.StaticInt64)
}

func attrFloat32Default(ad *attrDefault) defaults.Float32 {
	if ad == nil {
		return nil
	}
	return float32default.StaticFloat32(ad.StaticFloat)
}

func attrFloat64Default(ad *attrDefault) defaults.Float64 {
	if ad == nil {
		return nil
	}
	return float64default.StaticFloat64(ad.StaticDouble)
}

func (s BaseSchema) ResourceAttr() rschema.Attribute {
	if s.MdDesc == "" {
		s.MdDesc = s.Desc
	}

	if s.AttrType == TfString {
		return rschema.StringAttribute{
			Required:            s.Required,
			Optional:            s.Optional,
			Computed:            s.Computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
			Default:             attrStringDefault(s.Default),
		}
	}
	if s.AttrType == TfInt {
		return rschema.Int32Attribute{
			Required:            s.Required,
			Optional:            s.Optional,
			Computed:            s.Computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
			Default:             attrIntDefault(s.Default),
		}
	}
	if s.AttrType == TfInt64 {
		return rschema.Int64Attribute{
			Required:            s.Required,
			Optional:            s.Optional,
			Computed:            s.Computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
			Default:             attrInt64Default(s.Default),
		}
	}
	if s.AttrType == TfBoolean {
		return rschema.BoolAttribute{
			Required:            s.Required,
			Optional:            s.Optional,
			Computed:            s.Computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
			Default:             attrBoolDefault(s.Default),
		}
	}
	if s.AttrType == TFFloat {
		return rschema.Float32Attribute{
			Required:            s.Required,
			Optional:            s.Optional,
			Computed:            s.Computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
			Default:             attrFloat32Default(s.Default),
		}
	}
	if s.AttrType == TFFloat64 {
		return rschema.Float64Attribute{
			Required:            s.Required,
			Optional:            s.Optional,
			Computed:            s.Computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
			Default:             attrFloat64Default(s.Default),
		}
	}

	// List support (currently only list of strings is needed)
	if s.AttrType == TFList {
		return rschema.ListAttribute{
			ElementType:         types.StringType,
			Required:            s.Required,
			Computed:            s.Computed,
			Optional:            s.Optional,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}

	if s.AttrType == TFSet {
		return rschema.SetAttribute{
			ElementType:         types.StringType,
			Required:            s.Required,
			Computed:            s.Computed,
			Optional:            s.Optional,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}

	// Map support (currently only string element type used in project)
	if s.AttrType == TfMap {
		// default to string map; extend as needed
		return rschema.MapAttribute{
			ElementType:         types.StringType,
			Required:            s.Required,
			Computed:            s.Computed,
			Optional:            s.Optional,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	// JSON support using NormalizedType
	if s.AttrType == TfJSON {
		return rschema.StringAttribute{
			CustomType:          jsontypes.NormalizedType{},
			Required:            s.Required,
			Computed:            s.Computed,
			Optional:            s.Optional,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}

	panic(fmt.Sprintf("unreachable code for attr type = \n%+v", s))
}

func (s BaseSchema) DataSourceAttr() dschema.Attribute {
	// Default behavior for data sources: attributes are typically Computed,
	// except identifiers or when explicitly overridden via DS* flags.
	computed := true
	required := false

	// If any DS* override flags are set, honor them.
	if s.Name == "id" {
		// Special-case ID when not overridden: required input
		computed = false
		required = true
	}

	if s.AttrType == TfString {
		return dschema.StringAttribute{
			Required:            required,
			Computed:            computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	if s.AttrType == TfInt {
		return dschema.Int32Attribute{
			Required:            required,
			Computed:            computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	if s.AttrType == TfInt64 {
		return dschema.Int64Attribute{
			Required:            required,
			Computed:            computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	if s.AttrType == TfBoolean {
		return dschema.BoolAttribute{
			Required:            required,
			Computed:            computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	if s.AttrType == TFFloat {
		return dschema.Float32Attribute{
			Required:            required,
			Computed:            computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	if s.AttrType == TFFloat64 {
		return dschema.Float64Attribute{
			Required:            required,
			Computed:            computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	// List support (currently list of strings is needed)
	if s.AttrType == TFList {
		return dschema.ListAttribute{
			ElementType:         types.StringType,
			Required:            required,
			Computed:            computed,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	// Map support (currently only string element type used in project)
	if s.AttrType == TfMap {
		// default to string map; extend as needed
		return dschema.MapAttribute{
			ElementType:         types.StringType,
			Required:            required,
			Computed:            computed,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}
	// JSON support using NormalizedType
	if s.AttrType == TfJSON {
		return dschema.StringAttribute{
			CustomType:          jsontypes.NormalizedType{},
			Required:            required,
			Computed:            computed,
			Sensitive:           s.Sensitive,
			Description:         s.Desc,
			MarkdownDescription: s.MdDesc,
		}
	}

	return nil
}

func ResAttributes(attrs []BaseSchema) map[string]rschema.Attribute {
	result := make(map[string]rschema.Attribute)
	for _, attr := range attrs {
		result[attr.Name] = attr.ResourceAttr()
	}
	return result
}

// DSAttributes converts a slice of BaseSchema into a map for data source attributes.
func DSAttributes(attrs []BaseSchema) map[string]dschema.Attribute {
	result := make(map[string]dschema.Attribute)
	for _, attr := range attrs {
		result[attr.Name] = attr.DataSourceAttr()
	}
	return result
}
