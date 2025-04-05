package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateRandomKey(t *testing.T) {
	key1, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	if len(key1) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key1))
	}

	// Test uniqueness
	key2, _ := GenerateRandomKey()
	if bytes.Equal(key1, key2) {
		t.Error("Generated keys should be unique")
	}
}

func TestSaveAndLoadKey(t *testing.T) {
	// Setup
	tempKeyFile := filepath.Join(t.TempDir(), "test.key")
	originalKey, _ := GenerateRandomKey()

	// Test saving
	err := SaveKeyToFile(originalKey, tempKeyFile)
	if err != nil {
		t.Fatalf("Failed to save key: %v", err)
	}

	// Test loading
	loadedKey, err := LoadKeyFromFile(tempKeyFile)
	if err != nil {
		t.Fatalf("Failed to load key: %v", err)
	}

	if !bytes.Equal(originalKey, loadedKey) {
		t.Error("Loaded key doesn't match original key")
	}
}

func TestFileEncryptionDecryption(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	plainFile := filepath.Join(tempDir, "plain.txt")
	encryptedFile := filepath.Join(tempDir, "encrypted.enc")
	decryptedFile := filepath.Join(tempDir, "decrypted.txt")
	
	originalContent := []byte("Hello, this is a test message!")
	err := os.WriteFile(plainFile, originalContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Generate key
	key, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test encryption
	err = EncryptFile(plainFile, encryptedFile, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Verify encrypted file exists and is different from original
	encryptedContent, err := os.ReadFile(encryptedFile)
	if err != nil {
		t.Fatalf("Failed to read encrypted file: %v", err)
	}
	if bytes.Equal(encryptedContent, originalContent) {
		t.Error("Encrypted content should be different from original content")
	}

	// Test decryption
	err = DecryptFile(encryptedFile, decryptedFile, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Verify decrypted content matches original
	decryptedContent, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}
	if !bytes.Equal(decryptedContent, originalContent) {
		t.Error("Decrypted content doesn't match original content")
	}
}

func TestDirectoryEncryptionDecryption(t *testing.T) {
	// Setup test directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	encryptedDir := filepath.Join(tempDir, "encrypted")
	decryptedDir := filepath.Join(tempDir, "decrypted")

	// Create test directory structure
	dirs := []string{
		filepath.Join(inputDir, "subdir1"),
		filepath.Join(inputDir, "subdir2"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Create test files
	files := map[string]string{
		filepath.Join(inputDir, "file1.txt"):          "Content of file 1",
		filepath.Join(inputDir, "subdir1/file2.txt"):  "Content of file 2",
		filepath.Join(inputDir, "subdir2/file3.txt"):  "Content of file 3",
	}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Generate key
	key, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test directory encryption
	err = EncryptDirectory(inputDir, encryptedDir, key)
	if err != nil {
		t.Fatalf("Directory encryption failed: %v", err)
	}

	// Test directory decryption
	err = DecryptDirectory(encryptedDir, decryptedDir, key)
	if err != nil {
		t.Fatalf("Directory decryption failed: %v", err)
	}

	// Verify decrypted files match original files
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(inputDir, path)
		decryptedPath := filepath.Join(decryptedDir, relPath)

		originalContent, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read original file %s: %v", path, err)
			return nil
		}

		decryptedContent, err := os.ReadFile(decryptedPath)
		if err != nil {
			t.Errorf("Failed to read decrypted file %s: %v", decryptedPath, err)
			return nil
		}

		if !bytes.Equal(originalContent, decryptedContent) {
			t.Errorf("Content mismatch for file %s", relPath)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to verify decrypted files: %v", err)
	}
}

func TestHMACVerification(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	plainFile := filepath.Join(tempDir, "plain.txt")
	encryptedFile := filepath.Join(tempDir, "encrypted.enc")
	decryptedFile := filepath.Join(tempDir, "decrypted.txt")

	originalContent := []byte("Test content for HMAC verification")
	err := os.WriteFile(plainFile, originalContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Generate keys
	key1, _ := GenerateRandomKey()
	key2, _ := GenerateRandomKey()

	// Encrypt with key1
	err = EncryptFile(plainFile, encryptedFile, key1)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to decrypt with key2 (should fail due to HMAC verification)
	err = DecryptFile(encryptedFile, decryptedFile, key2)
	if err == nil {
		t.Error("Expected decryption to fail with wrong key")
	} else if !strings.Contains(err.Error(), "integrity check failed") {
		t.Errorf("Expected integrity check failure, got: %v", err)
	} else {
		t.Log("Successfully detected tampering attempt with wrong key")
	}

	// Verify successful decryption with correct key
	err = DecryptFile(encryptedFile, decryptedFile, key1)
	if err != nil {
		t.Errorf("Failed to decrypt with correct key: %v", err)
	}

	// Verify content
	decryptedContent, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}
	if !bytes.Equal(decryptedContent, originalContent) {
		t.Error("Decrypted content doesn't match original content")
	}
}
