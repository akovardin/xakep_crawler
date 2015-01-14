package crawler

import (
    "testing"
)

func TestCrawler(t *testing.T) {
    var links []string
    links = append(links, "http://xakep.ru/issues/xa/page/2")
    Crawler(links)
}
