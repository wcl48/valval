package valval

import (
	"errors"
	"fmt"
	"strings"
)

type ObjectFieldError struct {
	Name string
	Err  error
}

func (ofe ObjectFieldError) Error() string {
	return fmt.Sprintf("%s : %s", ofe.Name, ofe.Err.Error())
}

type ObjectError []ObjectFieldError

func (oe ObjectError) Error() string {
	buf := []string{}
	for _, in := range oe {
		buf = append(buf, in.Error())
	}
	return strings.Join(buf, "\n")
}

type M map[string]Validator

type ObjectValidatorFunc func(content map[string]interface{}) error

type ObjectValidator interface {
	Validator
	Self(...ObjectValidatorFunc) ObjectValidator
}

type objectValidator struct {
	vMap       M
	selfVfuncs []ObjectValidatorFunc
}

func (ov *objectValidator) Validate(val interface{}) error {
	valMap, err := obj2Map(val)
	if err != nil {
		return typeMissmatchError("object")
	}

	// inner
	if err := ov.checkInner(valMap); err != nil {
		return err
	}

	// self
	if err := ov.checkSelf(valMap); err != nil {
		return err
	}
	return nil
}

func (ov *objectValidator) checkInner(valMap map[string]interface{}) error {
	var errs []ObjectFieldError
	for k, fv := range ov.vMap {
		fValue := valMap[k]
		err := fv.Validate(fValue)
		if err != nil {
			errs = append(errs, ObjectFieldError{
				Name: k,
				Err:  err,
			})
		}
	}

	if errs != nil {
		return ObjectError(errs)
	}
	return nil
}

func (ov *objectValidator) checkSelf(valMap map[string]interface{}) error {
	errs := []error{}
	for _, svf := range ov.selfVfuncs {
		err := svf(valMap)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// TODO Detail
		return errors.New("error")
	}
	return nil
}

func (ov *objectValidator) Self(vfs ...ObjectValidatorFunc) ObjectValidator {
	// copy
	newOv := *ov
	newOv.selfVfuncs = vfs
	return &newOv
}

func Object(m M) ObjectValidator {
	return &objectValidator{
		vMap: m,
	}
}
