/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import (
	"sync"
)

// MTable is the table in which method data is stored for quick reference at
// method invocation. It consists of a map whose key is a string consisting of a
// concatenation of the class name, method name, and method type. For example,
//		java/io/PrintStream.println(Ljava/lang/String;)V
//
// The value consists of a byte identifying whether the method is a Java method
// ('J'), that is, a method that is executed by executing bytecodes, or a go
// method ('G"), which is a Go funciton that is being used as a stand-in for
// the named Java method. In most contexts, this would be called a native method,
// but that term is used in a different context in Java (see JNI), so avoided here.
//
// The second field in the value is an empty interface, which is Go's way of
// implementing generics. Ultimately, this mechanism supports two types of entries--
// one for each kind of method.
//
// When a function is invoked, the lookup mechanism first checks the MTable, and
// if an entry is found, that entry is what is executed. If no entry is found,
// the search goes to the class and faiing that to the superclass, etc. Once the
// method is located it's added to the MTable so that all future invocations will
// result in fast look-ups in the MTable.
var MTable = make(map[string]MTentry)

// MTentry is described in detail in the comments to MTable
type MTentry struct {
	meth  mData // the method data
	mType byte  // method type, G = Go method, J = Java method
}

type mData interface{}

// GmEntry is the entry in the MTable for Go functions. See MTable comments for details.
type GmEntry struct {
	ParamSlots int
	Fu         func([]interface{})
}

// Function is the generic-style function used for Go entries: a function that accepts a
// slice of empty interfaces and returns nothing (b/c all returns are pushed onto the
// stack rather than actually returned to a caller.
type Function func([]interface{})

// MTmutex is used for updates to the MTable because multiple threads could be
// updating it simultaneously.
var MTmutex sync.Mutex

// MTableLoadNatives loads the Go methods from files that contain them. It does this
// by calling the Load_* function in each of those files to load whatever Go functions
// they make available.
func MTableLoadNatives() {
	loadlib(Load_System_PrintStream())

}

func loadlib(libMeths map[string]GMeth) {
	for key, val := range libMeths {
		gme := GmEntry{}
		gme.ParamSlots = val.ParamSlots
		gme.Fu = val.GFunction

		tableEntry := MTentry{
			mType: 'G',
			meth:  gme,
		}

		MTmutex.Lock()
		MTable[key] = tableEntry
		MTmutex.Unlock()
	}
}
