package mdmail

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark-meta"
)

var mds = `---
From:
  - name: Charles Srisuwananukorn
    email: csrisuw@gmail.com
To:
  - name: Charles Srisuwananukorn
    email: csrisuw@gmail.com
Subject: Hello, world!
---
# Hello goldmark-meta

Sample text.

[link](http://example.com)
`

func mdToHTML(md []byte) (map[string]interface{}, []byte) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(md, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := meta.Get(context)

	return metaData, buf.Bytes()
}

func getAddresses(meta map[string]interface{}, key string) []*mail.Address {
	var addresses []*mail.Address
	for _, v := range meta[key].([]interface{}) {
		v := v.(map[interface{}]interface{})
		addresses = append(addresses, &mail.Address{v["name"].(string), v["email"].(string)})
	}
	return addresses
}

type frontmatter struct {
	From    []*mail.Address
	To      []*mail.Address
	Subject string
	Date    time.Time
}

// TODO: Make this return an error since parsing the frontmatter can fail
func newFrontmatterFromMeta(meta map[string]interface{}) *frontmatter {
	return &frontmatter{
		From:    getAddresses(meta, "From"),
		To:      getAddresses(meta, "To"),
		Subject: meta["Subject"].(string),
		Date:    time.Now(),
	}
}

type imapCredentials struct {
	Server   string
	Username string
	Password string
}

// TODO: NewImapCredentialsFromEnv should return an error if any of the env vars are missing
func newImapCredentialsFromEnv() *imapCredentials {
	return &imapCredentials{
		Server:   os.Getenv("IMAP_SERVER"),
		Username: os.Getenv("IMAP_USER"),
		Password: os.Getenv("IMAP_PASSWORD"),
	}
}

func createMail(frontmatter *frontmatter, html []byte) (bytes.Buffer) {
	var b bytes.Buffer

	// Create our mail header
	var h mail.Header
	h.SetAddressList("From", frontmatter.From)
	h.SetAddressList("To", frontmatter.To)
	h.SetSubject(frontmatter.Subject)
	h.SetDate(frontmatter.Date)

	// Create a new mail writer
	mw, err := mail.CreateWriter(&b, h)
	if err != nil {
		log.Fatal(err)
	}

	tw, err := mw.CreateInline()
	if err != nil {
		log.Fatal(err)
	}

	// Omitting a plain-text alternative beacuse Mail.app will create one anyway

	var htmlHeader mail.InlineHeader
	htmlHeader.Set("Content-Type", "text/html")
	hpw, err := tw.CreatePart(htmlHeader)
	if err != nil {
		log.Fatal(err)
	}
	hpw.Write(html)
	hpw.Close()

	tw.Close()
	mw.Close()

	return b
}

func createDraft(imapCredentials *imapCredentials, b bytes.Buffer) {
	// Connect to server
	c, err := client.DialTLS(imapCredentials.Server, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(imapCredentials.Username, imapCredentials.Password); err != nil {
		log.Fatal(err)
	}

	// Append it to Drafts
	if err := c.Append("[Gmail]/Drafts", nil, time.Now(), &b); err != nil {
		log.Fatal(err)
	}
}

func CreateDraftFromMarkdown(cmd *cobra.Command, args []string) {
	md, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	meta, html := mdToHTML(md)
	frontmatter := newFrontmatterFromMeta(meta)
	imapCredentials := newImapCredentialsFromEnv()

	b := createMail(frontmatter, html)
	createDraft(imapCredentials, b)
}
