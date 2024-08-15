package pb

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/VictoriaMetrics/easyproto"
)

// Uint64 represents uint64 type, which is encoded as a string in JSON and as uint64 in Prorobuf
// See https://protobuf.dev/programming-guides/proto3/#json
type Uint64 uint64

// UnmarshalJSON decodes t from data
func (u *Uint64) UnmarshalJSON(src []byte) error {
	vStr, err := strconv.Unquote(string(src))
	if err != nil {
		return fmt.Errorf("failed to unquote uint64: %w", err)
	}
	v, err := strconv.ParseUint(vStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse uint64: %w", err)
	}
	*u = Uint64(v)
	return nil
}

// Int64 represents int64 type, which is encoded as a string in JSON and as int64 in Prorobuf
// See https://protobuf.dev/programming-guides/proto3/#json
type Int64 int64

// UnmarshalJSON decodes t from data
func (i *Int64) UnmarshalJSON(src []byte) error {
	vStr, err := strconv.Unquote(string(src))
	if err != nil {
		return fmt.Errorf("failed to unquote int64: %w", err)
	}
	v, err := strconv.ParseInt(vStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse int64: %w", err)
	}
	*i = Int64(v)
	return nil
}

// Resource represents the corresponding OTEL protobuf message
type Resource struct {
	Attributes []*KeyValue
}

func (r *Resource) marshalProtobuf(mm *easyproto.MessageMarshaler) {
	for _, a := range r.Attributes {
		a.marshalProtobuf(mm.AppendMessage(1))
	}
}

// UnmarshalProtobuf unmarshals r from protobuf message at src.
func (r *Resource) UnmarshalProtobuf(src []byte) (err error) {
	// message Resource {
	//   repeated KeyValue attributes = 1;
	// }
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, err = fc.NextField(src)
		if err != nil {
			return fmt.Errorf("cannot read next field in Resource: %w", err)
		}
		switch fc.FieldNum {
		case 1:
			data, ok := fc.MessageData()
			if !ok {
				return fmt.Errorf("cannot read Attribute data")
			}
			r.Attributes = append(r.Attributes, &KeyValue{})
			a := r.Attributes[len(r.Attributes)-1]
			if err := a.UnmarshalProtobuf(data); err != nil {
				return fmt.Errorf("cannot unmarshal Attribute: %w", err)
			}
		}
	}
	return nil
}

// KeyValue represents the corresponding OTEL protobuf message
type KeyValue struct {
	Key   string
	Value *AnyValue
}

func (kv *KeyValue) marshalProtobuf(mm *easyproto.MessageMarshaler) {
	mm.AppendString(1, kv.Key)
	if kv.Value != nil {
		kv.Value.marshalProtobuf(mm.AppendMessage(2))
	}
}

// UnmarshalProtobuf unmarshals kv from protobuf message at src.
func (kv *KeyValue) UnmarshalProtobuf(src []byte) (err error) {
	// message KeyValue {
	//   string key = 1;
	//   AnyValue value = 2;
	// }
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, err = fc.NextField(src)
		if err != nil {
			return fmt.Errorf("cannot read next field in KeyValue: %w", err)
		}
		switch fc.FieldNum {
		case 1:
			key, ok := fc.String()
			if !ok {
				return fmt.Errorf("cannot read Key")
			}
			kv.Key = strings.Clone(key)
		case 2:
			data, ok := fc.MessageData()
			if !ok {
				return fmt.Errorf("cannot read Value")
			}
			kv.Value = &AnyValue{}
			if err := kv.Value.UnmarshalProtobuf(data); err != nil {
				return fmt.Errorf("cannot unmarshal Value: %w", err)
			}
		}
	}
	return nil
}

// AnyValue represents the corresponding OTEL protobuf message
type AnyValue struct {
	StringValue  *string
	BoolValue    *bool
	IntValue     *int64
	DoubleValue  *float64
	ArrayValue   *ArrayValue
	KeyValueList *KeyValueList
	BytesValue   *[]byte
}

func (av *AnyValue) marshalProtobuf(mm *easyproto.MessageMarshaler) {
	switch {
	case av.StringValue != nil:
		mm.AppendString(1, *av.StringValue)
	case av.BoolValue != nil:
		mm.AppendBool(2, *av.BoolValue)
	case av.IntValue != nil:
		mm.AppendInt64(3, *av.IntValue)
	case av.DoubleValue != nil:
		mm.AppendDouble(4, *av.DoubleValue)
	case av.ArrayValue != nil:
		av.ArrayValue.marshalProtobuf(mm.AppendMessage(5))
	case av.KeyValueList != nil:
		av.KeyValueList.marshalProtobuf(mm.AppendMessage(6))
	case av.BytesValue != nil:
		mm.AppendBytes(7, *av.BytesValue)
	}
}

// UnmarshalProtobuf unmarshals av from protobuf message at src.
func (av *AnyValue) UnmarshalProtobuf(src []byte) (err error) {
	// message AnyValue {
	//   oneof value {
	//     string string_value = 1;
	//     bool bool_value = 2;
	//     int64 int_value = 3;
	//     double double_value = 4;
	//     ArrayValue array_value = 5;
	//     KeyValueList kvlist_value = 6;
	//     bytes bytes_value = 7;
	//   }
	// }
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, err = fc.NextField(src)
		if err != nil {
			return fmt.Errorf("cannot read next field in AnyValue")
		}
		switch fc.FieldNum {
		case 1:
			stringValue, ok := fc.String()
			if !ok {
				return fmt.Errorf("cannot read StringValue")
			}
			stringValue = strings.Clone(stringValue)
			av.StringValue = &stringValue
		case 2:
			boolValue, ok := fc.Bool()
			if !ok {
				return fmt.Errorf("cannot read BoolValue")
			}
			av.BoolValue = &boolValue
		case 3:
			intValue, ok := fc.Int64()
			if !ok {
				return fmt.Errorf("cannot read IntValue")
			}
			av.IntValue = &intValue
		case 4:
			doubleValue, ok := fc.Double()
			if !ok {
				return fmt.Errorf("cannot read DoubleValue")
			}
			av.DoubleValue = &doubleValue
		case 5:
			data, ok := fc.MessageData()
			if !ok {
				return fmt.Errorf("cannot read ArrayValue")
			}
			av.ArrayValue = &ArrayValue{}
			if err := av.ArrayValue.UnmarshalProtobuf(data); err != nil {
				return fmt.Errorf("cannot unmarshal ArrayValue: %w", err)
			}
		case 6:
			data, ok := fc.MessageData()
			if !ok {
				return fmt.Errorf("cannot read KeyValueList")
			}
			av.KeyValueList = &KeyValueList{}
			if err := av.KeyValueList.UnmarshalProtobuf(data); err != nil {
				return fmt.Errorf("cannot unmarshal KeyValueList: %w", err)
			}
		case 7:
			bytesValue, ok := fc.Bytes()
			if !ok {
				return fmt.Errorf("cannot read BytesValue")
			}
			bytesValue = bytes.Clone(bytesValue)
			av.BytesValue = &bytesValue
		}
	}
	return nil
}

// ArrayValue represents the corresponding OTEL protobuf message
type ArrayValue struct {
	Values []*AnyValue
}

func (av *ArrayValue) marshalProtobuf(mm *easyproto.MessageMarshaler) {
	for _, v := range av.Values {
		v.marshalProtobuf(mm.AppendMessage(1))
	}
}

// UnmarshalProtobuf unmarshals av from protobuf message at src.
func (av *ArrayValue) UnmarshalProtobuf(src []byte) (err error) {
	// message ArrayValue {
	//   repeated AnyValue values = 1;
	// }
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, err = fc.NextField(src)
		if err != nil {
			return fmt.Errorf("cannot read next field in ArrayValue")
		}
		switch fc.FieldNum {
		case 1:
			data, ok := fc.MessageData()
			if !ok {
				return fmt.Errorf("cannot read Value data")
			}
			av.Values = append(av.Values, &AnyValue{})
			v := av.Values[len(av.Values)-1]
			if err := v.UnmarshalProtobuf(data); err != nil {
				return fmt.Errorf("cannot unmarshal Value: %w", err)
			}
		}
	}
	return nil
}

// KeyValueList represents the corresponding OTEL protobuf message
type KeyValueList struct {
	Values []*KeyValue
}

func (kvl *KeyValueList) marshalProtobuf(mm *easyproto.MessageMarshaler) {
	for _, v := range kvl.Values {
		v.marshalProtobuf(mm.AppendMessage(1))
	}
}

// UnmarshalProtobuf unmarshals kvl from protobuf message at src.
func (kvl *KeyValueList) UnmarshalProtobuf(src []byte) (err error) {
	// message KeyValueList {
	//   repeated KeyValue values = 1;
	// }
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, err = fc.NextField(src)
		if err != nil {
			return fmt.Errorf("cannot read next field in KeyValueList")
		}
		switch fc.FieldNum {
		case 1:
			data, ok := fc.MessageData()
			if !ok {
				return fmt.Errorf("cannot read Value data")
			}
			kvl.Values = append(kvl.Values, &KeyValue{})
			v := kvl.Values[len(kvl.Values)-1]
			if err := v.UnmarshalProtobuf(data); err != nil {
				return fmt.Errorf("cannot unmarshal Value: %w", err)
			}
		}
	}
	return nil
}
