package gomock

import (
	"reflect"
)

// TestingT matches the minimal subset of *testing.T used by gomock.
type TestingT interface {
	Helper()
	Fatalf(format string, args ...interface{})
}

// Controller coordinates recorded expectations and actual calls.
type Controller struct {
	T     TestingT
	calls map[string][]*Call
}

// NewController creates a new Controller bound to the provided test helper.
func NewController(t TestingT) *Controller {
	return &Controller{
		T:     t,
		calls: make(map[string][]*Call),
	}
}

// Finish asserts that all expected calls were satisfied.
func (c *Controller) Finish() {
	c.T.Helper()
	for method, calls := range c.calls {
		if len(calls) > 0 {
			c.T.Fatalf("expected call to %s not made", method)
		}
	}
}

// RecordCallWithMethodType stores an expected call for later matching.
func (c *Controller) RecordCallWithMethodType(receiver interface{}, method string, methodType reflect.Type, args ...interface{}) *Call {
	c.T.Helper()
	call := &Call{ctrl: c, method: method, args: args}
	c.calls[method] = append(c.calls[method], call)
	return call
}

// Call matches the incoming invocation to recorded expectations.
func (c *Controller) Call(receiver interface{}, method string, args ...interface{}) []interface{} {
	c.T.Helper()
	queued := c.calls[method]
	if len(queued) == 0 {
		c.T.Fatalf("unexpected call to %s", method)
	}
	call := queued[0]
	c.calls[method] = queued[1:]

	if len(call.args) != len(args) {
		c.T.Fatalf("unexpected number of arguments for %s: got %d want %d", method, len(args), len(call.args))
	}
	for i, exp := range call.args {
		if matcher, ok := exp.(Matcher); ok {
			if !matcher.Matches(args[i]) {
				c.T.Fatalf("argument %d for %s does not match: got %v want %s", i, method, args[i], matcher.String())
			}
			continue
		}
		if !reflect.DeepEqual(exp, args[i]) {
			c.T.Fatalf("argument %d for %s does not match: got %v want %v", i, method, args[i], exp)
		}
	}

	if call.doFunc.IsValid() {
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = reflect.ValueOf(arg)
		}
		out := call.doFunc.Call(in)
		res := make([]interface{}, len(out))
		for i, rv := range out {
			res[i] = rv.Interface()
		}
		return res
	}

	return call.returns
}

// Matcher allows custom argument matching.
type Matcher interface {
	Matches(x interface{}) bool
	String() string
}

// Any matches any value.
type anyMatcher struct{}

func (anyMatcher) Matches(interface{}) bool { return true }
func (anyMatcher) String() string           { return "is anything" }

// Any returns a matcher that accepts any value.
func Any() Matcher { return anyMatcher{} }

// Call holds a single expected invocation.
type Call struct {
	ctrl    *Controller
	method  string
	args    []interface{}
	returns []interface{}
	doFunc  reflect.Value
}

// Return configures the return values for the call.
func (c *Call) Return(values ...interface{}) *Call {
	c.returns = values
	return c
}

// DoAndReturn configures a function to execute when the call is matched.
func (c *Call) DoAndReturn(fn interface{}) *Call {
	rv := reflect.ValueOf(fn)
	if rv.Kind() != reflect.Func {
		c.ctrl.T.Fatalf("DoAndReturn requires a function")
	}
	c.doFunc = rv
	return c
}
