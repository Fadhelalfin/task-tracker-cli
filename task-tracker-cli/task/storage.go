package task

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"fmt"
)

// dataDir adalah direktori tempat menyimpan file JSON relatif terhadap executable
const dataDir = "data"
const tasksFile = "tasks.json"

var (
	// Mutex untuk menghindari race condition saat mengakses file
	mu sync.Mutex 
)

// getTasksFilePath mengembalikan jalur lengkap ke file tasks.json
func getTasksFilePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex) // Ini yang benar: filepath.Dir(string)
	
	// Lokasi file JSON: di dalam subdirektori 'data' di samping executable
	return filepath.Join(exPath, dataDir, tasksFile), nil // Ini yang benar: filepath.Join(string, ...)
}

// LoadTasks membaca tugas dari file JSON
func LoadTasks() ([]Task, error) {
	mu.Lock()
	defer mu.Unlock()

	filePath, err := getTasksFilePath() // Panggil fungsi ini untuk mendapatkan jalur
	if err != nil {
		return nil, err
	}

	// Dapatkan jalur direktori data
	dataDirPath := filepath.Dir(filePath) // Dapatkan direktori dari filePath

	// Pastikan direktori data ada
	if _, err := os.Stat(dataDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDirPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
		return []Task{}, nil // Kembalikan slice kosong jika direktori/file belum ada
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // File tidak ada, kembalikan slice kosong
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}
	return tasks, nil
}

// SaveTasks menulis tugas ke file JSON
func SaveTasks(tasks []Task) error {
	mu.Lock()
	defer mu.Unlock()

	filePath, err := getTasksFilePath() // Panggil fungsi ini untuk mendapatkan jalur
	if err != nil {
		return err
	}

	// Dapatkan jalur direktori data
	dataDirPath := filepath.Dir(filePath) // Dapatkan direktori dari filePath

	// Pastikan direktori data ada sebelum menulis file
	if _, err := os.Stat(dataDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	data, err := json.MarshalIndent(tasks, "", "  ") // Format JSON dengan indentasi
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	
	// os.WriteFile adalah alternatif yang lebih modern dari ioutil.WriteFile di Go 1.16+
	return os.WriteFile(filePath, data, 0644) // Menggunakan os.WriteFile
}

// GetNextID mendapatkan ID berikutnya untuk tugas baru
func GetNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}