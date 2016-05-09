package rucaptcha

import "testing"

var solver *CaptchaSolver
var lastCaptchaID *string

func TestNew(t *testing.T) {
	solver = New("rucaptcha api key")
	if solver == nil {
		t.Errorf("Could not create solver instance\n")
	}
}

func TestLoadFile(t *testing.T) {
	solver.ImagePath = "./test/captcha.jpg"
	data, err := solver.loadCaptchaImage()
	if err != nil {
		t.Errorf("Load file from disk error: %s\n", err)
	}

	if data != nil && len(*data) <= 0 {
		t.Errorf("Load file from disk: Data too small\n")
	}

	solver.ImagePath = "https://raw.githubusercontent.com/allush/go-rucaptha/master/test/captcha.jpg"
	data, err = solver.loadCaptchaImage()
	if err != nil {
		t.Errorf("Load file from url error: %s\n", err)
	}

	if data != nil && len(*data) <= 0 {
		t.Errorf("Load file from url: Data too small\n")
	}

}

func TestSolve(t *testing.T) {
	solver.IsRegsence = true
	answer, _, err := solver.Solve(solver.ImagePath)
	if err != nil {
		t.Errorf("Solve error: %s\n", err)
	}

	want := "8AnF"
	if *answer != want {
		t.Errorf("Answer(%s) not equal want(%s)\n", *answer, want)
	}
}

func TestComlpain(t *testing.T) {
	_, lastCaptchaID, err := solver.Solve(solver.ImagePath)
	if err != nil {
		t.Errorf("Solve error: %s\n", err)
	}

	err = solver.Complain(*lastCaptchaID)
	if err != nil {
		t.Errorf("Complain error: %s\n", err)
	}
}
