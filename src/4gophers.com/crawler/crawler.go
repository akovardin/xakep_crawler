package crawler

import (
    "4gophers.com/tor"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "io"
    // "log"
    "net/http"
    "os"
    "regexp"
    "time"
    "sync"
    "log"
)

var (
    re        *regexp.Regexp
    clientPtr *http.Client
    folder    string
)

func init() {
    folder = "./downloads/"
}

// Подготавливаем регулярку
func regex() *regexp.Regexp {
    if re == nil {
        re = regexp.MustCompile("http://[/A-Za-z0-9.-]*([XAxa][XAxa][_A-Za-z0-9-]*.pdf)[.pdf]*[?][a-z0-9]*")
    }
    return re
}

// Подготавливаем http клиент для выполнения
// запросов в сайту
func client() *http.Client {
    if clientPtr == nil {
        clientPtr = tor.PrepareProxyClient()
    }
    return clientPtr
}

// Собираем контент с подготовленных ссылок
func Crawler(links []string) {

    var wg sync.WaitGroup
    wg.Add(len(links))

    for _, link := range links { // ждем, пока не получим линки
        go worker(link, &wg)
    }

    wg.Wait()
}

// PassThru врапер над io.Reader.
//
// Эта структура нужна нам для отображения процесса
// скачивания файлов
type PassThru struct {
    io.Reader
    total      int64 // Total # of bytes transferred
    downloaded int64
}

// Read 'переопределенный' метод Read из io.Reader method.
// Именно этот метод вызывается, когда мы используем io.Copy().
// Это можно использовать для отображения процесса копирования файла
func (pt *PassThru) Read(p []byte) (int, error) {
    n, err := pt.Reader.Read(p)
    pt.downloaded += int64(n)

    if err == nil {
        fmt.Println("Скачано", pt.downloaded, "байтов из ", pt.total)
    }

    return n, err
}

// Функция в которой выполняется вся работа
// и которая запускается в рамках go-рутины
func worker(link string, wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Println("Start: " + link)

    response, err := tor.HttpGet(client(), link)
    if err != nil {
        log.Print(err)
        return
    }

    doc, err := goquery.NewDocumentFromResponse(response)
    if err != nil {
        log.Print(err)
        return
    }

    // Получаем все ссылки на файлы со страницы
    hrefs := doc.Find("a.download-button").Map(
        func(i int, s *goquery.Selection) string {
            fmt.Println("Начало обработки: ", link)
            href, _ := s.Attr("href")
            fmt.Println(href)
            return href
        })

    // Скачиваем все файлы по ссылкам
    for _, href := range hrefs {
        name := regex().ReplaceAllString(href, "${1}")
        fmt.Println("Файл: " + name)

        filename := folder + name

        // Если файл уже существует - переходим на следующую итерацию
        if _, err := os.Stat(filename); err == nil {
            continue
        }

        output, err := os.Create(filename)
        if err != nil {
            fmt.Println("Проблема при создании файла:", filename, "-", err)
            continue
        }
        defer output.Close()

        time.Sleep(3 * time.Second)

        // Запрашиваем файл через tor прокси
        response, err := tor.HttpGet(client(), href)
        if err != nil {
            fmt.Println("Проблема со скачиваним файла: ", href, "-", err)
            os.Remove(filename)
            continue
        }
        defer response.Body.Close()

        fmt.Println("размер:", response.ContentLength)
        src := &PassThru{
            Reader:     response.Body,
            total:      response.ContentLength,
            downloaded: 0,
        }
        n, err := io.Copy(output, src)

        if err != nil {
            fmt.Println("Проблема со скачиваним файла: ", href, "-", err)
            os.Remove(filename)
            continue
        }
        fmt.Println(n, "байтов скачено. файл: "+name)
    }

    fmt.Println("Finished: " + link)
    wg.Done()
    return
}
