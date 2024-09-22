package params

import (
	"fmt"
	"strconv"
)

type Params map[string]any

func FromAny(a any) Params {
	m, ok := a.(map[string]any)
	if !ok {
		return Params{}
	}
	return m
}

func Merge(pms ...Params) Params {
	result := Params{}
	for _, pm := range pms {
		for k, v := range pm {
			result[k] = v
		}
	}
	return result
}

func (p Params) Set(key string, value any) {
	p[key] = value
}

func (p Params) Map(key ...string) Params {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case map[string]any:
			return Params(v)
		case map[string]string:
			params := Params{}
			for k, v := range v {
				params[k] = v
			}
			return params
		case map[any]any:
			params := Params{}
			for k, v := range v {
				params[fmt.Sprintf("%v", k)] = v
			}
			return params
		default:
			return Params{}
		}
	}
	return Params{}
}

func (p Params) Slice(key ...string) []Params {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case []any:
			var result []Params
			for _, item := range v {
				result = append(result, FromAny(item))
			}
			return result
		case []map[string]any:
			var result []Params
			for _, item := range v {
				result = append(result, Params(item))
			}
			return result
		default:
			return []Params{}
		}
	}
	return []Params{}
}

func (p Params) Int(key ...string) int {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case string:
			i, err := strconv.Atoi(v)
			if err != nil {
				return 0
			}
			return i
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		default:
			return 0
		}
	}

	return 0
}

func (p Params) Float32(key ...string) float32 {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case string:
			i, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return 0
			}
			return float32(i)
		case int:
			return float32(v)
		case int64:
			return float32(v)
		case float32:
			return float32(v)
		case float64:
			return float32(v)
		default:
			return 0
		}
	}

	return 0
}

func (p Params) Float64(key ...string) float64 {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case string:
			i, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0
			}
			return float64(i)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return float64(v)
		default:
			return 0
		}
	}

	return 0
}

func (p Params) String(key ...string) string {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}
		switch v := n.(type) {
		case string:
			return v
		case int:
			return strconv.Itoa(v)
		case int64:
			return strconv.FormatInt(v, 10)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(v)
		default:
			return ""
		}
	}

	return ""
}

func (p Params) Bool(key ...string) bool {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case bool:
			return v
		case string:
			return v != "false"
		case int:
			return v != 0
		case int64:
			return v != 0
		case float64:
			return v != 0
		default:
			return false
		}
	}

	return false
}

func (p Params) IntSlice(key ...string) []int {
	return slice[int](p, key...)
}

func (p Params) FloatSlice(key ...string) []float64 {
	return slice[float64](p, key...)
}

func (p Params) StringSlice(key ...string) []string {
	return slice[string](p, key...)
}

func (p Params) BoolSlice(key ...string) []bool {
	return slice[bool](p, key...)
}

func (p Params) ByteSlice(key ...string) []byte {
	for _, k := range key {
		n, ok := p[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case string:
			return []byte(v)
		case []byte:
			return v
		default:
			return nil
		}
	}

	return nil
}

func slice[T any](m Params, key ...string) []T {
	for _, k := range key {
		n, ok := m[k]
		if !ok {
			continue
		}

		switch v := n.(type) {
		case []any:
			var a []T
			for _, n := range v {
				switch n := n.(type) {
				case T:
					a = append(a, n)
				}
			}
			return a
		default:
			return nil
		}
	}

	return nil
}
