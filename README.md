# 🔐 Secure File Encryption Tool

A robust and secure file encryption tool written in Go that provides military-grade encryption for files and directories using AES-256-GCM with HMAC integrity verification.

## ✨ Features

- 🔒 AES-256-GCM encryption for strong security
- 🛡️ HMAC-SHA256 for integrity verification
- 📁 Support for both single file and directory encryption
- 📝 Comprehensive audit logging
- 🔑 Secure key management
- 🖥️ Interactive CLI interface

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or later
- Make (optional, for using Makefile commands)

### Installation

1. Clone this repository:
```bash
git clone https://github.com/amulthantharate/secure_files.git
cd secure_files
```

2. Build using one of these methods:

Using Go directly:
```bash
go build -o secure-files
```

Using Make:
```bash
make build
```

### Quick Start

Run the application:
```bash
./secure-files
```

## 💫 Usage Guide

The interactive menu provides the following options:

1. **Encrypt a single file**
   - Select option 1
   - Enter source file path
   - Enter destination path for encrypted file

2. **Decrypt a single file**
   - Select option 2
   - Enter encrypted file path (.enc)
   - Enter destination path for decrypted file

3. **Encrypt a directory**
   - Select option 3
   - Enter source directory path
   - Enter destination directory for encrypted files

4. **Decrypt a directory**
   - Select option 4
   - Enter encrypted directory path
   - Enter destination directory for decrypted files

5. **View audit logs**
   - Select option 5 to view security logs

## 🛠️ Build Options

The included Makefile provides several useful commands:

- `make build` - Build the application
- `make clean` - Clean build artifacts
- `make test` - Run all tests
- `make run` - Build and run the application
- `make fmt` - Format the code
- `make all` - Clean, format, build, and test

## 🔒 Security Features

### Key Management
- Automatic generation of 256-bit encryption keys
- Secure key storage with 0600 permissions
- Unique nonce generation for each encryption operation

### Integrity Protection
- HMAC-SHA256 for tamper detection
- Automatic verification during decryption
- Immediate alert on detected tampering

### Audit Logging
- Detailed operation logging
- Timestamp and status tracking
- Success/failure recording

## ⚠️ Important Security Notes

1. **Backups**: Always maintain backups of original files before encryption
2. **Key Security**: Store encryption keys in a secure location
3. **File Extensions**: Encrypted files use the `.enc` extension
4. **Permissions**: Ensure proper file permissions on sensitive data

## 🐛 Troubleshooting

Common issues and solutions:

1. **Permission Denied**
   - Check file/directory permissions
   - Ensure you have write access to the output location

2. **Decryption Failed**
   - Verify you're using the correct encryption key
   - Check if the file is corrupted or tampered

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 🔐 Security Reports

For security issues, please email [amulthantharate@gmail.com] instead of using the issue tracker.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
