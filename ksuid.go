package uksuid

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

// Taken from: https://github.com/segmentio/ksuid/blob/master/ksuid.go
// KSUID's epoch starts more recently so that the 32-bit number space gives a
// significantly higher useful lifetime of around 136 years from March 2017.
// This number (14e8) was picked to be easy to remember.
var epochStamp uint32 = 1400000000

// KSUID is a wrapper around github.com/segmentio/ksuid.
// It implements some convenient Marshalling/Unmarshalling
// interfaces for better compatibility across-databases.
type KSUID struct {
	ksuid.KSUID
}

// MarshalCQL converts the KSUID into GoCql-compatible []byte.
func (k KSUID) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return k.MarshalBinary()
}

// UnmarshalCQL converts GoCQL bytes to local KSUID.
func (k *KSUID) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return k.UnmarshalBinary(data)
}

// Timestamp returns timestamp portion of the ID as a bare integer,
// which is corrected for KSUID's special epoch.
// Timestamp is in Milliseconds.
func (k KSUID) Timestamp() uint32 {
	t := k.TimestampUncorrected()
	return t + epochStamp
}

// TimestampUncorrected returns timestamp portion of the ID as a bare integer,
// which is uncorrected for KSUID's special epoch.
// Timestamp is in Milliseconds.
func (k KSUID) TimestampUncorrected() uint32 {
	return k.KSUID.Timestamp()
}

// New creates a new random KSUID using current time.
func New() (KSUID, error) {
	u, err := ksuid.NewRandom()
	if err != nil {
		err = errors.Wrap(err, "Error generating KSUID")
		return KSUID{}, err
	}

	return KSUID{u}, nil
}

// NewWithTime creates a new random KSUID using specified time.
func NewWithTime(time time.Time) (KSUID, error) {
	u, err := ksuid.NewRandomWithTime(time)
	if err != nil {
		err = errors.Wrapf(
			err,
			"Error generating KSUID with Time: %s", time.String(),
		)
		return KSUID{}, err
	}
	return KSUID{u}, nil
}

// Compare implements comparison for KSUID type.
// The comparison is done using bytes.Compare.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a KSUID, b KSUID) int {
	return ksuid.Compare(a.KSUID, b.KSUID)
}

// IsSorted checks whether a slice of KSUIDs is sorted
func IsSorted(ids []KSUID) bool {
	return ksuid.IsSorted(ToWrappedKSUID(ids...))
}

// Sort returns sorted KSUID
func Sort(ids []KSUID) []KSUID {
	ksids := ToWrappedKSUID(ids...)
	ksuid.Sort(ksids)
	ids = ToWrapperKSUID(ksids...)
	return ids
}
