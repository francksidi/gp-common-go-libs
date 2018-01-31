package structmatcher_test

import (
	"github.com/greenplum-db/gp-common-go-libs/structmatcher"

	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStructMatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "structmatcher tests")
}

var _ = Describe("structmatcher.StructMatchers", func() {
	type SimpleStruct struct {
		Field1 int
		Field2 string
	}
	type NestedStruct struct {
		Field1      int
		Field2      string
		NestedField []SimpleStruct
	}
	Describe("structmatcher.StructMatcher", func() {
		It("returns no failures for the same structs", func() {
			struct1 := SimpleStruct{Field1: 0, Field2: "test_schema"}
			struct2 := SimpleStruct{Field1: 0, Field2: "test_schema"}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, false, false)
			Expect(mismatches).To(BeEmpty())
		})
		It("returns mismatches with different structs", func() {
			struct1 := SimpleStruct{Field1: 0, Field2: "test_schema"}
			struct2 := SimpleStruct{Field1: 0, Field2: ""}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, false, false)
			Expect(mismatches).ToNot(BeEmpty())
		})
		It("returns mismatches in nested structs", func() {
			struct1 := NestedStruct{Field1: 0, Field2: "testrole", NestedField: []SimpleStruct{{Field1: 3}}}
			struct2 := NestedStruct{Field1: 0, Field2: "testrole", NestedField: []SimpleStruct{{Field1: 4}}}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, false, false)
			Expect(len(mismatches)).To(Equal(1))
			Expect(mismatches[0]).To(Equal("Mismatch on field Field1\nExpected\n    <int>: 3\nto equal\n    <int>: 4"))
		})
		It("returns mismatches including struct fields", func() {
			struct1 := NestedStruct{Field1: 0, Field2: "testrole", NestedField: []SimpleStruct{{Field1: 3}}}
			struct2 := NestedStruct{Field1: 0, Field2: "teststruct2", NestedField: []SimpleStruct{{Field1: 4}}}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, true, true, "Field2")
			Expect(len(mismatches)).To(Equal(1))
			Expect(mismatches[0]).To(Equal("Mismatch on field Field2\nExpected\n    <string>: testrole\nto equal\n    <string>: teststruct2"))
		})
		It("returns mismatches including nested struct fields", func() {
			struct1 := NestedStruct{Field1: 0, Field2: "testrole", NestedField: []SimpleStruct{{Field1: 3}}}
			struct2 := NestedStruct{Field1: 0, Field2: "teststruct2", NestedField: []SimpleStruct{{Field1: 4}}}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, true, true, "NestedField.Field1")
			Expect(len(mismatches)).To(Equal(1))
			Expect(mismatches[0]).To(Equal("Mismatch on field Field1\nExpected\n    <int>: 3\nto equal\n    <int>: 4"))
		})
		It("returns mismatches excluding struct fields", func() {
			struct1 := NestedStruct{Field1: 0, Field2: "testrole", NestedField: []SimpleStruct{{Field1: 3}}}
			struct2 := NestedStruct{Field1: 0, Field2: "teststruct2", NestedField: []SimpleStruct{{Field1: 4}}}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, true, false, "Field2")
			Expect(len(mismatches)).To(Equal(1))
			Expect(mismatches[0]).To(Equal("Mismatch on field Field1\nExpected\n    <int>: 3\nto equal\n    <int>: 4"))
		})
		It("returns mismatches excluding nested struct fields", func() {
			struct1 := NestedStruct{Field1: 0, Field2: "testrole", NestedField: []SimpleStruct{{Field1: 3}}}
			struct2 := NestedStruct{Field1: 0, Field2: "teststruct2", NestedField: []SimpleStruct{{Field1: 4}}}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, true, false, "NestedField.Field1")
			Expect(len(mismatches)).To(Equal(1))
			Expect(mismatches[0]).To(Equal("Mismatch on field Field2\nExpected\n    <string>: testrole\nto equal\n    <string>: teststruct2"))
		})
		It("formats a nice error message for mismatches", func() {
			struct1 := SimpleStruct{Field1: 0, Field2: "test_schema"}
			struct2 := SimpleStruct{Field1: 0, Field2: "another_schema"}
			mismatches := structmatcher.StructMatcher(&struct1, &struct2, false, false)
			Expect(mismatches).To(Equal([]string{"Mismatch on field Field2\nExpected\n    <string>: test_schema\nto equal\n    <string>: another_schema"}))
		})
	})
})
