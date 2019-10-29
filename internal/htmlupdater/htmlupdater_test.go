package htmlupdater

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFullHtml(t *testing.T) {
	htmlUpdater := NewHtmlUpdater("http://orcz.com/api.php")
	result, err := htmlUpdater.getFullHTML("User:ImDevinC")
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
}

func TestCheckIfKeyExists(t *testing.T) {
	htmlUpdater := NewHtmlUpdater("http://orcz.com/api.php")
	html, err := htmlUpdater.getFullHTML("User:ImDevinC")
	assert.Nil(t, err)
	results := findNewKeys([]string{"CZCTJ-CZ59T-HC35W-T3BJB-ZTZJC", "ABCDE-ABCDE-ABCDE-ABCDE-ABCDE"}, html)
	assert.Len(t, results, 1)
	assert.Equal(t, "ABCDE-ABCDE-ABCDE-ABCDE-ABCDE", results[0])
}

func TestUpdatePost(t *testing.T) {
	htmlUpdater := NewHtmlUpdater("http://orcz.com/api.php")
	err := htmlUpdater.login(os.Getenv("ORCZ_USERNAME"), os.Getenv("ORCZ_PASSWORD"))
	assert.Nil(t, err)
	htmlUpdater.updatePost([]string{"ABCDE-ABCDE-ABCDE-ABCDE-ABCDE"}, "Oct 28, 2019", "User:ImDevinC", "2", "", "")
}

func TestSolveCaptcha(t *testing.T) {
	question := "Type only the four missing letters from the following words: The USA Presidential Nominees are Hillary Clinton and Donald T____"
	answer := solveCaptcha(question)
	assert.Equal(t, "rump", answer)
}

func TestLogin(t *testing.T) {
	htmlUpdater := NewHtmlUpdater("http://orcz.com/api.php")
	err := htmlUpdater.login(os.Getenv("ORCZ_USERNAME"), os.Getenv("ORCZ_PASSWORD"))
	assert.Nil(t, err)
}

func TestAddKeys(t *testing.T) {
	AddKeys("User:ImDevinC", os.Getenv("ORCZ_USERNAME"), os.Getenv("ORCZ_PASSWORD"), []string{"ABCDE-ABCDE-ABCDE-ABCDE-ABCDF"}, "Oct 28, 2019")
}
