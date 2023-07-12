// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build !go1.20
// +build !go1.20

package bsoncodec

import "reflect"

// TODO: this might not be worth it
func zeroValue(val reflect.Value) {
	val.Set(reflect.Zero(val.Type()))
}
