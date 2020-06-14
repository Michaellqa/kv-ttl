package repository

import (
	"fmt"
	"io/ioutil"
	"kv-ttl/kv"
	"os"
	"reflect"
	"regexp"
	"sort"
	"testing"
	"time"
)

// Removes all the previously created *.json files from the current directory.
// Initiate a cache with backing up to file. Waits for scheduled backup happened.
// Creates another cache instance based on the same file.
// Sorts values received from the new cache and tests against the original data.
func TestCacheBackup(t *testing.T) {
	// clean up
	rx, _ := regexp.Compile(`.*\.json`)
	files, _ := ioutil.ReadDir(".")
	for _, f := range files {
		if rx.MatchString(f.Name()) {
			_ = os.Remove(f.Name())
		}
	}
	// init
	fileStorage := NewFileRepo("snap.json")
	config := kv.Configuration{
		BackupInterval: 1 * time.Second,
		Storage:        fileStorage,
	}
	// insert
	cache := kv.NewCache(config)
	values := []kv.T{
		{V: "one"},
		{V: "two"},
		{V: "three"},
	}
	for i, v := range values {
		cache.Add(fmt.Sprintf("%d", i), v)
	}
	// verify
	time.Sleep(1200 * time.Millisecond)
	newCache := kv.NewCache(kv.Configuration{Storage: fileStorage})
	storedValues := newCache.ListAll()

	sort.Slice(values, func(i, j int) bool { return values[i].V > values[j].V })
	sort.Slice(storedValues, func(i, j int) bool { return storedValues[i].V > storedValues[j].V })

	if !reflect.DeepEqual(values, storedValues) {
		t.Errorf("%v\n!=\n%v", values, storedValues)
	}
}
