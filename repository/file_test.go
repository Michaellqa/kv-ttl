package repository

import (
	"fmt"
	"kv-ttl/kv"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestCacheBackup(t *testing.T) {
	// clean up
	//rx, _ := regexp.Compile(`.*\.json`)
	//files, _ := ioutil.ReadDir(".")
	//for _, f := range files {
	//	if rx.MatchString(f.Name()) {
	//		_ = os.Remove(f.Name())
	//	}
	//}

	fileStorage := NewFileRepo("snap.json")
	config := kv.Configuration{
		BackupInterval: 1 * time.Second,
		Storage:        fileStorage,
	}

	cache := kv.NewCache(config)
	values := []kv.T{
		{V: "one"},
		{V: "two"},
		{V: "three"},
	}
	for i, v := range values {
		cache.Add(fmt.Sprintf("%d", i), v)
	}
	go func() {
		time.Sleep(5 * time.Second)
		cache.AddWithTtl("4", kv.T{V: "four"}, 3*time.Second)
	}()

	time.Sleep(10500 * time.Millisecond)

	newCache := kv.NewCache(kv.Configuration{
		Storage: fileStorage,
	})
	storedValues := newCache.GetAll()

	sort.Slice(values, func(i, j int) bool { return values[i].V > values[j].V })
	sort.Slice(storedValues, func(i, j int) bool { return storedValues[i].V > storedValues[j].V })

	if !reflect.DeepEqual(values, storedValues) {
		t.Errorf("%v\n!=\n%v", values, storedValues)
	}
}
