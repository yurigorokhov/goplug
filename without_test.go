package goplug_test

import (
	plug "github.com/yurigorokhov/goplug"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func (s *GoPlugSuite) Test_without_removes_parameters(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p.
		With("param1", "value1").
		With("param2", "value2").
		Without("param2")
	c.Assert(p.GetParam("param1"), Equals, "value1")
	c.Assert(p.GetParam("param2"), Equals, "")
}
