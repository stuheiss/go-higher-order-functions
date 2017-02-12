# go-higher-order-functions

Basic higher order functions in go: map, filter, remove foldl, foldr, take, drop
All functions work with arrays.  Some have variants that work with channels.
A few have concurrent variants.

Utility functions provided to convert arrays to channels and vice versa.

go is not polymorphic. All functions below use the single type T. This is a
bit of a cheat as some functions should really have multiple types.

Set type T to your preference or replace T in a given function with the type you need.

This is not idiomatic go. You may find it useful if you prefer functional style.
