package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
)

// smtpMailer sends via net/smtp, configured from MAIL_* env vars.
type smtpMailer struct {
	host, port, user, pass, from, encryption string
}

func newSMTP() *smtpMailer {
	return &smtpMailer{
		host:       os.Getenv("MAIL_HOST"),
		port:       envOr("MAIL_PORT", "587"),
		user:       os.Getenv("MAIL_USERNAME"),
		pass:       os.Getenv("MAIL_PASSWORD"),
		from:       os.Getenv("MAIL_FROM"),
		encryption: strings.ToLower(envOr("MAIL_ENCRYPTION", "starttls")),
	}
}

func (s *smtpMailer) Send(_ context.Context, m Message) error {
	if s.host == "" {
		return fmt.Errorf("mail: MAIL_HOST not set")
	}
	from := m.From
	if from == "" {
		from = s.from
	}
	addr := net.JoinHostPort(s.host, s.port)
	recipients := append(append(append([]string{}, m.To...), m.Cc...), m.Bcc...)
	raw := buildMessage(from, m)

	var auth smtp.Auth
	if s.user != "" {
		auth = smtp.PlainAuth("", s.user, s.pass, s.host)
	}

	// Implicit TLS (e.g. port 465).
	if s.encryption == "tls" || s.encryption == "ssl" {
		return s.sendTLS(addr, auth, from, recipients, raw)
	}
	// STARTTLS / plain handled by smtp.SendMail (it upgrades when offered).
	return smtp.SendMail(addr, auth, from, recipients, raw)
}

func (s *smtpMailer) sendTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12})
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return err
	}
	defer c.Quit()
	if auth != nil {
		if err := c.Auth(auth); err != nil {
			return err
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, r := range to {
		if err := c.Rcpt(r); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}

// buildMessage assembles an RFC 5322 message (multipart when HTML is present).
func buildMessage(from string, m Message) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", strings.Join(m.To, ", "))
	if len(m.Cc) > 0 {
		fmt.Fprintf(&b, "Cc: %s\r\n", strings.Join(m.Cc, ", "))
	}
	fmt.Fprintf(&b, "Subject: %s\r\n", m.Subject)
	b.WriteString("MIME-Version: 1.0\r\n")
	if m.HTML != "" {
		boundary := "togo-boundary-9f8a7b6c"
		fmt.Fprintf(&b, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary)
		fmt.Fprintf(&b, "--%s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n", boundary, m.Text)
		fmt.Fprintf(&b, "--%s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n", boundary, m.HTML)
		fmt.Fprintf(&b, "--%s--\r\n", boundary)
	} else {
		b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		b.WriteString(m.Text)
	}
	return []byte(b.String())
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
