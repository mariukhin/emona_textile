package api

import (
	"amifactory.team/sequel/coton-app-backend/app/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var FormError = formError{}

type formError struct {
	Cause   error
	message string
}

func NewFormError(message string, reason error) error {
	return &formError{message: message, Cause: reason}
}
func (e formError) Unwrap() error { return e.Cause }
func (e formError) Error() string { return e.message }

var uuidRegex = regexp.MustCompile(`^[a-f\d]{8}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{12}$`)

func IsValidUUID(value string) bool {
	return uuidRegex.MatchString(value)
}

func ClearString(value string) string {
	return strings.TrimSpace(regexp.QuoteMeta(value))
}

type BaseForm struct {
	ValidationErrors []error
}

func (f *BaseForm) HasValidationErrors() bool {
	return len(f.ValidationErrors) > 0
}

func (f *BaseForm) AddValidationError(message string, err error) {
	f.ValidationErrors = append(f.ValidationErrors, NewFormError(message, err))
}

//
// Page form
//
type PageForm struct {
	*BaseForm
	Page         int
	PageSize     int
	SearchPhrase *string
	minSearchLen int
}

func NewPageForm(r *http.Request, minSearchLen int) *PageForm {
	f := PageForm{
		BaseForm:     &BaseForm{},
		Page:         0,
		PageSize:     10,
		SearchPhrase: nil,
		minSearchLen: minSearchLen,
	}

	srcParam := r.URL.Query().Get("src")
	pageParam := r.URL.Query().Get("page")
	pageSizeParam := r.URL.Query().Get("page_size")

	if len(pageParam) > 0 {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			f.AddValidationError("page__invalid", err)
		} else {
			f.Page = page
		}
	}

	if len(pageSizeParam) > 0 {
		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			f.AddValidationError("page_size__invalid", err)
		} else {
			f.PageSize = pageSize
			if f.PageSize <= 0 || 100 < f.PageSize {
				f.AddValidationError("page_size__out_of_bounds", nil)
			}
		}
	}

	if len(srcParam) > 0 {
		srcEscaped := ClearString(srcParam)
		f.SearchPhrase = &srcEscaped
		if utf8.RuneCountInString(*f.SearchPhrase) < f.minSearchLen {
			f.AddValidationError("src__invalid", nil)
		}
	}

	return &f
}

func (f *PageForm) Offset() int {
	return f.Page * f.PageSize
}

func (f *PageForm) HasFilter() bool {
	return f.SearchPhrase != nil
}

//
// Contact form
//
type ContactForm struct {
	Phone       *model.Phone
	Vars        map[string]interface{}
	allowedVars []model.ContactVariable
}

func NewContactForm(group *model.ContactGroup, r *http.Request) (*ContactForm, error) {
	form := ContactForm{
		Vars: make(map[string]interface{}),
	}
	form.allowedVars = group.Vars
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (f *ContactForm) UnmarshalJSON(b []byte) error {
	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}

	for _, v := range f.allowedVars {
		switch v.Type {
		case model.VarTypeString:
			strVar, found := jsonMap[v.Name].(string)
			if found {
				f.Vars[v.Name] = strVar
			}

		case model.VarTypeNumber:
			intVar, found := jsonMap[v.Name].(int)
			if found {
				f.Vars[v.Name] = intVar
			}

			floatVar, found := jsonMap[v.Name].(float64)
			if found {
				f.Vars[v.Name] = floatVar
			}

		case model.VarTypePhone:
			phoneRaw, found := jsonMap[v.Name].(string)
			if found {
				phone, err := model.ParsePhone(phoneRaw)
				if err == nil {
					if v.Name == "phone" {
						f.Phone = phone
					} else {
						f.Vars[v.Name] = phone
					}
				}
			}

		default:
			return fmt.Errorf("unexpected var type '%s'", v.Type)
		}
	}

	return nil
}

func (f *ContactForm) Validate() []error {
	errs := make([]error, 0)

	if f.Phone == nil {
		errs = append(errs, errors.New("phone__required"))
	}

	return errs
}

//
// Contact Variable form
//
type ContactVariableForm struct {
	Title              string                    `json:"title"`
	Type               model.ContactVariableType `json:"type"`
	DefaultValueString *string                   `json:"default_value_string,omitempty"`
	DefaultValueNumber *int                      `json:"default_value_number,omitempty"`
	isNew              bool
}

func NewContactVariableForm(r *http.Request, isNew bool) (*ContactVariableForm, error) {
	var form ContactVariableForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		return nil, err
	}
	form.isNew = isNew

	return &form, nil
}

func (f *ContactVariableForm) Validate() []error {
	errs := make([]error, 0)

	if f.isNew {
		if len(f.Title) == 0 {
			errs = append(errs, errors.New("title__required"))
		} else if len(f.Title) > 20 {
			errs = append(errs, errors.New("title__invalid"))
		}

		if len(f.Type) == 0 {
			errs = append(errs, errors.New("type__required"))
		} else {
			allowedTypes := []model.ContactVariableType{model.VarTypeString, model.VarTypeNumber}

			if !model.Contains(allowedTypes, f.Type) {
				errs = append(errs, errors.New("type__invalid"))
			}
		}
	}

	if f.DefaultValueString != nil && len(*f.DefaultValueString) > 20 {
		errs = append(errs, errors.New("default_value_string__invalid"))
	}

	return errs
}
