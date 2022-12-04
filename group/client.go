package group

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Getter interface {
	Get(name string, key string) ([]byte, error)
}

type Picker interface {
	Pick(key string) (Getter, error)
}

type HTTPClient struct {
	url string
}

func (hc *HTTPClient) Get(name string, key string) ([]byte, error) {
	Url := fmt.Sprintf("%s%s/%s", hc.url, url.QueryEscape(name), url.QueryEscape(key))
	log.Println(Url)
	res, err := http.Get(Url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("错误的返回码")
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("获取数据错误！")
	}
	return bytes, nil
}

//var _ Getter = (*HTTPClient)(nil)
