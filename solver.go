package rucaptcha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// CaptchaSolver structure
type CaptchaSolver struct {
	IsPhrase   bool
	IsRegsence bool
	IsNumeric  int
	MinLength  int
	MaxLength  int
	Language   int

	ImagePath string
	APIKey    string

	RequestURL         string
	ResultURL          string
	CheckResultTimeout time.Duration
}

// New creates instance of solver
func New(key string) *CaptchaSolver {
	return &CaptchaSolver{
		RequestURL: "http://rucaptcha.com/in.php",
		ResultURL:  "http://rucaptcha.com/res.php",
		APIKey:     key,
	}
}

// Solve get image by path and redn request to rucaptcha service
// Returns captcha code or nil if errors occured
func (solver *CaptchaSolver) Solve(path string) (*string, error) {
	solver.ImagePath = path

	file, err := solver.loadCaptchaImage()
	if err != nil {
		return nil, err
	}

	captchaID, err := solver.getCaptchaID(*file)
	if err != nil {
		return nil, err
	}

	answer, err := solver.waitForReady(*captchaID)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (solver *CaptchaSolver) getCaptchaID(file []byte) (*string, error) {

	response, err := solver.sendRequest(file)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}

	hasError := regexp.
		MustCompile(`ERROR`).
		MatchString(string(body))

	if hasError {
		return nil, fmt.Errorf("Captcha service error: %s\n", string(body))
	}

	isOk := regexp.
		MustCompile(`OK`).
		MatchString(string(body))

	if !isOk {
		return nil, fmt.Errorf("Unknown response: %s\n", string(body))
	}

	results := strings.Split(string(body), "|")

	return &results[1], nil
}

func (solver *CaptchaSolver) waitForReady(captchaID string) (*string, error) {

	data := url.Values{}
	data.Add("key", solver.APIKey)
	data.Add("action", "get")
	data.Add("id", captchaID)

	url := solver.ResultURL + "?"

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

func (solver *CaptchaSolver) loadCaptchaImage() (*[]byte, error) {
	isHTTP := regexp.
		MustCompile(`(http://|https://)`).
		MatchString(solver.ImagePath)

	if !isHTTP {
		body, err := ioutil.ReadFile(solver.ImagePath)
		if err != nil {
			return nil, err
		}
		return &body, nil
	}

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
