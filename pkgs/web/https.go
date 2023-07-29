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

func doesFileExist(fileName string) bool {
	_, error := os.Stat(fileName)

	if os.IsNotExist(error) {
		log.Printf("%v file does not exist\n", fileName)
		return false
	} else {
		return true
	}
}
func createCertFiles() error {
	cmd := exec.Command(
		"openssl",
		"req", "-new", "-newkey", "rsa:2048", "-sha256", "-days", "365", "-nodes", "-x509",
		"-keyout", server_key,
		"-out", server_crt,
		"-subj", "/O=Global Security/OU=IT Department")

	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	} else {
		log.Print(out.String())
	}
	return err
}

func init() {
	if !(doesFileExist(server_key) || doesFileExist(server_crt)) {
		log.Printf("try to create certs %v and %v \n", server_crt, server_key)
		err := createCertFiles()
		if err != nil {
			log.Printf("certs %v and %v created successfully\n", server_crt, server_key)
		}

	}

}
