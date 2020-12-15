package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	TagMappingName   = "mapping"
	TagSeperatorName = "seperator"
)

/*
 *	 ["a","b"] -> "a,b"  : ok , "a,b" -> ["a","b"] : ok
 */
func ToCopyObject(source interface{}, desc interface{}) {
	descType := reflect.TypeOf(desc)
	if descType.Kind() != reflect.Ptr {
		panic("Desc Type Error !! Need Ptr !")
	}

	srcType := indirectType(reflect.TypeOf(source))

	switch srcType.Kind() {
	case reflect.Slice:
		srcValue := indirectValue(reflect.ValueOf(source))
		isPtr := false

		//desc = ptr 이므로 elem 사용.
		descSliceType := descType.Elem()
		//Slice 의 자식 Elem
		descElemType := descSliceType.Elem()
		if descElemType.Kind() == reflect.Ptr {
			isPtr = true
			descElemType = descElemType.Elem()
		}

		resultSlice := reflect.MakeSlice(descSliceType, 0, 0)
		for i := 0; i < srcValue.Len(); i++ {
			srcIndexValue := indirectValue(srcValue.Index(i))
			descIndexValue := reflect.New(descElemType) //new 를 하게 되면 pointer

			toCopyOnOtherType(srcIndexValue.Interface(), descIndexValue.Interface())
			if isPtr {
				resultSlice = reflect.Append(resultSlice, descIndexValue)
			} else {
				resultSlice = reflect.Append(resultSlice, descIndexValue.Elem())
			}
		}
		reflect.ValueOf(desc).Elem().Set(resultSlice)
	case reflect.Struct:
		toCopyOnOtherType(source, desc)
	case reflect.Map:
		srcValue := indirectValue(reflect.ValueOf(source))
		isPtr := false

		//desc = ptr 이므로 elem 사용.
		descSliceType := descType.Elem()
		//Slice 의 자식 Elem
		descElemType := descSliceType.Elem()
		if descElemType.Kind() == reflect.Ptr {
			isPtr = true
			descElemType = descElemType.Elem()
		}

		resultMap := reflect.MakeMap(descSliceType)
		for _, v := range srcValue.MapKeys() {
			srcIndexValue := indirectValue(srcValue.MapIndex(v))
			descIndexValue := reflect.New(descElemType) //new 를 하게 되면 pointer
			toCopyOnOtherType(srcIndexValue.Interface(), descIndexValue.Interface())
			if isPtr {
				resultMap.SetMapIndex(v, descIndexValue)
			} else {
				resultMap.SetMapIndex(v, descIndexValue.Elem())
			}
		}
		reflect.ValueOf(desc).Elem().Set(resultMap)
	}
}

func toCopyOnOtherType(src interface{}, desc interface{}) {
	descType := reflect.TypeOf(desc)
	if descType.Kind() != reflect.Ptr {
		panic("Desc Type Error !! Need Ptr !")
	}

	descVal := reflect.ValueOf(desc).Elem()
	srcVal := indirectValue(reflect.ValueOf(src))

	if !srcVal.IsValid() || !descVal.IsValid() {
		return
	}

	srcValType := srcVal.Type()
	if srcValType.Kind() != reflect.Struct {
		setValue(srcVal.Interface(), reflect.StructField{}, srcVal, descVal)
		return
	}

	length := srcVal.NumField()

	descMappingMap, isSrcPb := checkSrcPbAndGetMappingDesc(srcValType, descVal)

	for i := 0; i < length; i++ {
		srcTypeField := srcValType.Field(i)
		srcFieldName := srcTypeField.Tag.Get(TagMappingName)
		if srcFieldName == "" {
			if isSrcPb {
				if dfn, ok := descMappingMap[srcTypeField.Name]; ok && dfn != "" {
					srcFieldName = dfn
				}
			}

			if srcFieldName == "" {
				srcFieldName = srcTypeField.Name
			}
		}

		//if !descVal.IsValid() {
		//	continue
		//}

		var descField = descVal.FieldByName(srcFieldName)
		if !descField.IsValid() {
			continue
		}

		f := srcVal.Field(i)
		iv := indirectValue(f)
		if !iv.IsValid() {
			continue
		}

		setValue(iv.Interface(), srcTypeField, iv, descField)
	}
}

// value = indirect interface value
// f = indirect reflect.value
// descField = reflect.value
func setValue(value interface{}, srcTypeField reflect.StructField, f, descField reflect.Value) {
	//todo: 나중에 desc type 기준 으로 수정 필요
	//switch descField.Interface().(type) {
	//case *types.Struct:
	//	b, _ := json.Marshal(value)
	//	s := &types.Struct{}
	//	if err := jsonpb.Unmarshal(bytes.NewReader(b), s); err != nil {
	//		fmt.Println(err)
	//	} else {
	//		descField.Set(reflect.ValueOf(s))
	//	}
	//	return
	//}

	switch v := value.(type) {
	case string:
		switch descField.Interface().(type) {
		case time.Time:
			if t, ok := TryConvertTime(v); ok {
				descField.Set(reflect.ValueOf(t))
			} else {
				return
			}
		case *time.Time:
			if t, ok := TryConvertTime(v); ok {
				descField.Set(reflect.ValueOf(&t))
			} else {
				return
			}
		case []string:
			sep := srcTypeField.Tag.Get(TagSeperatorName)
			if sep == "" {
				sep = ","
			}

			if v == "" {
				descField.Set(reflect.ValueOf([]string{}))
			} else {
				descField.Set(reflect.ValueOf(strings.Split(v, sep)))
			}
		case *string:
			descField.Set(reflect.ValueOf(&v))
		case map[string]interface{}:
			var m map[string]interface{}
			if v == "" {
				descField.Set(reflect.ValueOf(map[string]interface{}{}))
			} else {
				if err := json.Unmarshal([]byte(v), &m); err == nil {
					descField.Set(reflect.ValueOf(m))
				} else {
					descField.Set(reflect.ValueOf(map[string]interface{}{}))
				}
			}
		case bool:
			descField.SetBool(convertStringToBoolean(v))
		case *bool:
			bo := convertStringToBoolean(v)
			descField.Set(reflect.ValueOf(&bo))
		default:
			descField.SetString(v)
		}
		break
	case []string:
		switch descField.Interface().(type) {
		case *string:
			sep := srcTypeField.Tag.Get(TagSeperatorName)
			if sep == "" {
				sep = ","
			}
			descField.Set(reflect.ValueOf(strings.Join(v, sep)).Addr())
		case string:
			sep := srcTypeField.Tag.Get(TagSeperatorName)
			if sep == "" {
				sep = ","
			}
			descField.Set(reflect.ValueOf(strings.Join(v, sep)))
		case []string:
			if f.Len() > 0 {
				descField.Set(f)
			} else {
				descField.Set(reflect.ValueOf([]string{}))
			}
		}
		break
	case time.Time:
		switch descField.Interface().(type) {
		case string:
			descField.SetString(v.Format(time.RFC3339))
			break
		case *time.Time:
			descField.Set(reflect.ValueOf(&v))
		default:
			descField.Set(reflect.ValueOf(v))
		}
		break
	case []byte:
		switch descField.Interface().(type) {
		case []byte:
			descField.Set(reflect.ValueOf(v))
		default:
			// struct
			targetIsPtr := false
			targetType := descField.Type()

			if descField.Kind() == reflect.Ptr {
				targetIsPtr = true
				targetType = targetType.Elem()
			}

			if targetType.Kind() == reflect.Struct {
				descNewValue := reflect.New(targetType) //new 를 하게 되면 pointer
				if err := json.Unmarshal(v, descNewValue.Interface()); err == nil {
					if targetIsPtr {
						descField.Set(descNewValue)
					} else {
						descField.Set(descNewValue.Elem())
					}
				}
			}
		}
	case int64:
		switch descField.Interface().(type) {
		case time.Duration:
			x := time.Duration(v)
			descField.Set(reflect.ValueOf(x))
		case *time.Duration:
			x := time.Duration(v)
			descField.Set(reflect.ValueOf(&x))
		default:
			convertAndSetNumValue(f, descField)
		}
	default:
		fKind := f.Kind()
		//todo: need refactoring
		if fKind == reflect.Slice || fKind == reflect.Struct || fKind == reflect.Map {
			targetType := descField.Type()
			targetIsPtr := false
			if descField.Kind() == reflect.Ptr {
				targetIsPtr = true
				targetType = targetType.Elem()
			}

			if fKind == reflect.Slice {
				sliceType := reflect.SliceOf(targetType.Elem())
				resultSlice := reflect.New(sliceType)

				ToCopyObject(value, resultSlice.Interface())

				if targetIsPtr {
					descField.Set(resultSlice)
				} else {
					descField.Set(resultSlice.Elem())
				}
			} else {
				descNewValue := reflect.New(targetType) //new 를 하게 되면 pointer
				ToCopyObject(value, descNewValue.Interface())

				if targetIsPtr {
					descField.Set(descNewValue)
				} else {
					descField.Set(descNewValue.Elem())
				}
			}
		} else {
			convertAndSetNumValue(f, descField)
		}
	}
}

func convertAndSetNumValue(src, dst reflect.Value) {
	dft := dst.Type()
	dfkIsPtr := dst.Kind() == reflect.Ptr
	if dfkIsPtr {
		dft = dft.Elem()
	}
	if src.Type() == dft { //* 땐 타입 비교
		if dfkIsPtr {
			dst.Set(src.Addr())
		} else {
			dst.Set(src)
		}
	} else {
		c := src.Convert(dft)

		var rv reflect.Value
		switch c.Kind() {
		case reflect.Int:
			v := c.Interface().(int)
			if dfkIsPtr {
				rv = reflect.ValueOf(&v)
			} else {
				rv = reflect.ValueOf(v)
			}
		case reflect.Int8:
			v := c.Interface().(int8)
			if dfkIsPtr {
				rv = reflect.ValueOf(&v)
			} else {
				rv = reflect.ValueOf(v)
			}
		case reflect.Int16:
			v := c.Interface().(int16)
			if dfkIsPtr {
				rv = reflect.ValueOf(&v)
			} else {
				rv = reflect.ValueOf(v)
			}
		case reflect.Int32:
			v := c.Interface().(int32)
			if dfkIsPtr {
				rv = reflect.ValueOf(&v)
			} else {
				rv = reflect.ValueOf(v)
			}
		case reflect.Int64:
			v := c.Interface().(int64)
			if dfkIsPtr {
				rv = reflect.ValueOf(&v)
			} else {
				rv = reflect.ValueOf(v)
			}
		case reflect.Float32:
			v := c.Interface().(float32)
			if dfkIsPtr {
				rv = reflect.ValueOf(&v)
			} else {
				rv = reflect.ValueOf(v)
			}
		case reflect.Float64:
			v := c.Interface().(float64)
			if dfkIsPtr {
				rv = reflect.ValueOf(&v)
			} else {
				rv = reflect.ValueOf(v)
			}
		//uint 는 쓴다면 또 추가.
		default:
			fmt.Println("HELLO ", c, src, dst, c.Kind())
			return
		}
		dst.Set(rv)
	}
}

func TryConvertTime(valueString string) (time.Time, bool) {
	resultTime, er := time.Parse(time.RFC3339Nano, valueString)
	if er != nil {
		resultTime, er = time.Parse(time.RFC3339, valueString)
		if er != nil {
			return resultTime, false
		}
	}

	return resultTime, true
}

func convertStringToBoolean(v string) bool {
	ret, err := strconv.ParseBool(v)
	if err == nil {
		return ret
	} else if v == "Y" || v == "y" {
		return true
	} else {
		return false
	}
}

func checkSrcPbAndGetMappingDesc(srcType reflect.Type, descVal reflect.Value) (map[string]string, bool) {
	//if _, isSrcPb := srcType.FieldByName(DefaultPbField); isSrcPb {
	descMappingMap := map[string]string{}
	descValType := descVal.Type()
	descLength := descVal.NumField()
	for i := 0; i < descLength; i++ {
		descTypeField := descValType.Field(i)
		descFieldName := descTypeField.Tag.Get(TagMappingName)
		if descFieldName != "" {
			descMappingMap[descFieldName] = descTypeField.Name
		}
	}

	return descMappingMap, true
	//}
	//
	//return nil, false
}

//elem = Pointer 를 값으로 변경
func indirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectValue reflect.Type) reflect.Type {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}
