// Copyright 2013
// Author: Christopher Van Arsdale

package main

import (
  "flag"
  "fmt"
)

var FLAGS_max_args = flag.Int("max_args", 7, "Number of arguments to auto-gen.")

func OutputHeader() {
  fmt.Println(
    `// Copyright 2013
// Author: Christopher Van Arsdale
//
// THIS CODE IS AUTO-GENERATED BY callback_gen.go, DO NOT MODIFY DIRECTLY
//
// Note: This was originally implemented with std::bind/std::function, but
// you cannot declare them ahead of time (for anything with an argument and an
// object) so passing the callbacks around is not tenable.

#ifndef _COMMON_BASE_CALLBACK_H__
#define _COMMON_BASE_CALLBACK_H__

#include "common/base/macros.h"`)
}

func OutputFooter() {
  fmt.Println("")
  fmt.Println("#endif  // _COMMON_BASE_CALLBACK_H__")
}

func GetTemplateArgs(num_args int, object bool, final_object bool) string {
  if num_args <= 0 && !object {
    return ""
  }

  template := "template <"
  if object {
    template += "typename Object"
  }
  for i := 0; i < num_args; i++ {
    if i > 0 || object { template += ", " }
    template += fmt.Sprintf("typename Arg%d", i)
  }
  if final_object {
    template += ", typename InputObject=Object"
  }
  return template + ">"
}

func GetFunctionArgsStart(start int, num_args int) string {
  output := ""
  for i := start; i < num_args; i++ {
    if i > start {
      output += ", "
    }
    output += fmt.Sprintf("Arg%d arg%d", i, i)
  }
  return output
}

func GetFunctionArgs(num_args int) string {
  return GetFunctionArgsStart(0, num_args)
}

func OutputBaseClass(num_args int, name string) {
  fmt.Println(GetTemplateArgs(num_args, false, false))
  fmt.Println("class " + name + " {")
  fmt.Println(" public:")
  fmt.Println("  " + name + "() {}")
  fmt.Println("  virtual ~" + name + "() {}")
  fmt.Println("  virtual void Run(" + GetFunctionArgs(num_args) + ") = 0;")
  fmt.Println("  virtual bool IsPermanentCallback() = 0;")
  fmt.Println(" private: ")
  fmt.Println("  DISALLOW_COPY_AND_ASSIGN(" + name + ");")
  fmt.Println("};")
}

func PrintConstructor(class string, num_args int, object bool) {
  // Start the constructor
  line := "  " + class + "(bool perm, Func func"
  if object {
    line += ", Object* obj"
  }
  if num_args > 0 {
    line += ", " + GetFunctionArgs(num_args)
  }
  line += ")"
  fmt.Println(line)

  // Fill in member initialization
  line = "    : perm_(perm), func_(func)"
  if object {
    line += ", object_(obj)"
  }
  for i := 0; i < num_args; i++ {
    line += fmt.Sprintf(", arg%d_(arg%d)", i, i)
  }
  line += "{"
  fmt.Println(line)
  fmt.Println("  }")
}

func PrintRun(num_total int, num_inputs int, object bool) {
  fmt.Println("  virtual void Run(" +
    GetFunctionArgsStart(num_inputs, num_total) + ") {")
  fmt.Println("    bool del = !perm_;")
  line := "    func_("
  if object {
    line = "    (object_->*func_)("
  }
  for i := 0; i < num_inputs; i++ {
    if i > 0 { line += ", " }
    line += fmt.Sprintf("arg%d_", i)
  }
  for i := num_inputs; i < num_total; i++ {
    if i > 0 { line += ", " }
    line += fmt.Sprintf("arg%d", i)
  }
  line += ");"
  fmt.Println(line)
  fmt.Println("    if (del) { delete this; }")
  fmt.Println("  }")
}

func PrintFunctionTypedef(num_total int, object bool) {
  if !object {
    fmt.Println("  typedef void (*Func)(" +
      GetFunctionArgs(num_total) + ");")
  } else {
    fmt.Println("  typedef void (Object::*Func)(" +
      GetFunctionArgs(num_total) + ");")
  }
}

func GetBaseClass(start_num int, num_args int) string {
  if num_args == 0 {
    return "Closure"
  }

  base_class := "Callback" + fmt.Sprintf("%d", num_args)
  base_class += "<"
  for i := 0; i < num_args; i++ {
    if i > 0 { base_class += ", " }
    base_class += fmt.Sprintf("Arg%d", i + start_num)
  }
  base_class += ">"
  return base_class
}

func GetClassName(num_inputs int, num_args int,
  object bool, templated bool) string {
  class := ""
  if object {
    class = "MemberCallback"
  } else {
    class = "FunctionCallback"
  }
  class += fmt.Sprintf("%d_%d", num_inputs, num_args)
  if templated && (object || num_inputs + num_args > 0) {
    template := ""
    if object {
      template += "Object"
    }
    for i := 0; i < num_inputs + num_args; i++ {
      if len(template) > 0 { template += "," }
      template += fmt.Sprintf("Arg%d", i)
    }
    class += "<" + template + ">"
  }
  return class
}

func OutputCallbackClass(num_total int, num_args int, object bool) {
  num_inputs := num_total - num_args
  base_class := GetBaseClass(0, num_args)
  class := GetClassName(num_inputs, num_args, object, false)

  fmt.Println(GetTemplateArgs(num_total, object, false))
  fmt.Println("class " + class + " : public " + base_class + " {")
  fmt.Println(" public:")
  PrintFunctionTypedef(num_total, object)
  PrintConstructor(class, num_inputs, object)
  fmt.Println("  virtual ~" + class + "() {}")
  fmt.Println("")
  PrintRun(num_total, num_inputs, object)
  fmt.Println("  virtual bool IsPermanentCallback() { return perm_; }")
  fmt.Println("")
  fmt.Println(" private:")
  fmt.Println("  DISALLOW_COPY_AND_ASSIGN(" + class + ");")
  fmt.Println("")
  fmt.Println("  bool perm_;")
  fmt.Println("  Func func_;")
  if object {
    fmt.Println("  Object* object_;")
  }
  for i := 0; i < num_inputs; i++ {
    fmt.Println(fmt.Sprintf("  Arg%d arg%d_;", i, i))
  }
  fmt.Println("};")
  fmt.Println("")
}

func OutputNewCallback(num_total int, num_args int, object bool, perm bool) {
  num_inputs := num_total - num_args

  // tempalte <...>
  fmt.Println("")
  fmt.Println(GetTemplateArgs(num_total, object, object))

  // Callback...* NewCallback(...) {
  line := GetBaseClass(num_inputs, num_args) + "* "
  if perm {
    line += "NewPermanentCallback("
  } else {
    line += "NewCallback("
  }
  if object {
    line += "InputObject* object, void(Object::*f)"
  } else {
    line += "void (*f)"
  }
  line += "(" + GetFunctionArgs(num_total) + ")"
  if num_inputs > 0 {
    line += ", " + GetFunctionArgsStart(0, num_inputs)
  }
  line += ") {"
  fmt.Println(line)

  line = "  return new " +
    GetClassName(num_inputs, num_args, object, true) + "(";
  fmt.Println(line)
  fmt.Println(fmt.Sprintf("    %t,", perm))
  line = "    f"
  if object {
    fmt.Println(line + ",")
    line = "    object"
  }
  if num_inputs > 0 {
    fmt.Println(line + ",")
    line = "    "
    for i := 0; i < num_inputs; i++ {
      if i > 0 { line += ", " }
      line += fmt.Sprintf("arg%d", i)
    }
  }
  fmt.Println(line + ");")
  fmt.Println("}")
}

func main() {
  flag.Parse()

  OutputHeader()

  // Base classes
  OutputBaseClass(0, "Closure")
  for i := 1; i <= *FLAGS_max_args; i++ {
    fmt.Println("")
    OutputBaseClass(i, "Callback" + fmt.Sprintf("%d", i))
  }

  // Output callback classes
  for i := 0; i <= *FLAGS_max_args; i++ {
    for j := 0; j <= i; j++ {
      OutputCallbackClass(i, j, false)
      OutputCallbackClass(i, j, true)
    }
  }

  // Output NewCallback functions
  for i := 0; i <= *FLAGS_max_args; i++ {
    for j := 0; j <= i; j++ {
      OutputNewCallback(i, j, false, false)
      OutputNewCallback(i, j, false, true)
      OutputNewCallback(i, j, true, false)
      OutputNewCallback(i, j, true, true)
    }
  }

  OutputFooter()
}
