// Copyright 2013
// Author: Christopher Van Arsdale

#include <string>
#include "strings/path.h"
#include "strings/strutil.h"
#include "strings/stringpiece.h"

namespace strings {

std::string JoinPath(const StringPiece& a, const StringPiece& b) {
  return CleanPath(a.as_string() + "/" + b.as_string());
}

std::string CleanPath(const StringPiece& input) {
  bool absolute = HasPrefix(input, "/");
  std::vector<StringPiece> pieces = Split(input, "/");
  std::vector<StringPiece> output;
  int num_deep = 0;
  for (int i = 0; i < pieces.size(); ++i) {
    if (pieces[i] == ".") {
      continue;
    }
    if (pieces[i] == "..") {
      if (!output.empty() && output.back() != "..") {
        output.resize(output.size() - 1);
      } else if (!absolute) {
        output.push_back(pieces[i]);
      }
    } else {
      output.push_back(pieces[i]);
    }
  }
  if (output.empty()) {
    return absolute ? "/" : ".";
  }
  return (absolute ? "/" : "") + Join(output, "/");
}

}  // namespace strings