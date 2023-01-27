/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/frames"
	"jacobin/globals"
	"testing"
)

func TestJdkArrayTypeToJacobinType(t *testing.T) {

	a := jdkArrayTypeToJacobinType(T_BOOLEAN)
	if a != BYTE {
		t.Errorf("Expected Jacobin type of %d, got: %d", BYTE, a)
	}

	b := jdkArrayTypeToJacobinType(T_CHAR)
	if b != INT {
		t.Errorf("Expected Jacobin type of %d, got: %d", INT, b)
	}

	c := jdkArrayTypeToJacobinType(T_DOUBLE)
	if c != FLOAT {
		t.Errorf("Expected Jacobin type of %d, got: %d", FLOAT, c)
	}

	d := jdkArrayTypeToJacobinType(999)
	if d != ERROR {
		t.Errorf("Expected Jacobin type of %d, got: %d", ERROR, d)
	}
}

// ARRAYLENGTH
// First, we create the array of 13 elements, then we push the reference
// to it and execute the ARRAYLENGTH bytecode using the address stored
// in the global array address list
func TestByteArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))             // make the array 13 elements big
	f.Meth = append(f.Meth, T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(ARRAYLENGTH)
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 13 {
		t.Errorf("Expecting array length of 13, got %d", size)
	}
}

func TestIntArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(22))            // make the array 22 elements big
	f.Meth = append(f.Meth, T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 22 {
		t.Errorf("Expecting array length of 13, got %d", size)
	}
}

func TestFloatArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(34))               // make the array 34 elements big
	f.Meth = append(f.Meth, T_DOUBLE) // make it an array of doubles

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 34 {
		t.Errorf("Expecting array length of 34, got %d", size)
	}
}

// NEWARRAY: creation of array for primitive values
func TestNewrray(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))             // make the array 13 elements big
	f.Meth = append(f.Meth, T_LONG) // make it an array of longs

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, test the length of the array, which should be 13
	element := g.ArrayAddressList.Front()
	ptr := element.Value.(*JacobinIntArray)
	if len(*ptr.Arr) != 13 {
		t.Errorf("Expecting array length of 13, got %d", len(*ptr.Arr))
	}
}