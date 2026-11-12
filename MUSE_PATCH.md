# Muse MyGo patch

This is a narrowly scoped local fork of `github.com/egoist/mygo` v0.8.3.
The upstream MIT license and `go.mod` are retained. The fork is based on the
module source published for v0.8.3; docs, examples, benchmark fixtures, package
examples and upstream tests were pruned. Build source and embedded runtime
assets are retained.

Muse needs to promote an already active in-window tree drag to a native
cross-window drag only after the pointer leaves its window. Calling
`Window.StartDataDrag` directly bypasses the native UI engine's completion
wrapper, leaving its local drag and pressed state active after the OS drag.
The patch adds `(*ui.Context).HandoffDataDrag`, which promotes the active
runtime drag and reuses the same cleanup behavior as `Element.DragData`.

Keep the fork synchronized with upstream deliberately: this patch touches the
native UI drag lifecycle and should be reviewed against the upstream
`ui/dragdrop.go` whenever the pinned MyGo version changes.
