package cache

import (
	"runtime"
	"sync"
)

type interned struct {
	val string
}

var internPool = struct {
	sync.RWMutex
	m map[string]*interned
}{m: make(map[string]*interned)}

func intern(s string) string {
	internPool.RLock()
	if obj, ok := internPool.m[s]; ok {
		internPool.RUnlock()
		return obj.val
	}
	internPool.RUnlock()

	internPool.Lock()
	if obj, ok := internPool.m[s]; ok {
		internPool.Unlock()
		return obj.val
	}
	obj := &interned{val: s}
	internPool.m[s] = obj
	// when obj is GC'ed, remove it from pool. Temp. Not Optimal.
	runtime.SetFinalizer(obj, func(i *interned) {
		internPool.Lock()
		delete(internPool.m, i.val)
		internPool.Unlock()
	})
	internPool.Unlock()

	return s
}
