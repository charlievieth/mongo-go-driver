package bsoncodec

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestKindCacheArray(t *testing.T) {
	// Make sure that reflect.UnsafePointer is the largest reflect.Type.
	//
	// The String() method of invalid reflect.Type types are of the format
	// "kind{NUMBER}".
	for rt := reflect.UnsafePointer + 1; rt < reflect.UnsafePointer+16; rt++ {
		s := rt.String()
		if !strings.Contains(s, strconv.Itoa(int(rt))) {
			t.Errorf("reflect.Type(%d) appears to be valid: %q", rt, s)
		}
	}
}

func TestKindCacheClone(t *testing.T) {
	e1 := new(kindEncoderCache)
	d1 := new(kindDecoderCache)
	for k := reflect.Invalid; k <= reflect.UnsafePointer; k++ {
		if k&1 == 0 {
			e1.Store(k, new(fakeCodec))
			d1.Store(k, new(fakeCodec))
		}
	}
	e2 := e1.Clone()
	for k := reflect.Invalid; k <= reflect.UnsafePointer; k++ {
		v1, ok1 := e1.Load(k)
		v2, ok2 := e2.Load(k)
		if !reflect.DeepEqual(v1, v2) || ok1 != ok2 {
			t.Errorf("Encoder: %v, %t != %v, %t", v1, ok1, v2, ok2)
		}
	}
	d2 := d1.Clone()
	for k := reflect.Invalid; k <= reflect.UnsafePointer; k++ {
		v1, ok1 := d1.Load(k)
		v2, ok2 := d2.Load(k)
		if !reflect.DeepEqual(v1, v2) || ok1 != ok2 {
			t.Errorf("Decoder: %v, %t != %v, %t", v1, ok1, v2, ok2)
		}
	}
}

func TestKindCacheEncoderNilEncoder(t *testing.T) {
	t.Run("Encoder", func(t *testing.T) {
		c := new(kindEncoderCache)
		c.Store(reflect.Invalid, ValueEncoder(nil))
		v, ok := c.Load(reflect.Invalid)
		if v != nil || ok {
			t.Errorf("Load of nil ValueEncoder should return: nil, false; got: %v, %t", v, ok)
		}
	})
	t.Run("Decoder", func(t *testing.T) {
		c := new(kindDecoderCache)
		c.Store(reflect.Invalid, ValueDecoder(nil))
		v, ok := c.Load(reflect.Invalid)
		if v != nil || ok {
			t.Errorf("Load of nil ValueDecoder should return: nil, false; got: %v, %t", v, ok)
		}
	})
}

var codecCacheTestTypes = [16]reflect.Type{
	reflect.TypeOf(uint8(0)),
	reflect.TypeOf(uint16(0)),
	reflect.TypeOf(uint32(0)),
	reflect.TypeOf(uint64(0)),
	reflect.TypeOf(uint(0)),
	reflect.TypeOf(uintptr(0)),
	reflect.TypeOf(int8(0)),
	reflect.TypeOf(int16(0)),
	reflect.TypeOf(int32(0)),
	reflect.TypeOf(int64(0)),
	reflect.TypeOf(int(0)),
	reflect.TypeOf(float32(0)),
	reflect.TypeOf(float64(0)),
	reflect.TypeOf(true),
	reflect.TypeOf(struct{ A int }{}),
	reflect.TypeOf(map[int]int{}),
}

func BenchmarkEncoderCacheLoad(b *testing.B) {
	typs := codecCacheTestTypes
	c := new(encoderCache)
	codec := new(fakeCodec)
	for _, t := range typs {
		c.Store(t, codec)
	}
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			c.Load(typs[i%len(typs)])
		}
	})
}

func BenchmarkEncoderCacheStore(b *testing.B) {
	typs := codecCacheTestTypes
	c := new(encoderCache)
	codec := new(fakeCodec)
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			c.Store(typs[i%len(typs)], codec)
		}
	})
}
