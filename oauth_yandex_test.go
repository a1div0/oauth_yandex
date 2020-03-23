package oauth_yandex

import (
    "testing"
)

func TestServiceName(t *testing.T) {
    s := OAuthYandex {}
    name := s.ServiceName()
    if (name != "yandex") {
        t.Error("ServiceName = '", name, "'\nа должно быть = 'yandex'")
    }
}

func TestGetLoginURL(t *testing.T) {

    s := OAuthYandex {
        ClientId: "123",
        ClientPsw: "456",
    }
    url := s.LoginURL("a","b")
    target_url := "https://oauth.yandex.ru/authorize?client_id=123&response_type=code&state=b";
    if (url != target_url) {
        t.Error("URL = ", url, "\nа должен быть = ", target_url)
    }
}
