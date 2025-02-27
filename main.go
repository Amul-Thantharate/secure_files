package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 📝 Initialize log file
var logFile = "audit.log"

func init() {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ Error opening log file:", err)
		os.Exit(1)
	}
	log.SetOutput(f)
}

// 🔍 Log events (Encryption/Decryption)
func logEvent(event, file string, success bool) {
	status := "✅ SUCCESS"
	if !success {
		status = "❌ FAILED"
	}
	logMsg := fmt.Sprintf("%s | %s | File: %s\n", time.Now().Format("2006-01-02 15:04:05"), status, event+" "+file)
	log.Println(logMsg)
	fmt.Println(logMsg) // Also print to console
}

// 🔐 Generate a random key (AES-256)
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

// 💾 Save key securely
func SaveKeyToFile(key []byte, filename string) error {
	return os.WriteFile(filename, []byte(hex.EncodeToString(key)), 0600)
}

// 🔑 Load key from file
func LoadKeyFromFile(filename string) ([]byte, error) {
	keyHex, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(keyHex)))
}

// 🛡️ Generate HMAC for integrity check
func generateHMAC(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// 🔏 Encrypt a file
func EncryptFile(inputFile, outputFile string, key []byte) error {
	plaintext, err := os.ReadFile(inputFile)
	if err != nil {
		logEvent("Encryption", inputFile, false)
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logEvent("Encryption", inputFile, false)
		return err
	}

	nonce := make([]byte, 12) // AES-GCM nonce
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logEvent("Encryption", inputFile, false)
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logEvent("Encryption", inputFile, false)
		return err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	hmacHash := generateHMAC(ciphertext, key)

	output := append(nonce, append(hmacHash, ciphertext...)...)
	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		logEvent("Encryption", inputFile, false)
		return err
	}

	logEvent("Encryption", inputFile, true)
	return nil
}

// 🔓 Decrypt a file
func DecryptFile(inputFile, outputFile string, key []byte) error {
	ciphertext, err := os.ReadFile(inputFile)
	if err != nil {
		logEvent("Decryption", inputFile, false)
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logEvent("Decryption", inputFile, false)
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logEvent("Decryption", inputFile, false)
		return err
	}

	nonceSize := aesGCM.NonceSize()
	hmacSize := sha256.Size

	if len(ciphertext) < nonceSize+hmacSize {
		logEvent("Decryption", inputFile, false)
		return fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	hmacReceived := ciphertext[nonceSize : nonceSize+hmacSize]
	ciphertext = ciphertext[nonceSize+hmacSize:]

	hmacCalculated := generateHMAC(ciphertext, key)
	if !hmac.Equal(hmacCalculated, hmacReceived) {
		logEvent("Tamper Detection", inputFile, false)
		return fmt.Errorf("integrity check failed! Possible tampering detected ❌")
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logEvent("Decryption", inputFile, false)
		return err
	}

	if err := os.WriteFile(outputFile, plaintext, 0644); err != nil {
		logEvent("Decryption", inputFile, false)
		return err
	}

	logEvent("Decryption", inputFile, true)
	return nil
}

// 📁 Encrypt all files in a directory
func EncryptDirectory(inputDir, outputDir string, key []byte) error {
	return filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, _ := filepath.Rel(inputDir, path)
			outputFile := filepath.Join(outputDir, relPath+".enc")
			os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
			return EncryptFile(path, outputFile, key)
		}
		return nil
	})
}

// 📂 Decrypt all files in a directory
func DecryptDirectory(inputDir, outputDir string, key []byte) error {
	return filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".enc") {
			relPath, _ := filepath.Rel(inputDir, path)
			outputFile := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".enc"))
			os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
			return DecryptFile(path, outputFile, key)
		}
		return nil
	})
}

// 🎛️ Interactive CLI Menu
func ShowMenu() {
	fmt.Println("\n🔐 ==== Secure File Encryption Menu ==== 🔐")
	fmt.Println("1️⃣ Encrypt a file")
	fmt.Println("2️⃣ Decrypt a file")
	fmt.Println("3️⃣ Encrypt a directory")
	fmt.Println("4️⃣ Decrypt a directory")
	fmt.Println("5️⃣ View security audit logs")
	fmt.Println("6️⃣ Exit 🚪")
	fmt.Print("👉 Enter your choice: ")

	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1, 2:
		fmt.Print("📄 Enter the file path: ")
		var inputFile string
		fmt.Scan(&inputFile)

		fmt.Print("📂 Enter the output file path: ")
		var outputFile string
		fmt.Scan(&outputFile)

		key, err := LoadKeyFromFile("encryption.key")
		if err != nil {
			key, _ = GenerateRandomKey()
			SaveKeyToFile(key, "encryption.key")
		}

		if choice == 1 {
			EncryptFile(inputFile, outputFile, key)
		} else {
			DecryptFile(inputFile, outputFile, key)
		}

	case 3, 4:
		fmt.Print("📁 Enter the directory path: ")
		var inputDir string
		fmt.Scan(&inputDir)

		fmt.Print("📂 Enter the output directory path: ")
		var outputDir string
		fmt.Scan(&outputDir)

		key, _ := LoadKeyFromFile("encryption.key")

		if choice == 3 {
			EncryptDirectory(inputDir, outputDir, key)
		} else {
			DecryptDirectory(inputDir, outputDir, key)
		}

	case 5:
		logData, _ := os.ReadFile(logFile)
		fmt.Println("\n📜 Security Audit Logs:\n" + string(logData))

	case 6:
		fmt.Println("👋 Exiting... Have a great day! 🚀")
		os.Exit(0)

	default:
		fmt.Println("❌ Invalid choice. Please try again.")
	}
}

func main() {
	for {
		ShowMenu()
	}
}

