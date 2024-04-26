// Package creator provides a common location for defining infra.Creator objects.
// Each such object is a factory for generating slog.Handler instances.
//
// The author keeps thinking that Creator objects may have some external use, so they are public.
// Declarations are kept in separate packages so that they can potentially be reused in
// other projects without dragging in all the various other handlers and logging packages.
// This applies to the utility packages as well, which are shared by creators.
package creator
