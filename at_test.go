package goplug_test

import (
	plug "github.com/yurigorokhov/goplug"
	. "gopkg.in/check.v1"
)

func (s *GoPlugSuite) Test_at(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p.At("a", "b", "c", "d")
	c.Assert(p.String(), Equals, "http://test.dev/a/b/c/d")
}

func (s *GoPlugSuite) Test_at_multiple(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p.At("a", "b", "c", "d")
	p.At("e", "f", "g", "h")
	c.Assert(p.String(), Equals, "http://test.dev/a/b/c/d/e/f/g/h")
}
