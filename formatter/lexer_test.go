package formatter_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/shceph/nasm-fmt/formatter"
)

func TestOnlyBasicInstructions(t *testing.T) {
	expected_tokens := []formatter.Token{
		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "eax"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "cr0"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "ebx"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "0x13"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "ecx"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "0x15"},
	}

	t.Log("Testing examples/just_instructions.asm ...")

	file, err := os.Open("../examples/just_instructions.asm")

	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	tokens := formatter.Tokenize(file)

	if !reflect.DeepEqual(tokens, &expected_tokens) {
		formatter.PrintTokens(&expected_tokens)
		fmt.Println()
		formatter.PrintTokens(tokens)

		t.Error("Tokens are not matching")
	}
}

func TestLabelsAndInstructions(t *testing.T) {
	expected_tokens := []formatter.Token{
		{TkType: formatter.TkLabel, TkValue: "_start"},

		{TkType: formatter.TkInstruction, TkValue: "extern"},
		{TkType: formatter.TkOperand, TkValue: "page_directory"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "eax"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "page_directory"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "cr3"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "eax"},

		{TkType: formatter.TkEmptyLine, TkValue: "\\n"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "eax"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "cr0"},

		{TkType: formatter.TkInstruction, TkValue: "or"},
		{TkType: formatter.TkOperand, TkValue: "eax"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "0x80000001"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "cr0"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "eax"},
	}

	t.Log("Testing examples/label_and_instructions.asm ...")

	file, err := os.Open("../examples/label_and_instructions.asm")

	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	tokens := formatter.Tokenize(file)

	if !reflect.DeepEqual(tokens, &expected_tokens) {
		formatter.PrintTokens(&expected_tokens)
		fmt.Println()
		formatter.PrintTokens(tokens)

		t.Error("Tokens are not matching")
	}
}

func TestWithComments(t *testing.T) {
	expected_tokens := []formatter.Token{
		{TkType: formatter.TkCommentNewLine, TkValue: "; Start label"},

		{TkType: formatter.TkLabel, TkValue: "_start"},

		{TkType: formatter.TkInstruction, TkValue: "extern"},
		{TkType: formatter.TkOperand, TkValue: "page_directory"},
		{TkType: formatter.TkCommentSameLine, TkValue: "; An extern variable"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "eax"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "page_directory"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "cr3"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "eax"},

		{TkType: formatter.TkEmptyLine, TkValue: "\\n"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "eax"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "cr0"},

		{TkType: formatter.TkInstruction, TkValue: "or"},
		{TkType: formatter.TkOperand, TkValue: "eax"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "0x80000001"},

		{TkType: formatter.TkInstruction, TkValue: "mov"},
		{TkType: formatter.TkOperand, TkValue: "cr0"},
		{TkType: formatter.TkComma, TkValue: ","},
		{TkType: formatter.TkOperand, TkValue: "eax"},
	}

	t.Log("Testing examples/with_comments.asm ...")

	file, err := os.Open("../examples/with_comments.asm")

	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	tokens := formatter.Tokenize(file)

	if !reflect.DeepEqual(tokens, &expected_tokens) {
		formatter.PrintTokens(&expected_tokens)
		fmt.Println()
		formatter.PrintTokens(tokens)

		t.Error("Tokens are not matching")
	}
}
