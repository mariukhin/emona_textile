package template

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

type EmailTemplates interface {
	AccountMemberInvitation(invitationLink, accountName, invitedUserName string) (string, error)
	AccountEmailConfirmation(confirmationLink, userName string) (string, error)

	UserSignUpEmail(confirmationLink, staffName string) (string, error)
	UserPasswordReset(passwordResetLink, staffName string) (string, error)
	UserEmailConfirmation(confirmationLink, customerName string) (string, error)

	StaffSignUpEmail(confirmationLink, loginLink, staffName, staffEmail, staffPass string) (string, error)
	StaffPasswordUpdated(loginLink, staffName, staffEmail, staffPass string) (string, error)
	StaffEmailConfirmation(confirmationLink, staffName string) (string, error)
}

func NewEmailTemplates(backOfficeHost, adminPanelHost string) (EmailTemplates, error) {
	return &emailTemplates{
		backOfficeHost: backOfficeHost,
		adminPanelHost: adminPanelHost,
		brand:          "CotonSMS",
	}, nil
}

type emailTemplates struct {
	backOfficeHost string
	adminPanelHost string
	brand          string
}

func (t *emailTemplates) AccountMemberInvitation(invitationLink, accountName, invitedUserName string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/account-invitation.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "account-invitation.html", struct {
		AccountName    string
		Brand          string
		InvitationLink string
	}{
		AccountName:    accountName,
		Brand:          t.brand,
		InvitationLink: invitationLink,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.backOfficeLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", invitedUserName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) AccountEmailConfirmation(confirmationLink, userName string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/account-email-confirmation.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "account-email-confirmation.html", struct {
		Brand            string
		ConfirmationLink string
	}{
		Brand:            t.brand,
		ConfirmationLink: confirmationLink,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.backOfficeLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", userName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) UserSignUpEmail(confirmationLink, staffName string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/customer-sign-up.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "customer-sign-up.html", struct {
		Brand            string
		ConfirmationLink string
	}{
		Brand:            t.brand,
		ConfirmationLink: confirmationLink,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.backOfficeLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", staffName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) UserPasswordReset(passwordResetLink, staffName string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/customer-password-reset.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "customer-password-reset.html", struct {
		Brand             string
		PasswordResetLink string
	}{
		Brand:             t.brand,
		PasswordResetLink: passwordResetLink,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.backOfficeLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", staffName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) UserEmailConfirmation(confirmationLink, customerName string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/customer-email-confirmation.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "customer-email-confirmation.html", struct {
		Brand            string
		ConfirmationLink string
	}{
		Brand:            t.brand,
		ConfirmationLink: confirmationLink,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.backOfficeLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", customerName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) StaffSignUpEmail(confirmationLink, loginLink, staffName, staffEmail, staffPass string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/staff-sign-up.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "staff-sign-up.html", struct {
		Brand            string
		ConfirmationLink string
		LoginUrl         string
		LoginEmail       string
		LoginPass        string
	}{
		Brand:            t.brand,
		ConfirmationLink: confirmationLink,
		LoginUrl:         loginLink,
		LoginEmail:       staffEmail,
		LoginPass:        staffPass,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.adminLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", staffName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) StaffPasswordUpdated(loginLink, staffName, staffEmail, staffPass string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/staff-password-updated.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "staff-password-updated.html", struct {
		Brand      string
		LoginUrl   string
		LoginEmail string
		LoginPass  string
	}{
		Brand:      t.brand,
		LoginUrl:   loginLink,
		LoginEmail: staffEmail,
		LoginPass:  staffPass,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.adminLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", staffName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) StaffEmailConfirmation(confirmationLink, staffName string) (string, error) {
	tmpl, err := template.ParseFiles("template/default.html", "template/staff-email-confirmation.html")
	if err != nil {
		return "", fmt.Errorf("fail to parse email template - %v", err)
	}

	var contentBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&contentBuff, "staff-email-confirmation.html", struct {
		Brand            string
		ConfirmationLink string
	}{
		Brand:            t.brand,
		ConfirmationLink: confirmationLink,
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	contentHtml := template.HTML(contentBuff.String())

	var emailBuff bytes.Buffer
	err = tmpl.ExecuteTemplate(&emailBuff, "default.html", struct {
		LogoURL string
		Brand   string
		Title   string
		Content template.HTML
		Year    string
	}{
		LogoURL: t.adminLogoUrl(),
		Brand:   t.brand,
		Title:   fmt.Sprintf("Dear %s", staffName),
		Content: contentHtml,
		Year:    t.currentYear(),
	})
	if err != nil {
		return "", fmt.Errorf("fail to execute template - %v", err)
	}

	return emailBuff.String(), nil
}

func (t *emailTemplates) backOfficeLogoUrl() string {
	return fmt.Sprintf("%s%s", t.backOfficeHost, "/assets/img/email-back-office-logo-dark.png")
}

func (t *emailTemplates) adminLogoUrl() string {
	return fmt.Sprintf("%s%s", t.adminPanelHost, "/assets/img/email-admin-logo-dark.png")
}

func (t *emailTemplates) currentYear() string {
	return fmt.Sprintf("%d", time.Now().Year())
}
