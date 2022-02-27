package getopt

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (opts *GetOpt) Marshal(target interface{}, argv []string, posix bool) ([]string, error) {
	targetValue := reflect.ValueOf(target).Elem()
	if targetValue.Kind() == reflect.Ptr {
		targetValue = reflect.Indirect(targetValue)
	}
	if targetValue.Kind() != reflect.Struct {
		opts.done = true
		return nil, errors.New("struct pointer expected, " + targetValue.Kind().String() + " receved")
	}
	if !targetValue.CanAddr() {
		opts.done = true
		return nil, errors.New("struct pointer expected, " + targetValue.Kind().String() + " receved")
	}
	for i, I := 0, targetValue.NumField(); i < I; i++ {
		fieldType := targetValue.Type().Field(i)
		if found, ok := fieldType.Tag.Lookup("flag"); ok {
			if !fieldType.IsExported() {
				return nil, errors.New("can't use flags for unexported fieldType " + fieldType.Name)
			}
			fieldValue := targetValue.Field(i)
			synonyms := strings.Split(found, ",")
			help := fieldType.Tag.Get("help")
			flags, longopts := opts.separateFlagsFromLognopts(synonyms)
			if len(flags) == 0 && len(longopts) == 0 {
				continue
			}
			var err error
			var callback func(string) error
			var trigger func() error
			switch value := fieldValue.Interface().(type) {
			case func() error:
				trigger = value
			case func(str string) error:
				callback = value
			case string:
				callback = func(strval string) error {
					fieldValue.Set(reflect.ValueOf(strval))
					return nil
				}
			case []string:
				callback = func(strval string) error {
					value = append(value, strval)
					fieldValue.Set(reflect.ValueOf(value))
					return nil
				}
			case map[string]string:
				callback = func(strval string) error {
					key, val := getKeyValue(strval)
					if value == nil {
						value = make(map[string]string)
					}
					value[key] = val
					fieldValue.Set(reflect.ValueOf(value))
					return nil
				}
			case uint64:
				callback = func(strval string) error {
					value, err := strconv.ParseUint(strval, 0, 64)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case []uint64:
				callback = func(strval string) error {
					val, err := strconv.ParseUint(strval, 0, 64)
					if err == nil {
						value = append(value, val)
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]uint64:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := strconv.ParseUint(sval, 0, 64); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]uint64)
						}
						value[key] = val
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			case uint:
				callback = func(strval string) error {
					value, err := strconv.ParseUint(strval, 0, 32)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(uint(value)))
					}
					return err
				}
			case []uint:
				callback = func(strval string) error {
					val, err := strconv.ParseUint(strval, 0, 32)
					if err == nil {
						value = append(value, uint(val))
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]uint:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := strconv.ParseUint(sval, 0, 32); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]uint)
						}
						value[key] = uint(val)
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			case int64:
				callback = func(strval string) error {
					value, err := strconv.ParseInt(strval, 0, 64)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case []int64:
				callback = func(strval string) error {
					val, err := strconv.ParseInt(strval, 0, 64)
					if err == nil {
						value = append(value, val)
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]int64:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := strconv.ParseInt(sval, 0, 64); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]int64)
						}
						value[key] = val
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			case int:
				callback = func(strval string) error {
					value, err := strconv.ParseInt(strval, 0, 32)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(int(value)))
					}
					return err
				}
			case []int:
				callback = func(strval string) error {
					val, err := strconv.ParseInt(strval, 0, 32)
					if err == nil {
						value = append(value, int(val))
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]int:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := strconv.ParseInt(sval, 0, 32); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]int)
						}
						value[key] = int(val)
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			case bool:
				err = opts.FlagFuncV(flags, longopts, func() error {
					fieldValue.Set(reflect.ValueOf(true))
					return nil
				}, help)
			case float64:
				callback = func(strval string) error {
					value, err := strconv.ParseFloat(strval, 64)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case []float64:
				callback = func(strval string) error {
					val, err := strconv.ParseFloat(strval, 64)
					if err == nil {
						value = append(value, val)
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]float64:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := strconv.ParseFloat(sval, 64); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]float64)
						}
						value[key] = val
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			case float32:
				callback = func(strval string) error {
					value, err := strconv.ParseFloat(strval, 32)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(float32(value)))
					}
					return err
				}
			case []float32:
				callback = func(strval string) error {
					val, err := strconv.ParseFloat(strval, 32)
					if err == nil {
						value = append(value, float32(val))
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]float32:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := strconv.ParseFloat(sval, 32); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]float32)
						}
						value[key] = float32(val)
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			case time.Time:
				callback = func(strval string) error {
					value, err := time.Parse(time.RFC3339, strval)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case []time.Time:
				callback = func(strval string) error {
					val, err := time.Parse(time.RFC3339, strval)
					if err == nil {
						value = append(value, val)
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]time.Time:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := time.Parse(time.RFC3339, sval); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]time.Time)
						}
						value[key] = val
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			case time.Duration:
				callback = func(strval string) error {
					value, err := time.ParseDuration(strval)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case []time.Duration:
				callback = func(strval string) error {
					val, err := time.ParseDuration(strval)
					if err == nil {
						value = append(value, val)
						fieldValue.Set(reflect.ValueOf(value))
					}
					return err
				}
			case map[string]time.Duration:
				callback = func(strval string) error {
					key, sval := getKeyValue(strval)
					if val, err := time.ParseDuration(sval); err != nil {
						return err
					} else {
						if value == nil {
							value = make(map[string]time.Duration)
						}
						value[key] = val
						fieldValue.Set(reflect.ValueOf(value))
						return nil
					}
				}
			default:
				return nil, errors.New("unsupported type " + fieldType.Type.Kind().String() + " for " + fieldType.Name)
			}
			if callback != nil {
				err = opts.ArgFuncV(flags, longopts, callback, help)
				if err == nil {
					var val *string
					if found, ok := fieldType.Tag.Lookup("default"); ok {
						val = &found
					}
					if found, ok := fieldType.Tag.Lookup("env"); ok {
						if found, ok := os.LookupEnv(found); ok {
							val = &found
						}
					}
					if val != nil {
						err = callback(*val)
					}
				}
			} else if trigger != nil {
				err = opts.FlagFuncV(flags, longopts, trigger, help)
				if err == nil {
					var val bool
					if found, ok := fieldType.Tag.Lookup("default"); ok {
						val, err = strconv.ParseBool(found)
					}
					if err == nil {
						if found, ok := fieldType.Tag.Lookup("env"); ok {
							if found, ok := os.LookupEnv(found); ok {
								val, err = strconv.ParseBool(found)
							}
						}
					}
					if err == nil && val {
						trigger()
					}
				}
			}
			if err != nil {
				opts.done = true
				return nil, err
			}
		}
	}
	return opts.Parse(argv, posix)
}

func getKeyValue(arg string) (key, value string) {
	vec := strings.SplitN(arg, ":", 2)
	if len(vec) > 1 {
		return vec[0], vec[1]
	} else {
		return key, ""
	}
}
