package out

import (
	"context"
	"fmt"
	"io"
	"os"
)

type OutputContext struct {
	base   *OutputContext
	stdout io.Writer
	stderr io.Writer
}

func New(out io.Writer, err io.Writer) *OutputContext {
	return &OutputContext{
		stdout: out,
		stderr: err,
	}
}

func (o *OutputContext) Printf(msg string, args ...interface{}) (int, error) {
	return fmt.Fprintf(o.Stdout(), msg, args...)
}

func (o *OutputContext) Print(args ...interface{}) (int, error) {
	return fmt.Fprint(o.Stdout(), args...)
}

func (o *OutputContext) Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(o.Stdout(), args...)
}

func (o *OutputContext) ErrPrintf(msg string, args ...interface{}) (int, error) {
	return fmt.Fprintf(o.Stderr(), msg, args...)
}

func (o *OutputContext) ErrPrint(args ...interface{}) (int, error) {
	return fmt.Fprint(o.Stderr(), args...)
}

func (o *OutputContext) ErrPrintln(args ...interface{}) (int, error) {
	return fmt.Fprintln(o.Stderr(), args...)
}

func (o *OutputContext) Stdout() io.Writer {
	if o.stdout != nil {
		return o.stdout
	}
	if o.base != nil {
		return o.base.Stdout()
	}
	return os.Stdout
}

func (o *OutputContext) Stderr() io.Writer {
	if o.stderr != nil {
		return o.stderr
	}
	if o.base != nil {
		return o.base.Stderr()
	}
	return os.Stderr
}

var def = New(os.Stdout, os.Stderr)

func With(ctx context.Context, o *OutputContext) context.Context {
	o.base = Get(ctx)
	return context.WithValue(ctx, "stdout", o)
}

func Get(ctx context.Context) *OutputContext {
	o := ctx.Value("stdout")
	if o == nil {
		return def
	}
	return o.(*OutputContext)
}

func Printf(ctx context.Context, msg string, args ...interface{}) (int, error) {
	return Get(ctx).Printf(msg, args...)
}

func Print(ctx context.Context, args ...interface{}) (int, error) {
	return Get(ctx).Print(args...)
}

func Println(ctx context.Context, args ...interface{}) (int, error) {
	return Get(ctx).Println(args...)
}

////////////////////////////////////////////////////////////////////////////////

func ErrPrintf(ctx context.Context, msg string, args ...interface{}) (int, error) {
	return Get(ctx).ErrPrintf(msg, args...)
}

func ErrPrint(ctx context.Context, args ...interface{}) (int, error) {
	return Get(ctx).ErrPrint(args...)
}

func ErrPrintln(ctx context.Context, args ...interface{}) (int, error) {
	return Get(ctx).ErrPrintln(args...)
}
