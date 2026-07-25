package kv

import "testing"

func TestPairKV(t *testing.T) {
	p := NewPair(42, "answer")

	key, value := p.KV()
	if key != 42 || value != "answer" {
		t.Fatalf("KV returned (%v, %v), want (42, answer)", key, value)
	}

	valuePtrKey, valuePtr := p.KVMut()
	if valuePtrKey != 42 || valuePtr == nil {
		t.Fatalf("KVMut returned (%v, %v), want key 42 and non-nil value pointer", valuePtrKey, valuePtr)
	}

	*valuePtr = "updated"
	if p.Value != "updated" {
		t.Fatalf("KVMut pointer did not update value, got %q", p.Value)
	}
}
