package gomail

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

const (
	testTo1  = "to1@example.com"
	testTo2  = "to2@example.com"
	testFrom = "from@example.com"
	testBody = "Test message"
	testMsg  = "To: " + testTo1 + ", " + testTo2 + "\r\n" +
		"From: " + testFrom + "\r\n" +
		"DKIM-Signature: v=1; a=rsa-sha256; q=dns/txt; c=relaxed/relaxed;\r\n" +
		" s=1592040826; d=example.com; h=from:subject:date:message-id;\r\n" +
		" bh=PkbvdhgKiEAAhc+GiwM2ZnxMu+StJ76lWGj3Y9phfFA=;\r\n" +
		" b=YsFEEu+AVAA5+Ccm18aF37Wl3f/CgxV9x12oCIA41JWQaseCmcqLD0DJCMGKEuyoJmdZ/D\r\n" +
		" L6HVSsiWIpDpgTEogOpQgqy86zEYW2KiZdjG3TBwxSrAVOI8y2ZBtW3KnanxpkBVoVEsTC\r\n" +
		" 7F0E3qk51MO5V6vCiR4HI9VaHbD9O6I=\r\n" +
		"Mime-Version: 1.0\r\n" +
		"Date: Wed, 25 Jun 2014 17:46:00 +0000\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: quoted-printable\r\n" +
		"\r\n" +
		testBody
)

type mockSender SendFunc

func (s mockSender) Send(from string, to []string, dkc *DKIMConfig, msg io.WriterTo) error {
	return s(from, to, dkc, msg)
}

type mockSendCloser struct {
	mockSender
	close func() error
}

func (s *mockSendCloser) Close() error {
	return s.close()
}

func TestSend(t *testing.T) {

	dkc := DKIMConfig{
		Selector: "1592040826",
		Domain:   "example.com",
		PubKey:   "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDbDWNDSezM0CLzzEvpM9dzw89DTDO+SMy4q6aZ63jTg3azolMZfhUcesDAd/4sRyPl+TnJ4Y60ULa67Z3wK61NyBOoCVzWCG9FvAO16RxAC11E6JAoj+DsusGjEGwYHq3fLPCgHhlprcOVLIgi3at5Zo9flh2K9+EuAgyWKanHQIDAQAB",
		PrivKey: `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDDbDWNDSezM0CLzzEvpM9dzw89DTDO+SMy4q6aZ63jTg3azolM
ZfhUcesDAd/4sRyPl+TnJ4Y60ULa67Z3wK61NyBOoCVzWCG9FvAO16RxAC11E6JA
oj+DsusGjEGwYHq3fLPCgHhlprcOVLIgi3at5Zo9flh2K9+EuAgyWKanHQIDAQAB
AoGAErJBlVMS30QiAr452HMOG815ib+/Ua3oPjANwFv2+O44yRxFane/AGU9tLXz
NZnMP7iqf6r6XpoyqTsv49kdXbIF0cjS+RLn/mfqlT77N1j5iJRuPmOmw4VGjZWr
naDHpSKCz+fWlGSxyLsSpHsYAUMjXX2q2bbAlWMY+DMAa0ECQQDMb8GQqmMmye+g
SWHvb5C1xVsBkQoxGR0u+pRkJuY5RYvyCB4VkGpz+i7d/nMcm6374NrYh/WpeKX9
a1zL09xlAkEA9LZwNWvMGCYE1O0uWpdsvF2/1tTAcxAPTWNi4JqyIRzrwob14uy1
Pw+d21S6BOyflzknz7EBypNzMo4AQj5oWQJAXEzUstEK5RdlFhQroGPZjQfmt8VZ
OaOiFnTSDIm3hgINViAuHQRP278H6/iW1kK/gaoahIqV8objQpB3nBsyNQJBAOwt
Q6CbWFAaKWGjQ7CVIpGt3V+m19J1Nn+XIy/ovXBt7DBDdv67O7YQCWdMn3fvM5uM
wxqVGEh+BJlPKXrFpokCQC1HVLroevl0SpRcNWpvi+ap/1f+FS9E9ZpC1M1bGZAn
v3TCqRxZuUGWPJkrNo0auVsxEVzmbjVAmTfLROprALc=
-----END RSA PRIVATE KEY-----`,
	}
	s := &mockSendCloser{
		mockSender: stubSend(t, testFrom, []string{testTo1, testTo2}, &dkc, testMsg),
		close: func() error {
			t.Error("Close() should not be called in Send()")
			return nil
		},
	}

	if err := Send(s, &dkc, getTestMessage()); err != nil {
		t.Errorf("Send(): %v", err)
	}
}

func getTestMessage() *Message {
	m := NewMessage()
	m.SetHeader("From", testFrom)
	m.SetHeader("To", testTo1, testTo2)
	m.SetBody("text/plain", testBody)

	return m
}

func stubSend(t *testing.T, wantFrom string, wantTo []string, dkc *DKIMConfig, wantBody string) mockSender {
	return func(from string, to []string, dkc *DKIMConfig, msg io.WriterTo) error {
		if from != wantFrom {
			t.Errorf("invalid from, got %q, want %q", from, wantFrom)
		}
		if !reflect.DeepEqual(to, wantTo) {
			t.Errorf("invalid to, got %v, want %v", to, wantTo)
		}

		buf := new(bytes.Buffer)
		_, err := msg.WriteTo(buf)
		if err != nil {
			t.Fatal(err)
		}
		compareBodies(t, buf.String(), wantBody)

		return nil
	}
}
