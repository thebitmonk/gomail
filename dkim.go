package gomail

type DKIMConfig struct {
	Domain   string
	PubKey   string
	PrivKey  string
	Selector string
	Headers  string
}
