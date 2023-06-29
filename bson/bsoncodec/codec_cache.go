package bsoncodec

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type encoderCache struct {
	cache sync.Map // map[reflect.Type]ValueEncoder
}

func (c *encoderCache) Store(rt reflect.Type, enc ValueEncoder) {
	c.cache.Store(rt, enc)
}

func (c *encoderCache) Load(rt reflect.Type) (ValueEncoder, bool) {
	if v, _ := c.cache.Load(rt); v != nil {
		return v.(ValueEncoder), true
	}
	return nil, false
}

func (c *encoderCache) LoadOrStore(rt reflect.Type, enc ValueEncoder) ValueEncoder {
	if v, loaded := c.cache.LoadOrStore(rt, enc); loaded {
		enc = v.(ValueEncoder)
	}
	return enc
}

func (c *encoderCache) Clone() *encoderCache {
	cc := new(encoderCache)
	c.cache.Range(func(k, v interface{}) bool {
		if k != nil && v != nil {
			cc.cache.Store(k, v)
		}
		return true
	})
	return cc
}

type decoderCache struct {
	cache sync.Map // map[reflect.Type]ValueDecoder
}

func (c *decoderCache) Store(rt reflect.Type, dec ValueDecoder) {
	c.cache.Store(rt, dec)
}

func (c *decoderCache) Load(rt reflect.Type) (ValueDecoder, bool) {
	if v, _ := c.cache.Load(rt); v != nil {
		return v.(ValueDecoder), true
	}
	return nil, false
}

func (c *decoderCache) LoadOrStore(rt reflect.Type, dec ValueDecoder) ValueDecoder {
	if v, loaded := c.cache.LoadOrStore(rt, dec); loaded {
		dec = v.(ValueDecoder)
	}
	return dec
}

func (c *decoderCache) Clone() *decoderCache {
	cc := new(decoderCache)
	c.cache.Range(func(k, v interface{}) bool {
		if k != nil && v != nil {
			cc.cache.Store(k, v)
		}
		return true
	})
	return cc
}

// atomic.Value requires that all calls to Store() have the same concrete
// type so we wrap the ValueEncoder with a kindEncoderCacheEntry to ensure
// the type is always the same (since different concrete types may implement
// the ValueEncoder interface).
type kindEncoderCacheEntry struct {
	enc ValueEncoder
}

type kindEncoderCache struct {
	entries [reflect.UnsafePointer]atomic.Value // *kindEncoderCacheEntry
}

func (c *kindEncoderCache) Store(rt reflect.Kind, enc ValueEncoder) {
	if enc != nil && rt < reflect.Kind(len(c.entries)) {
		c.entries[rt].Store(&kindEncoderCacheEntry{enc: enc})
	}
}

func (c *kindEncoderCache) Load(rt reflect.Kind) (ValueEncoder, bool) {
	if rt < reflect.Kind(len(c.entries)) {
		if ent, ok := c.entries[rt].Load().(*kindEncoderCacheEntry); ok {
			return ent.enc, ent.enc != nil
		}
	}
	return nil, false
}

func (c *kindEncoderCache) Clone() *kindEncoderCache {
	cc := new(kindEncoderCache)
	for i, v := range c.entries {
		if val := v.Load(); val != nil {
			cc.entries[i].Store(val)
		}
	}
	return cc
}

// atomic.Value requires that all calls to Store() have the same concrete
// type so we wrap the ValueDecoder with a kindDecoderCacheEntry to ensure
// the type is always the same (since different concrete types may implement
// the ValueDecoder interface).
type kindDecoderCacheEntry struct {
	dec ValueDecoder
}

type kindDecoderCache struct {
	entries [reflect.UnsafePointer]atomic.Value // *kindDecoderCacheEntry
}

func (c *kindDecoderCache) Store(rt reflect.Kind, dec ValueDecoder) {
	if rt < reflect.Kind(len(c.entries)) {
		c.entries[rt].Store(&kindDecoderCacheEntry{dec: dec})
	}
}

func (c *kindDecoderCache) Load(rt reflect.Kind) (ValueDecoder, bool) {
	if rt < reflect.Kind(len(c.entries)) {
		if ent, ok := c.entries[rt].Load().(*kindDecoderCacheEntry); ok {
			return ent.dec, ent.dec != nil
		}
	}
	return nil, false
}

func (c *kindDecoderCache) Clone() *kindDecoderCache {
	cc := new(kindDecoderCache)
	for i, v := range c.entries {
		if val := v.Load(); val != nil {
			cc.entries[i].Store(val)
		}
	}
	return cc
}
