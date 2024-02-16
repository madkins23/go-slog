// Package infra contains functionality shared between test and benchmark managers.
//
//   - Convenient definitions for:
//   - AttrFn defines the slog.HandlerOptions.ReplaceAttr function template.
//   - EmptyAttr() returns an empty attribute.
//   - Creator struct and instances thereof for specific slog.Handlers.
//   - Functions to return slog.HandlerOptions of general utility.
package infra
