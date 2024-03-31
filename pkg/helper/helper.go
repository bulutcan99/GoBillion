package helper

import (
	"encoding/hex"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	intType interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
			~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}
)

func Contains[T comparable](elems []T, v T, fn func(value T, element T) bool) bool {
	if fn == nil {
		fn = func(value T, element T) bool {
			return value == element
		}
	}
	for _, s := range elems {
		if fn(v, s) {
			return true
		}
	}
	return false
}

func GracefulShutdown(shutdown chan os.Signal) {
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func InDateRange(SelectsDay []int) bool {
	weekday := int(time.Now().Weekday())
	for _, day := range SelectsDay {
		if day == weekday {
			return true
		}
	}
	return false
}

func IsThresholdPassed(triggeredAt time.Time, threshold time.Duration) bool {
	return triggeredAt.Add(threshold).Before(time.Now())
}

func ConvertToFloat64(value any) float64 {
	if value == nil {
		return 0
	}

	typeOf := reflect.TypeOf(value).String()

	switch typeOf {
	case "float64":
		return value.(float64)
	case "float32":
		return float64(value.(float32))
	case "int":
		return float64(value.(int))
	case "int32":
		return float64(value.(int32))
	case "int64":
		return float64(value.(int64))
	case "string":
		f, _ := strconv.ParseFloat(value.(string), 64)
		return f
	case "bool":
		if value.(bool) {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ParseMongoDoubleToFloat64(value any) float64 {
	if value == nil {
		return 0
	}
	typeOfDecimal := reflect.TypeOf(value).String()
	switch typeOfDecimal {
	case "primitive.Decimal128":
		parsedValue, _ := strconv.ParseFloat(value.(primitive.Decimal128).String(), 64)
		return parsedValue
	case "float64":
		return value.(float64)
	case "int64":
		return float64(value.(int64))
	case "int32":
		return float64(value.(int32))
	case "int":
		return float64(value.(int))
	case "string":
		parsedValue, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return 0
		}
		return parsedValue
	default:
		return 0
	}
}

func ConvertAnyToUint16(value any) uint16 {
	if value == nil {
		return 0
	}

	typeOf := reflect.TypeOf(value).String()

	switch typeOf {
	case "int":
		return uint16(value.(int))
	case "int32":
		return value.(uint16)
	case "int64":
		return uint16(value.(int64))
	case "float32":
		return uint16(value.(float32))
	case "float64":
		return uint16(value.(float64))
	case "string":
		i, _ := strconv.ParseInt(value.(string), 10, 32)
		return uint16(i)
	case "bool":
		if value.(bool) {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ConvertAnyToInt8(value interface{}) int8 {
	if value == nil {
		return 0
	}

	typeOf := reflect.TypeOf(value).String()

	switch typeOf {
	case "int":
		return int8(value.(int))
	case "int32":
		return int8(value.(int32))
	case "int64":
		return int8(value.(int64))
	case "float32":
		return int8(value.(float32))
	case "float64":
		return int8(value.(float64))
	default:
		return -1
	}
}

func ConvertToUint16Ptr(value interface{}) *uint16 {
	switch v := value.(type) {
	case int:
		if v >= 0 && v <= math.MaxUint16 {
			result := uint16(v)
			return &result
		}
		zap.S().Error("int value %d is out of range for uint16", v)
		return nil
	case int16:
		if v >= 0 {
			result := uint16(v)
			return &result
		}
		zap.S().Error("int16 value %d is out of range for uint16", v)
		return nil
	case int32:
		if v >= 0 && v <= math.MaxUint16 {
			result := uint16(v)
			return &result
		}
		zap.S().Error("int16 value %d is out of range for uint16", v)
		return nil
	case int64:
		if v >= 0 && v <= math.MaxUint16 {
			result := uint16(v)
			return &result
		}
		zap.S().Error("int16 value %d is out of range for uint16", v)
		return nil
	case float64:
		if v >= 0 && v <= math.MaxUint16 && v == float64(int(v)) {
			result := uint16(int(v))
			return &result
		}
		zap.S().Error("int16 value %d is out of range for uint16", v)
		return nil
	case uint:
		if v <= math.MaxUint16 {
			result := uint16(v)
			return &result
		}
		zap.S().Error("int16 value %d is out of range for uint16", v)
		return nil
	case uint32:
		if v <= math.MaxUint16 {
			result := uint16(v)
			return &result
		}
		zap.S().Error("int16 value %d is out of range for uint16", v)
		return nil
	case uint64:
		if v <= math.MaxUint16 {
			result := uint16(v)
			return &result
		}
		zap.S().Error("int16 value %d is out of range for uint16", v)
		return nil
	default:
		zap.S().Error("value %v is not a valid type for uint16", value)
		return nil
	}
}

func ConvertToInt32(value any) int32 {
	if value == nil {
		return 0
	}

	typeOf := reflect.TypeOf(value).String()

	switch typeOf {
	case "int":
		return int32(value.(int))
	case "int32":
		return value.(int32)
	case "int64":
		return int32(value.(int64))
	case "float32":
		return int32(value.(float32))
	case "float64":
		return int32(value.(float64))
	case "string":
		i, _ := strconv.ParseInt(value.(string), 10, 32)
		return int32(i)
	case "bool":
		if value.(bool) {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ConvertToUInt16[Int intType](value Int) uint16 {
	castedValue := uint16(value)
	if 0 > value {
		zap.S().Error("Overflow occurred")
		return 0
	}

	return castedValue
}

func ConvertToPtrUInt16[Int intType](value Int) *uint16 {
	castedValue := uint16(value)
	if 0 > value {
		zap.S().Error("Overflow occurred")
		return nil
	}

	return &castedValue
}

func CalculateDuration(fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)

	return duration
}

func HexToASCII(hexStr string) string {
	str, err := hex.DecodeString(hexStr)
	if err != nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&str))
}

func ConvertProtoTimeToTime(protoTime *timestamppb.Timestamp) time.Time {
	if protoTime == nil {
		return time.Time{}
	}

	return protoTime.AsTime()
}
