package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/patnaikankit/KV-Store.git/pkg/utils"
)

type KVMapValue struct {
	Value    string    `json:"value"`
	ExpireAt time.Time `json:"expire_at"`
}

type kvStore struct {
	Store    map[string]KVMapValue
	Mu       sync.RWMutex
	FileName string
	Count    int64
}

var (
	KVStore *kvStore
	once    sync.Once
)

func NewKVStore(fileName string) *kvStore {
	return &kvStore{
		Store:    make(map[string]KVMapValue),
		Mu:       sync.RWMutex{},
		FileName: fileName,
		Count:    0,
	}
}

func InitKvStore(filename string) {
	once.Do(func() {
		KVStore = NewKVStore(filename)
		if err := KVStore.loadFromFile(); err != nil {
			fmt.Println("Error loading data from file: ", err)
		}
		go KVStore.periodicSaveData()
		go KVStore.removeExpiredKey()
	})
}

func (kv *kvStore) Get(key string) (string, error) {
	kv.Mu.RLock()
	defer kv.Mu.RLocker().Unlock()

	resp, ok := kv.Store[key]
	if !ok {
		utils.AddLog("READ", "FAILED", key)
		return "", errors.New("key not found")
	}

	utils.AddLog("READ", "SUCCESS", key)
	return resp.Value, nil
}

func (kv *kvStore) Set(key string, value KVMapValue) error {
	kv.Mu.Lock()
	defer kv.Mu.Unlock()

	if _, exists := kv.Store[key]; exists {
		utils.AddLog("SET", "FAILED", key)
		return errors.New("key already exists")
	}

	kv.Store[key] = value
	kv.Count++
	utils.AddLog("SET", "SUCCESS", key)
	return nil
}

func (kv *kvStore) Update(key string, value string) (bool, error) {
	kv.Mu.Lock()
	defer kv.Mu.Unlock()

	if val, exists := kv.Store[key]; exists {
		kv.Store[key] = KVMapValue{Value: value, ExpireAt: val.ExpireAt}
		kv.Count++
		utils.AddLog("UPDATE", "SUCCESS", key)
		return true, nil
	}

	utils.AddLog("UPDATE", "FAILED", key)
	return false, errors.New("key doesn't exist")
}

func (kv *kvStore) Delete(key string) error {
	kv.Mu.Lock()
	defer kv.Mu.Unlock()

	if _, exists := kv.Store[key]; exists {
		delete(kv.Store, key)
		kv.Count++
		utils.AddLog("DELETE", "SUCCESS", key)
		return nil
	}

	utils.AddLog("DELETE", "FAILED", key)
	return errors.New("key doesn't exist")
}

func (kv *kvStore) loadFromFile() error {
	resp, err := os.ReadFile(kv.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist, starting with an empty Store")
			return nil
		}
		return err
	}
	return json.Unmarshal(resp, &kv.Store)
}

func (kv *kvStore) periodicSaveData() {
	timer := time.NewTicker(1 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			kv.checkAndSave("Periodic")
		default:
			time.Sleep(1000 * time.Millisecond)
			kv.checkAndSave("Count")
		}
	}
}

func (kv *kvStore) checkAndSave(opType string) {
	kv.Mu.Lock()
	defer kv.Mu.Unlock()

	if opType == "Periodic" || kv.Count >= 5 {
		err := kv.saveDataToFile()
		if err != nil {
			fmt.Printf("error saving to file: %v\n", err)
		}
	}
}

func (kv *kvStore) saveDataToFile() error {
	jsonData, err := json.Marshal(kv.Store)
	if err != nil {
		return err
	}

	err = os.WriteFile(kv.FileName, []byte(jsonData), 0644)
	if err != nil {
		return err
	}

	kv.Count = 0
	return nil
}

func (kv *kvStore) removeExpiredKey() {
	timer := time.NewTicker(time.Second * 5)
	defer timer.Stop()

	for range timer.C {
		for key, val := range kv.Store {
			if time.Now().After(val.ExpireAt) {
				kv.Delete(key)
				utils.AddLog("CLEANUP", "Expired", key)
			}
		}
	}
}
