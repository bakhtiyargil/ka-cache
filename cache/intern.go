package cache

import (
	"runtime"
	"sync"
)

type interned struct {
	value string
}

var internPool = struct {
	sync.RWMutex
	internMap map[string]*interned
}{internMap: make(map[string]*interned)}

func intern(stringValue string) string {
	internPool.RLock()
	if internObj, ok := internPool.internMap[stringValue]; ok {
		internPool.RUnlock()
		return internObj.value
	}
	internPool.RUnlock()

	internPool.Lock()
	if obj, ok := internPool.internMap[stringValue]; ok {
		internPool.Unlock()
		return obj.value
	}
	obj := &interned{value: stringValue}
	internPool.internMap[stringValue] = obj

	// when internObj is GC'ed, remove it from pool. Temp. Not Optimal.
	runtime.SetFinalizer(obj, func(i *interned) {
		internPool.Lock()
		delete(internPool.internMap, i.value)
		internPool.Unlock()
	})
	internPool.Unlock()

	return stringValue
}
