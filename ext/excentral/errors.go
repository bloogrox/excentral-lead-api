package excentral

import "github.com/pkg/errors"

var (
	ErrRequiredField                 = errors.New("Required field missing")
	ErrFieldNotValid                 = errors.New("Field not valid")
	ErrInvalidData                   = errors.New("Invalid data")
	ErrAPIAccess                     = errors.New("Affiliate not found or no API access")
	ErrVerificationParametersMissing = errors.New("Verification parameters missing")
	ErrWrongChecksum                 = errors.New("Wrong Checksum")
	ErrUnknown                       = errors.New("UnknownError")
	ErrDuplicateEmail                = errors.New("Lead email is already exists in the system")
	ErrRestrictedCountry             = errors.New("We do not accept leads from specified country")
)

var errMap = map[int]error{
	2:  ErrRequiredField,
	3:  ErrFieldNotValid,
	4:  ErrInvalidData,
	5:  ErrAPIAccess,
	7:  ErrVerificationParametersMissing,
	8:  ErrWrongChecksum,
	9:  ErrUnknown,
	11: ErrDuplicateEmail,
	24: ErrRestrictedCountry,
}
