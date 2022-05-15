package common

import "errors"

var (
	ErrInvalidPKARecord         = errors.New("invalid pka record was found")
	ErrNoDNS                    = errors.New("could not contact dns")
	ErrInvalidMasterKey         = errors.New("invalid master key")
	ErrInvalidMasterFingerprint = errors.New("invalid fingerprint of master")
	ErrBadIndexSignature        = errors.New("bad signature of index")
)
