//go:build reckless

package options

// Reckless being true allows the build to do things like turning off
// synchronisation of global variables, and assume more undefined behaviour.
//
// This value is set at compile-time and is declared as a const to allow
// the compiler to remove dead code.
const Reckless = true
