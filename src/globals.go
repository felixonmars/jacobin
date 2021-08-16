/* Jacobin VM -- A Java virtual machine
 * (c) Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0
 */

package main

import (
	"time"
)

// Globals contains variables that need to be globally accessible,
// such as VM and program args, pointers to classloaders, etc.
type Globals struct {
	// ---- jacobin version number ----
	// note: all references to version number must come from this literal
	version string
	vmModel string // "client" or "server" (both the same acc. to JVM docs)

	// ---- processing stoppage? ----
	exitNow bool

	// ---- logging items ----
	logLevel  int
	startTime time.Time

	// ---- command-line items ----
	jacobinName string // name of the executing Jacobin executable
	args        []string
	commandLine string

	startingClass string
	startingJar   string
	appArgs       []string
	options       map[string]Option

	// ---- classloading items ----
	/*
		var bootstrapLoader = Classloader( name: "bootstrap", parent: "" )
		var systemLoader    = Classloader( name: "system", parent: "bootstrap" )
		var assertionStatus = true //default assertion status is that assertions are executed. This is only for start-up.
		var verifyBytecode  = verifyLevel.remote
	*/
}

// initialize the global values that are known at start-up
// listed in alpha order after the first two items
func initGlobals(progName string) *Globals {
	globals := new(Globals)
	globals.startTime = time.Now()

	globals.exitNow = false
	globals.jacobinName = progName
	globals.logLevel = WARNING
	globals.options = make(map[string]Option)
	globals.startingClass = ""
	globals.startingJar = ""
	globals.version = "0.1.0"
	globals.vmModel = "server"

	return globals
}

// the value portion of the globals.ptions table. This is described in more detail in
// option_table_loader.go introductory comments
type Option struct {
	supported bool
	set       bool
	argStyle  int16
	action    func(position int, name string) (int, error)
}
