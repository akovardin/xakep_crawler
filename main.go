package main

import (
    "4gophers.com/crawler"
    "strconv"
)

func main() {
    var links []string
    for i := 2; i <= 5; i++ {
        links = append(links, "http://xakep.ru/issues/xa/page/" +
                                    strconv.Itoa(i))
    }
    crawler.Crawler(links)
}
