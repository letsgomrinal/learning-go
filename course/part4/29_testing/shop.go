// Package lesson29: the code you'll write tests for. Read README.md first!
//
// You'll fix ONE bug in this file (in Slugify), after your tests find it. Leave the rest as it is.
package lesson29

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

// Slugify turns a title into a URL-friendly "slug":
//
//   - lower case
//
//   - letters and digits are kept
//
//   - every run of anything else (spaces, punctuation...) becomes ONE dash
//
//   - no dash at the start or the end
//
//     Slugify("Hello, World!")      == "hello-world"
//     Slugify("  Go -- is   fun  ") == "go-is-fun"
//     Slugify("Café au lait")       == "café-au-lait"
//
// There's a bug in here. Write tests (Ex1) until one of them shows it, then fix it.
func Slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

// Clock tells the time. Code that needs "now" takes a Clock, so tests can pass a fake one.
type Clock interface {
	Now() time.Time
}

// RealClock is the Clock a real program uses.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// Greeting depends on the hour: before 12 "Good morning", before 18 "Good afternoon",
// otherwise "Good evening".
func Greeting(c Clock) string {
	h := c.Now().Hour()
	switch {
	case h < 12:
		return "Good morning"
	case h < 18:
		return "Good afternoon"
	default:
		return "Good evening"
	}
}

// Mailer sends emails. A real program would use an SMTP server or an email API.
type Mailer interface {
	Send(to, body string) error
}

// Invoice is a bill a customer has to pay.
type Invoice struct {
	Customer string // the customer's email
	Amount   int    // in cents
	Due      time.Time
	Paid     bool
}

// NotifyOverdue emails every customer whose invoice is unpaid and past its due date (strictly before
// now). It returns how many emails were sent. If sending fails, it stops and returns the count so far
// and the error, wrapped as "notify <customer>: <err>".
func NotifyOverdue(m Mailer, invoices []Invoice, now time.Time) (int, error) {
	sent := 0
	for _, inv := range invoices {
		if inv.Paid || !inv.Due.Before(now) {
			continue
		}
		body := fmt.Sprintf("Your invoice of $%d.%02d was due on %s.", inv.Amount/100, inv.Amount%100, inv.Due.Format("2006-01-02"))
		if err := m.Send(inv.Customer, body); err != nil {
			return sent, fmt.Errorf("notify %s: %w", inv.Customer, err)
		}
		sent++
	}
	return sent, nil
}
