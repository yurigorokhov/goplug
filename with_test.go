package goplug_test

import (
	plug "github.com/yurigorokhov/goplug"
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type GoPlugSuite struct{}

var _ = Suite(&GoPlugSuite{})

func (s *GoPlugSuite) With_adds_parameters(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p.
		With("param1", "value1").
		With("param2", "value2")
	c.Assert(p.GetParam("param1"), Equals, "value1")
	c.Assert(p.GetParam("param2"), Equals, "value2")
}
