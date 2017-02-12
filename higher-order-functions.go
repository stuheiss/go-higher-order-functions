package main

// Basic higher order functions in go: map, filter, remove foldl, foldr, take, drop
// All functions work with arrays.  Some have variants that work with channels.
// A few have concurrent variants.

// Utility functions provided to convert arrays to channels and vice versa.

// go is not polymorphic. All functions below use the single type T. This is a
// bit of a cheat as some functions should really have multiple types.

// Set type T to your preference or replace T in a given function with the type you need.

// This is not idiomatic go. You may find it useful if you prefer functional style.

import (
	"fmt"
	"sync"
)

type T int

// return reversed copy of array of T
func reverse(in []T) []T {
	l := len(in)
	out := make([]T, l)
	for i, v := range in {
		out[(l-1)-i] = v
	}
	return out
}

// send array of T to channel, return channel
func to_chan(in []T) <-chan T {
	out := make(chan T)
	go func() {
		for _, n := range in {
			out <- n
		}
		close(out)
	}()
	return out
}

// read array of T from channel, return array
func from_chan(in <-chan T) []T {
	out := make([]T, 0)
	for n := range in {
		out = append(out, n)
	}
	return out
}

// map
func mapT(f func(T) T, from []T) []T {
	to := make([]T, len(from))
	for i, v := range from {
		to[i] = f(v)
	}
	return to
}

// parallel map
func pmapT(f func(T) T, from []T) []T {
	N := len(from)
	to := make([]T, N)
	var wg sync.WaitGroup
	wg.Add(N)
	for i, v := range from {
		go func(i int, v T) {
			defer wg.Done()
			to[i] = f(v)
		}(i, v)
	}
	wg.Wait()
	return to
}

// mapchan
func mapchanT(f func(T) T, from <-chan T) chan T {
	to := make(chan T)
	go func() {
		for {
			i, e := <-from
			if e == false {
				close(to)
				break
			}
			to <- f(i)
		}
	}()
	return to
}

// filter
func filterT(f func(T) bool, from []T) []T {
	to := make([]T, 0)
	for _, v := range from {
		if f(v) {
			to = append(to, v)
		}
	}
	return to
}

// filterchan
func filterchanT(f func(T) bool, from <-chan T) <-chan T {
	to := make(chan T)
	go func(to chan T) {
		for n := range from {
			if f(n) {
				to <- n
			}
		}
		close(to)
	}(to)
	return to
}

func removeT(f func(T) bool, from []T) []T {
	to := make([]T, 0)
	for _, n := range from {
		if !f(n) {
			to = append(to, n)
		}
	}
	return to
}

func removechanT(f func(T) bool, from <-chan T) <-chan T {
	to := make(chan T)
	go func(to chan T) {
		for n := range from {
			if !f(n) {
				to <- n
			}
		}
		close(to)
	}(to)
	return to
}

func take(n int, from []T) []T {
	to := make([]T, 0)
	for _, v := range from {
		if n <= 0 {
			break
		}
		to = append(to, v)
		n -= 1
	}
	return to
}

func drop(n int, from []T) []T {
	to := make([]T, 0)
	for _, v := range from {
		if n > 0 {
			n -= 1
			continue
		}
		to = append(to, v)
	}
	return to
}

// foldl :: (b -> a -> b) -> b -> [a] -> b
// foldl f z []     = z
// foldl f z (x:xs) = foldl f (f z x) xs
func foldlT(f func(T, T) T, z T, xs []T) T {
	if len(xs) == 0 {
		return z
	} else {
		x := xs[0]
		xs = xs[1:]
		return foldlT(f, f(z, x), xs)
	}
}

// foldr :: (a -> b -> b) -> b -> [a] -> b
// foldr f z []     = z
// foldr f z (x:xs) = f x (foldr f z xs)
func foldrT(f func(T, T) T, z T, xs []T) T {
	if len(xs) == 0 {
		return z
	} else {
		x := xs[0]
		xs = xs[1:]
		return f(x, foldrT(f, z, xs))
	}
}

func main() {
	t := []T{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("dataset", t)
	fmt.Println("to/from channel", from_chan(to_chan(t)))
	fmt.Println("array reverse", reverse(t))
	fmt.Println("array filter < 5", filterT(func(i T) bool { return i < 5 }, t))
	fmt.Println("array filter even", filterT(func(i T) bool { return i%2 == 0 }, t))
	fmt.Println("array remove even", removeT(func(i T) bool { return i%2 == 0 }, t))
	fmt.Println("array take 3", take(3, t))
	fmt.Println("array drop 3", drop(3, t))
	fmt.Println("array map double", mapT(func(i T) T { return i * 2 }, t))
	fmt.Println("array parallel map double", pmapT(func(i T) T { return i * 2 }, t))
	fmt.Println("channel map double", from_chan(mapchanT(func(i T) T { return i * 2 }, to_chan(t))))
	fmt.Println("channel filter odd", from_chan(filterchanT(func(i T) bool { return i%2 != 0 }, to_chan(t))))
	fmt.Println("channel remove odd", from_chan(removechanT(func(i T) bool { return i%2 != 0 }, to_chan(t))))
	fmt.Println("array foldl sum", foldlT(func(x, y T) T { return x + y }, 0, t))
	fmt.Println("array foldl sub", foldlT(func(x, y T) T { return x - y }, 0, t))
	fmt.Println("array foldl mult", foldlT(func(x, y T) T { return x * y }, 1, t))
	fmt.Println("array foldr sum", foldrT(func(x, y T) T { return x + y }, 0, t))
	fmt.Println("array foldr sub", foldrT(func(x, y T) T { return x - y }, 0, t))
	fmt.Println("array foldr mult", foldrT(func(x, y T) T { return x * y }, 1, t))
}
