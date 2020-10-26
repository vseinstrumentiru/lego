package server

const (
	rtNoWait = "nowait"
)

func NoWait(r *Runtime) {
	r.options[rtNoWait] = true
}
