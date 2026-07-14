package cache

import (
	"time"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	c := NewCache(1*time.Second, 100*time.Millisecond, 10)
	defer c.Stop()

	c.Set("key1", "value1", 0)
	val, ok := c.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("вместо ожидаемого значения получено %v", val)
	}

	_, ok = c.Get("not-exist")
	if ok {
		t.Errorf("ключ не может существовать")
	}
}

func TestTTL(t *testing.T) {
	c := NewCache(1*time.Second, 100*time.Millisecond, 10)	
	defer c.Stop()

	c.Set("key", "value", 100*time.Millisecond)

	_, ok := c.Get("key")
	if !ok {
		t.Error("ключ еще должен существовать")
	}
	
	time.Sleep(150 * time.Millisecond)
	
	_, ok = c.Get("key")
	if ok {
		t.Error("ключ не должен существовать после TTL")
	}
}

func TestDelete(t *testing.T) {
	c := NewCache(5*time.Minute, 10*time.Second, 10)
	defer c.Stop()

	c.Set("key", "value", 0)
	c.Delete("key")

	_, ok := c.Get("key")
	if ok {
		t.Errorf("ключ существует после удаления")
	}	
}

func TestMaxSize(t *testing.T) {
	c := NewCache(5*time.Minute, 10*time.Second, 2)
	defer c.Stop()

	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0)

	if c.Size() != 2 {
		t.Errorf("не соблюден ожидаемый размер кэша")
	}
}

func TestKeys(t *testing.T) {
	c := NewCache(5*time.Minute, 10*time.Second, 10)
	defer c.Stop()

	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0)
	
	keys := c.Keys()
	if len(keys) != 3 {
		t.Errorf("ожидалось 3 ключа, получено: %d", len(keys))
	}
}


