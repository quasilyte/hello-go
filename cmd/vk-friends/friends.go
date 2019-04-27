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
)

// Утилита, демонстрирующая использование VK API без каких-либо
// сторонних библиотек на примере friends методом.
//
// Пример использования:
//	$ vk-friends -token $TOKEN online
//	$ vk-friends -token $TOKEN list
//
// Для более простого примера смотри cmd/vk-storage.

func main() {
	var args arguments
	ctxt := newContext(&args)

	// Все шаги программы, в последовательности выполнения.
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
		// Единственное место для обработки ошибок в main функции.
		if err := step.fn(); err != nil {
			log.Printf("%s: %v", step.name, err)
			return
		}
	}
}

// arguments - аргументы командной строки.
type arguments struct {
	// Все поля описаны внутри метода parse.

	token      string
	apiVersion string
	command    string
	verbose    bool
}

// parse связывает аргументы командной строки с объектом args.
// Не производит детальной валидации.
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

// context хранит состояние выполнения программы.
// Для создания экземпляра следует использовать newContext.
type context struct {
	args *arguments

	// commands хранит все зарегистрированные обработчики.
	// Заполняется в newContext.
	commands map[string]func() error

	// requests является счётчиком выполненного количества запросов к API.
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
	// Мы могли бы просто сформировать строку, но url.URL
	// можно использовать как простой билдер для URL'ов.
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

// debugf подобен log.Printf, но печатает только в verbose режиме.
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

	// Проверяем, что подкоманда определена.
	if _, ok := ctxt.commands[ctxt.args.command]; !ok {
		// Если пользователь использует неправильную команду,
		// соберём список доступных комманд и подскажем ему их.
		var hints []string
		for command := range ctxt.commands {
			hints = append(hints, command)
		}
		// Поскольку map в Go не отсортирован, после получения
		// ключей их следует отсортировать, иначе каждый раз
		// подсказки будут печататься в разном порядке.
		sort.Strings(hints)
		return fmt.Errorf("unrecognized command %q (supported commands: %s)",
			ctxt.args.command, strings.Join(hints, ", "))
	}

	return nil
}

func (ctxt *context) execCommand() error {
	// Делегируем выполнение зарегистрированной команде.
	return ctxt.commands[ctxt.args.command]()
}

func (ctxt *context) printStats() error {
	log.Printf("made %d API requests", ctxt.requests)
	return nil
}

// printFriends печатает список друзей.
func (ctxt *context) printFriends(friends []interface{}) error {
	var ids []string
	for _, id := range friends {
		ids = append(ids, fmt.Sprint(int(id.(float64))))
	}

	idsParam := strings.Join(ids, ",")
	resp, err := ctxt.apiGet("method/users.get", "user_ids="+idsParam)
	if err != nil {
		return err
	}
	users := resp["response"].([]interface{})

	for i, user := range users {
		user := user.(map[string]interface{})
		fmt.Printf("\t%4d %s %s\n", i+1, user["first_name"], user["last_name"])
	}

	return nil
}
