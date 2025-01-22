package email

import (
	"bytes"
	"html/template"
	"os"
)

// RenderTemplate loads an HTML template and injects dynamic data
func RenderTemplate(templatePath string, data map[string]interface{}) (string, error) {
	// Open the template file
	file, err := os.Open(templatePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	// Render the template with dynamic data
	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, data)
	if err != nil {
		return "", err
	}

	return rendered.String(), nil
}
