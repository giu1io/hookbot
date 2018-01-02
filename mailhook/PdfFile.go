package mailhook

import "os"

// PdfFile class
type PdfFile struct {
	path string
}

// Path method
func (p *PdfFile) Path() string {
	return p.path
}

// Destroy method
func (p *PdfFile) Destroy() error {
	return os.Remove(p.path)
}
