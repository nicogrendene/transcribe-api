package utils

import (
	"strconv"
	"strings"
)

// cleanPointerFormat limpia el formato de puntero "&{valor}" que puede venir de Pinecone
func CleanPointerFormat(s string) string {
	if len(s) > 3 && s[:2] == "&{" && s[len(s)-1] == '}' {
		return s[2 : len(s)-1]
	}
	return s
}

// validateFilename valida que el nombre de archivo sea seguro
func ValidateFilename(filename string) bool {
	if filename == "" {
		return false
	}

	// Verificar que no contenga caracteres peligrosos
	dangerousChars := []string{"..", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range dangerousChars {
		if strings.Contains(filename, char) {
			return false
		}
	}

	return true
}

// parseFloatFromString convierte una cadena a float64 de forma segura
func ParseFloatFromString(s string) (float64, error) {
	cleanVal := CleanPointerFormat(s)
	return strconv.ParseFloat(cleanVal, 64)
}

// sanitizeString limpia una cadena de caracteres no deseados
func SanitizeString(s string) string {
	// Remover caracteres de control y espacios extra
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")

	// Remover espacios múltiples
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	return s
}

// formatFileSize formatea el tamaño de archivo en formato legible
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return strconv.FormatFloat(float64(bytes)/float64(div), 'f', 1, 64) + " " + string([]byte{'K', 'M', 'G', 'T', 'P', 'E'}[exp]) + "B"
}

// containsString verifica si un slice contiene una cadena
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getFileExtension obtiene la extensión de un archivo
func GetFileExtension(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		return ""
	}
	return strings.ToLower(filename[dotIndex:])
}

// isVideoFile verifica si un archivo es de video
func IsVideoFile(filename string) bool {
	videoExtensions := []string{".mp4", ".avi", ".mov", ".mkv", ".webm", ".flv"}
	ext := GetFileExtension(filename)
	return ContainsString(videoExtensions, ext)
}

// isSubtitleFile verifica si un archivo es de subtítulos
func IsSubtitleFile(filename string) bool {
	subtitleExtensions := []string{".vtt", ".srt", ".ass", ".ssa"}
	ext := GetFileExtension(filename)
	return ContainsString(subtitleExtensions, ext)
}
