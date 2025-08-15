package formatter

import (
	"os"
)

type FormatOpts = struct {
	NoTab            bool
	TabWidth         int
	OperandAlignDist int
}

var DefaultFormatOpts FormatOpts = FormatOpts{
	NoTab:            true,
	TabWidth:         4,
	OperandAlignDist: 10,
}

func Format(filePath string, opts FormatOpts) (string, error) {
	tokens, err := Tokenize(filePath)

	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	output := generateFromTokens(tokens, DefaultFormatOpts)
	_, err = file.WriteString(output)

	return output, err
}

func bufAddTab(buf *[]byte, tabCount int, opts *FormatOpts) {
	for range tabCount {
		if !opts.NoTab {
			*buf = append(*buf, '\t')
			continue
		}

		for range opts.TabWidth {
			*buf = append(*buf, ' ')
		}
	}
}

func bufAddStr(buf *[]byte, values ...string) {
	for _, val := range values {
		for _, ch := range val {
			*buf = append(*buf, byte(ch))
		}
	}
}

func bufAddChr(buf *[]byte, char byte) {
	*buf = append(*buf, char)
}

func bufAddWhitespaces(buf *[]byte, count int) {
	for range count {
		*buf = append(*buf, ' ')
	}
}

func generateFromTokens(tokens []Token, opts FormatOpts) string {
	if len(tokens) == 0 {
		return ""
	}

	buf := make([]byte, 0, len(tokens)*5)
	indentLevel := 0
	instructionLen := 0
	dontAddSpaceInFrontOfOperand := false

	for _, tk := range tokens {
		switch tk.TkType {
		case TkLabel:
			bufAddStr(&buf, "\n", tk.TkValue, ":")
			indentLevel = 1
		case TkBracketedDirective:
			bufAddStr(&buf, "\n[", tk.TkValue, "]")
			indentLevel = 0
		case TkInstruction:
			bufAddChr(&buf, '\n')
			bufAddTab(&buf, indentLevel, &opts)
			bufAddStr(&buf, tk.TkValue)
			instructionLen = len(tk.TkValue)
		case TkOperand:
			if instructionLen != 0 && indentLevel != 0 {
				bufAddWhitespaces(&buf, opts.OperandAlignDist-instructionLen)
				instructionLen = 0
			}
			if dontAddSpaceInFrontOfOperand {
				dontAddSpaceInFrontOfOperand = false
			} else {
				bufAddChr(&buf, ' ')
			}
			bufAddStr(&buf, tk.TkValue)
		case TkComma:
			bufAddChr(&buf, ',')
		case TkColon:
			bufAddChr(&buf, ':')
			dontAddSpaceInFrontOfOperand = true
		case TkCommentSameLine:
			bufAddStr(&buf, " ; ", tk.TkValue)
		case TkCommentNewLine:
			bufAddChr(&buf, '\n')
			bufAddTab(&buf, indentLevel, &opts)
			bufAddStr(&buf, "; ", tk.TkValue)
		case TkEmptyLine:
			bufAddChr(&buf, '\n')
		case TkPushIndentLevel:
			indentLevel = int(tk.TkValue[0] - '0')
		}
	}

	startIndex := 0

	if buf[0] == '\n' {
		startIndex = 1
	}

	return string(buf[startIndex:])
}
