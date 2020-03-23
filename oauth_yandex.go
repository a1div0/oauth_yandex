// klenov
// 2019.11.02

// Регистрация веб-приложения
// https://yandex.ru/dev/oauth/doc/dg/tasks/register-client-docpage/

package oauth_yandex

import (
    "fmt"
    "net/http"
    "net/url"
    "encoding/json"
    "time"
    "io/ioutil"
    "github.com/a1div0/oauth"
)

type OAuthYandex struct {
    ClientId string
    ClientPsw string
    token string
    token_dt_start time.Time
    token_dt_end time.Time
    refresh_token string
}

func (s *OAuthYandex) ServiceName() (string) {
    return "yandex"
}

func (s *OAuthYandex) LoginURL(verification_code_callback_url string, state string) (string) {

    data := url.Values{}
    data.Set("response_type", "code")
    data.Set("client_id"    , s.ClientId)
    data.Set("state"        , state)

    return "https://oauth.yandex.ru/authorize?" + data.Encode()
}

func (s *OAuthYandex) OnRecieveVerificationCode(code string, u *oauth.UserData) (error) {

    err := s.code_to_token(code)
    if err != nil {
		return err
	}
    err = s.token_to_userdata(u)
    if err != nil {
		return err
	}
    return nil
}

func (s *OAuthYandex) code_to_token(code string) (error) {

    formData := url.Values{
		"grant_type": {"authorization_code"},
        "code": {code},
        "client_id": {s.ClientId},
        "client_secret": {s.ClientPsw},
	}

    resp, err := http.PostForm("https://oauth.yandex.ru/token", formData)
    defer resp.Body.Close()
	if err != nil {
		return err
    }

    type YaTokenAnswerStruct struct {
        Token_type string
        Access_token string
        Expires_in int64
        Refresh_token string
        Error string
        Error_description string
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return err
    }

    var YaTokenAnswer YaTokenAnswerStruct

    err = json.Unmarshal(body, &YaTokenAnswer)
    if err != nil {
		return err
    }

    if (YaTokenAnswer.Error != "") {
        return fmt.Errorf("Error - %s: %s", YaTokenAnswer.Error, YaTokenAnswer.Error_description)
    }

    s.token = YaTokenAnswer.Access_token
    s.token_dt_start = time.Now()
    //s.token_dt_end = s.token_dt_start.Add(YaTokenAnswer.expires_in)
    s.refresh_token = YaTokenAnswer.Refresh_token

    return nil
}

func (s *OAuthYandex) token_to_userdata(u *oauth.UserData) (error) {

    req, err := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
	if err != nil {
		return err
	}
	// Получаем и устанавливаем тип контента
	req.Header.Set("Authorization", "OAuth " + s.token)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
    defer resp.Body.Close()
	if err != nil {
		return err
	}

    type YaUserAnswerStruct struct {
        First_name string
        Last_name string
        Display_name string
        Real_name string
        Login string
        Default_email string
        Id string
        Client_id string
        Emails []string
        Default_avatar_id string
        Is_avatar_empty bool
        Sex string
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return err
    }

    var YaUserAnswer YaUserAnswerStruct
    err = json.Unmarshal(body, &YaUserAnswer)
    if err != nil {
		return err
    }

    u.ExtId = YaUserAnswer.Id
    u.Name = YaUserAnswer.Real_name
    u.Email = YaUserAnswer.Default_email

    return nil
}
