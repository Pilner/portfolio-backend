package http_api

// NormalizeEmailAlias strips email aliases from the local part of the email address.
// For example, "foo+bar@test.com" becomes "foo@test.com".
func NormalizeEmailAlias(email string) (string, error) {
	// Find the position of '@' to separate the local part from the domain.
	at := -1
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			at = i
			break
		}
	}
	// If '@' is not found, return as an invalid email.
	if at == -1 {
		// Not a valid email, return as-is with an error.
		return email, ErrInvalidEmail
	}

	// local is the part before '@', domain is the part including and after '@'
	local := email[:at]
	domain := email[at:]

	// Look for a '+' in the local part, which signifies an alias.
	plusIdx := -1
	for j := 0; j < len(local); j++ {
		if local[j] == '+' {
			plusIdx = j
			break
		}
	}

	// If a '+' is found, remove everything from '+' to '@'.
	if plusIdx != -1 {
		return local[:plusIdx] + domain, nil
	}

	// If there is no alias, return the email unchanged.
	return email, nil
}
