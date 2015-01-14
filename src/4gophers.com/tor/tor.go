package tor

import (
    "github.com/hailiang/socks"
    "io/ioutil"
    "log"
    "net/http"
)

func PrepareProxyClient() *http.Client {
    dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:9050")
    transport := &http.Transport{Dial: dialSocksProxy}
    return &http.Client{Transport: transport}
}

func HttpGet(httpClient *http.Client, url string) (response *http.Response, err error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    // req.Header.Set("User-Agent", "curl/7.21.4 (universal-apple-darwin11.0) libcurl/7.21.4 OpenSSL/0.9.8x zlib/1.2.5")
    req.Header.Set("User-Agent", `Mozilla/5.0 (X11; Linux i686)
        AppleWebKit/537.36 (KHTML, like Gecko)
        Chrome/36.0.1985.125 Safari/537.36`)
    response, err = httpClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    return
}

func HttpGetBody(httpClient *http.Client, url string) (body string, err error) {
    resp, err := HttpGet(httpClient, url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    bodyb, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    body = string(bodyb)
    return
}
