package utils

import "strings"

// MaskToken masks a git token while showing the prefix and first few characters of the secret
func MaskToken(token string) string {
	if token == "" {
		return "<empty>"
	}

	// Common git forge token prefixes
	prefixes := []string{
		"github_pat_", // GitHub fine-grained PAT
		"ghp_",        // GitHub personal access token
		"gho_",        // GitHub OAuth token
		"ghu_",        // GitHub user-to-server token
		"ghs_",        // GitHub server-to-server token
		"ghr_",        // GitHub refresh token
		"glpat-",      // GitLab personal access token
		"gloas-",      // GitLab OAuth application secret
		"glptt-",      // GitLab project access token
	}

	// Find matching prefix
	var prefix string
	secretStart := 0

	for _, p := range prefixes {
		if strings.HasPrefix(token, p) {
			prefix = p
			secretStart = len(p)
			break
		}
	}

	// If no known prefix, treat entire token as secret
	if prefix == "" {
		if len(token) <= 7 {
			return "***"
		}

		return token[:7] + "***"
	}

	// Show prefix + first 7 chars of secret
	secretPart := token[secretStart:]
	if len(secretPart) <= 7 {
		return prefix + "***"
	}

	return prefix + secretPart[:7] + "***"
}

// MaskTokensInMap recursively masks token values in a map
func MaskTokensInMap(m map[string]interface{}) map[string]interface{} {
	masked := make(map[string]interface{})
	for k, v := range m {
		if k == "token" {
			if strValue, ok := v.(string); ok && strValue != "" {
				masked[k] = MaskToken(strValue)
			} else {
				masked[k] = v
			}
		} else if mapValue, ok := v.(map[string]interface{}); ok {
			masked[k] = MaskTokensInMap(mapValue)
		} else {
			masked[k] = v
		}
	}
	return masked
}
