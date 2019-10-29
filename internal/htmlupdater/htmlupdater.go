package htmlupdater

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

var captchas = [...]string{
	"Type only the four missing letters from the following words: The USA Presidential Nominees are Hillary Clinton and Donald Trump",
	"Type only the four missing letters from the following words: Happy Birthday.",
}

type HtmlUpdater struct {
	rootURL   string
	csrfToken string
	client    *http.Client
}

type QueryResponseWrapper struct {
	Query QueryResponse `json:"query"`
}

type QueryResponse struct {
	Export QueryExportResponse `json:"export"`
}

type QueryExportResponse struct {
	Body string `json:"*"`
}

type EditResponseWrapper struct {
	Edit EditResponse `json:"edit"`
}

type EditResponse struct {
	Result       string          `json:"result"`
	PageID       int             `json:"pageid"`
	Title        string          `json:"title"`
	ContentModel string          `json:"contentmodel"`
	OldRevID     int             `json:"oldrevid"`
	NewRevID     int             `json:"newrevid"`
	NewTimestamp string          `json:"newtimestamp"`
	Captcha      CaptchaResponse `json:"captcha,omitempty"`
}

type CaptchaResponse struct {
	Type     string `json:"type"`
	Mime     string `json:"mime"`
	ID       string `json:"id"`
	Question string `json:"question"`
}

type TokenQueryResponseWrapper struct {
	Query TokenQueryResponse `json:"query"`
}

type TokenQueryResponse struct {
	Tokens TokenResponse `json:"tokens"`
}

type TokenResponse struct {
	CSRF string `json:"csrftoken"`
}

type LoginResponseWrapper struct {
	Login LoginResponse `json:"login"`
}

type LoginResponse struct {
	LGUserID     int    `json:"lguserid,omitempty"`
	LGUserName   string `json:"lgusername,omitempty"`
	Result       string `json:"result,omitempty"`
	Token        string `json:"token,omitempty"`
	CookiePrefix string `json:"cookieprefix,omitempty"`
	SessionID    string `json:"sessionid,omitempty"`
}

func NewHtmlUpdater(rootURL string) HtmlUpdater {
	cookieJar, _ := cookiejar.New(nil)
	return HtmlUpdater{rootURL: rootURL, client: &http.Client{
		Jar: cookieJar,
	}}
}

func AddKeys(title string, username string, password string, keys []string, timestamp string) {
	h := NewHtmlUpdater("http://orcz.com/api.php")
	fullHTML, err := h.getFullHTML(title)
	if err != nil {
		log.Println(err)
		return
	}
	newKeys := findNewKeys(keys, fullHTML)
	if len(newKeys) == 0 {
		log.Println("No new keys found")
		return
	}

	err = h.login(username, password)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.updatePost(keys, timestamp, title, "2", "", "")
	if err != nil {
		log.Println(err)
	}
}

func (h *HtmlUpdater) login(username string, password string) error {
	err := h.updateCSRFToken()
	if err != nil {
		return err
	}

	base, err := url.Parse(h.rootURL)
	if err != nil {
		log.Printf("Failed to create baseURL with error %v\n", err)
		return err
	}

	form := url.Values{}
	form.Add("action", "login")
	form.Add("format", "json")
	form.Add("lgname", "username")
	form.Add("lgpassword", "password")
	form.Add("lgtoken", h.csrfToken)

	resp, err := h.client.PostForm(base.String(), form)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResponse LoginResponseWrapper
	err = json.Unmarshal([]byte(body), &loginResponse)
	if err != nil {
		return err
	}

	return nil
}

func (h *HtmlUpdater) getFullHTML(title string) (string, error) {
	base, err := url.Parse(h.rootURL)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Add("format", "json")
	q.Add("action", "query")
	q.Add("titles", title)
	q.Add("export", "true")

	base.RawQuery = q.Encode()
	resp, err := h.client.Get(base.String())
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var queryResponse QueryResponseWrapper
	err = json.Unmarshal([]byte(bytes), &queryResponse)
	if err != nil {
		return "", err
	}

	return queryResponse.Query.Export.Body, nil
}

func (h *HtmlUpdater) updatePost(keys []string, timestamp string, title string, section string, captchaID string, captchaWord string) error {
	var append string
	for _, key := range keys {
		append = append + buildKeyHTML(key, timestamp)
	}
	fullHTML, err := h.getFullHTML(title)
	if err != nil {
		return err
	}
	shiftCodeIndex := strings.Index(fullHTML, "===Shift Codes===")
	if shiftCodeIndex == -1 {
		return errors.New("Failed to find shift code index")
	}

	lastRowIndex := strings.LastIndex(fullHTML, "\n|-\nLegend ... ")
	if lastRowIndex == -1 {
		return errors.New("Unable to find last row, not continuing")
	}
	lastRowIndex += 4

	bookmarkIndex := strings.LastIndex(fullHTML, "==Bookmark Here==")
	if bookmarkIndex == -1 {
		return errors.New("Failed to find bookmark index")
	}
	bookmarkIndex -= 2

	finalText := fullHTML[shiftCodeIndex:lastRowIndex] + append + fullHTML[lastRowIndex:bookmarkIndex]
	decoded := html.UnescapeString(finalText)

	base, err := url.Parse(h.rootURL)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Add("action", "edit")
	form.Add("title", title)
	form.Add("section", section)
	form.Add("text", decoded)
	form.Add("format", "json")
	form.Add("bot", "true")
	form.Add("token", h.csrfToken)
	if len(captchaID) > 0 && len(captchaWord) > 0 {
		form.Add("captchaid", captchaID)
		form.Add("captchaword", captchaWord)
	}

	resp, err := h.client.PostForm(base.String(), form)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll((resp.Body))
	if err != nil {
		return err
	}

	var editResponse EditResponseWrapper
	err = json.Unmarshal([]byte(bytes), &editResponse)
	if err != nil {
		return err
	}
	if strings.ToLower(editResponse.Edit.Result) != "success" {
		log.Println(string(bytes))
		log.Printf("%v\n", editResponse)
		if len(captchaID) == 0 && len(captchaWord) == 0 {
			answer := solveCaptcha(editResponse.Edit.Captcha.Question)
			return h.updatePost(keys, timestamp, title, section, editResponse.Edit.Captcha.ID, answer)
		}
		return errors.New(editResponse.Edit.Result)
	}

	return nil
}

func (h *HtmlUpdater) updateCSRFToken() error {
	base, err := url.Parse(h.rootURL)
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Add("action", "query")
	q.Add("meta", "tokens")
	q.Add("format", "json")

	base.RawQuery = q.Encode()
	resp, err := h.client.Get(base.String())
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tokenResponse TokenQueryResponseWrapper
	err = json.Unmarshal([]byte(bytes), &tokenResponse)
	if err != nil {
		return err
	}
	h.csrfToken = tokenResponse.Query.Tokens.CSRF
	return nil
}

func findNewKeys(keys []string, html string) []string {
	results := []string{}
	re := regexp.MustCompile(`(\w{5}\-){4}\w{5}`)
	matches := re.FindAll([]byte(html), -1)
	for _, new := range keys {
		keyFound := false
		for _, existing := range matches {
			if new == string(existing) {
				keyFound = true
				break
			}
		}
		if keyFound {
			continue
		}
		results = append(results, new)
	}
	return results
}

func buildKeyHTML(key string, timestamp string) string {
	return fmt.Sprintf("||@DuvalMagic || 1 Golden Key || %s || Unknown || %s || ❓ || ❓ || ❓ || ❓ || ❓\n|-\n", timestamp, key)
}

func solveCaptcha(question string) string {
	var captcha string
	for _, val := range captchas {
		if len(question) == len(val) {
			captcha = val
			break
		}
	}
	if len(captcha) == 0 {
		log.Printf("Failed to find matching captcha for question %s\n", question)
		return ""
	}

	index := strings.Index(question, "_")
	answer := captcha[index : index+4]

	return answer
}
