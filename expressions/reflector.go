package expressions

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetAttributeValueMemberType(val reflect.Value) types.AttributeValue {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", val.Interface())}
	case reflect.Bool:
		return &types.AttributeValueMemberBOOL{Value: val.Bool()}
	//case reflect.Map:

	case reflect.Slice, reflect.Array:
		return &types.AttributeValueMemberSS{Value: val.Interface().([]string)}
	default:
		return &types.AttributeValueMemberS{Value: val.String()}
	}
}
