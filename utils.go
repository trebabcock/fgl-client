package main

import (
	"runtime"
	"strings"

	"github.com/c-bata/go-prompt"
)

func baseURL(path string) string {
	return "http://" + server + path
}

func chatURL() string {
	return "ws://127.0.0.1:2814/chat"
}

func emptyCompleter(d prompt.Document) []prompt.Suggest {
	s := make([]prompt.Suggest, 0)
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := make([]prompt.Suggest, 0)
	s = append(s, prompt.Suggest{Text: "cd", Description: "change directory"})
	s = append(s, prompt.Suggest{Text: "clear", Description: "clear console"})
	s = append(s, prompt.Suggest{Text: "exit", Description: "exit the shell"})
	s = append(s, prompt.Suggest{Text: "pwd", Description: "print working directory"})
	s = append(s, prompt.Suggest{Text: "cat", Description: "output the contents of a file"})
	s = append(s, prompt.Suggest{Text: "echo", Description: "output the supplied value"})

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func makeMainPrompt() string {
	return currentUser + " > "
}

func makeShellPrompt() string {
	return currentUser + "@fgldb " + currentDir.Name + " $ "
}

func clear() {
	value, ok := mClear[runtime.GOOS]
	if ok {
		value()
	} else {
		mClear["linux"]()
	}
}

func wordWrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped
}

func wordWrap2(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n        " + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped
}
