package helpers

import (
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ProtoMessageToMap converts a protobuf message to map[string]string
func ProtoMessageToMap(msg proto.Message) map[string]string {
	result := make(map[string]string)

	m := msg.ProtoReflect()
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		// Skip empty values
		if reflect.ValueOf(v.Interface()).IsZero() {
			return true
		}

		// Get the field name (in snake_case)
		fieldName := string(fd.Name())

		// Convert to camelCase for Google API params
		fieldName = toCamelCase(fieldName)

		// Convert the value to string based on its type
		var strValue string
		switch fd.Kind() {
		case protoreflect.StringKind:
			strValue = v.String()
		case protoreflect.Int32Kind, protoreflect.Int64Kind,
			protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			strValue = fmt.Sprintf("%d", v.Interface())
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			strValue = fmt.Sprintf("%f", v.Interface())
		case protoreflect.BoolKind:
			strValue = fmt.Sprintf("%t", v.Interface())
		default:
			// Skip complex types or handle them specifically if needed
			return true
		}

		// Special handling for certain fields
		if fieldName == "lr" && strValue != "" && !strings.HasPrefix(strValue, "lang_") {
			strValue = fmt.Sprintf("lang_%s", strValue)
		}

		result[fieldName] = strValue
		return true
	})

	return result
}

// toCamelCase converts snake_case to camelCase
func toCamelCase(s string) string {
	// Handle special cases for Google API params
	switch s {
	case "c2coff":
		return "c2coff"
	case "cr":
		return "cr"
	case "cx":
		return "cx"
	case "gl":
		return "gl"
	case "hl":
		return "hl"
	case "hq":
		return "hq"
	case "lr":
		return "lr"
	case "num":
		return "num"
	case "q":
		return "q"
	}

	// General snake_case to camelCase conversion
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
