# 🔐 Secure File Encryption Tool

A robust and secure file encryption tool written in Go that provides military-grade encryption for files and directories using AES-256-GCM with HMAC integrity verification.

## ✨ Features

- 🔒 AES-256-GCM encryption for strong security
- 🛡️ HMAC-SHA256 for integrity verification
- 📁 Support for both single file and directory encryption
- 📝 Comprehensive audit logging
- 🔑 Secure key management
- 🖥️ Interactive CLI interface

## 🚀 Installation

1. Ensure you have Go 1.16 or later installed on your system.
2. Clone this repository:
```bash
git clone https://github.com/amulthantharate/secure_files.git
cd secure_files
```
3. Build the application:
```bash
go build -o secure-files
```

## 💫 Usage

Run the application:
```bash
./secure-files
```

The interactive menu provides the following options:
1. Encrypt a single file
2. Decrypt a single file
3. Encrypt an entire directory
4. Decrypt an entire directory

### 🔑 Key Management

- The application automatically generates and manages encryption keys
- Keys are stored securely with appropriate file permissions (0600)
- Each encryption operation uses a unique nonce for added security

### 📝 Audit Logging

All operations are logged to `audit.log` with timestamps and status information for security auditing and troubleshooting.

## 🛡️ Security Features

- AES-256-GCM for authenticated encryption
- Secure random key generation
- HMAC-SHA256 for integrity verification
- Tamper detection
- Secure file permissions for sensitive files
- Comprehensive audit logging

## ⚠️ Important Notes

1. Always keep backups of your original files before encryption
2. Store your encryption keys securely
3. The `.enc` extension is added to encrypted files automatically

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.


## 🔒 Security

If you discover any security-related issues, please email [amulthantharate@gmail.com] instead of using the issue tracker.
