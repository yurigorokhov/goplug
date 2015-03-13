package goplug_test

import (
	plug "github.com/yurigorokhov/goplug"
	. "gopkg.in/check.v1"
)

func (s *GoPlugSuite) Test_withParams_adds_parameters(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p.WithParams(map[string]string{"param1": "value1", "param2": "value2"})
	c.Assert(p.GetParam("param1"), Equals, "value1")
	c.Assert(p.GetParam("param2"), Equals, "value2")
}
