package burgoking

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

const (
	baseURL = "https://www.bkvousecoute.fr"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
)

var (
	ErrInvalidAPIResponse = errors.New("invalid response from the website")
)

func NewCodeGenerator() *CodeGenerator {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	return &CodeGenerator{&client}
}

type CodeGenerator struct {
	client *http.Client
}

func (c *CodeGenerator) Generate() (code string, err error) {
	req, err := buildFirstRequest()
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrInvalidAPIResponse
		return
	}

	action, err := parseAction(resp)
	if err != nil {
		return
	}

	for i := 0; i < 5; i++ {
		req, err = buildActionRequest(action)
		if err != nil {
			return
		}

		resp, err = c.client.Do(req)
		if err != nil {
			return
		}

		if resp.StatusCode != http.StatusOK {
			err = ErrInvalidAPIResponse
			return
		}

		action, err = parseAction(resp)
		if err != nil {
			return
		}
	}

	return
}

func addUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", userAgent)
}

func parseAction(resp *http.Response) (nextAction string, err error) {
	doc, err := goquery.NewDocumentFromReader(io.TeeReader(resp.Body, os.Stdout))
	if err != nil {
		return
	}

	nextAction, exists := doc.Find("#surveyEntryForm").Attr("action")
	if !exists {
		err = ErrInvalidAPIResponse
		return
	}

	err = resp.Body.Close()
	if err != nil {
		err = ErrInvalidAPIResponse
		return
	}

	return
}

func buildFirstRequest() (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return
	}
	addUserAgent(req)
	return
}

func buildUrl(action string) string {
	return fmt.Sprintf("%s/%s", baseURL, action)
}

func buildActionRequest(action string) (req *http.Request, err error) {
	data := url.Values{}
	data.Set("JavaScriptEnabled", "1")
	data.Set("AcceptCookies", "Y")

	req, err = http.NewRequest(http.MethodPost, buildUrl(action), strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	addUserAgent(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return
}