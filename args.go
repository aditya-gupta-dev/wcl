package main

type ArgsModel struct {
	WorkDir string `arg:"positional"`
	Verbose bool   `arg:"-v"`
	Size    bool   `arg:"-s"`
	Lines   bool   `arg:"-l"`
}
