package main

type OutputFormat string

const (
	FormatAuto     OutputFormat = "Auto"
	FormatMOBI     OutputFormat = "MOBI"
	FormatEPUB     OutputFormat = "EPUB"
	FormatCBZ      OutputFormat = "CBZ"
	FormatKFX      OutputFormat = "KFX"
	FormatMOBIEPUB OutputFormat = "MOBI+EPUB"
)
