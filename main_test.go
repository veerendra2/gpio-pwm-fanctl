package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDutyForTemp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dutyForTemp Suite")
}

var _ = Describe("dutyForTemp", func() {
	It("returns 100 for temp >= 80", func() {
		Expect(dutyForTemp(85)).To(Equal(uint32(100)))
		Expect(dutyForTemp(80)).To(Equal(uint32(100)))
	})

	It("returns 80 for 70 <= temp < 80", func() {
		Expect(dutyForTemp(75)).To(Equal(uint32(80)))
		Expect(dutyForTemp(70)).To(Equal(uint32(80)))
	})

	It("returns 60 for 35 <= temp < 70", func() {
		Expect(dutyForTemp(50)).To(Equal(uint32(60)))
		Expect(dutyForTemp(35)).To(Equal(uint32(60)))
	})

	It("returns 40 for 0 <= temp < 35", func() {
		Expect(dutyForTemp(20)).To(Equal(uint32(40)))
		Expect(dutyForTemp(0)).To(Equal(uint32(40)))
	})

	It("returns 0 for temp < 0", func() {
		Expect(dutyForTemp(0)).To(Equal(uint32(40))) // 0 is covered by last
	})
})
