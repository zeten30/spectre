package main

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/zeten30/spectre"
)

func generateRandomBytes(nbytes int) ([]byte, error) {
	uuid := make([]byte, nbytes)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return []byte{}, err
	}

	return uuid, nil
}

var base32Encoder = base32.NewEncoding("abcdefghjkmnopqrstuvwxyz23456789")

func generateRandomBase32String(outlen int) (string, error) {
	nbytes := (outlen * 5 / 8) + 1
	uuid, err := generateRandomBytes(nbytes)
	if err != nil {
		return "", err
	}

	s := base32Encoder.EncodeToString(uuid)
	if outlen == -1 {
		outlen = len(s)
	}

	return s[0:outlen], nil
}

type mockPaste struct {
	ID, LanguageName, Title string
	ExpirationTime          *time.Time
	Body                    string
}

func (m *mockPaste) GetID() spectre.PasteID {
	return spectre.PasteID(m.ID)
}

func (m *mockPaste) GetLanguageName() string {
	return m.LanguageName
}

func (m *mockPaste) GetExpirationTime() *time.Time {
	return m.ExpirationTime
}

func (m *mockPaste) GetTitle() string {
	return m.Title
}

func (m *mockPaste) IsEncrypted() bool {
	return false
}

func (m *mockPaste) GetEncryptionMethod() spectre.EncryptionMethod {
	return 0
}

func (m *mockPaste) GetModificationTime() time.Time {
	return time.Now()
}

type irc struct {
	io.Reader
}

func (*irc) Close() error {
	return nil
}

func (m *mockPaste) Reader() (io.ReadCloser, error) {
	logrus.Infof("Reader(%v)", m.ID)
	return &irc{Reader: strings.NewReader(m.Body)}, nil
}

func (m *mockPaste) Update(u spectre.PasteUpdate) error {
	logrus.Infof("Update(%v, %v)", m.ID, u)
	if u.Title != nil {
		m.Title = *u.Title
	}

	if u.LanguageName != nil {
		m.LanguageName = *u.LanguageName
	}

	if u.Body != nil {
		m.Body = *u.Body
	}

	if u.ExpirationTime != nil {
		m.ExpirationTime = u.ExpirationTime
	}

	return nil
}

func (m *mockPaste) Erase() error {
	return nil
}

type mockPasteService struct {
	o sync.Once
	m map[spectre.PasteID]*mockPaste
}

func (m *mockPasteService) init() {
	m.o.Do(func() {
		m.m = make(map[spectre.PasteID]*mockPaste)
		t := time.Now().Add(2 * time.Hour)
		m.m[spectre.PasteID("abcde")] = &mockPaste{
			ID:             "abcde",
			Title:          "Hello World",
			LanguageName:   "text",
			ExpirationTime: &t,
			Body:           "I am a real paste; I promise!",
		}
	})

}

func (m *mockPasteService) CreatePaste(context.Context, spectre.Cryptor) (spectre.Paste, error) {
	m.init()
	logrus.Infof("CreatePaste()")
	i, _ := generateRandomBase32String(5)
	id := spectre.PasteID(i)
	logrus.Infof("-> %v", id)
	p := &mockPaste{ID: i, LanguageName: "text"}
	m.m[id] = p
	return p, nil
}

func (m *mockPasteService) GetPaste(c context.Context, cr spectre.Cryptor, id spectre.PasteID) (spectre.Paste, error) {
	m.init()
	logrus.Infof("GetPaste(%v)", id)
	p, ok := m.m[id]
	if !ok {
		logrus.Errorf("-> not found")
		return nil, spectre.ErrNotFound
	}
	logrus.Infof("-> found")
	return p, nil
}

func (m *mockPasteService) GetPastes(context.Context, []spectre.PasteID) ([]spectre.Paste, error) {
	return nil, spectre.ErrNotFound
}

func (m *mockPasteService) DestroyPaste(c context.Context, id spectre.PasteID) (bool, error) {
	logrus.Infof("DestroyPaste(%v)", id)
	_, ok := m.m[id]
	delete(m.m, id)
	return ok, nil
}
