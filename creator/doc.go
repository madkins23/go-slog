// Package creator provides a common location for defining infra.Creator objects.
// Each object is a factory for generating slog.Handler instances.
//
// Declarations are kept in separate packages so that they can be reused in
// other projects without dragging in all the various handlers and logging packages.
// This applies to the utility packages as well, which are shared by creators.
package creator
