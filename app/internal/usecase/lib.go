package usecase

import "strings"

func buildLogoUrl(baseUrl string, path string) string {
    if path == "" {
        return ""
    }

    return strings.Join([]string{baseUrl, path}, "/")
}
