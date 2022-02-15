/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"jacobin/globals"
	"jacobin/log"
	"strconv"
	"testing"
)

// This test has an entire class (HaveInterface.class) in an array. This class is then parsed in
// this test and a variety of tests about the class verify the successful parse. This class was
// chosen because it's comparatively small and known to be well-formed. The bytes are laid out
// 16 to a row. The source code for the class appears here:
//
// // simple class for testing of parser
// import java.io.Serializable;
// import java.io.ObjectStreamException;
// import java.io.IOException;
// import java.io.InvalidClassException;
// import java.lang.Runnable;
//
// public class HaveInterface implements Serializable, Runnable {
//
//   private void writeObject(java.io.ObjectOutputStream out) throws IOException {
//      throw new IOException();
//   }
//
//   private void readObject(java.io.ObjectInputStream in) throws IOException, ClassNotFoundException {
// 	    var i = 3;
// 	    if( i > 2)
// 		    throw new IOException();
// 	    else
// 		    throw new ClassNotFoundException();
//   }
//
//   public void run(){
//      return;
//   }
//
//   private void readObjectNoData() throws ObjectStreamException {
// 	        throw new InvalidClassException("test");
//   }
// }

func TestASimpleValidClass(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	classBytes := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x37, 0x00, 0x28, 0x0A, 0x00, 0x0A, 0x00, 0x1C, 0x07,
		0x00, 0x1D, 0x0A, 0x00, 0x02, 0x00, 0x1C, 0x07, 0x00, 0x1E, 0x0A, 0x00, 0x04, 0x00, 0x1C, 0x07,
		0x00, 0x1F, 0x08, 0x00, 0x20, 0x0A, 0x00, 0x06, 0x00, 0x21, 0x07, 0x00, 0x22, 0x07, 0x00, 0x23,
		0x07, 0x00, 0x24, 0x07, 0x00, 0x25, 0x01, 0x00, 0x06, 0x3C, 0x69, 0x6E, 0x69, 0x74, 0x3E, 0x01,
		0x00, 0x03, 0x28, 0x29, 0x56, 0x01, 0x00, 0x04, 0x43, 0x6F, 0x64, 0x65, 0x01, 0x00, 0x0F, 0x4C,
		0x69, 0x6E, 0x65, 0x4E, 0x75, 0x6D, 0x62, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6C, 0x65, 0x01, 0x00,
		0x0B, 0x77, 0x72, 0x69, 0x74, 0x65, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x01, 0x00, 0x1F, 0x28,
		0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F, 0x2F, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x4F,
		0x75, 0x74, 0x70, 0x75, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x3B, 0x29, 0x56, 0x01, 0x00,
		0x0A, 0x45, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6F, 0x6E, 0x73, 0x01, 0x00, 0x0A, 0x72, 0x65,
		0x61, 0x64, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x01, 0x00, 0x1E, 0x28, 0x4C, 0x6A, 0x61, 0x76,
		0x61, 0x2F, 0x69, 0x6F, 0x2F, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x49, 0x6E, 0x70, 0x75, 0x74,
		0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x3B, 0x29, 0x56, 0x01, 0x00, 0x0D, 0x53, 0x74, 0x61, 0x63,
		0x6B, 0x4D, 0x61, 0x70, 0x54, 0x61, 0x62, 0x6C, 0x65, 0x01, 0x00, 0x03, 0x72, 0x75, 0x6E, 0x01,
		0x00, 0x10, 0x72, 0x65, 0x61, 0x64, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x4E, 0x6F, 0x44, 0x61,
		0x74, 0x61, 0x07, 0x00, 0x26, 0x01, 0x00, 0x0A, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x46, 0x69,
		0x6C, 0x65, 0x01, 0x00, 0x12, 0x48, 0x61, 0x76, 0x65, 0x49, 0x6E, 0x74, 0x65, 0x72, 0x66, 0x61,
		0x63, 0x65, 0x2E, 0x6A, 0x61, 0x76, 0x61, 0x0C, 0x00, 0x0D, 0x00, 0x0E, 0x01, 0x00, 0x13, 0x6A,
		0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F, 0x2F, 0x49, 0x4F, 0x45, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69,
		0x6F, 0x6E, 0x01, 0x00, 0x20, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x43,
		0x6C, 0x61, 0x73, 0x73, 0x4E, 0x6F, 0x74, 0x46, 0x6F, 0x75, 0x6E, 0x64, 0x45, 0x78, 0x63, 0x65,
		0x70, 0x74, 0x69, 0x6F, 0x6E, 0x01, 0x00, 0x1D, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F, 0x2F,
		0x49, 0x6E, 0x76, 0x61, 0x6C, 0x69, 0x64, 0x43, 0x6C, 0x61, 0x73, 0x73, 0x45, 0x78, 0x63, 0x65,
		0x70, 0x74, 0x69, 0x6F, 0x6E, 0x01, 0x00, 0x04, 0x74, 0x65, 0x73, 0x74, 0x0C, 0x00, 0x0D, 0x00,
		0x27, 0x01, 0x00, 0x0D, 0x48, 0x61, 0x76, 0x65, 0x49, 0x6E, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
		0x65, 0x01, 0x00, 0x10, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x4F, 0x62,
		0x6A, 0x65, 0x63, 0x74, 0x01, 0x00, 0x14, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F, 0x2F, 0x53,
		0x65, 0x72, 0x69, 0x61, 0x6C, 0x69, 0x7A, 0x61, 0x62, 0x6C, 0x65, 0x01, 0x00, 0x12, 0x6A, 0x61,
		0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x52, 0x75, 0x6E, 0x6E, 0x61, 0x62, 0x6C, 0x65,
		0x01, 0x00, 0x1D, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F, 0x2F, 0x4F, 0x62, 0x6A, 0x65, 0x63,
		0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x45, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6F, 0x6E,
		0x01, 0x00, 0x15, 0x28, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53,
		0x74, 0x72, 0x69, 0x6E, 0x67, 0x3B, 0x29, 0x56, 0x00, 0x21, 0x00, 0x09, 0x00, 0x0A, 0x00, 0x02,
		0x00, 0x0B, 0x00, 0x0C, 0x00, 0x00, 0x00, 0x05, 0x00, 0x01, 0x00, 0x0D, 0x00, 0x0E, 0x00, 0x01,
		0x00, 0x0F, 0x00, 0x00, 0x00, 0x1D, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x2A, 0xB7,
		0x00, 0x01, 0xB1, 0x00, 0x00, 0x00, 0x01, 0x00, 0x10, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x08, 0x00, 0x02, 0x00, 0x11, 0x00, 0x12, 0x00, 0x02, 0x00, 0x0F, 0x00, 0x00, 0x00,
		0x20, 0x00, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0xBB, 0x00, 0x02, 0x59, 0xB7, 0x00, 0x03,
		0xBF, 0x00, 0x00, 0x00, 0x01, 0x00, 0x10, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x0A, 0x00, 0x13, 0x00, 0x00, 0x00, 0x04, 0x00, 0x01, 0x00, 0x02, 0x00, 0x02, 0x00, 0x14, 0x00,
		0x15, 0x00, 0x02, 0x00, 0x0F, 0x00, 0x00, 0x00, 0x47, 0x00, 0x02, 0x00, 0x03, 0x00, 0x00, 0x00,
		0x17, 0x06, 0x3D, 0x1C, 0x05, 0xA4, 0x00, 0x0B, 0xBB, 0x00, 0x02, 0x59, 0xB7, 0x00, 0x03, 0xBF,
		0xBB, 0x00, 0x04, 0x59, 0xB7, 0x00, 0x05, 0xBF, 0x00, 0x00, 0x00, 0x02, 0x00, 0x10, 0x00, 0x00,
		0x00, 0x12, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x02, 0x00, 0x0E, 0x00, 0x07, 0x00, 0x0F,
		0x00, 0x0F, 0x00, 0x11, 0x00, 0x16, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0xFC, 0x00, 0x0F, 0x01,
		0x00, 0x13, 0x00, 0x00, 0x00, 0x06, 0x00, 0x02, 0x00, 0x02, 0x00, 0x04, 0x00, 0x01, 0x00, 0x17,
		0x00, 0x0E, 0x00, 0x01, 0x00, 0x0F, 0x00, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x01, 0xB1, 0x00, 0x00, 0x00, 0x01, 0x00, 0x10, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x15, 0x00, 0x02, 0x00, 0x18, 0x00, 0x0E, 0x00, 0x02, 0x00, 0x0F, 0x00, 0x00, 0x00,
		0x22, 0x00, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x0A, 0xBB, 0x00, 0x06, 0x59, 0x12, 0x07, 0xB7,
		0x00, 0x08, 0xBF, 0x00, 0x00, 0x00, 0x01, 0x00, 0x10, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x19, 0x00, 0x13, 0x00, 0x00, 0x00, 0x04, 0x00, 0x01, 0x00, 0x19, 0x00, 0x01, 0x00,
		0x1A, 0x00, 0x00, 0x00, 0x02, 0x00, 0x1B,
	}

	// the following tests work against hard-coded values for this class,
	// which are known to be valid. (Compare with javap output for this class.)

	klass, err := parse(classBytes)
	if err != nil {
		fmt.Println(err)
		t.Error("Previous line shows unexpected error in parsing of HaveInterface.class")
	}

	if klass.javaVersion != 55 {
		t.Error("Expected Java version # of 55. Got: " + strconv.Itoa(klass.javaVersion))
	}

	if klass.cpCount != 40 {
		t.Error("Expected constant pool cont of 40. Got: " + strconv.Itoa(klass.cpCount))
	}

	if klass.accessFlags != 0x21 {
		t.Error("Expected access flags of 0x21 (33d). Got: " + strconv.Itoa(klass.accessFlags))
	}

	if klass.classIsPublic != true {
		t.Error("Expected class to be public, but it's not.")
	}

	if klass.classIsSuper != true {
		t.Error("Expected class to be a superclass, but it isn't.")
	}
	if klass.classIsInterface != false {
		t.Error("classIsInterface() returns true, when it should return false.")
	}

	if klass.className != "HaveInterface" {
		t.Error("Expected class to be named 'HaveInterface'. Got: " + klass.className)
	}

	if klass.superClass != "java/lang/Object" {
		t.Error("Expected superclass to be 'java/lang/Object'. Got: " + klass.superClass)
	}

	if klass.interfaceCount != 2 {
		t.Error("Expected 2 interfaces in this class. Got: " + strconv.Itoa(klass.interfaceCount))
	}

	if klass.fieldCount != 0 || len(klass.fields) != 0 {
		t.Error("Expected 0 fields, but got: " + strconv.Itoa(klass.fieldCount) + " and " +
			strconv.Itoa(len(klass.fields)))
	}

	if klass.methodCount != 5 || len(klass.methods) != 5 {
		t.Error("Expected 5 methods, but got: " + strconv.Itoa(klass.methodCount) + " and " +
			strconv.Itoa(len(klass.methods)))
	}

	meth3 := klass.methods[2]
	if klass.utf8Refs[meth3.name].content != "readObject" {
		t.Error("Expected a method name of 'readObject'. Got: " + klass.utf8Refs[meth3.name].content)
	}

	if klass.utf8Refs[meth3.description].content != "(Ljava/io/ObjectInputStream;)V" {
		t.Error("Expected readObject() to have a descriptor of '(Ljava/io/ObjectInputStream;)V'. Got: " +
			klass.utf8Refs[meth3.description].content)
	}

	if len(meth3.attributes) != 2 {
		t.Error("Expected method readObject() to have 2 attributes. Got: " + strconv.Itoa(len(meth3.attributes)))
	}

	attribName := klass.utf8Refs[meth3.attributes[0].attrName].content
	if attribName != "Code" {
		t.Error("Expected name of first method attribute in readObject() to be 'Code'. Got: " +
			attribName)
	}

	if meth3.deprecated != false {
		t.Error("Expected method readObject() not to be deprecated, but it was.")
	}

	if klass.attribCount != 1 || len(klass.attributes) != 1 {
		t.Error("Expected 1 attribute, but got: " + strconv.Itoa(klass.attribCount) + " and " +
			strconv.Itoa(len(klass.attributes)))
	}

	if klass.utf8Refs[klass.attributes[0].attrName].content != "SourceFile" {
		t.Error("Expected a class attribute named 'SourceFile'. Got: " +
			klass.utf8Refs[klass.attributes[0].attrName].content)
	}

	byte1 := klass.attributes[0].attrContent[0]
	byte2 := klass.attributes[0].attrContent[1]
	attrIndex := int(byte1)*256 + int(byte2)
	attrName := klass.utf8Refs[klass.cpIndex[attrIndex].slot]
	if attrName.content != "HaveInterface.java" {
		t.Error("Expected SourceFile attribute to be 'HaveInterface.java'. Got: " + attrName.content)
	}

	err = formatCheckClass(&klass)
	if err != nil {
		t.Error("HaveInterface.class failed format check")
	}
}
