package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	PrivateKey string `json:"privateKey"`
	OtpFolder  string `json:"otpFolder"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Config) create() {
	fmt.Println("Creating a new config file...")
	reader := bufio.NewReader((os.Stdin))

	// Ask for the name
	fmt.Print("Name: ")
	c.Name, _ = reader.ReadString('\n')
	c.Name = strings.TrimSuffix(c.Name, "\n")

	// Ask for the email
	fmt.Print("Email: ")
	c.Email, _ = reader.ReadString('\n')
	c.Email = strings.TrimSuffix(c.Email, "\n")

	// Ask for the private key
	fmt.Print("Private Key: ")
	c.PrivateKey, _ = reader.ReadString('\n')
	c.PrivateKey = strings.TrimSuffix(c.PrivateKey, "\n")

	// Ask for the otp folder
	fmt.Print("OTP Folder: ")
	c.OtpFolder, _ = reader.ReadString('\n')
	c.OtpFolder = strings.TrimSuffix(c.OtpFolder, "\n")

	// Save struct to the config file as JSON
	file, _ := json.MarshalIndent(c, "", "    ")
	_ = os.WriteFile("config.json", file, 0644)

	fmt.Println("Config file created.")
}

// Get the public key of the email address
// from keys.openpgp.org/
func getPublicKey(email string) {
	// URLENCODE email
	escapedEmail := url.QueryEscape(email)
	resp, err := http.Get("https://keys.openpgp.org/vks/v1/by-email/" + escapedEmail)
	check(err)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
}

func main() {
	fmt.Println("One-Time Pad")
	fmt.Println("------------")
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	configCreate := configCmd.Bool("create", false, "Create a new config file")

	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	encryptTo := encryptCmd.String("to", "", "Encrypt the message to this email address")
	encryptFile := encryptCmd.String("file", "", "Encrypt the message from this file")

	if len(os.Args) < 2 {
		fmt.Println("Please specify a command.")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "config":
		configCmd.Parse(os.Args[2:])
		if *configCreate {
			var c Config
			c.create()
			return
		}

	case "encrypt":
		encryptCmd.Parse(os.Args[2:])
		if *encryptTo == "" || *encryptFile == "" {
			fmt.Println("Please specify the email address and the file to encrypt.")
			os.Exit(0)
		}

		// Get the public key of the email address
		// from keys.openpgp.org/
		getPublicKey(*encryptTo)
	}
}

// Usage: otp [OPTIONS]
//   -e, --encrypt			Encrypt the message
//   -d, --decrypt			Decrypt the message
//   -k, --keyfile string	Path to the keyfile

// Examples:
// otp --keyfolder ./keys --inputfile ./message --email kaangiray@buzl.uk

// otp config --create
