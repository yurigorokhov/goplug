package goplug_test

import (
	plug "github.com/yurigorokhov/goplug"
	. "gopkg.in/check.v1"
)

func (s *GoPlugSuite) Test_atPath(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p.AtPath("a/b/c/d")
	c.Assert(p.String(), Equals, "http://test.dev/a/b/c/d")
}

func (s *GoPlugSuite) Test_atPath_starting_with_forward_slash(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p.AtPath("/a/b/c/d")
	c.Assert(p.String(), Equals, "http://test.dev/a/b/c/d")
}
