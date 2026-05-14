package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateService(t *testing.T) {
	svc := NewService()

	t.Run("Render invitation template", func(t *testing.T) {
		data := map[string]interface{}{
			"UserName":      "John Doe",
			"Title":         "Test Survey",
			"InvitationURL": "http://example.com/invite",
			"PlatformName":  "CustomForms",
		}
		result, err := svc.Render("invitation.html", data)
		assert.NoError(t, err)
		assert.Contains(t, result, "John Doe")
		assert.Contains(t, result, "Test Survey")
		assert.Contains(t, result, "http://example.com/invite")
	})

	t.Run("Render non-existent template", func(t *testing.T) {
		_, err := svc.Render("non-existent.html", nil)
		assert.Error(t, err)
	})
}
