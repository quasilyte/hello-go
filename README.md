# Hello, Go!

## Способы изучать Go

* Если нравится формат FAQ, читайте дальше и пропустите эту секцию.
* Если любите решать задачки на [leetcode](https://leetcode.com/problemset/all/?difficulty=Easy), [Codewars](https://www.codewars.com/), [CodinGame](https://www.codingame.com/start) или [HackerRank](https://www.hackerrank.com/), то можете попробовать решать их на Go (они поддерживают решения на этом языке).
* Если любите смотреть уже готовые примеры решений и сравнивать их с решениями
  на знакомых вам языках программирования, загляните в [rosettacode](http://www.rosettacode.org/wiki/Category:Go).
  Аналогичным ресурсом является [Go by example](https://gobyexample.com/).
* Для rosettacode есть сайт с side-by-side сравнением решений на двух языках.
  Вот, например, страница для [Go<->Python](https://rosetta.alhur.es/compare/Go/Python/).
* Для уверенных в себе есть [learnxinyminutes](https://learnxinyminutes.com/docs/ru-ru/go-ru/). Качество подачи не слишком высокое, но это один из самых быстрых способов изучить самые базовые возможности языка.
* Если есть опыт с другими языками программирования, можно попробовать [go tour](https://tour.golang.org/welcome/1).

## Как установить Go?

Скачать нужный дистрибутив по ссылке: https://golang.org/dl/.
Есть версии под Windows, Linux и macOS.

Качать последнюю стабильную версию (1.13.3).

## Какой текстовой редактор использовать для Go?

- Visual studio code + [плагин для Go](https://marketplace.visualstudio.com/items?itemName=golang.go)
- Если знакомы продукты JetBrains Intellij, то Goland IDE
- Если всё выше звучит непонятно, то используйте https://play.golang.org/

## Как проверить, что Go установлен?

Набрать в терминале `go version`.
Если не работает, то нужно добавить путь к папке с
исполняемым файлом `go` в переменную окружения `PATH`.

## Как запускать программы на Go?

Создайте файл `hello.go` следующего содержания:

```go
package main

import "fmt"

func main() {
  fmt.Println("Hello, World!")
}
```

Для запуска нужно сначала скомпилировать программу, а затем её запустить:

```bash
$ go build -o hello.exe hello.go
$ ./hello.exe
Hello, World!
```

Но для таких простых случаев есть команда `run`, выполняющая эти два шага за вас:

```bash
$ go run hello.go
Hello, World!
```

## Что такое GOPATH? Как узнать его значение?

`GOPATH` указывает на директорию, куда будут устанавливаться пакеты
и в которой будут искаться импортируемые (подключаемые) пакеты.

Если переменной окружения `GOPATH` нет, Go всё равно будет
иметь некоторое значение по умолчанию. Узнать текущее
значение `GOPATH` проще всего командой:

```bash
go env GOPATH
```

На системах типа Linux директория по умолчанию `~/go`.

> Внимание: после Go 1.13 важность GOPATH понизилась. Теперь нужно работать с [модулями](https://github.com/golang/go/wiki/Modules).

## Что такое пакет?

Пакет - это набор файлов, который образует логическую группу.
Можно называть пакет словом "библиотека" (хотя библиотека может состоять из
нескольких пакетов).

## Что такое модуль?

Модуль может содержать в себе один или более пакетов, он же ассоциирует с ними версию.

Например, пакет `foo` может быть частью модуля `github.com/someuser/foo` с версией `v0.5.0`.

Для версий используется подход семантического версионирования.

## Где найти документацию по стандартной библиотеке Go?

Документация по пакетам: https://golang.org/pkg/.

Можно установить godoc и смотреть документацию оффлайн:

```bash
go get -v golang.org/x/tools/cmd/godoc
```

Теперь можно запустить godoc:

```bash
godoc -http=:8080
```

Если открыть в браузере адрес <http://localhost:8080/pkg/>, то вы
увидите ту же документацию, что была доступна онлайн.

Сайт [godoc.org](https://godoc.org/) можно использовать для поиска Go пакетов и/или их документации.

> Внимание: вместо `godoc.org` теперь стоит использовать [pkg.go.dev](https://pkg.go.dev/).

## Какие ещё есть полезные ресурсы?

Большинство ссылок легко найти в гугле по запросу "golang learning resources".<br>
Самое главное правило - всегда искать по слову `golang`, а не `go`.<br>

Ниже наиболее стоящие результаты с описаниями:

* [С чего начать новичку в Go](http://dev.tulu.la/post/go-newbies/)
* [Golang book (перевод на русский)](http://golang-book.ru)
* [Resources for new Go programmers](https://dave.cheney.net/resources-for-new-go-programmers) - статья [Dave Cheney](https://dave.cheney.net/about), одного из ведущих разработчиков Go.
* [golang/go/wiki/Learn](https://github.com/golang/go/wiki/Learn) - много учебного материала.
* [Go videos](https://github.com/hH39797J/golang-videos-ru) - собрание видеозаписей докладов про Go.
* [Go webdev examples](https://gowebexamples.com/) - аналог Go by example, но с уклоном в веб разработку.

## Книги по Go

Многие книги имеют переводы на русский язык.

* [Get programming with Go](https://www.manning.com/books/get-programming-with-go) - хорошая книга если Go один из первых ваших языков программирования.
* [The Go Programming Language](http://www.gopl.io/) - очень известная книга, довольно хороша, но подойдёт только тем, кто уже более-менее комфортно программирует на одном или более языках программирования.
* [Go in practice](https://www.manning.com/books/go-in-practice) - книга, которая может дополнить книги, перечисленные выше.

## Что такое "сообщество Go"?

[GolangShow](http://golangshow.com) - русскоязычный подкаст о Go. Крутые ведущие, интересные гости.

Сообщество стоит понимать как "группа людей со схожими интересами и/или целями".

[golang-ru Slack](http://slack.golang-ru.com) - русскоязычное Go сообщество.
Там можно задавать вопросы, обсуждать Go, библиотеки под него и прочее.

Для вопросов лучше всего подходит канал `#school` (при формулировке вопроса можно
опираться на [How To Ask Questions The Smart Way](http://www.catb.org/esr/faqs/smart-questions.html)).

Всем участникам следует соблюдать [кодекс норм поведения](https://golang.org/conduct).

Для Казани есть группа [GolangKazan](https://vk.com/golangkazan).
