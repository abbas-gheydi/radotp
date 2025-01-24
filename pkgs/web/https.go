package web

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

var (
	server_crt                = "/etc/radotp/server.crt"
	server_key                = "/etc/radotp/server.key"
	HTTPSListenAddr           = "0.0.0.0:8081"
	RedirectToHTTPS           = true
	RedirectToHTTPSPortNumber = "443"
)

// doesFileExist checks if a file exists at the given path.
func doesFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Printf("File %v does not exist.\n", fileName)
		return false
	}
	return true
}

// createCertFiles generates self-signed SSL certificate and private key files
// using OpenSSL. This is useful for development environments.
func createCertFiles() error {
	cmd := exec.Command(
		"openssl",
		"req", "-new", "-newkey", "rsa:2048", "-sha256", "-days", "365", "-nodes", "-x509",
		"-keyout", server_key,
		"-out", server_crt,
		"-subj", "/O=Global Security/OU=IT Department",
	)

	// Capture the command's error output
	var errOutput bytes.Buffer
	cmd.Stderr = &errOutput

	err := cmd.Run()
	if err != nil {
		log.Printf("Error creating certificate files: %v\n", err)
		log.Printf("OpenSSL Output: %s\n", errOutput.String())
		return err
	}

	log.Println("Certificate files created successfully.")
	return nil
}

// init function initializes the application by checking for the existence of
// the SSL certificate and key files. If they do not exist, it creates them.
func init() {
	// Check if both the certificate and key files exist
	if !(doesFileExist(server_key) && doesFileExist(server_crt)) {
		log.Printf("Certificate or key file missing. Attempting to create: %v and %v\n", server_crt, server_key)
		err := createCertFiles()
		if err != nil {
			log.Printf("Failed to create certificate files: %v\n", err)
		} else {
			log.Printf("Certificate files created successfully: %v and %v\n", server_crt, server_key)
		}
	} else {
		log.Println("Certificate and key files already exist.")
	}
}
