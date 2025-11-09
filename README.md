# Command Line Utility Library

This library provides two support frameworks
- Handling of Option Sets
- Handling of list-based command output using the [streaming library](https://github.com/mandelsoft/streaming).

Under the folder `examples` you can find two complete examples 
describing how to use both library parts in combination:

- [Listing of files](examples/files)
- [Simple Graph Traversal](examples/graph)

## Option Sets

An `OptionSet` represents a set of objects implementing the `Options` interface.
Such an object is able to configure a [`pflag.FlagSet`](https://github.com/spf13/pflag) flag set.

Options implementing the `Options` interface are typically implemented in a dedicated 
package providing soe standard functions like `New`and `From`. If you follow
this pattern, options sets can be used as follows.

You can configure sets of options like this

```go
  opts:= flagutils.DefaultOptionSet{}
  opts.Add(otype1.New(), otype2.New(), ...)
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
They follow the standard convention for option objects. Every `Options` object 
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

It implements the `flagutils.Validatable` interface.

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

It implements the `flagutils.Validatable` and `flagutils.Finalizable`interface.

When finalized, the manged processing pool is closed again.

#### Output Mode Option

The package `output` provides an output mode option usable to request
one of multiple possible output modes (like `-o wide` or `-o tree`)
(value type `string`).

Default values:
- *Long Option*: `mode`
- *Short Option*: `o`

Configuration:
- `WithNames(long,short)`
- `WithDescription(desc)`

The `New` function takes an `output.OutputsFactory` defining the available
output modes (see [list-based-output](#list-based-output)).

It implements the `output.FieldNameProvider` and `flagutils.Validatable` interface.

#### Table Output Options

The package `tableoutput` provides an output mode for [table-based list output](#table-output).

It also offers an `Options` object for configuring the behavior via
command line options.

- list of filtered output column names (value type `[]string`).

  Default values:
  - *Long Option*: `columns`
  - *Short Option*: none

  Configuration:
  - `WithColumnsNames(long,short)`
  - `WithColumnsDescription(desc)`

- request all fields if created with optimized mode. (value type `bool`).

  Default values:
  - *Long Option*: `all-columns`
  - *Short Option*: none

  Configuration:
  - `WithAllColumnsNames(long,short)`
  - `WithAllColumnsDescription(desc)`

### Option Type Support

There are some types supporting the creation of options.

`flagutils.SimpleOption` is a standard implementation for an
`Options` object implementing a single option. It can also be used
to be aggregated to implement a multi-option `Options` object.

It offers a default configuration and the methods to adapt
the names when creating such an `Options` object (`WithNames`and `WithDescription`).

The type `flagutils.VarPFunc[T]` is the type for a function
usable to add a flag to a `pflag.FlagSet` for the value type `T`.
Implementations can be achieved directly from the `pflag.FlagSet`
type, like

```go
(*pflag.FlagSet).StringVarP
```

It used by the simple option to implement the `AddFlags` method.
With `NewSimpleOption[T]` an option for the value type `T` is created,
it uses the type `T` to implicitly determine the flag setter function.
With `NewSimpleOptionWithSetter[T]` the setter can explicitly be given.

## Output destinations

The package `utils.out` offers a simple output redirection bound to a `context.Context`. 

It can be set by `out.With(context.Context, OutputContext)` and retrieved
by `out.Get(context.Context)` It always provides an output context. The default context is reflecting `os.Stdout` and `os.Stderr`.

This package also provides functions for printing using a context, which implicitly evaluate the configured `OutputContext`.

The [outputs](#predefined-output-modes) provided for the [list-based output](#list-based-output) use this functionality 
to support context-specific output redirection.

Output destinations are configured by an `out.OutputContext` object.

## List-Based Output

A common use case for some reporting command line interface 
is to provide commands taking some element specifications and listing
attributes for those elements (see kubectl) potentially with different
output modes. 

The steps required to fulfill this task are always similar:
- first the input specification is mapped to some root elements.
- this initial set is enriched by other objects, for example, following dependencies.
- The elements are mapped out a set of attributes which should be displayed
- And finally, the elements in the given order are formatted to be displayed on the output stream.

This part of the library provides some support for those commands, based
on the [streaming library](https://github.com/mandelsoft/streaming).
This library supports the execution of a processing chain consisting of multiple steps, like mapping elements and substituting elements by a set of other elements.

The basic functionality can be found in package `output`. 
The central interface is `OutputsFactory`, It shields a set of
available output mode described by elements of type `OutputFactory` and
is used to configure an `output.Options` describing the command line flags to select the desired output mode.

An `OutputFactory` is able to create an object implementing the `output.Output` interface. It can be used to process a slice of input elements and provide the desired output.

Additionally, it supports the `output.FieldNameProvider` interface to support the [sort](#sort-option) option. It should at least support the `sort.FIELD_MODE_SORT` and `output.FIELD_MODE_OUTPUT` field mode.
The first one describes the field names and order for the sort step, and 
the second one describes the fields available for the final output formatting.

All those factories get access to the option set used to configure the output on the command line. This way, they can adapt their processing to
the desires of the user.


### Predefined Output Modes

The package provides some default output mode implementations. They
are based on the streaming library used to implement the various steps
required to map the initial input to the final output.

- [*Manifest Outputs*](#manifest-output): Map the elements to a textual format like JSON or YAML.
- [*Table Output*](#table-output): Show the elements as table with a particular column per value field.
- [*Tree Output*](#tree-output): Like a table output but shows the attributes as a tree. This is applicable if selected elements feature dependencies among each other.

Every mode creates a chain of processing steps, potentially influenced by
options of an `OptionSet`. 

Input is an iterator provided by a source object. It is used to feed some
processing steps, which may include
- an *explode* step used to build the transitive closure.
- map the elements to a slice of field values
- sort those elements according to some sor function (provided by the `sort`option)
- and finally, processing the provided elements to generate the desired output

### Table Output

The package `tableoutput` offers an output mode displaying a sequence of elements as table, one column per attribute and one rwo per element.

It offers some [formatting options](#table-output-options).
It is defined by a mapping function able to map elements to a slice of attribute fields. Those elements can then be fed into a sort step (which can be configured by the [`sort`](#sort-option)), if it is present in the given option set.
It also observes the [`closure` option](#closure-option). If required, an appropriate *explode* step is processed before the mapping. 

The mapping can either be defined by directly giving a mapper, or an `output.MappingProvider`, which is able to provide a mapper based on an `OptionSet`. For example, is a transitive  output the path should be added to a *name field value*, but not for non-transitive processing.

A sample output may look like this:

```
MODE       NAME                                 SIZE ERROR
drwxrwxrwx output                               4096 
-rw-rw-rw- output\interface.go                   653 
-rw-rw-rw- output\options.go                    1513 
-rw-rw-rw- output\output.go                      347 
-rw-rw-rw- output\outputs.go                    1375 
-rw-rw-rw- output\utils.go                      1249 
drwxrwxrwx output\internal                         0 
drwxrwxrwx output\manifest                         0 
-rw-rw-rw- output\internal\impl.go              1041 
-rw-rw-rw- output\internal\interface.go         1691 
-rw-rw-rw- output\manifest\factory.go           1162 
-rw-rw-rw- output\manifest\manifest.go          2237 
drwxrwxrwx output\tableoutput                   4096 
drwxrwxrwx output\treeoutput                    4096 
-rw-rw-rw- output\manifest\output.go            1167 
-rw-rw-rw- output\tableoutput\factory.go        2674 
-rw-rw-rw- output\treeoutput\factory.go         3397 
-rw-rw-rw- output\tableoutput\options.go        1760 
-rw-rw-rw- output\treeoutput\output_test.go     2131 
-rw-rw-rw- output\tableoutput\output.go         3573 
-rw-rw-rw- output\treeoutput\suite_test.go       201 
-rw-rw-rw- output\tableoutput\utils.go          2600 
-rw-rw-rw- output\treeoutput\treeoptions.go     2303 
drwxrwxrwx output\treeoutput\test                  0 
drwxrwxrwx output\treeoutput\topo               4096 
-rw-rw-rw- output\treeoutput\test\a                5 
-rw-rw-rw- output\treeoutput\test\b                3 
-rw-rw-rw- output\treeoutput\topo\sort.go       3067 
-rw-rw-rw- output\treeoutput\topo\sort_test.go  2461 
drwxrwxrwx output\treeoutput\test\dir              0 
-rw-rw-rw- output\treeoutput\topo\suite_test.go  183 
-rw-rw-rw- output\treeoutput\test\dir\a            5 
-rw-rw-rw- output\treeoutput\topo\topo.go       1442 
-rw-rw-rw- output\treeoutput\test\dir\c            6 
drwxrwxrwx output\treeoutput\test\dir\sub          0 
-rw-rw-rw- output\treeoutput\test\dir\sub\d        6 
-rw-rw-rw- output\treeoutput\test\dir\sub\e        3 
drwxrwxrwx examples                                0 
drwxrwxrwx examples\graph                          0 
drwxrwxrwx examples\files                          0 
drwxrwxrwx examples\graph\graph                    0 
drwxrwxrwx examples\files\files                    0 
-rw-rw-rw- examples\graph\graph\closure.go      1584 
-rw-rw-rw- examples\files\files\closure.go      2465 
-rw-rw-rw- examples\graph\graph\graph.go        1087 
-rw-rw-rw- examples\files\files\options.go       465 
drwxrwxrwx examples\graph\app                      0 
-rw-rw-rw- examples\graph\graph\outputs.go      1181 
-rw-rw-rw- examples\graph\graph\source.go       2098 
-rw-rw-rw- examples\files\files\outputs.go      1806 
-rw-rw-rw- examples\graph\app\main.go           1549 
-rw-rw-rw- examples\files\files\sort.go          390 
drwxrwxrwx examples\files\app                      0 
-rw-rw-rw- examples\files\files\source.go       2335 
-rw-rw-rw- examples\files\app\main.go           1725 
processed 55 files
```

### Manifest Output

The package `manifest` offers an output mode displaying a sequence of elements as textual structured data, like JSON or YAML.
Therefore, the elements must implement the `Manifest` interface.

The elements are mapped to `Manifest` providing objects, which are then passed to a formatter for the final output.

Before this mapping, optionally the [`closure` option](#closure-option)
is observed to enrich the chain by an appropriate *explode* step.

The package provides formatters for JSON and YAML. 

With the function `AddManifestOutputs` the known modes can be added to an existing `OutputsFactory`:
- `json` compressed JSON
- `JSON` pretty printed JSON
- `yaml` elements as a YAML list
- `YAML` elements as a sequence of YAML documents.


### Tree Output

The package `manifest` offers an output mode displaying a sequence of elements as a table of attributes preceeded with a column visualizing a tree structure. This visualization if generated using  the `tree` package.
It is able to map a sequence of elements providing some tree-relevant
information, like the nesting hierarchy to a sequence `tree.TreeObject`.

The inbound elements provided by the element source must implement the 
`treeoutput.Element` interface, providing some standard node information and topology information.

This sequence is handled by a [`table output`](#table-output) with some
intermediate processing steps, doing
- a topological sort, observing the order of elements on every level as found in the inbound sequence.
- a mapping to elements providing visualization information
- and a mapping to enriched field value slices 

The last step is then fed into the table output.

The topological sort is defined by a compare function created
by a `topo.ComparerFactory`.
The sub package `topo` provides a standard implementation by providing
a factory for creating such a comparer based on an initial sorting order
( `topo.NewDefaultComparerFactory`).

An example how to use it can be found in [`exampled/graph`](examples/graph).

An output of a table output could look like this:

```
            NAME VALUE  ERROR
└─ ⊗        c    charly 
   ├─ ⊗     e    eve    
   │  └─    ...         already shown
   ├─ ⊗     b    bob    
   │  ├─    d    david  
   │  └─ ⊗  c    charly 
   │     └─ ...         cycle
   └─ ⊗     a    alice  
      ├─ ⊗  e    eve    
      │  └─ d    david  
      └─    d    david  
processed 11 nodes
```