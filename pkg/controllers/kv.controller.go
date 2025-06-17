package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/patnaikankit/KV-Store.git/pkg/store"
)

type SetBody struct {
	Key string `json:"key"`
}

type RequestBody struct {
	Key   string `json:"Key"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}

func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not found", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Please Provide a key", http.StatusNotImplemented)
		return
	}

	val := store.KVStore.Get(key)
	if val == "" {
		http.Error(w, "Key does not exits", http.StatusNotFound)
		return
	}

	resp := map[string]string{
		"data": val,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func Set(w http.ResponseWriter, r *http.Request) {
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

	store.KVStore.Set(
		data.Key,
		store.KVMapValue{Value: data.Value, ExpireAt: time.Now().Add(ttl)},
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]string{
		"status": "Success", 
		"message": "Key-value pair set successfully"
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return 
	}
}

func Update(w *http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch || r.Method != http.MethodPut {
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

	isUpdated := store.KVStore.Update(data.Key, data.Value)
	w.Header().Set("Content-Type", "application/json")
	var resp map[string]string

	if isUpdated {
		resp = map[string]string{
			"status": "Success",
			"message": "Key Updated successfully"
		}
		w.WriteHeader(http.StatusOK)
	}
	else {
		resp = map[string]string{
			"status": "failure",
			"message": "Couldn't update successfully"
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("Key")
	if key == "" {
		http.Error(w, "Please provide key", http.StatusBadRequest)
		return
	}

	store.KVStore.Delete(key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{
		"status": "Success",
		"message": "Key-value deleted successfully"
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}