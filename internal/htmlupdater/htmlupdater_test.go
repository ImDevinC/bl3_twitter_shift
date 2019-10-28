package htmlupdater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFullHtml(t *testing.T) {
	htmlUpdater := NewHtmlUpdater("http://orcz.com")
	result := htmlUpdater.GetFullHtml("Borderlands_3:_Shift_Codes", "2")
	assert.NotEmpty(t, result)
}

func TestCheckIfKeyExists(t *testing.T) {
	htmlUpdater := NewHtmlUpdater("http://orcz.com")
	html := htmlUpdater.GetFullHtml("Borderlands_3:_Shift_Codes", "2")
	results := FindNewKeys([]string{"CZCTJ-CZ59T-HC35W-T3BJB-ZTZJC", "ABCDE-ABCDE-ABCDE-ABCDE-ABCDE"}, html)
	assert.Len(t, results, 1)
	assert.Equal(t, "ABCDE-ABCDE-ABCDE-ABCDE-ABCDE", results[0])
}

func TestUpdatePost(t *testing.T) {
	htmlUpdater := NewHtmlUpdater("http://orcz.com")
	htmlUpdater.UpdatePost([]string{"ABCDE-ABCDE-ABCDE-ABCDE-ABCDE"}, "Borderlands_3:_Shift_Codes", "2")
}
