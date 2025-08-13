package formatter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TokenType int

const (
	TkLabel TokenType = iota
	TkInstruction
	TkOperand
	TkComma
	TkCommentSameLine
	TkCommentNewLine
	TkEmptyLine
)

type Token = struct {
	TkType  TokenType
	TkValue string
}

func TokenTypeToStr(tkType TokenType) string {
	switch tkType {
	case TkLabel:
		return "TkLabel"
	case TkInstruction:
		return "TkInstruction"
	case TkOperand:
		return "TkOperand"
	case TkComma:
		return "TkComma"
	case TkCommentSameLine:
		return "TkCommentSameLine"
	case TkCommentNewLine:
		return "TkCommentNewLine"
	case TkEmptyLine:
		return "TlEmptyLine"
	}

	return "type not added"
}

func PrintTokens(tokens *[]Token) {
	for _, tk := range *tokens {
		const formatStr string = "type: %-20v =====          value: %v\n"
		fmt.Printf(formatStr, TokenTypeToStr(tk.TkType), tk.TkValue)
	}
}

func Tokenize(file *os.File) *[]Token {
	const tokensCapacity int = 1024
	tokens := make([]Token, 0, tokensCapacity)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokenizeLine(&tokens, line)
	}

	return &tokens
}

func appendTokens(tokens *[]Token, tk Token, addComma bool) {
	*tokens = append(*tokens, tk)

	if addComma {
		*tokens = append(*tokens, Token{TkType: TkComma, TkValue: ","})
	}
}

func tokenizeLine(tokens *[]Token, line string) {
	if len(line) == 0 {
		tk := Token{TkType: TkEmptyLine, TkValue: "\\n"}
		appendTokens(tokens, tk, false)
		return
	}

	const breakChars = " \t\n,:"
	const specialChars = ",:"
	buf := make([]byte, 0, len(line))
	instructionFound := false

	shouldAddToken := false
	specialChar := ' '

	for i := 0; i < len(line); i++ {
		if strings.ContainsRune(breakChars, rune(line[i])) {
			if strings.ContainsRune(specialChars, rune(line[i])) {
				specialChar = rune(line[i])
			}

			shouldAddToken = true

			if i != len(line)-1 {
				continue
			}
		} else if i == len(line)-1 {
			buf = append(buf, line[i])
			shouldAddToken = true
		}

		if len(buf) == 0 {
			shouldAddToken = false
		}

		if shouldAddToken {
			shouldAddToken = false

			tkType := TkOperand

			if specialChar == ':' {
				tkType = TkLabel
			} else if !instructionFound {
				tkType = TkInstruction
				instructionFound = true
			}

			tk := Token{TkType: tkType, TkValue: string(buf)}
			*tokens = append(*tokens, tk)

			if specialChar == ',' {
				tk := Token{TkType: TkComma, TkValue: ","}
				*tokens = append(*tokens, tk)
			}

			buf = buf[:0]
			specialChar = ' '
		}

		buf = append(buf, line[i])
	}
}
