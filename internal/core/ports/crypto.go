package ports

type HashManager interface {
	Hash(text string) (string, error)
	Compare(hashedText, textToCompare string) error
}
