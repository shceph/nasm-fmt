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
	TkBracketedDirective
	TkInstruction
	TkOperand
	TkComma
	TkColon
	TkCommentSameLine
	TkCommentNewLine
	TkEmptyLine
	TkNoType
)

type Token = struct {
	TkType  TokenType
	TkValue string
}

func TokenTypeToStr(tkType TokenType) string {
	switch tkType {
	case TkLabel:
		return "TkLabel"
	case TkBracketedDirective:
		return "TkBracketedDirective"
	case TkInstruction:
		return "TkInstruction"
	case TkOperand:
		return "TkOperand"
	case TkComma:
		return "TkComma"
	case TkColon:
		return "TkColon"
	case TkCommentSameLine:
		return "TkCommentSameLine"
	case TkCommentNewLine:
		return "TkCommentNewLine"
	case TkEmptyLine:
		return "TkEmptyLine"
	}

	return "type not added"
}

func PrintTokens(tokens []Token) {
	for _, tk := range tokens {
		const formatStr string = "type: %-20v =====          value: %v\n"
		fmt.Printf(formatStr, TokenTypeToStr(tk.TkType), tk.TkValue)
	}
}

func Tokenize(filePath string) ([]Token, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	const tokensCapacity int = 1024
	tokens := make([]Token, 0, tokensCapacity)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokenizeLine(&tokens, line)
	}

	return tokens, nil
}

func addToken(tokens *[]Token, tkType TokenType, tkValue string) {
	tk := Token{TkType: tkType, TkValue: tkValue}
	*tokens = append(*tokens, tk)
}

func handleComment(tokens *[]Token, commentType TokenType, line string, index int) {
	if index >= len(line) {
		panic("index greater than line lenght in trimComment")
	}

	// line[index] is ';', we don't want that included in value
	index++
	tkValue := ""

	if index < len(line) && line[index] == ' ' {
		index++
	}

	if index < len(line) {
		tkValue = strings.Join(strings.Fields(line[index:]), " ")
	}

	addToken(tokens, commentType, tkValue)
}

func handleBracketedDirective(tokens *[]Token, index int, line string) int {
	buf := make([]byte, 0, len(line))

	i := 0
	for i = index + 1; line[i] != ']' && i < len(line); i++ {
		buf = append(buf, line[i])
	}

	addToken(tokens, TkBracketedDirective, string(buf))
	return i
}

var (
	whitespaceChars = [256]bool{
		' ':  true,
		'\t': true,
		'\n': true,
	}

	specialChars = [256]bool{
		',': true,
		':': true,
	}
)

func isWhitespaceChar(ch byte) bool {
	return int(ch) < len(whitespaceChars) && whitespaceChars[ch]
}

func isSpecialChar(ch byte) bool {
	return int(ch) < len(specialChars) && specialChars[ch]
}

func flushPendingToken(
	tokens *[]Token,
	bufValue *[]byte,
	pendingType *TokenType,
	instructionFound, addComma *bool) {

	if len(*bufValue) != 0 {
		if !*instructionFound && *pendingType == TkOperand {
			*pendingType = TkInstruction
			*instructionFound = true
		}
		addToken(tokens, *pendingType, string(*bufValue))
	}

	if *addComma {
		addToken(tokens, TkComma, ",")
		*addComma = false
	}

	*bufValue = (*bufValue)[:0]
	*pendingType = TkNoType
}

func tokenizeLine(tokens *[]Token, line string) {
	if len(line) == 0 {
		addToken(tokens, TkEmptyLine, "\\n")
		return
	}

	bufValue := make([]byte, 0, len(line))

	instructionFound := false
	addComma := false
	pendingType := TkNoType

	for i := 0; i < len(line); i++ {
		ch := line[i]

		switch {
		case ch == ';':
			commentType := TkCommentSameLine
			if len(bufValue) == 0 {
				commentType = TkCommentNewLine
			}

			flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, &addComma)
			handleComment(tokens, commentType, line, i)
			return
		case ch == '[' && !instructionFound:
			i = handleBracketedDirective(tokens, i, line)
		case isWhitespaceChar(ch):
			if len(bufValue) != 0 {
				pendingType = TkOperand
			}
		case isSpecialChar(ch):
			if len(bufValue) != 0 {
				pendingType = TkOperand
			}

			switch ch {
			case ':':
				if !instructionFound {
					pendingType = TkLabel
				} else {
					flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, &addComma)
					addToken(tokens, TkColon, ":")
				}
			case ',':
				addComma = true
			}
		case pendingType != TkNoType && len(bufValue) != 0:
			flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, &addComma)
			fallthrough
		default:
			bufValue = append(bufValue, ch)
		}
	}

	if len(bufValue) != 0 {
		if pendingType == TkNoType {
			pendingType = TkOperand
		}

		addToken(tokens, pendingType, string(bufValue))
		if addComma {
			addToken(tokens, TkComma, ",")
		}
	}
}
