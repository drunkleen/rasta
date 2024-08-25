package emailPkg

import (
	"bytes"
	"errors"
	"github.com/drunkleen/rasta/config"
	commonerrors "github.com/drunkleen/rasta/internal/common/errors"
	newslettermodel "github.com/drunkleen/rasta/internal/models/newsletter"
	"github.com/drunkleen/rasta/internal/models/user"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"time"
)

type OtpEmailData struct {
	Otp               string
	FirstName         string
	Username          string
	HelpCenterEmail   string
	HelpCenterAddress string
	IssuerName        string
	DateNow           time.Time
}

type NewsletterEmailData struct {
	Body              string
	HelpCenterEmail   string
	HelpCenterAddress string
	IssuerName        string
	DateNow           time.Time
}

func SendEmail(htmlPathFile string, targetEmail string, subject string, EmailData any) error {
	tmpl, err := template.ParseFiles(htmlPathFile)
	if err != nil {
		return errors.New("internal server error")
	}

	var data any
	var ok bool
	switch EmailData.(type) {
	case *OtpEmailData:
		data, ok = EmailData.(*OtpEmailData)
		if !ok {
			return errors.New("internal server error")
		}
	case *NewsletterEmailData:
		data, ok = EmailData.(*NewsletterEmailData)
		if !ok {
			return errors.New("internal server error")
		}
	default:
		return errors.New("internal server error")
	}
	var body bytes.Buffer
	if err = tmpl.Execute(&body, data); err != nil {
		return errors.New("internal server error")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.GetEmailUsername())
	m.SetHeader("To", targetEmail)
	m.SetHeader("Subject", config.GetJwtIssuer()+" - "+subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(config.GetEmailHost(), config.GetEmailPort(), config.GetEmailUsername(), config.GetEmailPassword())

	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmailVerify(user *usermodel.User) error {
	data := &OtpEmailData{
		Otp:               user.OtpEmail.Code,
		FirstName:         user.FirstName,
		Username:          user.Username,
		HelpCenterEmail:   config.GetHelpCenterEmail(),
		HelpCenterAddress: config.GetHelpCenterAddress(),
		IssuerName:        config.GetJwtIssuer(),
		DateNow:           time.Now().Truncate(24 * time.Hour),
	}
	return SendEmail(
		"pkg/email/email_templates/welcome_and_verify.html",
		user.Email,
		"Verify your E-mail address",
		data,
	)
}

func SendEmailResetPassword(user *usermodel.User) error {
	data := &OtpEmailData{
		Otp:               user.ResetPwd.Code,
		FirstName:         user.FirstName,
		Username:          user.Username,
		HelpCenterEmail:   config.GetHelpCenterEmail(),
		HelpCenterAddress: config.GetHelpCenterAddress(),
		IssuerName:        config.GetJwtIssuer(),
		DateNow:           time.Now().Truncate(24 * time.Hour),
	}
	return SendEmail(
		"pkg/email/email_templates/reset_password.html",
		user.Email,
		"Reset password",
		data,
	)
}

func SendNewsletter(targetEmails *[]newslettermodel.Newsletter, EmailBody *string) error {
	subject := "Newsletter"
	for _, email := range *targetEmails {
		data := NewsletterEmailData{
			Body:              *EmailBody,
			HelpCenterEmail:   config.GetHelpCenterEmail(),
			HelpCenterAddress: config.GetHelpCenterAddress(),
			IssuerName:        config.GetJwtIssuer(),
			DateNow:           time.Now().Truncate(24 * time.Hour),
		}
		err := SendEmail("pkg/email/email_templates/news_letter.html", email.Email, subject, &data)
		if err != nil {
			log.Printf("failed to send email: %v", err)
			return errors.New(commonerrors.ErrInternalServer)
		}
	}
	return nil
}
