package builtin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestContext(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}

type ContextSuite struct {
	suite.Suite
}

var (
	_empty1 = struct{}{}
	_empty2 = struct{}{}
)

func (su *ContextSuite) TestEmptyStructKey() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, _empty1, "Key1")
	ctx = context.WithValue(ctx, _empty2, "Replaced")

	su.EqualValues(ctx.Value(_empty1), "Replaced")
	su.EqualValues(ctx.Value(_empty2), "Replaced")
}

type (
	definedStruct1 struct{}
	definedStruct2 struct{}
)

var (
	_defined1 = definedStruct1{}
	_defined2 = definedStruct2{}
)

func (su *ContextSuite) TestDefinedStructKey() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, _defined1, "Key1")
	ctx = context.WithValue(ctx, _defined2, "Key2")

	su.EqualValues(ctx.Value(_defined1), "Key1")
	su.EqualValues(ctx.Value(_defined2), "Key2")
}
