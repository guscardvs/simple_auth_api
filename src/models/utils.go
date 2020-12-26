package models

import (
	"reflect"

	"github.com/fatih/structs"
)

// EditModel assigns new value to modelStruct if values are notZero
func EditModel(oldModel, newModel interface{}) {

	oldValues := reflect.ValueOf(oldModel).Elem()
	newValues := reflect.ValueOf(newModel).Elem()

	fields := structs.Names(newModel)

	for index := range fields {
		newValue := newValues.FieldByName(fields[index])
		oldValue := oldValues.FieldByName(fields[index])

		if !newValue.IsZero() {
			oldValue.Set(newValue)
		}
	}

}
