package store

import (
	"errors"
	"sync"
	"time"

	"github.com/patnaikankit/KV-Store.git/pkg/utils"
)

type kvStore struct {
	Store    map[string]kvMapValue
	Mu       sync.RWMutex
	FileName string
	Count    int64
}

type kvMapValue struct {
	Value    string    `json:"value"`
	ExpireAt time.Time `json:"expire_at"`
}

var (
	KVStore *kvStore
	once    sync.Once
)

func NewKVStore(fileName string) *kvStore {
	return &kvStore{
		Store:    make(map[string]kvMapValue),
		Mu:       sync.RWMutex{},
		FileName: fileName,
		Count:    0,
	}

}

func (kv *kvStore) Get(key string) (string, error) {
	kv.Mu.RLock()
	defer kv.Mu.RLocker().Unlock()
	resp, ok := kv.Store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	utils.AddLog("READ", "SUCCESS", key)
	return resp.Value, nil
}

func (kv *kvStore) Set(key string, value kvMapValue) {
	kv.Mu.Lock()
	defer kv.Mu.Unlock()

	kv.Store[key] = value
	kv.Count++
	utils.AddLog("SET", "SUCCESS", key)
}

func (kv *kvStore) Update(key string, value string) (bool, error) {
	kv.Mu.Lock()
	defer kv.Mu.Unlock()

	if val, exists := kv.Store[key]; exists {
		kv.Store[key] := kvMapValue{Value: value, ExpireAt: val.ExpireAt}
		kv.Count++
		utils.AddLog("UPDATE", "SUCCESS", key)
		return true, nil
	}
	utils.AddLog("UPDATE", "FAILED", key)
	return false, errors.New("Failed to update key")
}

func (kv *kvStore) Delete(key string) {
	kv.Mu.Lock()
	defer kv.Mu.Unlock()

	delete(kv.Store, key)
	kv.Count++
	utils.AddLog("DELETE", "SUCCESS", key)
}
