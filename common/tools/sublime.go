package tools

const SUBLIME_APP = "open -a \"/Applications/Sublime Text.app\" "

func SublimeOpenFile(path string) {
	RunCommandPrintError(SUBLIME_APP + path)
}
