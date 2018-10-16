package uksuid

import (
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

// codec contains functions for converting to and from KSUID.

// ToWrappedKSUID converts slice of KSUID to slice is wrapped ksuid.KSUID.
func ToWrappedKSUID(ids ...KSUID) []ksuid.KSUID {
	ksids := []ksuid.KSUID{}
	for _, id := range ids {
		ksids = append(ksids, id.KSUID)
	}
	return ksids
}

// ToWrapperKSUID is inverse of ToWrappedKSUID. Converts slice of ksuid.KSUID to wrapper KSUID.
func ToWrapperKSUID(ksids ...ksuid.KSUID) []KSUID {
	ids := []KSUID{}
	for _, ksid := range ksids {
		ids = append(ids, KSUID{ksid})
	}
	return ids
}

// FromBytes returns a KSUID generated from the raw byte slice input.
func FromBytes(input []byte) (KSUID, error) {
	u, err := ksuid.FromBytes(input)
	if err != nil {
		err = errors.Wrap(err, "Error converting bytes to KSUID")
		return KSUID{}, err
	}
	return KSUID{u}, nil
}

// FromString decodes a string-encoded representation of a KSUID object
func FromString(input string) (KSUID, error) {
	u, err := ksuid.Parse(input)
	if err != nil {
		err = errors.Wrap(err, "Error converting string to KSUID")
		return KSUID{}, err
	}
	return KSUID{u}, nil
}

// FromParts returns a KSUID generated from provided time nad payload.
func FromParts(time time.Time, payload []byte) (KSUID, error) {
	u, err := ksuid.FromParts(time, payload)
	if err != nil {
		err = errors.Wrap(err, "Error converting specified parts to KSUID")
		return KSUID{}, err
	}
	return KSUID{u}, nil
}
