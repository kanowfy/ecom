package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidCreditCard = errors.New("Credit card info is invalid.")
	ErrExpiredCreditCard = errors.New("Your credit card has been expired.")
	UnacceptedCreditCard = errors.New("Sorry, we only accept VISA or MasterCard.")
	MasterCardRegx       = regexp.MustCompile("^(5[1-5][0-9]{14}|2(22[1-9][0-9]{12}|2[3-9][0-9]{13}|[3-6][0-9]{14}|7[0-1][0-9]{13}|720[0-9]{12}))$")
	VisaRegx             = regexp.MustCompile("^4[0-9]{12}(?:[0-9]{3})?$")
)

func VerifyCard(number string, cvv uint32, expiration_year uint32, expiration_month uint32) error {
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")
	if len(number) < 13 || len(number) > 16 {
		return ErrInvalidCreditCard
	}

	cvvStr := strconv.FormatUint(uint64(cvv), 10)
	if len(cvvStr) != 3 {
		return ErrInvalidCreditCard
	}

	if !MasterCardRegx.MatchString(number) && !VisaRegx.MatchString(number) {
		return UnacceptedCreditCard
	}

	if expiration_year <= uint32(time.Now().Year()) {
		if expiration_year == uint32(time.Now().Year()) {
			if expiration_month < uint32(time.Now().Month()) {
				return ErrExpiredCreditCard
			}
		} else {
			return ErrExpiredCreditCard
		}
	}

	return nil
}
