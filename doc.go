/*
Package cmdline provides a way to parse command line arguments.

This package fixes a few problems that come with using the flag package.

 - Don't hog the name flag, which is a boolean option
 - Use appropriate names for arguments, options and flags
 - Self documenting arguments and options are preferred
 - Multiname options, e.g. -n, --dry-run map to same flag
 - Skip pointer variations
 - Include required arguments

*/
package cmdline
