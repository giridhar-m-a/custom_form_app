package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

type Service interface {
	Render(templateName string, data any) (string, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Render(templateName string, data any) (string, error) {
	tmpl, err := s.getTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}

func (s *service) getTemplate(templateName string) (*template.Template, error) {
	cacheMutex.RLock()
	tmpl, ok := templateCache[templateName]
	cacheMutex.RUnlock()

	if ok {
		return tmpl, nil
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check locking
	if tmpl, ok := templateCache[templateName]; ok {
		return tmpl, nil
	}

	parsed, err := template.ParseFS(
		templateFS,
		"files/"+templateName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	templateCache[templateName] = parsed
	return parsed, nil
}
