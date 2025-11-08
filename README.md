# Command Line Utility Library

This library provides two support frameworks
- Handling of Option Sets
- Handling of list-based command output using the [streaming library](https://github.com/mandelsoft/streaming).

## Option Sets

An `OptionSet` represents a set of objects implementing the `Options` interface.
Such an object is able to configure a [`pflag.FlagSet`](https://github.com/spf13/pflag) flag set.

Options implementing the `Options` interface are typically implemented in a dedicated 
package providing soe standard functions like `New`and `From`. If you follow
this pattern, options sets can be used as follows.

You can configure sets of options like this

```go
  opts:= flagutils.DefaultOptionSet{}
  opts.Add(otype1.New(), optype2.New(), ...)
```

To add the options to a `pflag.FlagSet` just use

```go
  fs :=  pflag.NewFlagSet("test", pflag.ContinueOnError)
  opts.AddFlags(fs)
```

The option set can be passed around and 
if your code wants to access configured values for a dedicated option type
the appropriate `Options` object can be retrieved by using

```go
  myopt := otype1.From(opts)
  myopt.Value() // or any other method provided by your option type. 
```

A third interface is the `OptionSetProvider` interface. It is used to describe access to
an `OptionSet` for and other kind of object (for example, a command).

Option sets may be cascaded, they again implement the `Options` as well as the `OptionSetProvider`
interface and can therefore
be added to another `OptionSet`.

All involved (transitively) `Options` object can be iterated using the `Options`
method, which is an `iter.Seq[Options]`. An option set just provides access to options.
If it is extendable it should implement the `ExtendableOptionSet` interface.

A default implementation for an option set provided by the type `DefaultOptionSet`.
It also implements the `ExtendableOptionSet`interface and supports the `Add`
method to aggregate `Options`.

With `GetFrom[T]` it is possible to retrieve the option in an option set
implementing the interface `T`. Similarly, `Filter[T]` provides a slice
with all options implementing the interface `T`. `T` might be 
a pointer to a concreate option type (`*otypepkg.Options`), or any interface implemented by an option type.

### Option Completion and Validation

An `Options` object may optionally implement the `Validatable` interface.
If implemented it will be called whenever an option set is validated using
the `flagutils.Validate` function. The `Validate(ctx context.Context, opts OptionSet, v ValidationSet) error`
methods gets access to the used context, the actually validated option set and
a validation set.

The validation set can be used to recursively get access to other validated
options. THis might be required for different correlated options if
the validation of one option requires the state of another option.
The validation set keeps track of already validated options to assure 
that every option is validated only once. Cyclic dependencies among
options should be avoided but do not lead to an error.
The validation set figures out whether an option is already validated or validating
and returns the requested option without further recursive calls.
This way the initial order options are added to the `OptionSet`
determines the order resolution for cyclic validation dependencies.

The same way works a `Finalizable` interface. It can be used to clean up
external state after the processing based on an option set. Finalization
should be done in the opposite order than the validation.
If an `OptionSet` implements the `Validatable` or `Finalizable` interface,
it gets control over the handling of the included option objects.

### Predefined Option Types

Additionally, some common option types are defined.
They follow the standard convention for option objects. Every options object 
is implemented in a separate package, always following the same layout:
- A struct type `Options` defines the option variables for the correlated set of commandline
  options bundled by this option object. It implements the `Options` interface.
- A function `New()` and optionally more special functions or additional
  parameters are provided to create a new Options object. It can then be added to
  a `DefaultOptionSet`.
- Every such `Options`object supports the configuration methods (if it represents
  a single flag)
  - `WithNames` to configure the long and short option names
  - `WithDescription` to configure the option description. The string is
    potentially used as format for an `Sprintf` call fed with option-type
    specific values.
- A function `From(OptionSetProvider) *Options` retrieving the option
  from an option set, if it is available in this set.

#### Closure Option

The package `closure` provides a closure option usable to request recursive processing a list of initial elements (value type `bool`).

Default values:
- *Long Option*: `closure`
- *Short Option*: `c`

Configuration:
- `WithNames(long,short)`
- `WithDescription(desc)`

This option supports the element processing by being able to
provide an `Exploder` for a processing chain.

Therefore, there are two constructors taking some info for generating such
an exploder:
- `New(chain.Exploder)`: The exploder code to use
- `NewByFactory(ExploderFactory)`: a factory able to create an exploder based on other options.

To support types elements for those processing chains, the 
option type is parameterized with the element type.

### Sort Option

The package `sort` provides a sort option usable to request sorting of field-based output.
It accepts a list of sort fields names (value type `[]string`).

Default values:
- *Long Option*: `sort`
- *Short Option*: `s`

Configuration:
- `WithNames(long,short)`
- `WithDescription(desc)`
- `WithComparator(field, cmp)`

This option supports the element processing by being able to
provide an `CompateFunc` for a processing chain to sort elements
offering a field value slice. Field values are always strings.

If a field name is prefixed by `-` the sort order is reversed.
Possible field names are taken from another option in the 
used `OptionSet` offering a field name slice for the state name
FIELD_MODE_SORT. Such a slice is, for example, offered by the 
[output mode option](#output-mode-option) by implementing
the `output.FieldNameProvider` interface.

#### Parallel Option

The package `parallel` provides a parallel option usable to request 
parallel processing of elements with a limited degree of parallelity
(value type `bool`).

Default values:
- *Long Option*: `parallel`
- *Short Option*: `p`

Configuration:
- `WithNames(long,short)`
- `WithDescription(desc)`
- `WithPoolProvider(PoolProvider)`

This option works together with the `pool` package to support parallel
processing with limited degree of parallelity. It again works together
with the streaming package used to for the [list-based processing](#list-based-output).

If enabled, the option object provides access to a processor pool
able to handle processing requests. The used pool provider can be configured
for the option. By default, the `simplepool` provider is used.

#### Output Mode Option

#### Table Output Options

### Option Type Support

## List-Based Output