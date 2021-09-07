/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"fmt"
	"jacobin/log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Classloaders hold the parsed bytecode in classes, where they can be retrieved
// and moved to an execution role. Most of the comments and code presuppose some
// familiarity with the role of classloaders. More information can be found at:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-5.html#jvms-5.3
type Classloader struct {
	Name    string
	Parent  string
	Classes map[string]parsedClass
}

var AppCL Classloader
var BootstrapCL Classloader
var ExtensionCL Classloader

// the parsed class
type parsedClass struct {
	javaVersion    int
	className      string   // name of class without path and without .class
	superClass     string   // name of superclass for this class
	interfaceCount int      // number of interfaces this class implements
	interfaces     []string // the interfaces this class implements

	// ---- constant pool data items ----
	cpCount       int             // count of constant pool entries
	cpIndex       []cpEntry       // the constant pool index to entries
	classRefs     []classRefEntry // this and next slices hold CP entries
	fieldRefs     []fieldRefEntry
	intConsts     []intConst
	interfaceRefs []interfaceRefEntry
	methodRefs    []methodRefEntry
	nameAndTypes  []nameAndTypeEntry
	stringRefs    []stringConstantEntry
	utf8Refs      []utf8Entry

	// ---- access flags items ----
	accessFlags       int // the following booleans interpret the access flags
	classIsPublic     bool
	classIsFinal      bool
	classIsSuper      bool
	classIsInterface  bool
	classIsAbstract   bool
	classIsSynthetic  bool
	classIsAnnotation bool
	classIsEnum       bool
	classIsModule     bool
}

// cfe = class format error, which is the error thrown by the parser for most
// of the errors arising from malformed bytecode. Prints out file and line# where
// the call to cfe() occurred.
func cfe(msg string) error {
	errMsg := "Class Format Error: " + msg

	// get the filename and line# of the function where the error occurred
	// implementation note: Caller(0) would be this function. (1) is the
	// previous function on the stack (so, the one calling this error routine)
	// To traverse all the way back to the start of the program, set up a loop
	// and exit when ok is no longer true.
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		fileName, fileLine := fn.FileLine(pc)
		errMsg = errMsg + "\n  dectected by file: " + filepath.Base(fileName) +
			", line: " + strconv.Itoa(fileLine)
	}
	log.Log(errMsg, log.SEVERE)
	return errors.New(errMsg)
}

// first canonicalizes the filename, checks whether the class is already loaded,
// and if not, then parses the class and loads it.
/*
// 1 TODO: canonicalize class name
// 2 TODO: search through classloaders for this class
// 3 TODO: determine which classloader should load the class, then
// 4 TODO: have *it* parse and load the class.
*/
func (cl Classloader) LoadClassFromFile(filename string) error {
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Log("Could not read file: "+filename+". Exiting.", log.SEVERE)
		return fmt.Errorf("file I/O error")
	}

	log.Log(filename+" read", log.FINE)

	fullyParsedClass, err := parse(rawBytes)
	if err != nil {
		log.Log("error parsing "+filename+". Exiting.", log.SEVERE)
		return fmt.Errorf("parsing error")
	}

	return insert(fullyParsedClass)

}

// Init() simply initializes the three classloaders and points them to each other
// in the proper order. This function might be substantially revised later.
func Init() error {
	BootstrapCL.Name = "bootstrap"
	BootstrapCL.Parent = ""
	BootstrapCL.Classes = make(map[string]parsedClass)

	ExtensionCL.Name = "extension"
	ExtensionCL.Parent = "bootstrap"
	ExtensionCL.Classes = make(map[string]parsedClass)

	AppCL.Name = "app"
	AppCL.Parent = "system"
	AppCL.Classes = make(map[string]parsedClass)
	return nil
}

// insert the fully parsed class into the classloader
func insert(class parsedClass) error {
	return nil //TODO: fill out after finishing parser
}