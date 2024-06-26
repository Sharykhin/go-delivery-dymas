// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     order_validation_message.avsc
 */
package v1

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

type UnionStringNullTypeEnum int

const (
	UnionStringNullTypeEnumString UnionStringNullTypeEnum = 0
)

type UnionStringNull struct {
	String    string
	Null      *types.NullVal
	UnionType UnionStringNullTypeEnum
}

func writeUnionStringNull(r *UnionStringNull, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(1, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionStringNullTypeEnumString:
		return vm.WriteString(r.String, w)
	}
	return fmt.Errorf("invalid value for *UnionStringNull")
}

func NewUnionStringNull() *UnionStringNull {
	return &UnionStringNull{}
}

func (r *UnionStringNull) Serialize(w io.Writer) error {
	return writeUnionStringNull(r, w)
}

func DeserializeUnionStringNull(r io.Reader) (*UnionStringNull, error) {
	t := NewUnionStringNull()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func DeserializeUnionStringNullFromSchema(r io.Reader, schema string) (*UnionStringNull, error) {
	t := NewUnionStringNull()
	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func (r *UnionStringNull) Schema() string {
	return "[{\"logicalType\":\"uuid\",\"type\":\"string\"},\"null\"]"
}

func (_ *UnionStringNull) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionStringNull) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionStringNull) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionStringNull) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionStringNull) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionStringNull) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionStringNull) SetLong(v int64) {

	r.UnionType = (UnionStringNullTypeEnum)(v)
}

func (r *UnionStringNull) Get(i int) types.Field {

	switch i {
	case 0:
		return &types.String{Target: (&r.String)}
	case 1:
		return r.Null
	}
	panic("Unknown field index")
}
func (_ *UnionStringNull) NullField(i int)                  { panic("Unsupported operation") }
func (_ *UnionStringNull) HintSize(i int)                   { panic("Unsupported operation") }
func (_ *UnionStringNull) SetDefault(i int)                 { panic("Unsupported operation") }
func (_ *UnionStringNull) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *UnionStringNull) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *UnionStringNull) Finalize()                        {}

func (r *UnionStringNull) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionStringNullTypeEnumString:
		return json.Marshal(map[string]interface{}{"string": r.String})
	}
	return nil, fmt.Errorf("invalid value for *UnionStringNull")
}

func (r *UnionStringNull) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["string"]; ok {
		r.UnionType = 0
		return json.Unmarshal([]byte(value), &r.String)
	}
	return fmt.Errorf("invalid value for *UnionStringNull")
}
