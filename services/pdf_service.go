package services

import (
	"bytes"
	"fmt"
	"html/template"
	"notification-server/models"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFService struct{}

// NewPDFService creates a new PDF service instance
func NewPDFService() *PDFService {
	return &PDFService{}
}

// GenerateUserProfilePDF generates a PDF document with user profile details
func (s *PDFService) GenerateUserProfilePDF(user *models.User) ([]byte, error) {
	// Generate HTML content from template
	htmlContent, err := s.generateUserProfileHTML(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML template: %w", err)
	}

	// Create PDF from HTML
	pdfBytes, err := s.htmlToPDF(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBytes, nil
}

// generateUserProfileHTML generates HTML content for user profile PDF
func (s *PDFService) generateUserProfileHTML(user *models.User) (string, error) {
	tmpl, err := template.ParseFiles("./templates/user_profile.html")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, user); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// htmlToPDF converts HTML string to PDF bytes using wkhtmltopdf
func (s *PDFService) htmlToPDF(htmlContent string) ([]byte, error) {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	// Create a new page from HTML string
	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlContent)))
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	// Return PDF bytes
	return pdfg.Bytes(), nil
}
