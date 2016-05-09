package rucapcha

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

func (solver *CaptchaSolver) createRequest(
	file []byte,
	data map[string]interface{},
) (*http.Request, error) {

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	fw, err := writer.CreateFormFile("file", "captcha")
	if err != nil {
		return nil, err
	}
	if _, err = fw.Write(file); err != nil {
		return nil, err
	}

	if fw, err = writer.CreateFormField("key"); err != nil {
		return nil, err
	}
	if _, err = fw.Write([]byte(solver.APIKey)); err != nil {
		return nil, err
	}

	for field, value := range data {
		result, err := solver.parseValue(value)
		if err != nil {
			return nil, err
		}

		if fw, err = writer.CreateFormField(field); err != nil {
			return nil, err
		}

		if _, err = fw.Write([]byte(result)); err != nil {
			return nil, err
		}
	}

	writer.Close()

	request, err := http.NewRequest("POST", solver.RequestURL, &buffer)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	return request, nil
}

func (solver *CaptchaSolver) parseValue(value interface{}) (string, error) {

	resultString, okString := value.(string)
	if okString {
		return resultString, nil
	}

	resultInt, okInt := value.(int)
	if okInt {
		return strconv.Itoa(resultInt), nil
	}
	resultBool, okBool := value.(bool)
	if okBool {
		boolVal := 0
		if resultBool == true {
			boolVal = 1
		}
		return strconv.Itoa(boolVal), nil
	}

	return "", fmt.Errorf("option could not be converted to string\n")
}

func (solver *CaptchaSolver) sendRequest(file []byte) (io.ReadCloser, error) {
	data := map[string]interface{}{
		"phrase":   solver.IsPhrase,
		"regsense": solver.IsRegsence,
		"numeric":  solver.IsNumeric,
		"min_len":  solver.MinLength,
		"max_len":  solver.MaxLength,
		"language": solver.Language,
	}

	request, err := solver.createRequest(file, data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status of captcha request: %s", res.Status)
	}

	return res.Body, nil
}

func (solver *CaptchaSolver) complainRequest(captchaID string) error {

	data := url.Values{}
	data.Add("key", solver.APIKey)
	data.Add("action", "reportbad")
	data.Add("id", captchaID)

	url := solver.ResultURL + "?" + data.Encode()

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	isOk := regexp.
		MustCompile(`OK_REPORT_RECORDED`).
		MatchString(string(body))

	if !isOk {
		return fmt.Errorf("Captcha complain error: %s\n", string(body))
	}

	return nil
}
