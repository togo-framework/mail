// Package mail is togo's email subsystem: a Mailer contract with an SMTP driver
// (default) and a dev "log" driver. Additional providers (AWS SES, Resend, …)
// ship as driver plugins that call mail.RegisterDriver and depend on this
// package. Used by auth for OTP/verification/reset flows.
//
// Install: `togo install togo-framework/mail` (blank-import registers it).
package mail

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/togo-framework/togo"
)

// Message is an email to send.
type Message struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Text    string // plain-text body
	HTML    string // optional HTML body
}

// Mailer sends messages. Drivers (smtp, resend, ses, …) implement it.
type Mailer interface {
	Send(ctx context.Context, m Message) error
}

// DriverFactory builds a Mailer from the kernel (env-configured).
type DriverFactory func(k *togo.Kernel) (Mailer, error)

var (
	regMu   sync.RWMutex
	drivers = map[string]DriverFactory{}
)

// RegisterDriver registers a mail driver by name (call from a plugin's init()).
func RegisterDriver(name string, f DriverFactory) {
	regMu.Lock()
	drivers[name] = f
	regMu.Unlock()
}

func init() {
	RegisterDriver("smtp", func(k *togo.Kernel) (Mailer, error) { return newSMTP(), nil })
	RegisterDriver("log", func(k *togo.Kernel) (Mailer, error) { return &logMailer{k: k}, nil })

	togo.RegisterProviderFunc("mail", togo.PriorityService, func(k *togo.Kernel) error {
		name := os.Getenv("MAIL_DRIVER")
		if name == "" {
			if os.Getenv("MAIL_HOST") != "" {
				name = "smtp"
			} else {
				name = "log" // safe dev default: don't send, just log
			}
		}
		regMu.RLock()
		f, ok := drivers[name]
		regMu.RUnlock()
		if !ok {
			return fmt.Errorf("mail: unknown driver %q (install its plugin?)", name)
		}
		m, err := f(k)
		if err != nil {
			return err
		}
		k.Set("mail", &Service{mailer: m, driver: name})
		return nil
	})
}

// Service is the mail runtime stored on the kernel (k.Get("mail")).
type Service struct {
	mailer Mailer
	driver string
}

// Send dispatches a message via the configured driver.
func (s *Service) Send(ctx context.Context, m Message) error { return s.mailer.Send(ctx, m) }

// Driver returns the active driver name.
func (s *Service) Driver() string { return s.driver }

// FromKernel fetches the mail service from the kernel container.
func FromKernel(k *togo.Kernel) (*Service, bool) {
	v, ok := k.Get("mail")
	if !ok {
		return nil, false
	}
	s, ok := v.(*Service)
	return s, ok
}

// logMailer logs messages instead of sending — the safe default for dev/tests.
type logMailer struct{ k *togo.Kernel }

func (l *logMailer) Send(_ context.Context, m Message) error {
	if l.k != nil && l.k.Log != nil {
		l.k.Log.Info("mail (log driver)", "to", m.To, "subject", m.Subject)
	}
	return nil
}
