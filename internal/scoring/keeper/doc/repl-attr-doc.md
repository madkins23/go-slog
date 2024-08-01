The `ReplAttr` scoring algorithm is the same as the `Default` algorithm.

The chart compares the behavior of the `madkins/flash` handler
configured with both `flash.Extras` customization
(to mimic a badly behaved handler) and
a set of `ReplaceAttr` functions to remove that customization.
The goal is to compare the performance with and without
the `ReplaceAttr` functions defined in this repository.
