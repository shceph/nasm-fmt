package formatter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TokenType int

const (
	TkNoType TokenType = iota
	TkLabel
	TkBracketedDirective
	TkInstruction
	TkOperand
	TkComma
	TkColon
	TkPlus
	TkMinus
	TkAsterisk
	TkSlash
	TkCommentInline
	TkCommentNewLine
	TkEmptyLine
	TkPushIndentLevel
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
	case TkPlus:
		return "TkPlus"
	case TkMinus:
		return "TkMinus"
	case TkAsterisk:
		return "TkAsterisk"
	case TkSlash:
		return "TkSlash"
	case TkCommentInline:
		return "TkCommentSameLine"
	case TkCommentNewLine:
		return "TkCommentNewLine"
	case TkEmptyLine:
		return "TkEmptyLine"
	case TkPushIndentLevel:
		return "TkPushIndentLevel"
	case TkNoType:
		return "=====TkNoType====="
	}

	return "=====ValueNotInEnum======"
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

func handleBracketedDirective(tokens *[]Token, line string, index int) int {
	buf := make([]byte, 0, len(line))

	i := index
	for ; i < len(line); i++ {
		buf = append(buf, line[i])

		if line[i] == ']' {
			i++
			break
		}
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

	arithOperators = [256]TokenType{
		'+': TkPlus,
		'-': TkMinus,
		'*': TkAsterisk,
		'/': TkSlash,
	}
)

func isWhitespaceChar(ch byte) bool {
	return whitespaceChars[ch]
}

func isSpecialChar(ch byte) bool {
	return specialChars[ch]
}

func isArithOperator(ch byte) bool {
	return arithOperators[ch] != TkNoType
}

func operatorToToken(ch byte) Token {
	return Token{TkType: arithOperators[ch], TkValue: string(ch)}
}

func flushPendingToken(
	tokens *[]Token,
	bufValue *[]byte,
	pendingType *TokenType,
	instructionFound *bool,
	line string) {

	if len(*bufValue) == 0 {
		*pendingType = TkNoType
		return
	} else if *pendingType == TkNoType {
		*pendingType = TkOperand
	}

	if !*instructionFound && *pendingType == TkOperand {
		*pendingType = TkInstruction
		*instructionFound = true

		// If first character is not whitespace, that means that the
		// instruction is at the beginning of the line, which means that it
		// should not be indented like instructions usually are under labels
		if !isWhitespaceChar(line[0]) {
			addToken(tokens, TkPushIndentLevel, "0")
		}
	}

	addToken(tokens, *pendingType, string(*bufValue))

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
	pendingType := TkNoType

	for i := 0; i < len(line); i++ {
		ch := line[i]

		switch {
		case ch == ';':
			commentType := TkCommentInline
			if len(bufValue) == 0 {
				commentType = TkCommentNewLine
			}

			flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, line)
			handleComment(tokens, commentType, line, i)
			return
		case ch == '[' && !instructionFound && len(bufValue) == 0:
			i = handleBracketedDirective(tokens, line, i)
		case isWhitespaceChar(ch):
			if len(bufValue) != 0 && pendingType == TkNoType {
				pendingType = TkOperand
			}
		case isSpecialChar(ch):
			switch ch {
			case ':':
				if !instructionFound {
					pendingType = TkLabel
					break
				}
				flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, line)
				addToken(tokens, TkColon, ":")
			case ',':
				flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, line)
				addToken(tokens, TkComma, ",")
			}
		case isArithOperator(ch):
			flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, line)
			token := operatorToToken(ch)
			addToken(tokens, token.TkType, token.TkValue)
		case pendingType != TkNoType && len(bufValue) != 0:
			flushPendingToken(tokens, &bufValue, &pendingType, &instructionFound, line)
			fallthrough
		default:
			bufValue = append(bufValue, ch)
		}
	}

	if len(bufValue) != 0 {
		if pendingType == TkNoType {
			pendingType = TkOperand
		}

		if pendingType == TkOperand && !instructionFound {
			pendingType = TkInstruction
		}

		addToken(tokens, pendingType, string(bufValue))
	}
}
