package burgoking

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	baseURL   = "https://www.bkvousecoute.fr"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"

	surveyEntryId = "surveyEntryForm"
	surveyId      = "surveyForm"

	requiredRequests = 20
)

var (
	ErrInvalidAPIResponse = errors.New("invalid response from the website")
	ErrFormNotFound       = errors.New("cannot found form from the response")
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

	var (
		resp              *http.Response
		doc               *goquery.Document
		nextAction, index string
	)

	requestCount := 0
	for {
		if requestCount > requiredRequests {
			err = ErrInvalidAPIResponse
			break
		}

		resp, err = c.client.Do(req)
		if err != nil {
			break
		}

		if resp.StatusCode != http.StatusOK {
			err = ErrInvalidAPIResponse
			break
		}

		doc, err = buildDocument(resp)
		if err != nil {
			break
		}

		requestCount++
		body := doc.Find("body")

		if body.HasClass("Finish") {
			code = strings.Split(doc.Find(".ValCode").Text(), " : ")[1]
			break
		}

		nextAction, err = parseAction(doc)
		if err != nil {
			break
		}

		if body.HasClass("CookieSplashPage") {
			req, err = buildEntryRequest(nextAction)
			if err != nil {
				break
			}
		} else if body.HasClass("CouponEntryPage") {
			req, err = buildEntryRequest(nextAction)
			if err != nil {
				break
			}
		} else {
			index, err = parseIndex(doc)
			if err != nil {
				break
			}

			req, err = buildSurveyRequest(nextAction, index)
			if err != nil {
				break
			}
		}
	}

	return
}

func addUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", userAgent)
}

func buildDocument(resp *http.Response) (doc *goquery.Document, err error) {
	doc, err = goquery.NewDocumentFromReader(resp.Body)

	err = resp.Body.Close()
	if err != nil {
		return
	}

	return
}

func parseAction(doc *goquery.Document) (nextAction string, err error) {
	nextAction, exists := doc.Find("#" + surveyId).Attr("action")
	if exists {
		return
	}

	nextAction, exists = doc.Find("#" + surveyEntryId).Attr("action")
	if exists {
		return
	}

	err = ErrFormNotFound
	return
}

func parseIndex(doc *goquery.Document) (index string, err error) {
	index, exists := doc.Find("#IoNF").Attr("value")
	if !exists {
		err = ErrInvalidAPIResponse
	}

	return
}

func buildUrl(action string) string {
	return fmt.Sprintf("%s/%s", baseURL, action)
}

func buildFirstRequest() (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return
	}
	addUserAgent(req)
	return
}

func commonParams() *url.Values {
	data := url.Values{}
	data.Set("JavaScriptEnabled", "1")
	data.Set("AcceptCookies", "Y")
	return &data
}

func buildCommonRequest(action string, data *url.Values) (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodPost, buildUrl(action), strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	addUserAgent(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return
}

func buildEntryRequest(action string) (req *http.Request, err error) {
	data := commonParams()
	data.Set("FIP", "True")

	return buildCommonRequest(action, data)
}

func buildSurveyRequest(action string, questionIndex string) (req *http.Request, err error) {
	data := commonParams()
	data.Set("IoNF", questionIndex)

	return buildCommonRequest(action, data)
}
