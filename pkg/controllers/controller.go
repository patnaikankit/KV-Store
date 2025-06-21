package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/patnaikankit/KV-Store.git/pkg/store"
)

type RequestBody struct {
	Key   string `json:"Key"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not found", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Please Provide a key", http.StatusNotImplemented)
		return
	}

	val, err := store.KVStore.Get(key)
	if err != nil {
		http.Error(w, "Key does not exist", http.StatusNotFound)
		return
	}

	resp := map[string]string{
		"data": val,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func SetKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var data RequestBody
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	if data.Key == "" {
		http.Error(w, "Key cannot be empty", http.StatusBadRequest)
		return
	}

	var ttl time.Duration
	ttl, err = time.ParseDuration(data.TTL)
	if err != nil {
		ttl = time.Hour * 24
	}

	err = store.KVStore.Set(
		data.Key,
		store.KVMapValue{Value: data.Value, ExpireAt: time.Now().Add(ttl)},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]string{
		"status":  "Success",
		"message": "Key-value pair set successfully",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func UpdateKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var data RequestBody
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	if data.Key == "" {
		http.Error(w, "Key cannot be empty", http.StatusBadRequest)
		return
	}

	isUpdated, err := store.KVStore.Update(data.Key, data.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var resp map[string]string

	if isUpdated {
		resp = map[string]string{
			"status":  "Success",
			"message": "Key Updated successfully",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		resp = map[string]string{
			"status":  "Failure",
			"message": "Couldn't update successfully",
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func DeleteKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("Key")
	if key == "" {
		http.Error(w, "Please provide key", http.StatusBadRequest)
		return
	}

	if err := store.KVStore.Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{
		"status":  "Success",
		"message": "Key-value deleted successfully",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
