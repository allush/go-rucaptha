package rucapcha

import "testing"

var solver *CaptchaSolver

func TestNew(t *testing.T) {
	// solver = New("rucaptcha api key")
	solver = New("928c5c92741e747019c932be14da7762")
	if solver == nil {
		t.Errorf("Could not create solver instance\n")
	}
}

func TestLoadFile(t *testing.T) {
	solver.ImagePath = "./test/captcha.jpg"
	data, err := solver.loadCaptchaImage()
	if err != nil {
		t.Errorf("Load file error %s\n", err)
	}

	if data != nil && len(*data) <= 0 {
		t.Errorf("Data too small\n")
	}
}

func TestGetCaptchaID(t *testing.T) {

}

func TestSolve(t *testing.T) {

}
