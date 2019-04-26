package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Утилита, демонстрирующая использование VK API без каких-либо
// сторонних библиотек на примере friends методом.

func main() {
	var args arguments
	ctxt := newContext(&args)

	steps := []struct {
		name string
		fn   func() error
	}{
		{"parse args", args.parse},
		{"validate args", ctxt.validateArgs},
		{"exec command", ctxt.execCommand},
		{"print stats", ctxt.printStats},
	}

	for _, step := range steps {
		ctxt.debugf("start %q step", step.name)
		if err := step.fn(); err != nil {
			log.Printf("%s: %v", step.name, err)
			return
		}
	}
}

type arguments struct {
	token      string
	apiVersion string
	command    string
	verbose    bool
}

func (args *arguments) parse() error {
	flag.StringVar(&args.token, "token", "",
		`A token for VK API access_token parameter`)
	flag.StringVar(&args.apiVersion, "api", "5.95",
		`Which VK API version to use`)
	flag.BoolVar(&args.verbose, "verbose", false,
		`Whether to print debug information`)

	flag.Parse()

	if n := len(flag.Args()); n != 1 {
		return fmt.Errorf("expected exactly 1 positional argument, got %d", n)
	}

	args.command = flag.Args()[0]

	return nil
}

type context struct {
	args     *arguments
	commands map[string]func() error
	requests int
}

func newContext(args *arguments) *context {
	ctxt := &context{args: args}

	ctxt.commands = map[string]func() error{
		"online": ctxt.onlineCommand,
		"list":   ctxt.listCommand,
	}

	return ctxt
}

func (ctxt *context) apiURL(path string, params ...string) *url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "api.vk.com",
		Path:   path,
	}
	query := u.Query()
	query.Set("access_token", ctxt.args.token)
	query.Set("version", ctxt.args.apiVersion)
	for _, p := range params {
		kv := strings.Split(p, "=")
		query.Set(kv[0], kv[1])
	}
	u.RawQuery = query.Encode()
	return &u
}

type vkResponse map[string]interface{}

func (ctxt *context) debugf(format string, args ...interface{}) {
	if ctxt.args.verbose {
		log.Printf("debug: "+format, args...)
	}
}

func (ctxt *context) apiGet(path string, params ...string) (vkResponse, error) {
	ctxt.requests++

	targetURL := ctxt.apiURL(path, params...).String()
	ctxt.debugf("GET %q", targetURL)
	rawResp, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}
	defer rawResp.Body.Close()

	data, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp: %v", err)
	}
	ctxt.debugf("API response: %s", data)

	var resp vkResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("json decode: %v", err)
	}

	if apiError := resp["error"]; apiError != nil {
		return resp, fmt.Errorf("api error: %v", apiError)
	}

	return resp, nil
}

func (ctxt *context) getUserString(id int) (string, error) {
	resp, err := ctxt.apiGet("method/users.get", fmt.Sprintf("user_id=%d", id))
	if err != nil {
		return "", err
	}
	user := resp["response"].([]interface{})[0].(map[string]interface{})
	return fmt.Sprintf("%s %s", user["first_name"], user["last_name"]), nil
}

func (ctxt *context) printFriends(friends []interface{}) error {
	for i, id := range friends {
		id := int(id.(float64))
		name, err := ctxt.getUserString(id)

		// Для пользовательского токена есть лимит на 3 запроса в секунду.
		time.Sleep(time.Second / 2)

		if err != nil {
			return err
		}
		fmt.Printf("\t%4d %s (ID=%d)\n", i+1, name, id)
	}
	return nil
}

func (ctxt *context) listCommand() error {
	resp, err := ctxt.apiGet("method/friends.get")
	if err != nil {
		return err
	}

	friends := resp["response"].([]interface{})
	fmt.Printf("friends (%d):\n", len(friends))
	return ctxt.printFriends(friends)
}

func (ctxt *context) onlineCommand() error {
	resp, err := ctxt.apiGet("method/friends.getOnline")
	if err != nil {
		return err
	}

	friends := resp["response"].([]interface{})
	fmt.Printf("friends online (%d):\n", len(friends))
	return ctxt.printFriends(friends)
}

func (ctxt *context) validateArgs() error {
	checkNonEmpty := []struct {
		name  string
		value string
	}{
		{"-token", ctxt.args.token},
		{"-api", ctxt.args.apiVersion},
	}

	for _, check := range checkNonEmpty {
		if check.value != "" {
			continue
		}
		return fmt.Errorf("%s argument can't be empty", check.name)
	}

	if _, ok := ctxt.commands[ctxt.args.command]; !ok {
		var hints []string
		for command := range ctxt.commands {
			hints = append(hints, command)
		}
		sort.Strings(hints)
		return fmt.Errorf("unrecognized command %q (supported commands: %s)",
			ctxt.args.command, strings.Join(hints, ", "))
	}

	return nil
}

func (ctxt *context) execCommand() error {
	return ctxt.commands[ctxt.args.command]()
}

func (ctxt *context) printStats() error {
	log.Printf("made %d API requests", ctxt.requests)
	return nil
}
