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

var _ = fmt.Printf

// this event describes result validation from different services and accepts order if we pass validation
type OrderValidationMessage struct {
	Order_id string `json:"order_id"`

	Service_name string `json:"service_name"`

	Created_at int64 `json:"created_at"`

	Is_successful bool `json:"is_successful"`

	Payload PayloadMessageValidation `json:"Payload"`
}

const OrderValidationMessageAvroCRC64Fingerprint = "\xf4\xbc\xd0O\xe0\xaa 2"

func NewOrderValidationMessage() OrderValidationMessage {
	r := OrderValidationMessage{}
	r.Payload = NewPayloadMessageValidation()

	return r
}

func DeserializeOrderValidationMessage(r io.Reader) (OrderValidationMessage, error) {
	t := NewOrderValidationMessage()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeOrderValidationMessageFromSchema(r io.Reader, schema string) (OrderValidationMessage, error) {
	t := NewOrderValidationMessage()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeOrderValidationMessage(r OrderValidationMessage, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Order_id, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Service_name, w)
	if err != nil {
		return err
	}
	err = vm.WriteLong(r.Created_at, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.Is_successful, w)
	if err != nil {
		return err
	}
	err = writePayloadMessageValidation(r.Payload, w)
	if err != nil {
		return err
	}
	return err
}

func (r OrderValidationMessage) Serialize(w io.Writer) error {
	return writeOrderValidationMessage(r, w)
}

func (r OrderValidationMessage) Schema() string {
	return "{\"doc\":\"this event describes result validation from different services and accepts order if we pass validation\",\"fields\":[{\"logicalType\":\"UUID\",\"name\":\"order_id\",\"type\":\"string\"},{\"name\":\"service_name\",\"type\":\"string\"},{\"name\":\"created_at\",\"type\":{\"logicalType\":\"timestamp-millis\",\"type\":\"long\"}},{\"name\":\"is_successful\",\"type\":\"boolean\"},{\"name\":\"Payload\",\"type\":{\"fields\":[{\"name\":\"courier_id\",\"type\":[{\"logicalType\":\"uuid\",\"type\":\"string\"},\"null\"]}],\"name\":\"PayloadMessageValidation\",\"type\":\"record\"}}],\"name\":\"OrderValidationMessage\",\"type\":\"record\"}"
}

func (r OrderValidationMessage) SchemaName() string {
	return "OrderValidationMessage"
}

func (_ OrderValidationMessage) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ OrderValidationMessage) SetInt(v int32)       { panic("Unsupported operation") }
func (_ OrderValidationMessage) SetLong(v int64)      { panic("Unsupported operation") }
func (_ OrderValidationMessage) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ OrderValidationMessage) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ OrderValidationMessage) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ OrderValidationMessage) SetString(v string)   { panic("Unsupported operation") }
func (_ OrderValidationMessage) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *OrderValidationMessage) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Order_id}

		return w

	case 1:
		w := types.String{Target: &r.Service_name}

		return w

	case 2:
		w := types.Long{Target: &r.Created_at}

		return w

	case 3:
		w := types.Boolean{Target: &r.Is_successful}

		return w

	case 4:
		r.Payload = NewPayloadMessageValidation()

		w := types.Record{Target: &r.Payload}

		return w

	}
	panic("Unknown field index")
}

func (r *OrderValidationMessage) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *OrderValidationMessage) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ OrderValidationMessage) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ OrderValidationMessage) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ OrderValidationMessage) HintSize(int)                     { panic("Unsupported operation") }
func (_ OrderValidationMessage) Finalize()                        {}

func (_ OrderValidationMessage) AvroCRC64Fingerprint() []byte {
	return []byte(OrderValidationMessageAvroCRC64Fingerprint)
}

func (r OrderValidationMessage) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["order_id"], err = json.Marshal(r.Order_id)
	if err != nil {
		return nil, err
	}
	output["service_name"], err = json.Marshal(r.Service_name)
	if err != nil {
		return nil, err
	}
	output["created_at"], err = json.Marshal(r.Created_at)
	if err != nil {
		return nil, err
	}
	output["is_successful"], err = json.Marshal(r.Is_successful)
	if err != nil {
		return nil, err
	}
	output["Payload"], err = json.Marshal(r.Payload)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *OrderValidationMessage) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["order_id"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Order_id); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for order_id")
	}
	val = func() json.RawMessage {
		if v, ok := fields["service_name"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Service_name); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for service_name")
	}
	val = func() json.RawMessage {
		if v, ok := fields["created_at"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Created_at); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for created_at")
	}
	val = func() json.RawMessage {
		if v, ok := fields["is_successful"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Is_successful); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for is_successful")
	}
	val = func() json.RawMessage {
		if v, ok := fields["Payload"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Payload); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for Payload")
	}
	return nil
}
