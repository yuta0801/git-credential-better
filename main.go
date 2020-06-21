package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
)

const Service = "git"

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	credential := read()

	switch os.Args[1] {
	case "get":
		get(&credential)
	case "store":
		store(&credential)
	case "erase":
		erase(&credential)
	}
}

type Credential struct {
	username string
	password string
	protocol string
	host     string
	path     string
	url      string
	// quit     bool
}

func get(credential *Credential) {
	user := makeURL(credential)

	password, err := keyring.Get(Service, user)
	if err != nil {
		panic(err)
	}

	fmt.Println(password)
}

func store(credential *Credential) {
	if credential.protocol == "" ||
		(credential.host == "" && credential.path == "") ||
		credential.username == "" ||
		credential.password == "" {
		return
	}

	user := makeURL(credential)
	fmt.Println(user)

	err := keyring.Set(Service, user, credential.password)
	if err != nil {
		panic(err)
	}
}

func erase(credential *Credential) {
	if credential.protocol == "" ||
		(credential.host == "" && credential.path == "") ||
		credential.username == "" {
		return
	}

	user := makeURL(credential)

	err := keyring.Delete(Service, user)
	if err != nil {
		panic(err)
	}
}

func read() Credential {
	scanner := bufio.NewScanner(os.Stdin)
	var credential Credential

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		if !strings.Contains(line, "=") {
			fmt.Println("warning: invalid credential line:", line)
			fmt.Println("fatal: unable to read credential")
			os.Exit(1)
		}

		slice := strings.Split(line, "=")
		key := slice[0]
		value := slice[1]

		switch key {
		case "username":
			credential.username = value
		case "password":
			credential.password = value
		case "protocol":
			credential.protocol = value
		case "host":
			credential.host = value
		case "path":
			credential.path = value
		case "url":
			credential.url = value
		// case "quit":
		// 	credential.quit = value == "1" ||
		// 		value == "true"
		}
	}

	return credential
}

func makeURL(credential *Credential) string {
	url :=
		credential.protocol + "://" +
			credential.username + "@"

	if credential.host != "" {
		url += credential.host
	}
	if credential.path != "" {
		url += "/" + credential.path
	}

	return url
}
