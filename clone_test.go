package goplug_test

import (
	plug "github.com/yurigorokhov/goplug"
	. "gopkg.in/check.v1"
)

func (s *GoPlugSuite) Test_clone(c *C) {
	p, err := plug.New("http://test.dev")
	c.Assert(err, IsNil)
	p2 := p.Clone()
	p2.At("mypath")
	p.At("a", "b", "c", "d")
	c.Assert(p.String(), Equals, "http://test.dev/a/b/c/d")
	c.Assert(p2.String(), Equals, "http://test.dev/mypath")
}
