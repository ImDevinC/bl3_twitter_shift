package htmlupdater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

type HtmlUpdater struct {
	rootURL string
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
	Result       string `json:"result"`
	PageID       int    `json:"pageid"`
	Title        string `json:"title"`
	ContentModel string `json:"contentmodel"`
	OldRevID     int    `json:"oldrevid"`
	NewRevID     int    `json:"newrevid"`
	NewTimestamp string `json:"newtimestamp"`
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

func NewHtmlUpdater(rootURL string) HtmlUpdater {
	return HtmlUpdater{rootURL: rootURL}
}

func (h *HtmlUpdater) GetFullHtml(title string, section string) string {
	base, err := url.Parse(h.rootURL)
	if err != nil {
		log.Printf("Failed to create baseURL with error %v\n", err)
		return ""
	}
	base.Path = path.Join(base.Path, "api.php")

	q := url.Values{}
	q.Add("format", "json")
	q.Add("action", "query")
	q.Add("titles", title)
	q.Add("export", "true")

	base.RawQuery = q.Encode()
	resp, err := http.Get(base.String())
	if err != nil {
		log.Println(err)
		return ""
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	var queryResponse QueryResponseWrapper
	err = json.Unmarshal([]byte(bytes), &queryResponse)
	if err != nil {
		log.Println(err)
		return ""
	}

	return queryResponse.Query.Export.Body
}

func (h *HtmlUpdater) UpdatePost(keys []string, title string, section string) bool {
	token := h.GetCSRFToken()

	var append string
	for _, key := range keys {
		append = append + BuildKeyHtml(key)
	}

	html := h.GetFullHtml(title, section)
	lastRowIndex := strings.LastIndex(html, "\n|-\nLegend ... ")
	if lastRowIndex == -1 {
		log.Println("Unable to find last row, not continuing")
		return false
	}
	lastRowIndex += 4

	finalText := html[:lastRowIndex] + append + html[lastRowIndex:]

	base, err := url.Parse(h.rootURL)
	if err != nil {
		log.Printf("Failed to create baseURL with error %v\n", err)
		return false
	}
	base.Path = path.Join(base.Path, "api.php")

	form := url.Values{}
	form.Add("action", "edit")
	form.Add("title", title)
	form.Add("section", section)
	form.Add("text", finalText)
	form.Add("bot", "true")
	form.Add("format", "json")
	form.Add("token", token)

	resp, err := http.PostForm(base.String(), form)
	if err != nil {
		log.Println(err)
		return false
	}
	bytes, err := ioutil.ReadAll((resp.Body))
	if err != nil {
		log.Println(err)
		return false
	}

	var editResponse EditResponseWrapper
	err = json.Unmarshal([]byte(bytes), &editResponse)
	if err != nil {
		log.Println(err)
		return false
	}
	if strings.ToLower(editResponse.Edit.Result) != "success" {
		log.Println(string(bytes))
		log.Printf("%v\n", editResponse)
		return false
	}

	return true
}

func (h *HtmlUpdater) GetCSRFToken() string {
	base, err := url.Parse(h.rootURL)
	if err != nil {
		log.Printf("Failed to create baseURL with error %v\n", err)
		return ""
	}
	base.Path = path.Join(base.Path, "api.php")

	q := url.Values{}
	q.Add("action", "query")
	q.Add("meta", "tokens")
	q.Add("format", "json")

	base.RawQuery = q.Encode()
	resp, err := http.Get(base.String())
	if err != nil {
		log.Println(err)
		return ""
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	var tokenResponse TokenQueryResponseWrapper
	err = json.Unmarshal([]byte(bytes), &tokenResponse)
	if err != nil {
		log.Println(err)
		return ""
	}
	return tokenResponse.Query.Tokens.CSRF
}

func FindNewKeys(keys []string, html string) []string {
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

func BuildKeyHtml(key string) string {
	return fmt.Sprintf("||@DuvalMagic || 1 Golden Key || TODAY_DATE || Unknown || %s || ❓ || ❓ || ❓ || ❓ || ❓\n|-\n", key)
}
