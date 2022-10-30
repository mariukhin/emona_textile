package api

import "fmt"

type LinkBuilder interface {
	AccountMemberInvitation(invitationToken string) string
	AccountEmailConfirmation(confirmationToken string) string

	UserEmailConfirmation(confirmationToken string) string
	UserPasswordReset(passwordResetToken string) string

	StaffEmailConfirmation(confirmationToken string) string
	StaffLogin() string
}

type linkBuilder struct {
	backofficeHost string
	adminPanelHost string
}

func NewLinkBuilder(backofficeHost, adminPanelHost string) LinkBuilder {
	return &linkBuilder{
		backofficeHost: backofficeHost,
		adminPanelHost: adminPanelHost,
	}
}

func (linkBuilder *linkBuilder) AccountMemberInvitation(invitationToken string) string {
	return fmt.Sprintf("%s/account-invitation?t=%s", linkBuilder.backofficeHost, invitationToken)
}

func (linkBuilder *linkBuilder) AccountEmailConfirmation(confirmationToken string) string {
	return fmt.Sprintf("%s/account-email-confirmation?t=%s", linkBuilder.backofficeHost, confirmationToken)
}

func (linkBuilder *linkBuilder) UserEmailConfirmation(confirmationToken string) string {
	return fmt.Sprintf("%s/email-confirmation?t=%s", linkBuilder.backofficeHost, confirmationToken)
}

func (linkBuilder *linkBuilder) UserPasswordReset(passwordResetToken string) string {
	return fmt.Sprintf("%s/password-reset?t=%s", linkBuilder.backofficeHost, passwordResetToken)
}

func (linkBuilder *linkBuilder) StaffEmailConfirmation(confirmationToken string) string {
	return fmt.Sprintf("%s/email-confirmation?t=%s", linkBuilder.adminPanelHost, confirmationToken)
}

func (linkBuilder *linkBuilder) StaffLogin() string {
	return fmt.Sprintf("%s/sign-in", linkBuilder.adminPanelHost)
}
