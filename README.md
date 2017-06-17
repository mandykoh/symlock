# symlock

[![GoDoc](https://godoc.org/github.com/mandykoh/symlock?status.svg)](https://godoc.org/github.com/mandykoh/symlock)
[![Go Report Card](https://goreportcard.com/badge/github.com/mandykoh/symlock)](https://goreportcard.com/report/github.com/mandykoh/symlock)
[![Build Status](https://travis-ci.org/mandykoh/symlock.svg?branch=master)](https://travis-ci.org/mandykoh/symlock)

A symbolic lock implementation for Go.

## Introduction

SymLocks provide mutual exclusion on a string value rather than a specific lock object. This can be useful in situations where the complete set of mutex points isn’t immediately known, or may be too large for up front setup to be practical.


## Example

Use a SymLock like this:

```go
var s symlock.SymLock

s.WithMutex("some string value symbolising a mutex point", func() {
    // Do some things which require mutexing on something represented by the provided string
})
```

All code using the same SymLock with the same string will be mutexed from each other:

```go
go s.WithMutex("apple", func() {
    // Do some stuff
})

go s.WithMutex("apple", func() {
    // Do some more stuff, with mutually exclusive locking from the above
})

go s.WithMutex("pear", func() {
    // Do some stuff that can run concurrently with code using "apple" as
    // the mutex symbol, but not with other code that might be using "pear"
})
```
