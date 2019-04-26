package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Простейший пример утилиты командной строки, которая использует
// методы storage: get и set.
//
// Пример использования:
//	$ vk-storage -token $TOKEN set mykey 123
//	$ vk-storage -token $TOKEN get mykey

func main() {
	var args arguments

	if err := args.parse(); err != nil {
		log.Panicf("parse args error: %v", err)
	}

	switch args.command {
	case "get":
		storageGet(&args)
	case "set":
		storageSet(&args)
	default:
		log.Panicf("unknown storage method: %q", args.command)
	}
}

type arguments struct {
	token       string
	apiVersion  string
	command     string
	commandArgs []string
}

func (args *arguments) parse() error {
	flag.StringVar(&args.token, "token", "",
		`A token for VK API access_token parameter`)
	flag.StringVar(&args.apiVersion, "api", "5.95",
		`Which VK API version to use`)

	flag.Parse()

	if n := len(flag.Args()); n < 1 {
		return fmt.Errorf("expected at least 1 positional argument, got 0")
	}

	args.command = flag.Args()[0]
	args.commandArgs = flag.Args()[1:]

	return nil
}

func storageGet(args *arguments) {
	resp, err := apiGet(args, "method/storage.get", "key="+args.commandArgs[0])
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(resp["response"])
}

func storageSet(args *arguments) {
	resp, err := apiGet(args,
		"method/storage.set",
		"key="+args.commandArgs[0],
		"value="+args.commandArgs[1])
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(resp["response"])
}

func apiURL(args *arguments, path string, params ...string) *url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "api.vk.com",
		Path:   path,
	}
	query := u.Query()
	query.Set("access_token", args.token)
	query.Set("version", args.apiVersion)
	for _, p := range params {
		kv := strings.Split(p, "=")
		query.Set(kv[0], kv[1])
	}
	u.RawQuery = query.Encode()
	return &u
}

type vkResponse map[string]interface{}

func apiGet(args *arguments, path string, params ...string) (vkResponse, error) {
	targetURL := apiURL(args, path, params...).String()

	rawResp, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}
	defer rawResp.Body.Close()

	data, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp: %v", err)
	}

	var resp vkResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("json decode: %v", err)
	}

	if apiError := resp["error"]; apiError != nil {
		return resp, fmt.Errorf("api error: %v", apiError)
	}

	return resp, nil
}
