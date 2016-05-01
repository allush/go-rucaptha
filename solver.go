package rucapcha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type CaptchaSolver struct {
	IsPhrase   bool
	IsRegsence bool
	IsNumeric  int
	MinLength  int
	MaxLength  int
	Language   int

	ImagePath string
	ApiKey    string

	RequestUrl string
	ResultUrl  string
	CheckResultTimeout    int
}

func New(key string) *CaptchaSolver {
	return &CaptchaSolver{
		RequestUrl: "http://rucaptcha.com/in.php",
		ResultUrl:  "http://rucaptcha.com/res.php",
		ApiKey:     key,
	}
}

func (solver *CaptchaSolver) Solve(path string) (*string, error) {
	solver.ImagePath = path

	file, err := solver.getFile()
	if err != nil {
		return nil, err
	}

	response, err := solver.sendRequest(*file)
	if err != nil {
		return nil, err
	}

	captchaId, err := solver.getCaptchaId(response)
	if err != nil {
		return nil, err
	}

	answer, err := solver.waitForReady(captchaId)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (solver *CaptchaSolver) getCaptchaId(response []byte) (string, error){
		body, err := ioutil.ReadAll(response)
		if err != nil {
			return "", err
		}

		hasError := regexp.
			MustCompile(`ERROR`).
			MatchString(string(body))

		if hasError {
			return "", fmt.Errorf("Captcha service error: %s\n", string(body))
		}

		isOk := regexp.
			MustCompile(`OK`).
			MatchString(string(body))

		if !isOk {
			return "", fmt.Errorf("Unknown response: %s\n", string(body))
		}

		results := strings.Split(string(body), "|")

		return results[1], nil
}

func (solver *CaptchaSolver) waitForReady(captchaId string) (*string, error) {

	data := url.Values{}
	data.Add("key", solver.ApiKey)
	data.Add("action", "get")
	data.Add("id", captchaId)

	url := solver.ResultUrl + "?"

	var response *http.Response

	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	var answer *string
	for {
		time.Sleep(solver.CheckResultTimeout * time.Second)

		response, err := http.Get(url + data.Encode())
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		isOk := regexp.
			MustCompile(`OK`).
			MatchString(string(body))

		if isOk {
			results := strings.Split(string(body), "|")
			answer = &results[1]
			break
		}

		notReady := regexp.
			MustCompile(`CAPCHA_NOT_READY`).
			MatchString(string(body))

		if notReady {
			continue
		}

		hasError := regexp.
			MustCompile(`ERROR`).
			MatchString(string(body))

		if hasError {
			return nil, fmt.Errorf("Error response: %s", string(body))
		}
	}

	return answer, nil
}

func (solver *CaptchaSolver) getFile() (*[]byte, error) {
	response, err := http.Get(solver.ImagePath)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
