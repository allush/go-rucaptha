package rucapcha

import (
	"net/http"
	"bytes"
	"mime/multipart"
	"fmt"
	"io"
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
	if _, err = fw.Write([]byte(solver.ApiKey)); err != nil {
		return nil, err
	}

	for field, value := range data {
		result, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("option \"" + field + "\" should be string\n")
		}
		if fw, err = writer.CreateFormField(field); err != nil {
			return nil, err
		}
		if _, err = fw.Write([]byte(result)); err != nil {
			return nil, err
		}
	}

	writer.Close()

	request, err := http.NewRequest("POST", solver.RequestUrl, &buffer)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	return request, nil
}

func (solver *CaptchaSolver) sendRequest(file []byte) (io.ReadCloser, error) {
	data := map[string]interface{}{
		"phrase": solver.IsPhrase,
		"regsense": solver.IsRegsence,
		"numeric": solver.IsNumeric,
		"min_len": solver.MinLength,
		"max_len": solver.MaxLength,
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
