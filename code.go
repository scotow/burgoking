package burgoking

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"crypto/tls"
)

const (
	baseURL = "https://www.bkvousecoute.fr"

	userAgentHeader, userAgent     = "User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
	contentTypeHeader, contentType = "Content-Type", "application/x-www-form-urlencoded"

	startBodyClass  = "CookieSplashPage"
	entryBodyClass  = "CouponEntryPage"
	finishBodyClass = "Finish"

	surveyEntryId = "surveyEntryForm"
	surveyId      = "surveyForm"
	indexId       = "IoNF"
	codeClass     = "ValCode"

	formAttr  = "action"
	indexAttr = "value"

	fipKey, fipValue         = "FIP", "True"
	jsKey, jsValue           = "JavaScriptEnabled", "1"
	cookiesKey, cookiesValue = "AcceptCookies", "Y"

	indexKey      = "IoNF"
	surveyCodeKey = "SurveyCode"
	dayKey        = "InputDay"
	monthKey      = "InputMonth"
	yearKey       = "InputYear"
	hourKey       = "InputHour"
	minuteKey     = "InputMinute"

	requiredRequests = 19
)

var (
	ErrInvalidAPIResponse = errors.New("invalid response from the website")
	ErrFormNotFound       = errors.New("cannot find form from the response")
	ErrInvalidCode        = errors.New("cannot parse code from the page")
	ErrTooManyRequest     = errors.New("too many requests")
	ErrInvalidIndex       = errors.New("cannot find index from value")
)

func GenerateCode(meal *Meal) (code string, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client := http.Client{
		Transport: &http.Transport{
        		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    		},
		Jar: jar,
	}

	if meal == nil {
		meal = RandomMeal()
	}

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
		resp, err = client.Do(req)
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

		if body.HasClass(finishBodyClass) {
			code, err = parseCode(doc)
			break
		}

		nextAction, err = parseAction(doc)
		if err != nil {
			break
		}

		if body.HasClass(startBodyClass) {
			req, err = buildStartRequest(nextAction)
			if err != nil {
				break
			}
		} else if body.HasClass(entryBodyClass) {
			req, err = buildEntryRequest(nextAction, meal)
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

		if requestCount >= requiredRequests {
			err = ErrTooManyRequest
			break
		}
	}

	return
}

func addUserAgent(req *http.Request) {
	req.Header.Set(userAgentHeader, userAgent)
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
	nextAction, exists := doc.Find("#" + surveyId).Attr(formAttr)
	if exists {
		return
	}

	nextAction, exists = doc.Find("#" + surveyEntryId).Attr(formAttr)
	if exists {
		return
	}

	err = ErrFormNotFound
	return
}

func parseIndex(doc *goquery.Document) (index string, err error) {
	index, exists := doc.Find("#" + indexId).Attr(indexAttr)
	if !exists {
		err = ErrInvalidIndex
	}

	return
}

func parseCode(doc *goquery.Document) (code string, err error) {
	parts := strings.Split(doc.Find("."+codeClass).Text(), " : ")
	if len(parts) != 2 {
		err = ErrInvalidCode
		return
	}

	code = parts[1]
	return
}

func buildUrl(action string) string {
	return fmt.Sprintf("%s/%s", baseURL, action)
}

func padTo2(value int) string {
	return fmt.Sprintf("%02d", value)
}

func commonParams() *url.Values {
	data := url.Values{}
	data.Set(jsKey, jsValue)
	data.Set(cookiesKey, cookiesValue)
	return &data
}

func buildFirstRequest() (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return
	}
	addUserAgent(req)
	return
}

func buildCommonRequest(action string, data *url.Values) (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodPost, buildUrl(action), strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	addUserAgent(req)
	req.Header.Set(contentTypeHeader, contentType)

	return
}

func buildStartRequest(action string) (req *http.Request, err error) {
	data := commonParams()
	data.Set(fipKey, fipValue)

	return buildCommonRequest(action, data)
}

func buildEntryRequest(action string, meal *Meal) (req *http.Request, err error) {
	data := commonParams()
	data.Set(fipKey, fipValue)
	data.Set(surveyCodeKey, strconv.Itoa(meal.Restaurant))
	data.Set(dayKey, padTo2(meal.Date.Day()))
	data.Set(monthKey, padTo2(int(meal.Date.Month())))
	data.Set(yearKey, strconv.Itoa(meal.Date.Year())[2:])
	data.Set(hourKey, padTo2(meal.Date.Hour()))
	data.Set(minuteKey, padTo2(meal.Date.Minute()))

	return buildCommonRequest(action, data)
}

func buildSurveyRequest(action string, questionIndex string) (req *http.Request, err error) {
	data := commonParams()
	data.Set(indexKey, questionIndex)

	return buildCommonRequest(action, data)
}
