package cache

import (
	"time"
	"testing"
)

func TestCache_Expiry(t *testing.T) {

	cache, err := NewMinimalCache(3, 1) // cachesize, timeout in sec
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	cache.Add("key1", "value1")
	cache.Add("key2", "value2") // should evict


	result1, _:= cache.Get("key1")
	if result1 == false {
		t.Errorf("should not have an eviction")
	}

	result2, _:= cache.Get("key2")
	if result2 == false {
		t.Errorf("should not have an eviction")
	}

	time.Sleep(time.Second)
	result2, _= cache.Get("key2")
	if result2 == true {
		t.Errorf("After sleep, testK2 should have an eviction")
	}
}

func TestCache_Evict(t *testing.T) {

	cache, err := NewMinimalCache(1, 10) // cachesize, timeout in sec
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	cache.Add("key1", "value1")
	cache.Add("key2", "value2") // should evict


	result1, _:= cache.Get("key1")
	if result1 == true {
		t.Errorf("should have an eviction")
	}

	result2, _:= cache.Get("key2")
	if result2 == false {
		t.Errorf("should not have an eviction")
	}
}
