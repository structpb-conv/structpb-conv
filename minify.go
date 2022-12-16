package structpbConv

import "google.golang.org/protobuf/types/known/structpb"

func MinifyStruct(input *structpb.Struct) *structpb.Struct {
	if nil == input || len(input.GetFields()) == 0 {
		return nil
	}

	output := &structpb.Struct{
		Fields: make(map[string]*structpb.Value, len(input.GetFields())),
	}

	for key, value := range input.GetFields() {
		value = MinifyValue(value)
		if nil != value {
			output.Fields[key] = value
		}
	}

	if len(output.Fields) == 0 {
		return nil
	}
	return output
}

func MinifyList(input *structpb.ListValue) *structpb.ListValue {
	if nil == input || len(input.GetValues()) == 0 {
		return nil
	}

	output := &structpb.ListValue{Values: make([]*structpb.Value, 0, len(input.GetValues()))}
	for _, value := range input.GetValues() {
		value = MinifyValue(value)
		if nil != value {
			output.Values = append(output.Values, value)
		}
	}

	if len(output.Values) == 0 {
		return nil
	}
	return output
}

func MinifyValue(input *structpb.Value) *structpb.Value {
	if nil == input || nil == input.Kind {
		return nil
	}
	switch v := input.GetKind().(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_StructValue:
		if minValue := MinifyStruct(v.StructValue); nil != minValue {
			return structpb.NewStructValue(minValue)
		}
		return nil
	case *structpb.Value_ListValue:
		if minValue := MinifyList(v.ListValue); nil != minValue {
			return structpb.NewListValue(minValue)
		}
		return nil
	}
	return input
}
