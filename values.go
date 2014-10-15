package valval

import "reflect"

type valType int

const (
	t_number valType = iota
	t_string
	t_bool
)

var vtStrings = map[valType]string{
	t_number: "number",
	t_string: "string",
	t_bool:   "bool",
}

func (vt *valType) isKindOfType(val interface{}) bool {
	uvt := unwrapPtr(val)
	if uvt == nil {
		return false
	}
	rv := reflect.ValueOf(val)
	kind := rv.Kind()
	switch *vt {
	case t_number:
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr, reflect.Float32, reflect.Float64:
			return true
		}
	case t_string:
		return kind == reflect.String
	case t_bool:
		return kind == reflect.Bool
	}
	return false
}

type valueValidator struct {
	vt     valType
	vfuncs []ValidatorFunc
}

func (v *valueValidator) Validate(val interface{}) error {
	if val == nil {
		return nil
	}
	if v.vt.isKindOfType(val) {
		return validateByFuncs(v.vfuncs, val)
	}
	return typeMissmatchError(vtStrings[v.vt])
}

func Number(vfuncs ...ValidatorFunc) Validator {
	return &valueValidator{
		vt:     t_number,
		vfuncs: vfuncs,
	}
}

func String(vfuncs ...ValidatorFunc) Validator {
	return &valueValidator{
		vt:     t_string,
		vfuncs: vfuncs,
	}
}

func Bool(vfuncs ...ValidatorFunc) Validator {
	return &valueValidator{
		vt:     t_bool,
		vfuncs: vfuncs,
	}
}
