package uksuid

import (
	"log"
	"testing"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/gocql/gocql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKSUID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KSUID Suite")
}

// ksuidType is a mock gocql-type for testing CQL marshalling/Unmarshalling
type ksuidType struct{}

func (k ksuidType) Type() gocql.Type {
	return gocql.TypeCustom
}

func (k ksuidType) Version() byte {
	return 1
}
func (k ksuidType) Custom() string {
	return "custom"
}

func (k ksuidType) New() interface{} {
	return ksuidType{}
}

var _ = Describe("KSUID", func() {
	// Taken from: https://github.com/segmentio/ksuid/blob/master/ksuid.go
	// KSUID's epoch starts more recently so that the 32-bit number space gives a
	// significantly higher useful lifetime of around 136 years from March 2017.
	// This number (14e8) was picked to be easy to remember.
	var epochStamp uint32 = 1400000000

	Context("new KSUID is requested", func() {
		It("should return new KSUID", func() {
			u, err := New()
			Expect(err).ToNot(HaveOccurred())

			uid := u.String()
			Expect(uid).ToNot(BeEmpty())

			_, err = ksuid.Parse(uid)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("new KSUID WithTime is requested", func() {
		It("should return new KSUID", func() {
			uidTime := time.Now().AddDate(0, -2, -5)
			u, err := NewWithTime(uidTime)
			Expect(err).ToNot(HaveOccurred())

			log.Println(u.Timestamp())

			uid := u.String()
			Expect(uid).ToNot(BeEmpty())

			kuid, err := ksuid.Parse(uid)
			Expect(err).ToNot(HaveOccurred())

			t := int64(kuid.Timestamp() + epochStamp)
			Expect(t).To(Equal(uidTime.Unix()))
		})
	})

	Describe("Marshalling", func() {
		It("should marshal to CQL when time is not specified", func() {
			uid, err := New()
			Expect(err).ToNot(HaveOccurred())

			marshalled, err := gocql.Marshal(ksuidType{}, uid)
			Expect(err).ToNot(HaveOccurred())

			newKSUID, err := ksuid.FromBytes(marshalled)
			Expect(err).ToNot(HaveOccurred())
			Expect(newKSUID.String()).To(Equal(uid.String()))
		})

		It("should marshal to CQL when time is specified", func() {
			uidTime := time.Now().AddDate(0, -9, 4)
			uid, err := NewWithTime(uidTime)
			Expect(err).ToNot(HaveOccurred())

			marshalled, err := gocql.Marshal(ksuidType{}, uid)
			Expect(err).ToNot(HaveOccurred())

			newKSUID, err := ksuid.FromBytes(marshalled)
			Expect(err).ToNot(HaveOccurred())
			Expect(newKSUID.String()).To(Equal(uid.String()))

			t := int64(newKSUID.Timestamp() + epochStamp)
			Expect(t).To(Equal(uidTime.Unix()))
		})
	})

	Describe("Unmarshalling", func() {
		It("should unmarshal to CQL when time is not specified", func() {
			uid, err := New()
			Expect(err).ToNot(HaveOccurred())

			uidBytes := uid.Bytes()

			unmarshal := KSUID{}
			err = gocql.Unmarshal(ksuidType{}, uidBytes, &unmarshal)
			Expect(err).ToNot(HaveOccurred())
			Expect(unmarshal.String()).To(Equal(uid.String()))
		})

		It("should unmarshal to CQL when time not specified", func() {
			uidTime := time.Now().AddDate(0, 10, 4)
			uid, err := NewWithTime(uidTime)
			Expect(err).ToNot(HaveOccurred())

			uidBytes := uid.Bytes()

			unmarshal := KSUID{}
			err = gocql.Unmarshal(ksuidType{}, uidBytes, &unmarshal)
			Expect(err).ToNot(HaveOccurred())
			Expect(unmarshal.String()).To(Equal(uid.String()))
			Expect(unmarshal.Time().UnixNano()).To(Equal(uid.Time().UnixNano()))
		})
	})

	Describe("Timestamp", func() {
		It("should return corrected timestamp", func() {
			t := time.Now()
			ks, err := NewWithTime(t)
			Expect(err).ToNot(HaveOccurred())

			kt := int64(ks.Timestamp())
			Expect(kt).To(Equal(t.Unix()))
		})
	})

	Describe("TimestampUncorrected", func() {
		It("should return uncorrected timestamp", func() {
			t := time.Now()
			ks, err := NewWithTime(t)
			Expect(err).ToNot(HaveOccurred())

			kt := int64(ks.TimestampUncorrected())
			ut := t.Unix() - int64(epochStamp)
			Expect(kt).To(Equal(ut))
		})
	})

	Context("Sort operation is called", func() {
		It("should sort provided values", func() {
			ks := []KSUID{}
			for i := 0; i < 10; i++ {
				k, err := New()
				Expect(err).ToNot(HaveOccurred())
				ks = append(ks, k)
			}
			Expect(IsSorted(ks)).To(BeFalse())

			sorted := Sort(ks)
			Expect(IsSorted(sorted)).To(BeTrue())
		})
	})

	Context("parsing to KSUID", func() {
		Describe("FromBytes", func() {
			It("should parse bytes to KSUID", func() {
				u, err := New()
				Expect(err).ToNot(HaveOccurred())
				bytes := u.Bytes()

				uid, err := FromBytes(bytes)
				Expect(err).ToNot(HaveOccurred())
				_, err = ksuid.FromBytes(uid.Bytes())
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return any errors that occur", func() {
				invalidKSUID := "invalid"

				_, err := FromBytes([]byte(invalidKSUID))
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("FromString", func() {
			It("should parse string to KSUID", func() {
				u, err := New()
				Expect(err).ToNot(HaveOccurred())
				str := u.String()

				uid, err := FromString(str)
				Expect(uid).ToNot(BeNil())
				_, err = ksuid.Parse(uid.String())
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return any errors that occur", func() {
				invalidKSUID := "invalid"

				_, err := FromString(invalidKSUID)
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("FromParts", func() {
			It("should parse string to KSUID", func() {
				uidTime := time.Now().AddDate(0, 10, 4)

				// 16 bytes
				payload := []byte{
					229, 69, 95, 119, 215, 29, 196, 255, 186, 191, 89, 24, 30, 239, 132, 38,
				}
				uid, err := FromParts(uidTime, payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(uid.Payload()).To(Equal(payload))
				// Expect(uid.Time)
			})
		})

		Describe("ToWrappedKSUID", func() {
			It("should convert provided KSUID to ksuid.KSUID", func() {
				ks := []KSUID{}
				for i := 0; i < 10; i++ {
					k, err := New()
					Expect(err).ToNot(HaveOccurred())
					ks = append(ks, k)
				}

				ksids := ToWrappedKSUID(ks...)
				for i, ksid := range ksids {
					Expect(ksid).To(Equal(ks[i].KSUID))
				}
			})
		})

		Describe("ToWrapperKSUID", func() {
			It("should convert provided ksuid.KSUID to KSUID", func() {
				ks := []ksuid.KSUID{}
				for i := 0; i < 10; i++ {
					k, err := ksuid.NewRandom()
					Expect(err).ToNot(HaveOccurred())
					ks = append(ks, k)
				}

				ids := ToWrapperKSUID(ks...)
				for i, id := range ids {
					Expect(id.KSUID).To(Equal(ks[i]))
				}
			})
		})
	})
})
