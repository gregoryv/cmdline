/*
Package cmdline provides a way to parse command line arguments.

This package fixes a few problems that come with using the flag package.

  1. Don't hog the name flag, which is a boolean option
  2. Use appropriate names for arguments, options and flags
  3. Optional documentation, self documenting options are preferred
  4. Simplify multiname options, e.g. -n, --dry-run map to same flag
  5. Skip pointer variations
  6. Include required options TODO

*/
package cmdline
