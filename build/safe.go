// +build !reckless

package build

// Reckless being true allows the build to do things like turning off
// synchronisation of global variables, and assume more undefined behaviour.
const Reckless = false
