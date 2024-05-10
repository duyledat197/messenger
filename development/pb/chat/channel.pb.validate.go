// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: chat/channel.proto

package chat

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Channel with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Channel) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Channel with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ChannelMultiError, or nil if none found.
func (m *Channel) ValidateAll() error {
	return m.validate(true)
}

func (m *Channel) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ChannelId

	// no validation rules for Name

	// no validation rules for Description

	if len(errors) > 0 {
		return ChannelMultiError(errors)
	}

	return nil
}

// ChannelMultiError is an error wrapping multiple validation errors returned
// by Channel.ValidateAll() if the designated constraints aren't met.
type ChannelMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChannelMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChannelMultiError) AllErrors() []error { return m }

// ChannelValidationError is the validation error returned by Channel.Validate
// if the designated constraints aren't met.
type ChannelValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChannelValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChannelValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChannelValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChannelValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChannelValidationError) ErrorName() string { return "ChannelValidationError" }

// Error satisfies the builtin error interface
func (e ChannelValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChannel.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChannelValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChannelValidationError{}

// Validate checks the field values on SearchChannelByNameRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SearchChannelByNameRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SearchChannelByNameRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SearchChannelByNameRequestMultiError, or nil if none found.
func (m *SearchChannelByNameRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SearchChannelByNameRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	// no validation rules for Offset

	// no validation rules for Limit

	if len(errors) > 0 {
		return SearchChannelByNameRequestMultiError(errors)
	}

	return nil
}

// SearchChannelByNameRequestMultiError is an error wrapping multiple
// validation errors returned by SearchChannelByNameRequest.ValidateAll() if
// the designated constraints aren't met.
type SearchChannelByNameRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SearchChannelByNameRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SearchChannelByNameRequestMultiError) AllErrors() []error { return m }

// SearchChannelByNameRequestValidationError is the validation error returned
// by SearchChannelByNameRequest.Validate if the designated constraints aren't met.
type SearchChannelByNameRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SearchChannelByNameRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SearchChannelByNameRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SearchChannelByNameRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SearchChannelByNameRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SearchChannelByNameRequestValidationError) ErrorName() string {
	return "SearchChannelByNameRequestValidationError"
}

// Error satisfies the builtin error interface
func (e SearchChannelByNameRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSearchChannelByNameRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SearchChannelByNameRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SearchChannelByNameRequestValidationError{}

// Validate checks the field values on SearchChannelByNameResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SearchChannelByNameResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SearchChannelByNameResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SearchChannelByNameResponseMultiError, or nil if none found.
func (m *SearchChannelByNameResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SearchChannelByNameResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetChannels() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, SearchChannelByNameResponseValidationError{
						field:  fmt.Sprintf("Channels[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, SearchChannelByNameResponseValidationError{
						field:  fmt.Sprintf("Channels[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SearchChannelByNameResponseValidationError{
					field:  fmt.Sprintf("Channels[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return SearchChannelByNameResponseMultiError(errors)
	}

	return nil
}

// SearchChannelByNameResponseMultiError is an error wrapping multiple
// validation errors returned by SearchChannelByNameResponse.ValidateAll() if
// the designated constraints aren't met.
type SearchChannelByNameResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SearchChannelByNameResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SearchChannelByNameResponseMultiError) AllErrors() []error { return m }

// SearchChannelByNameResponseValidationError is the validation error returned
// by SearchChannelByNameResponse.Validate if the designated constraints
// aren't met.
type SearchChannelByNameResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SearchChannelByNameResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SearchChannelByNameResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SearchChannelByNameResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SearchChannelByNameResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SearchChannelByNameResponseValidationError) ErrorName() string {
	return "SearchChannelByNameResponseValidationError"
}

// Error satisfies the builtin error interface
func (e SearchChannelByNameResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSearchChannelByNameResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SearchChannelByNameResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SearchChannelByNameResponseValidationError{}
