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

	creds := read()

	if len(os.Args) < 2 {
		panic("usage: git credential [fill|approve|reject]")
	}

	switch os.Args[1] {
	case "get":
		get(&creds)
	case "store":
		store(&creds)
	case "erase":
		erase(&creds)
	}
}

type credential struct {
	username string
	password string
	protocol string
	host     string
	path     string
	url      string
	// quit     bool
}

func get(creds *credential) {
	user := makeURL(creds)

	password, err := keyring.Get(Service, user)
	if err != nil {
		panic(err)
	}

	fmt.Println(password)
}

func store(creds *credential) {
	if creds.protocol == "" ||
		(creds.host == "" && creds.path == "") ||
		creds.username == "" ||
		creds.password == "" {
		return
	}

	user := makeURL(creds)
	fmt.Println(user)

	err := keyring.Set(Service, user, creds.password)
	if err != nil {
		panic(err)
	}
}

func erase(creds *credential) {
	if creds.protocol == "" ||
		(creds.host == "" && creds.path == "") ||
		creds.username == "" {
		return
	}

	user := makeURL(creds)

	err := keyring.Delete(Service, user)
	if err != nil {
		panic(err)
	}
}

func read() credential {
	scanner := bufio.NewScanner(os.Stdin)
	var creds credential

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
			creds.username = value
		case "password":
			creds.password = value
		case "protocol":
			creds.protocol = value
		case "host":
			creds.host = value
		case "path":
			creds.path = value
		case "url":
			creds.url = value
			// case "quit":
			// 	creds.quit = value == "1" ||
			// 		value == "true"
		}
	}

	return creds
}

func makeURL(creds *credential) string {
	url :=
		creds.protocol + "://" +
			creds.username + "@"

	if creds.host != "" {
		url += creds.host
	}
	if creds.path != "" {
		url += "/" + creds.path
	}

	return url
}
