// Code generated from bpftrace.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // bpftrace

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type bpftraceParser struct {
	*antlr.BaseParser
}

var BpftraceParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func bpftraceParserInit() {
	staticData := &BpftraceParserStaticData
	staticData.LiteralNames = []string{
		"", "':'", "'/'", "'{'", "';'", "'}'", "'='", "'@'", "'['", "','", "']'",
		"'('", "')'", "'+'", "'-'", "'*'", "'%'", "'!'", "'~'", "'.'", "'|'",
		"'sprintf'", "'system'", "'count'", "'sum'", "'avg'", "'min'", "'max'",
		"'stats'", "'hist'", "'lhist'", "'cat'", "'join'", "'time'", "'strftime'",
		"'str'", "'strerror'", "'kaddr'", "'uaddr'", "'ntop'", "'pton'", "'reg'",
		"'kstack'", "'ustack'", "'ksym'", "'usym'", "'cgroupid'", "'macaddr'",
		"'nsecs'", "'elapsed'", "'cpu'", "'pid'", "'tid'", "'uid'", "'gid'",
		"'comm'", "'curtask'", "'rand'", "'ctx'", "'args'", "'retval'", "'probe'",
		"'username'", "", "", "", "", "", "'if'", "'else'", "'while'", "'for'",
		"'in'", "'return'", "'clear'", "'delete'", "'exit'", "'print'", "'printf'",
		"'END'", "'=='", "'!='", "'<'", "'>'", "'<='", "'>='", "'&&'", "'||'",
		"'<<'", "'>>'", "'+='", "'-='", "'*='", "'/='", "'%='", "'&='", "'|='",
		"'^='", "'<<='", "'>>='", "'++'", "'--'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "SHEBANG", "COMMENT",
		"WS", "NEWLINE", "PROBE_TYPE", "IF", "ELSE", "WHILE", "FOR", "IN", "RETURN",
		"CLEAR", "DELETE", "EXIT", "PRINT", "PRINTF", "END", "EQ", "NE", "LT",
		"GT", "LE", "GE", "AND", "OR", "SHL", "SHR", "ADD_ASSIGN", "SUB_ASSIGN",
		"MUL_ASSIGN", "DIV_ASSIGN", "MOD_ASSIGN", "AND_ASSIGN", "OR_ASSIGN",
		"XOR_ASSIGN", "SHL_ASSIGN", "SHR_ASSIGN", "INCR", "DECR", "NUMBER",
		"STRING", "IDENTIFIER", "VARIABLE", "MAP_NAME",
	}
	staticData.RuleNames = []string{
		"program", "shebang_section", "content", "shebang", "probe", "probe_def",
		"predicate", "block", "statement", "assignment", "map_assign", "var_assign",
		"function_call", "if_statement", "while_statement", "for_statement",
		"return_statement", "clear_statement", "delete_statement", "exit_statement",
		"print_statement", "printf_statement", "expression", "logical_or_expression",
		"logical_and_expression", "equality_expression", "relational_expression",
		"shift_expression", "additive_expression", "multiplicative_expression",
		"unary_expression", "postfix_expression", "primary_expression", "variable",
		"map_access", "expr_list", "output_redirection", "function_name", "builtin_name",
		"comment", "string",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 106, 499, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 1, 0, 3, 0, 84,
		8, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 5, 1, 91, 8, 1, 10, 1, 12, 1, 94, 9,
		1, 1, 2, 1, 2, 1, 2, 5, 2, 99, 8, 2, 10, 2, 12, 2, 102, 9, 2, 1, 3, 1,
		3, 1, 4, 1, 4, 3, 4, 108, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 114, 8, 4,
		1, 5, 1, 5, 1, 5, 5, 5, 119, 8, 5, 10, 5, 12, 5, 122, 9, 5, 1, 6, 1, 6,
		1, 6, 1, 6, 1, 7, 1, 7, 5, 7, 130, 8, 7, 10, 7, 12, 7, 133, 9, 7, 1, 7,
		1, 7, 5, 7, 137, 8, 7, 10, 7, 12, 7, 140, 9, 7, 1, 7, 3, 7, 143, 8, 7,
		5, 7, 145, 8, 7, 10, 7, 12, 7, 148, 9, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 164, 8, 8,
		1, 9, 1, 9, 3, 9, 168, 8, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10,
		1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 5, 10, 181, 8, 10, 10, 10, 12, 10, 184,
		9, 10, 1, 10, 1, 10, 3, 10, 188, 8, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1,
		10, 1, 10, 1, 10, 5, 10, 197, 8, 10, 10, 10, 12, 10, 200, 9, 10, 1, 10,
		1, 10, 3, 10, 204, 8, 10, 1, 10, 1, 10, 3, 10, 208, 8, 10, 1, 11, 1, 11,
		1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 216, 8, 11, 1, 12, 1, 12, 1, 12, 3,
		12, 221, 8, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 3, 12, 228, 8, 12, 1,
		12, 1, 12, 3, 12, 232, 8, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13,
		1, 13, 3, 13, 241, 8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1,
		15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15,
		259, 8, 15, 1, 15, 1, 15, 3, 15, 263, 8, 15, 1, 15, 1, 15, 3, 15, 267,
		8, 15, 1, 15, 1, 15, 3, 15, 271, 8, 15, 1, 16, 1, 16, 3, 16, 275, 8, 16,
		1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 5, 17, 285, 8,
		17, 10, 17, 12, 17, 288, 9, 17, 1, 17, 1, 17, 3, 17, 292, 8, 17, 3, 17,
		294, 8, 17, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 5, 18, 302, 8, 18,
		10, 18, 12, 18, 305, 9, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1,
		18, 1, 18, 5, 18, 315, 8, 18, 10, 18, 12, 18, 318, 9, 18, 1, 18, 1, 18,
		3, 18, 322, 8, 18, 1, 19, 1, 19, 3, 19, 326, 8, 19, 1, 20, 1, 20, 3, 20,
		330, 8, 20, 1, 20, 3, 20, 333, 8, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21,
		5, 21, 340, 8, 21, 10, 21, 12, 21, 343, 9, 21, 1, 21, 1, 21, 1, 22, 1,
		22, 1, 23, 1, 23, 1, 23, 5, 23, 352, 8, 23, 10, 23, 12, 23, 355, 9, 23,
		1, 24, 1, 24, 1, 24, 5, 24, 360, 8, 24, 10, 24, 12, 24, 363, 9, 24, 1,
		25, 1, 25, 1, 25, 5, 25, 368, 8, 25, 10, 25, 12, 25, 371, 9, 25, 1, 26,
		1, 26, 1, 26, 5, 26, 376, 8, 26, 10, 26, 12, 26, 379, 9, 26, 1, 27, 1,
		27, 1, 27, 5, 27, 384, 8, 27, 10, 27, 12, 27, 387, 9, 27, 1, 28, 1, 28,
		1, 28, 5, 28, 392, 8, 28, 10, 28, 12, 28, 395, 9, 28, 1, 29, 1, 29, 1,
		29, 5, 29, 400, 8, 29, 10, 29, 12, 29, 403, 9, 29, 1, 30, 1, 30, 1, 30,
		1, 30, 1, 30, 3, 30, 410, 8, 30, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1,
		31, 5, 31, 418, 8, 31, 10, 31, 12, 31, 421, 9, 31, 1, 31, 1, 31, 1, 31,
		1, 31, 3, 31, 427, 8, 31, 1, 31, 1, 31, 1, 31, 3, 31, 432, 8, 31, 1, 32,
		1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 3, 32, 443, 8,
		32, 1, 33, 1, 33, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34, 5, 34, 452, 8, 34,
		10, 34, 12, 34, 455, 9, 34, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34, 1,
		34, 5, 34, 464, 8, 34, 10, 34, 12, 34, 467, 9, 34, 1, 34, 1, 34, 3, 34,
		471, 8, 34, 3, 34, 473, 8, 34, 1, 35, 1, 35, 1, 35, 5, 35, 478, 8, 35,
		10, 35, 12, 35, 481, 9, 35, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 3,
		36, 489, 8, 36, 1, 37, 1, 37, 1, 38, 1, 38, 1, 39, 1, 39, 1, 40, 1, 40,
		1, 40, 0, 0, 41, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28,
		30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64,
		66, 68, 70, 72, 74, 76, 78, 80, 0, 10, 1, 0, 90, 99, 1, 0, 80, 81, 1, 0,
		82, 85, 1, 0, 88, 89, 1, 0, 13, 14, 2, 0, 2, 2, 15, 16, 2, 0, 13, 14, 17,
		18, 1, 0, 100, 101, 1, 0, 104, 105, 2, 0, 21, 62, 74, 78, 535, 0, 83, 1,
		0, 0, 0, 2, 88, 1, 0, 0, 0, 4, 100, 1, 0, 0, 0, 6, 103, 1, 0, 0, 0, 8,
		113, 1, 0, 0, 0, 10, 115, 1, 0, 0, 0, 12, 123, 1, 0, 0, 0, 14, 127, 1,
		0, 0, 0, 16, 163, 1, 0, 0, 0, 18, 167, 1, 0, 0, 0, 20, 207, 1, 0, 0, 0,
		22, 215, 1, 0, 0, 0, 24, 231, 1, 0, 0, 0, 26, 233, 1, 0, 0, 0, 28, 242,
		1, 0, 0, 0, 30, 270, 1, 0, 0, 0, 32, 272, 1, 0, 0, 0, 34, 293, 1, 0, 0,
		0, 36, 321, 1, 0, 0, 0, 38, 323, 1, 0, 0, 0, 40, 327, 1, 0, 0, 0, 42, 334,
		1, 0, 0, 0, 44, 346, 1, 0, 0, 0, 46, 348, 1, 0, 0, 0, 48, 356, 1, 0, 0,
		0, 50, 364, 1, 0, 0, 0, 52, 372, 1, 0, 0, 0, 54, 380, 1, 0, 0, 0, 56, 388,
		1, 0, 0, 0, 58, 396, 1, 0, 0, 0, 60, 409, 1, 0, 0, 0, 62, 411, 1, 0, 0,
		0, 64, 442, 1, 0, 0, 0, 66, 444, 1, 0, 0, 0, 68, 472, 1, 0, 0, 0, 70, 474,
		1, 0, 0, 0, 72, 488, 1, 0, 0, 0, 74, 490, 1, 0, 0, 0, 76, 492, 1, 0, 0,
		0, 78, 494, 1, 0, 0, 0, 80, 496, 1, 0, 0, 0, 82, 84, 3, 2, 1, 0, 83, 82,
		1, 0, 0, 0, 83, 84, 1, 0, 0, 0, 84, 85, 1, 0, 0, 0, 85, 86, 3, 4, 2, 0,
		86, 87, 5, 0, 0, 1, 87, 1, 1, 0, 0, 0, 88, 92, 3, 6, 3, 0, 89, 91, 5, 66,
		0, 0, 90, 89, 1, 0, 0, 0, 91, 94, 1, 0, 0, 0, 92, 90, 1, 0, 0, 0, 92, 93,
		1, 0, 0, 0, 93, 3, 1, 0, 0, 0, 94, 92, 1, 0, 0, 0, 95, 99, 3, 8, 4, 0,
		96, 99, 3, 78, 39, 0, 97, 99, 5, 66, 0, 0, 98, 95, 1, 0, 0, 0, 98, 96,
		1, 0, 0, 0, 98, 97, 1, 0, 0, 0, 99, 102, 1, 0, 0, 0, 100, 98, 1, 0, 0,
		0, 100, 101, 1, 0, 0, 0, 101, 5, 1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 103,
		104, 5, 63, 0, 0, 104, 7, 1, 0, 0, 0, 105, 107, 3, 10, 5, 0, 106, 108,
		3, 12, 6, 0, 107, 106, 1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108, 109, 1, 0,
		0, 0, 109, 110, 3, 14, 7, 0, 110, 114, 1, 0, 0, 0, 111, 112, 5, 79, 0,
		0, 112, 114, 3, 14, 7, 0, 113, 105, 1, 0, 0, 0, 113, 111, 1, 0, 0, 0, 114,
		9, 1, 0, 0, 0, 115, 120, 5, 67, 0, 0, 116, 117, 5, 1, 0, 0, 117, 119, 5,
		67, 0, 0, 118, 116, 1, 0, 0, 0, 119, 122, 1, 0, 0, 0, 120, 118, 1, 0, 0,
		0, 120, 121, 1, 0, 0, 0, 121, 11, 1, 0, 0, 0, 122, 120, 1, 0, 0, 0, 123,
		124, 5, 2, 0, 0, 124, 125, 3, 44, 22, 0, 125, 126, 5, 2, 0, 0, 126, 13,
		1, 0, 0, 0, 127, 131, 5, 3, 0, 0, 128, 130, 5, 66, 0, 0, 129, 128, 1, 0,
		0, 0, 130, 133, 1, 0, 0, 0, 131, 129, 1, 0, 0, 0, 131, 132, 1, 0, 0, 0,
		132, 146, 1, 0, 0, 0, 133, 131, 1, 0, 0, 0, 134, 142, 3, 16, 8, 0, 135,
		137, 5, 66, 0, 0, 136, 135, 1, 0, 0, 0, 137, 140, 1, 0, 0, 0, 138, 136,
		1, 0, 0, 0, 138, 139, 1, 0, 0, 0, 139, 143, 1, 0, 0, 0, 140, 138, 1, 0,
		0, 0, 141, 143, 5, 4, 0, 0, 142, 138, 1, 0, 0, 0, 142, 141, 1, 0, 0, 0,
		143, 145, 1, 0, 0, 0, 144, 134, 1, 0, 0, 0, 145, 148, 1, 0, 0, 0, 146,
		144, 1, 0, 0, 0, 146, 147, 1, 0, 0, 0, 147, 149, 1, 0, 0, 0, 148, 146,
		1, 0, 0, 0, 149, 150, 5, 5, 0, 0, 150, 15, 1, 0, 0, 0, 151, 164, 3, 18,
		9, 0, 152, 164, 3, 24, 12, 0, 153, 164, 3, 26, 13, 0, 154, 164, 3, 28,
		14, 0, 155, 164, 3, 30, 15, 0, 156, 164, 3, 32, 16, 0, 157, 164, 3, 34,
		17, 0, 158, 164, 3, 36, 18, 0, 159, 164, 3, 38, 19, 0, 160, 164, 3, 40,
		20, 0, 161, 164, 3, 42, 21, 0, 162, 164, 3, 44, 22, 0, 163, 151, 1, 0,
		0, 0, 163, 152, 1, 0, 0, 0, 163, 153, 1, 0, 0, 0, 163, 154, 1, 0, 0, 0,
		163, 155, 1, 0, 0, 0, 163, 156, 1, 0, 0, 0, 163, 157, 1, 0, 0, 0, 163,
		158, 1, 0, 0, 0, 163, 159, 1, 0, 0, 0, 163, 160, 1, 0, 0, 0, 163, 161,
		1, 0, 0, 0, 163, 162, 1, 0, 0, 0, 164, 17, 1, 0, 0, 0, 165, 168, 3, 20,
		10, 0, 166, 168, 3, 22, 11, 0, 167, 165, 1, 0, 0, 0, 167, 166, 1, 0, 0,
		0, 168, 19, 1, 0, 0, 0, 169, 170, 5, 106, 0, 0, 170, 171, 5, 6, 0, 0, 171,
		208, 3, 44, 22, 0, 172, 173, 5, 106, 0, 0, 173, 174, 7, 0, 0, 0, 174, 208,
		3, 44, 22, 0, 175, 187, 5, 7, 0, 0, 176, 177, 5, 8, 0, 0, 177, 182, 3,
		44, 22, 0, 178, 179, 5, 9, 0, 0, 179, 181, 3, 44, 22, 0, 180, 178, 1, 0,
		0, 0, 181, 184, 1, 0, 0, 0, 182, 180, 1, 0, 0, 0, 182, 183, 1, 0, 0, 0,
		183, 185, 1, 0, 0, 0, 184, 182, 1, 0, 0, 0, 185, 186, 5, 10, 0, 0, 186,
		188, 1, 0, 0, 0, 187, 176, 1, 0, 0, 0, 187, 188, 1, 0, 0, 0, 188, 189,
		1, 0, 0, 0, 189, 190, 5, 6, 0, 0, 190, 208, 3, 44, 22, 0, 191, 203, 5,
		7, 0, 0, 192, 193, 5, 8, 0, 0, 193, 198, 3, 44, 22, 0, 194, 195, 5, 9,
		0, 0, 195, 197, 3, 44, 22, 0, 196, 194, 1, 0, 0, 0, 197, 200, 1, 0, 0,
		0, 198, 196, 1, 0, 0, 0, 198, 199, 1, 0, 0, 0, 199, 201, 1, 0, 0, 0, 200,
		198, 1, 0, 0, 0, 201, 202, 5, 10, 0, 0, 202, 204, 1, 0, 0, 0, 203, 192,
		1, 0, 0, 0, 203, 204, 1, 0, 0, 0, 204, 205, 1, 0, 0, 0, 205, 206, 7, 0,
		0, 0, 206, 208, 3, 44, 22, 0, 207, 169, 1, 0, 0, 0, 207, 172, 1, 0, 0,
		0, 207, 175, 1, 0, 0, 0, 207, 191, 1, 0, 0, 0, 208, 21, 1, 0, 0, 0, 209,
		210, 5, 105, 0, 0, 210, 211, 5, 6, 0, 0, 211, 216, 3, 44, 22, 0, 212, 213,
		5, 105, 0, 0, 213, 214, 7, 0, 0, 0, 214, 216, 3, 44, 22, 0, 215, 209, 1,
		0, 0, 0, 215, 212, 1, 0, 0, 0, 216, 23, 1, 0, 0, 0, 217, 218, 3, 74, 37,
		0, 218, 220, 5, 11, 0, 0, 219, 221, 3, 70, 35, 0, 220, 219, 1, 0, 0, 0,
		220, 221, 1, 0, 0, 0, 221, 222, 1, 0, 0, 0, 222, 223, 5, 12, 0, 0, 223,
		232, 1, 0, 0, 0, 224, 225, 3, 76, 38, 0, 225, 227, 5, 11, 0, 0, 226, 228,
		3, 70, 35, 0, 227, 226, 1, 0, 0, 0, 227, 228, 1, 0, 0, 0, 228, 229, 1,
		0, 0, 0, 229, 230, 5, 12, 0, 0, 230, 232, 1, 0, 0, 0, 231, 217, 1, 0, 0,
		0, 231, 224, 1, 0, 0, 0, 232, 25, 1, 0, 0, 0, 233, 234, 5, 68, 0, 0, 234,
		235, 5, 11, 0, 0, 235, 236, 3, 44, 22, 0, 236, 237, 5, 12, 0, 0, 237, 240,
		3, 14, 7, 0, 238, 239, 5, 69, 0, 0, 239, 241, 3, 14, 7, 0, 240, 238, 1,
		0, 0, 0, 240, 241, 1, 0, 0, 0, 241, 27, 1, 0, 0, 0, 242, 243, 5, 70, 0,
		0, 243, 244, 5, 11, 0, 0, 244, 245, 3, 44, 22, 0, 245, 246, 5, 12, 0, 0,
		246, 247, 3, 14, 7, 0, 247, 29, 1, 0, 0, 0, 248, 249, 5, 71, 0, 0, 249,
		250, 5, 11, 0, 0, 250, 251, 5, 105, 0, 0, 251, 252, 5, 72, 0, 0, 252, 253,
		5, 106, 0, 0, 253, 254, 5, 12, 0, 0, 254, 271, 3, 14, 7, 0, 255, 256, 5,
		71, 0, 0, 256, 258, 5, 11, 0, 0, 257, 259, 3, 18, 9, 0, 258, 257, 1, 0,
		0, 0, 258, 259, 1, 0, 0, 0, 259, 260, 1, 0, 0, 0, 260, 262, 5, 4, 0, 0,
		261, 263, 3, 44, 22, 0, 262, 261, 1, 0, 0, 0, 262, 263, 1, 0, 0, 0, 263,
		264, 1, 0, 0, 0, 264, 266, 5, 4, 0, 0, 265, 267, 3, 18, 9, 0, 266, 265,
		1, 0, 0, 0, 266, 267, 1, 0, 0, 0, 267, 268, 1, 0, 0, 0, 268, 269, 5, 12,
		0, 0, 269, 271, 3, 14, 7, 0, 270, 248, 1, 0, 0, 0, 270, 255, 1, 0, 0, 0,
		271, 31, 1, 0, 0, 0, 272, 274, 5, 73, 0, 0, 273, 275, 3, 44, 22, 0, 274,
		273, 1, 0, 0, 0, 274, 275, 1, 0, 0, 0, 275, 33, 1, 0, 0, 0, 276, 277, 5,
		74, 0, 0, 277, 294, 5, 106, 0, 0, 278, 279, 5, 74, 0, 0, 279, 291, 5, 7,
		0, 0, 280, 281, 5, 8, 0, 0, 281, 286, 3, 44, 22, 0, 282, 283, 5, 9, 0,
		0, 283, 285, 3, 44, 22, 0, 284, 282, 1, 0, 0, 0, 285, 288, 1, 0, 0, 0,
		286, 284, 1, 0, 0, 0, 286, 287, 1, 0, 0, 0, 287, 289, 1, 0, 0, 0, 288,
		286, 1, 0, 0, 0, 289, 290, 5, 10, 0, 0, 290, 292, 1, 0, 0, 0, 291, 280,
		1, 0, 0, 0, 291, 292, 1, 0, 0, 0, 292, 294, 1, 0, 0, 0, 293, 276, 1, 0,
		0, 0, 293, 278, 1, 0, 0, 0, 294, 35, 1, 0, 0, 0, 295, 296, 5, 75, 0, 0,
		296, 297, 5, 106, 0, 0, 297, 298, 5, 8, 0, 0, 298, 303, 3, 44, 22, 0, 299,
		300, 5, 9, 0, 0, 300, 302, 3, 44, 22, 0, 301, 299, 1, 0, 0, 0, 302, 305,
		1, 0, 0, 0, 303, 301, 1, 0, 0, 0, 303, 304, 1, 0, 0, 0, 304, 306, 1, 0,
		0, 0, 305, 303, 1, 0, 0, 0, 306, 307, 5, 10, 0, 0, 307, 322, 1, 0, 0, 0,
		308, 309, 5, 75, 0, 0, 309, 310, 5, 7, 0, 0, 310, 311, 5, 8, 0, 0, 311,
		316, 3, 44, 22, 0, 312, 313, 5, 9, 0, 0, 313, 315, 3, 44, 22, 0, 314, 312,
		1, 0, 0, 0, 315, 318, 1, 0, 0, 0, 316, 314, 1, 0, 0, 0, 316, 317, 1, 0,
		0, 0, 317, 319, 1, 0, 0, 0, 318, 316, 1, 0, 0, 0, 319, 320, 5, 10, 0, 0,
		320, 322, 1, 0, 0, 0, 321, 295, 1, 0, 0, 0, 321, 308, 1, 0, 0, 0, 322,
		37, 1, 0, 0, 0, 323, 325, 5, 76, 0, 0, 324, 326, 3, 44, 22, 0, 325, 324,
		1, 0, 0, 0, 325, 326, 1, 0, 0, 0, 326, 39, 1, 0, 0, 0, 327, 329, 5, 77,
		0, 0, 328, 330, 3, 44, 22, 0, 329, 328, 1, 0, 0, 0, 329, 330, 1, 0, 0,
		0, 330, 332, 1, 0, 0, 0, 331, 333, 3, 72, 36, 0, 332, 331, 1, 0, 0, 0,
		332, 333, 1, 0, 0, 0, 333, 41, 1, 0, 0, 0, 334, 335, 5, 78, 0, 0, 335,
		336, 5, 11, 0, 0, 336, 341, 5, 103, 0, 0, 337, 338, 5, 9, 0, 0, 338, 340,
		3, 44, 22, 0, 339, 337, 1, 0, 0, 0, 340, 343, 1, 0, 0, 0, 341, 339, 1,
		0, 0, 0, 341, 342, 1, 0, 0, 0, 342, 344, 1, 0, 0, 0, 343, 341, 1, 0, 0,
		0, 344, 345, 5, 12, 0, 0, 345, 43, 1, 0, 0, 0, 346, 347, 3, 46, 23, 0,
		347, 45, 1, 0, 0, 0, 348, 353, 3, 48, 24, 0, 349, 350, 5, 87, 0, 0, 350,
		352, 3, 48, 24, 0, 351, 349, 1, 0, 0, 0, 352, 355, 1, 0, 0, 0, 353, 351,
		1, 0, 0, 0, 353, 354, 1, 0, 0, 0, 354, 47, 1, 0, 0, 0, 355, 353, 1, 0,
		0, 0, 356, 361, 3, 50, 25, 0, 357, 358, 5, 86, 0, 0, 358, 360, 3, 50, 25,
		0, 359, 357, 1, 0, 0, 0, 360, 363, 1, 0, 0, 0, 361, 359, 1, 0, 0, 0, 361,
		362, 1, 0, 0, 0, 362, 49, 1, 0, 0, 0, 363, 361, 1, 0, 0, 0, 364, 369, 3,
		52, 26, 0, 365, 366, 7, 1, 0, 0, 366, 368, 3, 52, 26, 0, 367, 365, 1, 0,
		0, 0, 368, 371, 1, 0, 0, 0, 369, 367, 1, 0, 0, 0, 369, 370, 1, 0, 0, 0,
		370, 51, 1, 0, 0, 0, 371, 369, 1, 0, 0, 0, 372, 377, 3, 54, 27, 0, 373,
		374, 7, 2, 0, 0, 374, 376, 3, 54, 27, 0, 375, 373, 1, 0, 0, 0, 376, 379,
		1, 0, 0, 0, 377, 375, 1, 0, 0, 0, 377, 378, 1, 0, 0, 0, 378, 53, 1, 0,
		0, 0, 379, 377, 1, 0, 0, 0, 380, 385, 3, 56, 28, 0, 381, 382, 7, 3, 0,
		0, 382, 384, 3, 56, 28, 0, 383, 381, 1, 0, 0, 0, 384, 387, 1, 0, 0, 0,
		385, 383, 1, 0, 0, 0, 385, 386, 1, 0, 0, 0, 386, 55, 1, 0, 0, 0, 387, 385,
		1, 0, 0, 0, 388, 393, 3, 58, 29, 0, 389, 390, 7, 4, 0, 0, 390, 392, 3,
		58, 29, 0, 391, 389, 1, 0, 0, 0, 392, 395, 1, 0, 0, 0, 393, 391, 1, 0,
		0, 0, 393, 394, 1, 0, 0, 0, 394, 57, 1, 0, 0, 0, 395, 393, 1, 0, 0, 0,
		396, 401, 3, 60, 30, 0, 397, 398, 7, 5, 0, 0, 398, 400, 3, 60, 30, 0, 399,
		397, 1, 0, 0, 0, 400, 403, 1, 0, 0, 0, 401, 399, 1, 0, 0, 0, 401, 402,
		1, 0, 0, 0, 402, 59, 1, 0, 0, 0, 403, 401, 1, 0, 0, 0, 404, 405, 7, 6,
		0, 0, 405, 410, 3, 60, 30, 0, 406, 407, 7, 7, 0, 0, 407, 410, 3, 66, 33,
		0, 408, 410, 3, 62, 31, 0, 409, 404, 1, 0, 0, 0, 409, 406, 1, 0, 0, 0,
		409, 408, 1, 0, 0, 0, 410, 61, 1, 0, 0, 0, 411, 431, 3, 64, 32, 0, 412,
		432, 7, 7, 0, 0, 413, 414, 5, 8, 0, 0, 414, 419, 3, 44, 22, 0, 415, 416,
		5, 9, 0, 0, 416, 418, 3, 44, 22, 0, 417, 415, 1, 0, 0, 0, 418, 421, 1,
		0, 0, 0, 419, 417, 1, 0, 0, 0, 419, 420, 1, 0, 0, 0, 420, 422, 1, 0, 0,
		0, 421, 419, 1, 0, 0, 0, 422, 423, 5, 10, 0, 0, 423, 432, 1, 0, 0, 0, 424,
		426, 5, 11, 0, 0, 425, 427, 3, 70, 35, 0, 426, 425, 1, 0, 0, 0, 426, 427,
		1, 0, 0, 0, 427, 428, 1, 0, 0, 0, 428, 432, 5, 12, 0, 0, 429, 430, 5, 19,
		0, 0, 430, 432, 5, 104, 0, 0, 431, 412, 1, 0, 0, 0, 431, 413, 1, 0, 0,
		0, 431, 424, 1, 0, 0, 0, 431, 429, 1, 0, 0, 0, 431, 432, 1, 0, 0, 0, 432,
		63, 1, 0, 0, 0, 433, 443, 5, 102, 0, 0, 434, 443, 3, 80, 40, 0, 435, 443,
		3, 66, 33, 0, 436, 443, 3, 68, 34, 0, 437, 438, 5, 11, 0, 0, 438, 439,
		3, 44, 22, 0, 439, 440, 5, 12, 0, 0, 440, 443, 1, 0, 0, 0, 441, 443, 3,
		24, 12, 0, 442, 433, 1, 0, 0, 0, 442, 434, 1, 0, 0, 0, 442, 435, 1, 0,
		0, 0, 442, 436, 1, 0, 0, 0, 442, 437, 1, 0, 0, 0, 442, 441, 1, 0, 0, 0,
		443, 65, 1, 0, 0, 0, 444, 445, 7, 8, 0, 0, 445, 67, 1, 0, 0, 0, 446, 447,
		5, 106, 0, 0, 447, 448, 5, 8, 0, 0, 448, 453, 3, 44, 22, 0, 449, 450, 5,
		9, 0, 0, 450, 452, 3, 44, 22, 0, 451, 449, 1, 0, 0, 0, 452, 455, 1, 0,
		0, 0, 453, 451, 1, 0, 0, 0, 453, 454, 1, 0, 0, 0, 454, 456, 1, 0, 0, 0,
		455, 453, 1, 0, 0, 0, 456, 457, 5, 10, 0, 0, 457, 473, 1, 0, 0, 0, 458,
		470, 5, 7, 0, 0, 459, 460, 5, 8, 0, 0, 460, 465, 3, 44, 22, 0, 461, 462,
		5, 9, 0, 0, 462, 464, 3, 44, 22, 0, 463, 461, 1, 0, 0, 0, 464, 467, 1,
		0, 0, 0, 465, 463, 1, 0, 0, 0, 465, 466, 1, 0, 0, 0, 466, 468, 1, 0, 0,
		0, 467, 465, 1, 0, 0, 0, 468, 469, 5, 10, 0, 0, 469, 471, 1, 0, 0, 0, 470,
		459, 1, 0, 0, 0, 470, 471, 1, 0, 0, 0, 471, 473, 1, 0, 0, 0, 472, 446,
		1, 0, 0, 0, 472, 458, 1, 0, 0, 0, 473, 69, 1, 0, 0, 0, 474, 479, 3, 44,
		22, 0, 475, 476, 5, 9, 0, 0, 476, 478, 3, 44, 22, 0, 477, 475, 1, 0, 0,
		0, 478, 481, 1, 0, 0, 0, 479, 477, 1, 0, 0, 0, 479, 480, 1, 0, 0, 0, 480,
		71, 1, 0, 0, 0, 481, 479, 1, 0, 0, 0, 482, 483, 5, 83, 0, 0, 483, 489,
		3, 44, 22, 0, 484, 485, 5, 89, 0, 0, 485, 489, 3, 44, 22, 0, 486, 487,
		5, 20, 0, 0, 487, 489, 3, 44, 22, 0, 488, 482, 1, 0, 0, 0, 488, 484, 1,
		0, 0, 0, 488, 486, 1, 0, 0, 0, 489, 73, 1, 0, 0, 0, 490, 491, 5, 104, 0,
		0, 491, 75, 1, 0, 0, 0, 492, 493, 7, 9, 0, 0, 493, 77, 1, 0, 0, 0, 494,
		495, 5, 64, 0, 0, 495, 79, 1, 0, 0, 0, 496, 497, 5, 103, 0, 0, 497, 81,
		1, 0, 0, 0, 56, 83, 92, 98, 100, 107, 113, 120, 131, 138, 142, 146, 163,
		167, 182, 187, 198, 203, 207, 215, 220, 227, 231, 240, 258, 262, 266, 270,
		274, 286, 291, 293, 303, 316, 321, 325, 329, 332, 341, 353, 361, 369, 377,
		385, 393, 401, 409, 419, 426, 431, 442, 453, 465, 470, 472, 479, 488,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// bpftraceParserInit initializes any static state used to implement bpftraceParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewbpftraceParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func BpftraceParserInit() {
	staticData := &BpftraceParserStaticData
	staticData.once.Do(bpftraceParserInit)
}

// NewbpftraceParser produces a new parser instance for the optional input antlr.TokenStream.
func NewbpftraceParser(input antlr.TokenStream) *bpftraceParser {
	BpftraceParserInit()
	this := new(bpftraceParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &BpftraceParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "bpftrace.g4"

	return this
}

// bpftraceParser tokens.
const (
	bpftraceParserEOF        = antlr.TokenEOF
	bpftraceParserT__0       = 1
	bpftraceParserT__1       = 2
	bpftraceParserT__2       = 3
	bpftraceParserT__3       = 4
	bpftraceParserT__4       = 5
	bpftraceParserT__5       = 6
	bpftraceParserT__6       = 7
	bpftraceParserT__7       = 8
	bpftraceParserT__8       = 9
	bpftraceParserT__9       = 10
	bpftraceParserT__10      = 11
	bpftraceParserT__11      = 12
	bpftraceParserT__12      = 13
	bpftraceParserT__13      = 14
	bpftraceParserT__14      = 15
	bpftraceParserT__15      = 16
	bpftraceParserT__16      = 17
	bpftraceParserT__17      = 18
	bpftraceParserT__18      = 19
	bpftraceParserT__19      = 20
	bpftraceParserT__20      = 21
	bpftraceParserT__21      = 22
	bpftraceParserT__22      = 23
	bpftraceParserT__23      = 24
	bpftraceParserT__24      = 25
	bpftraceParserT__25      = 26
	bpftraceParserT__26      = 27
	bpftraceParserT__27      = 28
	bpftraceParserT__28      = 29
	bpftraceParserT__29      = 30
	bpftraceParserT__30      = 31
	bpftraceParserT__31      = 32
	bpftraceParserT__32      = 33
	bpftraceParserT__33      = 34
	bpftraceParserT__34      = 35
	bpftraceParserT__35      = 36
	bpftraceParserT__36      = 37
	bpftraceParserT__37      = 38
	bpftraceParserT__38      = 39
	bpftraceParserT__39      = 40
	bpftraceParserT__40      = 41
	bpftraceParserT__41      = 42
	bpftraceParserT__42      = 43
	bpftraceParserT__43      = 44
	bpftraceParserT__44      = 45
	bpftraceParserT__45      = 46
	bpftraceParserT__46      = 47
	bpftraceParserT__47      = 48
	bpftraceParserT__48      = 49
	bpftraceParserT__49      = 50
	bpftraceParserT__50      = 51
	bpftraceParserT__51      = 52
	bpftraceParserT__52      = 53
	bpftraceParserT__53      = 54
	bpftraceParserT__54      = 55
	bpftraceParserT__55      = 56
	bpftraceParserT__56      = 57
	bpftraceParserT__57      = 58
	bpftraceParserT__58      = 59
	bpftraceParserT__59      = 60
	bpftraceParserT__60      = 61
	bpftraceParserT__61      = 62
	bpftraceParserSHEBANG    = 63
	bpftraceParserCOMMENT    = 64
	bpftraceParserWS         = 65
	bpftraceParserNEWLINE    = 66
	bpftraceParserPROBE_TYPE = 67
	bpftraceParserIF         = 68
	bpftraceParserELSE       = 69
	bpftraceParserWHILE      = 70
	bpftraceParserFOR        = 71
	bpftraceParserIN         = 72
	bpftraceParserRETURN     = 73
	bpftraceParserCLEAR      = 74
	bpftraceParserDELETE     = 75
	bpftraceParserEXIT       = 76
	bpftraceParserPRINT      = 77
	bpftraceParserPRINTF     = 78
	bpftraceParserEND        = 79
	bpftraceParserEQ         = 80
	bpftraceParserNE         = 81
	bpftraceParserLT         = 82
	bpftraceParserGT         = 83
	bpftraceParserLE         = 84
	bpftraceParserGE         = 85
	bpftraceParserAND        = 86
	bpftraceParserOR         = 87
	bpftraceParserSHL        = 88
	bpftraceParserSHR        = 89
	bpftraceParserADD_ASSIGN = 90
	bpftraceParserSUB_ASSIGN = 91
	bpftraceParserMUL_ASSIGN = 92
	bpftraceParserDIV_ASSIGN = 93
	bpftraceParserMOD_ASSIGN = 94
	bpftraceParserAND_ASSIGN = 95
	bpftraceParserOR_ASSIGN  = 96
	bpftraceParserXOR_ASSIGN = 97
	bpftraceParserSHL_ASSIGN = 98
	bpftraceParserSHR_ASSIGN = 99
	bpftraceParserINCR       = 100
	bpftraceParserDECR       = 101
	bpftraceParserNUMBER     = 102
	bpftraceParserSTRING     = 103
	bpftraceParserIDENTIFIER = 104
	bpftraceParserVARIABLE   = 105
	bpftraceParserMAP_NAME   = 106
)

// bpftraceParser rules.
const (
	bpftraceParserRULE_program                   = 0
	bpftraceParserRULE_shebang_section           = 1
	bpftraceParserRULE_content                   = 2
	bpftraceParserRULE_shebang                   = 3
	bpftraceParserRULE_probe                     = 4
	bpftraceParserRULE_probe_def                 = 5
	bpftraceParserRULE_predicate                 = 6
	bpftraceParserRULE_block                     = 7
	bpftraceParserRULE_statement                 = 8
	bpftraceParserRULE_assignment                = 9
	bpftraceParserRULE_map_assign                = 10
	bpftraceParserRULE_var_assign                = 11
	bpftraceParserRULE_function_call             = 12
	bpftraceParserRULE_if_statement              = 13
	bpftraceParserRULE_while_statement           = 14
	bpftraceParserRULE_for_statement             = 15
	bpftraceParserRULE_return_statement          = 16
	bpftraceParserRULE_clear_statement           = 17
	bpftraceParserRULE_delete_statement          = 18
	bpftraceParserRULE_exit_statement            = 19
	bpftraceParserRULE_print_statement           = 20
	bpftraceParserRULE_printf_statement          = 21
	bpftraceParserRULE_expression                = 22
	bpftraceParserRULE_logical_or_expression     = 23
	bpftraceParserRULE_logical_and_expression    = 24
	bpftraceParserRULE_equality_expression       = 25
	bpftraceParserRULE_relational_expression     = 26
	bpftraceParserRULE_shift_expression          = 27
	bpftraceParserRULE_additive_expression       = 28
	bpftraceParserRULE_multiplicative_expression = 29
	bpftraceParserRULE_unary_expression          = 30
	bpftraceParserRULE_postfix_expression        = 31
	bpftraceParserRULE_primary_expression        = 32
	bpftraceParserRULE_variable                  = 33
	bpftraceParserRULE_map_access                = 34
	bpftraceParserRULE_expr_list                 = 35
	bpftraceParserRULE_output_redirection        = 36
	bpftraceParserRULE_function_name             = 37
	bpftraceParserRULE_builtin_name              = 38
	bpftraceParserRULE_comment                   = 39
	bpftraceParserRULE_string                    = 40
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Content() IContentContext
	EOF() antlr.TerminalNode
	Shebang_section() IShebang_sectionContext

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_program
	return p
}

func InitEmptyProgramContext(p *ProgramContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_program
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) Content() IContentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContentContext)
}

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(bpftraceParserEOF, 0)
}

func (s *ProgramContext) Shebang_section() IShebang_sectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShebang_sectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShebang_sectionContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (p *bpftraceParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, bpftraceParserRULE_program)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(83)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == bpftraceParserSHEBANG {
		{
			p.SetState(82)
			p.Shebang_section()
		}

	}
	{
		p.SetState(85)
		p.Content()
	}
	{
		p.SetState(86)
		p.Match(bpftraceParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IShebang_sectionContext is an interface to support dynamic dispatch.
type IShebang_sectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Shebang() IShebangContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsShebang_sectionContext differentiates from other interfaces.
	IsShebang_sectionContext()
}

type Shebang_sectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShebang_sectionContext() *Shebang_sectionContext {
	var p = new(Shebang_sectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_shebang_section
	return p
}

func InitEmptyShebang_sectionContext(p *Shebang_sectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_shebang_section
}

func (*Shebang_sectionContext) IsShebang_sectionContext() {}

func NewShebang_sectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Shebang_sectionContext {
	var p = new(Shebang_sectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_shebang_section

	return p
}

func (s *Shebang_sectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Shebang_sectionContext) Shebang() IShebangContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShebangContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShebangContext)
}

func (s *Shebang_sectionContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserNEWLINE)
}

func (s *Shebang_sectionContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserNEWLINE, i)
}

func (s *Shebang_sectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Shebang_sectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Shebang_sectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterShebang_section(s)
	}
}

func (s *Shebang_sectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitShebang_section(s)
	}
}

func (p *bpftraceParser) Shebang_section() (localctx IShebang_sectionContext) {
	localctx = NewShebang_sectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, bpftraceParserRULE_shebang_section)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(88)
		p.Shebang()
	}
	p.SetState(92)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(89)
				p.Match(bpftraceParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(94)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IContentContext is an interface to support dynamic dispatch.
type IContentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllProbe() []IProbeContext
	Probe(i int) IProbeContext
	AllComment() []ICommentContext
	Comment(i int) ICommentContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsContentContext differentiates from other interfaces.
	IsContentContext()
}

type ContentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyContentContext() *ContentContext {
	var p = new(ContentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_content
	return p
}

func InitEmptyContentContext(p *ContentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_content
}

func (*ContentContext) IsContentContext() {}

func NewContentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContentContext {
	var p = new(ContentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_content

	return p
}

func (s *ContentContext) GetParser() antlr.Parser { return s.parser }

func (s *ContentContext) AllProbe() []IProbeContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IProbeContext); ok {
			len++
		}
	}

	tst := make([]IProbeContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IProbeContext); ok {
			tst[i] = t.(IProbeContext)
			i++
		}
	}

	return tst
}

func (s *ContentContext) Probe(i int) IProbeContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IProbeContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IProbeContext)
}

func (s *ContentContext) AllComment() []ICommentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICommentContext); ok {
			len++
		}
	}

	tst := make([]ICommentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICommentContext); ok {
			tst[i] = t.(ICommentContext)
			i++
		}
	}

	return tst
}

func (s *ContentContext) Comment(i int) ICommentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommentContext)
}

func (s *ContentContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserNEWLINE)
}

func (s *ContentContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserNEWLINE, i)
}

func (s *ContentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ContentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterContent(s)
	}
}

func (s *ContentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitContent(s)
	}
}

func (p *bpftraceParser) Content() (localctx IContentContext) {
	localctx = NewContentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, bpftraceParserRULE_content)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&32781) != 0 {
		p.SetState(98)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case bpftraceParserPROBE_TYPE, bpftraceParserEND:
			{
				p.SetState(95)
				p.Probe()
			}

		case bpftraceParserCOMMENT:
			{
				p.SetState(96)
				p.Comment()
			}

		case bpftraceParserNEWLINE:
			{
				p.SetState(97)
				p.Match(bpftraceParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(102)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IShebangContext is an interface to support dynamic dispatch.
type IShebangContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SHEBANG() antlr.TerminalNode

	// IsShebangContext differentiates from other interfaces.
	IsShebangContext()
}

type ShebangContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShebangContext() *ShebangContext {
	var p = new(ShebangContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_shebang
	return p
}

func InitEmptyShebangContext(p *ShebangContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_shebang
}

func (*ShebangContext) IsShebangContext() {}

func NewShebangContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShebangContext {
	var p = new(ShebangContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_shebang

	return p
}

func (s *ShebangContext) GetParser() antlr.Parser { return s.parser }

func (s *ShebangContext) SHEBANG() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHEBANG, 0)
}

func (s *ShebangContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShebangContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShebangContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterShebang(s)
	}
}

func (s *ShebangContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitShebang(s)
	}
}

func (p *bpftraceParser) Shebang() (localctx IShebangContext) {
	localctx = NewShebangContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, bpftraceParserRULE_shebang)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Match(bpftraceParserSHEBANG)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IProbeContext is an interface to support dynamic dispatch.
type IProbeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Probe_def() IProbe_defContext
	Block() IBlockContext
	Predicate() IPredicateContext
	END() antlr.TerminalNode

	// IsProbeContext differentiates from other interfaces.
	IsProbeContext()
}

type ProbeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProbeContext() *ProbeContext {
	var p = new(ProbeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_probe
	return p
}

func InitEmptyProbeContext(p *ProbeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_probe
}

func (*ProbeContext) IsProbeContext() {}

func NewProbeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProbeContext {
	var p = new(ProbeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_probe

	return p
}

func (s *ProbeContext) GetParser() antlr.Parser { return s.parser }

func (s *ProbeContext) Probe_def() IProbe_defContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IProbe_defContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IProbe_defContext)
}

func (s *ProbeContext) Block() IBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *ProbeContext) Predicate() IPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPredicateContext)
}

func (s *ProbeContext) END() antlr.TerminalNode {
	return s.GetToken(bpftraceParserEND, 0)
}

func (s *ProbeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProbeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProbeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterProbe(s)
	}
}

func (s *ProbeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitProbe(s)
	}
}

func (p *bpftraceParser) Probe() (localctx IProbeContext) {
	localctx = NewProbeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, bpftraceParserRULE_probe)
	var _la int

	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case bpftraceParserPROBE_TYPE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(105)
			p.Probe_def()
		}
		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == bpftraceParserT__1 {
			{
				p.SetState(106)
				p.Predicate()
			}

		}
		{
			p.SetState(109)
			p.Block()
		}

	case bpftraceParserEND:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(111)
			p.Match(bpftraceParserEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(112)
			p.Block()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IProbe_defContext is an interface to support dynamic dispatch.
type IProbe_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPROBE_TYPE() []antlr.TerminalNode
	PROBE_TYPE(i int) antlr.TerminalNode

	// IsProbe_defContext differentiates from other interfaces.
	IsProbe_defContext()
}

type Probe_defContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProbe_defContext() *Probe_defContext {
	var p = new(Probe_defContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_probe_def
	return p
}

func InitEmptyProbe_defContext(p *Probe_defContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_probe_def
}

func (*Probe_defContext) IsProbe_defContext() {}

func NewProbe_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Probe_defContext {
	var p = new(Probe_defContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_probe_def

	return p
}

func (s *Probe_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Probe_defContext) AllPROBE_TYPE() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserPROBE_TYPE)
}

func (s *Probe_defContext) PROBE_TYPE(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserPROBE_TYPE, i)
}

func (s *Probe_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Probe_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Probe_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterProbe_def(s)
	}
}

func (s *Probe_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitProbe_def(s)
	}
}

func (p *bpftraceParser) Probe_def() (localctx IProbe_defContext) {
	localctx = NewProbe_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, bpftraceParserRULE_probe_def)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(115)
		p.Match(bpftraceParserPROBE_TYPE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(120)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == bpftraceParserT__0 {
		{
			p.SetState(116)
			p.Match(bpftraceParserT__0)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(117)
			p.Match(bpftraceParserPROBE_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(122)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPredicateContext is an interface to support dynamic dispatch.
type IPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsPredicateContext differentiates from other interfaces.
	IsPredicateContext()
}

type PredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPredicateContext() *PredicateContext {
	var p = new(PredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_predicate
	return p
}

func InitEmptyPredicateContext(p *PredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_predicate
}

func (*PredicateContext) IsPredicateContext() {}

func NewPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PredicateContext {
	var p = new(PredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_predicate

	return p
}

func (s *PredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PredicateContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *PredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterPredicate(s)
	}
}

func (s *PredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitPredicate(s)
	}
}

func (p *bpftraceParser) Predicate() (localctx IPredicateContext) {
	localctx = NewPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, bpftraceParserRULE_predicate)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(123)
		p.Match(bpftraceParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(124)
		p.Expression()
	}
	{
		p.SetState(125)
		p.Match(bpftraceParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBlockContext is an interface to support dynamic dispatch.
type IBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext

	// IsBlockContext differentiates from other interfaces.
	IsBlockContext()
}

type BlockContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBlockContext() *BlockContext {
	var p = new(BlockContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_block
	return p
}

func InitEmptyBlockContext(p *BlockContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_block
}

func (*BlockContext) IsBlockContext() {}

func NewBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlockContext {
	var p = new(BlockContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_block

	return p
}

func (s *BlockContext) GetParser() antlr.Parser { return s.parser }

func (s *BlockContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserNEWLINE)
}

func (s *BlockContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserNEWLINE, i)
}

func (s *BlockContext) AllStatement() []IStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatementContext); ok {
			len++
		}
	}

	tst := make([]IStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatementContext); ok {
			tst[i] = t.(IStatementContext)
			i++
		}
	}

	return tst
}

func (s *BlockContext) Statement(i int) IStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *BlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterBlock(s)
	}
}

func (s *BlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitBlock(s)
	}
}

func (p *bpftraceParser) Block() (localctx IBlockContext) {
	localctx = NewBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, bpftraceParserRULE_block)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(127)
		p.Match(bpftraceParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == bpftraceParserNEWLINE {
		{
			p.SetState(128)
			p.Match(bpftraceParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(133)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&9223372036853098624) != 0) || ((int64((_la-68)) & ^0x3f) == 0 && ((int64(1)<<(_la-68))&545460848621) != 0) {
		{
			p.SetState(134)
			p.Statement()
		}
		p.SetState(142)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case bpftraceParserT__4, bpftraceParserT__6, bpftraceParserT__10, bpftraceParserT__12, bpftraceParserT__13, bpftraceParserT__16, bpftraceParserT__17, bpftraceParserT__20, bpftraceParserT__21, bpftraceParserT__22, bpftraceParserT__23, bpftraceParserT__24, bpftraceParserT__25, bpftraceParserT__26, bpftraceParserT__27, bpftraceParserT__28, bpftraceParserT__29, bpftraceParserT__30, bpftraceParserT__31, bpftraceParserT__32, bpftraceParserT__33, bpftraceParserT__34, bpftraceParserT__35, bpftraceParserT__36, bpftraceParserT__37, bpftraceParserT__38, bpftraceParserT__39, bpftraceParserT__40, bpftraceParserT__41, bpftraceParserT__42, bpftraceParserT__43, bpftraceParserT__44, bpftraceParserT__45, bpftraceParserT__46, bpftraceParserT__47, bpftraceParserT__48, bpftraceParserT__49, bpftraceParserT__50, bpftraceParserT__51, bpftraceParserT__52, bpftraceParserT__53, bpftraceParserT__54, bpftraceParserT__55, bpftraceParserT__56, bpftraceParserT__57, bpftraceParserT__58, bpftraceParserT__59, bpftraceParserT__60, bpftraceParserT__61, bpftraceParserNEWLINE, bpftraceParserIF, bpftraceParserWHILE, bpftraceParserFOR, bpftraceParserRETURN, bpftraceParserCLEAR, bpftraceParserDELETE, bpftraceParserEXIT, bpftraceParserPRINT, bpftraceParserPRINTF, bpftraceParserINCR, bpftraceParserDECR, bpftraceParserNUMBER, bpftraceParserSTRING, bpftraceParserIDENTIFIER, bpftraceParserVARIABLE, bpftraceParserMAP_NAME:
			p.SetState(138)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == bpftraceParserNEWLINE {
				{
					p.SetState(135)
					p.Match(bpftraceParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(140)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		case bpftraceParserT__3:
			{
				p.SetState(141)
				p.Match(bpftraceParserT__3)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(148)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(149)
		p.Match(bpftraceParserT__4)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Assignment() IAssignmentContext
	Function_call() IFunction_callContext
	If_statement() IIf_statementContext
	While_statement() IWhile_statementContext
	For_statement() IFor_statementContext
	Return_statement() IReturn_statementContext
	Clear_statement() IClear_statementContext
	Delete_statement() IDelete_statementContext
	Exit_statement() IExit_statementContext
	Print_statement() IPrint_statementContext
	Printf_statement() IPrintf_statementContext
	Expression() IExpressionContext

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_statement
	return p
}

func InitEmptyStatementContext(p *StatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_statement
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) Assignment() IAssignmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignmentContext)
}

func (s *StatementContext) Function_call() IFunction_callContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunction_callContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunction_callContext)
}

func (s *StatementContext) If_statement() IIf_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIf_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIf_statementContext)
}

func (s *StatementContext) While_statement() IWhile_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhile_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhile_statementContext)
}

func (s *StatementContext) For_statement() IFor_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFor_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFor_statementContext)
}

func (s *StatementContext) Return_statement() IReturn_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReturn_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReturn_statementContext)
}

func (s *StatementContext) Clear_statement() IClear_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClear_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClear_statementContext)
}

func (s *StatementContext) Delete_statement() IDelete_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDelete_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDelete_statementContext)
}

func (s *StatementContext) Exit_statement() IExit_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExit_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExit_statementContext)
}

func (s *StatementContext) Print_statement() IPrint_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrint_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrint_statementContext)
}

func (s *StatementContext) Printf_statement() IPrintf_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrintf_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrintf_statementContext)
}

func (s *StatementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (p *bpftraceParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, bpftraceParserRULE_statement)
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(151)
			p.Assignment()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(152)
			p.Function_call()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(153)
			p.If_statement()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(154)
			p.While_statement()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(155)
			p.For_statement()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(156)
			p.Return_statement()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(157)
			p.Clear_statement()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(158)
			p.Delete_statement()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(159)
			p.Exit_statement()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(160)
			p.Print_statement()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(161)
			p.Printf_statement()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(162)
			p.Expression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAssignmentContext is an interface to support dynamic dispatch.
type IAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Map_assign() IMap_assignContext
	Var_assign() IVar_assignContext

	// IsAssignmentContext differentiates from other interfaces.
	IsAssignmentContext()
}

type AssignmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAssignmentContext() *AssignmentContext {
	var p = new(AssignmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_assignment
	return p
}

func InitEmptyAssignmentContext(p *AssignmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_assignment
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_assignment

	return p
}

func (s *AssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignmentContext) Map_assign() IMap_assignContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMap_assignContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMap_assignContext)
}

func (s *AssignmentContext) Var_assign() IVar_assignContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVar_assignContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVar_assignContext)
}

func (s *AssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (p *bpftraceParser) Assignment() (localctx IAssignmentContext) {
	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, bpftraceParserRULE_assignment)
	p.SetState(167)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case bpftraceParserT__6, bpftraceParserMAP_NAME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(165)
			p.Map_assign()
		}

	case bpftraceParserVARIABLE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(166)
			p.Var_assign()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMap_assignContext is an interface to support dynamic dispatch.
type IMap_assignContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MAP_NAME() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	ADD_ASSIGN() antlr.TerminalNode
	SUB_ASSIGN() antlr.TerminalNode
	MUL_ASSIGN() antlr.TerminalNode
	DIV_ASSIGN() antlr.TerminalNode
	MOD_ASSIGN() antlr.TerminalNode
	AND_ASSIGN() antlr.TerminalNode
	OR_ASSIGN() antlr.TerminalNode
	XOR_ASSIGN() antlr.TerminalNode
	SHL_ASSIGN() antlr.TerminalNode
	SHR_ASSIGN() antlr.TerminalNode

	// IsMap_assignContext differentiates from other interfaces.
	IsMap_assignContext()
}

type Map_assignContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMap_assignContext() *Map_assignContext {
	var p = new(Map_assignContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_map_assign
	return p
}

func InitEmptyMap_assignContext(p *Map_assignContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_map_assign
}

func (*Map_assignContext) IsMap_assignContext() {}

func NewMap_assignContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Map_assignContext {
	var p = new(Map_assignContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_map_assign

	return p
}

func (s *Map_assignContext) GetParser() antlr.Parser { return s.parser }

func (s *Map_assignContext) MAP_NAME() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMAP_NAME, 0)
}

func (s *Map_assignContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Map_assignContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Map_assignContext) ADD_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserADD_ASSIGN, 0)
}

func (s *Map_assignContext) SUB_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSUB_ASSIGN, 0)
}

func (s *Map_assignContext) MUL_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMUL_ASSIGN, 0)
}

func (s *Map_assignContext) DIV_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserDIV_ASSIGN, 0)
}

func (s *Map_assignContext) MOD_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMOD_ASSIGN, 0)
}

func (s *Map_assignContext) AND_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserAND_ASSIGN, 0)
}

func (s *Map_assignContext) OR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserOR_ASSIGN, 0)
}

func (s *Map_assignContext) XOR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserXOR_ASSIGN, 0)
}

func (s *Map_assignContext) SHL_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHL_ASSIGN, 0)
}

func (s *Map_assignContext) SHR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHR_ASSIGN, 0)
}

func (s *Map_assignContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Map_assignContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Map_assignContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterMap_assign(s)
	}
}

func (s *Map_assignContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitMap_assign(s)
	}
}

func (p *bpftraceParser) Map_assign() (localctx IMap_assignContext) {
	localctx = NewMap_assignContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, bpftraceParserRULE_map_assign)
	var _la int

	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(169)
			p.Match(bpftraceParserMAP_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(170)
			p.Match(bpftraceParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(171)
			p.Expression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(172)
			p.Match(bpftraceParserMAP_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(173)
			_la = p.GetTokenStream().LA(1)

			if !((int64((_la-90)) & ^0x3f) == 0 && ((int64(1)<<(_la-90))&1023) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(174)
			p.Expression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(175)
			p.Match(bpftraceParserT__6)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(187)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == bpftraceParserT__7 {
			{
				p.SetState(176)
				p.Match(bpftraceParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(177)
				p.Expression()
			}
			p.SetState(182)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == bpftraceParserT__8 {
				{
					p.SetState(178)
					p.Match(bpftraceParserT__8)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(179)
					p.Expression()
				}

				p.SetState(184)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(185)
				p.Match(bpftraceParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(189)
			p.Match(bpftraceParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(190)
			p.Expression()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(191)
			p.Match(bpftraceParserT__6)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(203)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == bpftraceParserT__7 {
			{
				p.SetState(192)
				p.Match(bpftraceParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(193)
				p.Expression()
			}
			p.SetState(198)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == bpftraceParserT__8 {
				{
					p.SetState(194)
					p.Match(bpftraceParserT__8)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(195)
					p.Expression()
				}

				p.SetState(200)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(201)
				p.Match(bpftraceParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(205)
			_la = p.GetTokenStream().LA(1)

			if !((int64((_la-90)) & ^0x3f) == 0 && ((int64(1)<<(_la-90))&1023) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(206)
			p.Expression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVar_assignContext is an interface to support dynamic dispatch.
type IVar_assignContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VARIABLE() antlr.TerminalNode
	Expression() IExpressionContext
	ADD_ASSIGN() antlr.TerminalNode
	SUB_ASSIGN() antlr.TerminalNode
	MUL_ASSIGN() antlr.TerminalNode
	DIV_ASSIGN() antlr.TerminalNode
	MOD_ASSIGN() antlr.TerminalNode
	AND_ASSIGN() antlr.TerminalNode
	OR_ASSIGN() antlr.TerminalNode
	XOR_ASSIGN() antlr.TerminalNode
	SHL_ASSIGN() antlr.TerminalNode
	SHR_ASSIGN() antlr.TerminalNode

	// IsVar_assignContext differentiates from other interfaces.
	IsVar_assignContext()
}

type Var_assignContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVar_assignContext() *Var_assignContext {
	var p = new(Var_assignContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_var_assign
	return p
}

func InitEmptyVar_assignContext(p *Var_assignContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_var_assign
}

func (*Var_assignContext) IsVar_assignContext() {}

func NewVar_assignContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Var_assignContext {
	var p = new(Var_assignContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_var_assign

	return p
}

func (s *Var_assignContext) GetParser() antlr.Parser { return s.parser }

func (s *Var_assignContext) VARIABLE() antlr.TerminalNode {
	return s.GetToken(bpftraceParserVARIABLE, 0)
}

func (s *Var_assignContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Var_assignContext) ADD_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserADD_ASSIGN, 0)
}

func (s *Var_assignContext) SUB_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSUB_ASSIGN, 0)
}

func (s *Var_assignContext) MUL_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMUL_ASSIGN, 0)
}

func (s *Var_assignContext) DIV_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserDIV_ASSIGN, 0)
}

func (s *Var_assignContext) MOD_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMOD_ASSIGN, 0)
}

func (s *Var_assignContext) AND_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserAND_ASSIGN, 0)
}

func (s *Var_assignContext) OR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserOR_ASSIGN, 0)
}

func (s *Var_assignContext) XOR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserXOR_ASSIGN, 0)
}

func (s *Var_assignContext) SHL_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHL_ASSIGN, 0)
}

func (s *Var_assignContext) SHR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHR_ASSIGN, 0)
}

func (s *Var_assignContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Var_assignContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Var_assignContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterVar_assign(s)
	}
}

func (s *Var_assignContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitVar_assign(s)
	}
}

func (p *bpftraceParser) Var_assign() (localctx IVar_assignContext) {
	localctx = NewVar_assignContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, bpftraceParserRULE_var_assign)
	var _la int

	p.SetState(215)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(209)
			p.Match(bpftraceParserVARIABLE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(210)
			p.Match(bpftraceParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(211)
			p.Expression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(212)
			p.Match(bpftraceParserVARIABLE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(213)
			_la = p.GetTokenStream().LA(1)

			if !((int64((_la-90)) & ^0x3f) == 0 && ((int64(1)<<(_la-90))&1023) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(214)
			p.Expression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunction_callContext is an interface to support dynamic dispatch.
type IFunction_callContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Function_name() IFunction_nameContext
	Expr_list() IExpr_listContext
	Builtin_name() IBuiltin_nameContext

	// IsFunction_callContext differentiates from other interfaces.
	IsFunction_callContext()
}

type Function_callContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunction_callContext() *Function_callContext {
	var p = new(Function_callContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_function_call
	return p
}

func InitEmptyFunction_callContext(p *Function_callContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_function_call
}

func (*Function_callContext) IsFunction_callContext() {}

func NewFunction_callContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Function_callContext {
	var p = new(Function_callContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_function_call

	return p
}

func (s *Function_callContext) GetParser() antlr.Parser { return s.parser }

func (s *Function_callContext) Function_name() IFunction_nameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunction_nameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunction_nameContext)
}

func (s *Function_callContext) Expr_list() IExpr_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpr_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpr_listContext)
}

func (s *Function_callContext) Builtin_name() IBuiltin_nameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBuiltin_nameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBuiltin_nameContext)
}

func (s *Function_callContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Function_callContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Function_callContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterFunction_call(s)
	}
}

func (s *Function_callContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitFunction_call(s)
	}
}

func (p *bpftraceParser) Function_call() (localctx IFunction_callContext) {
	localctx = NewFunction_callContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, bpftraceParserRULE_function_call)
	var _la int

	p.SetState(231)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case bpftraceParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(217)
			p.Function_name()
		}
		{
			p.SetState(218)
			p.Match(bpftraceParserT__10)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(220)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&9223372036853098624) != 0) || ((int64((_la-74)) & ^0x3f) == 0 && ((int64(1)<<(_la-74))&8522825759) != 0) {
			{
				p.SetState(219)
				p.Expr_list()
			}

		}
		{
			p.SetState(222)
			p.Match(bpftraceParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case bpftraceParserT__20, bpftraceParserT__21, bpftraceParserT__22, bpftraceParserT__23, bpftraceParserT__24, bpftraceParserT__25, bpftraceParserT__26, bpftraceParserT__27, bpftraceParserT__28, bpftraceParserT__29, bpftraceParserT__30, bpftraceParserT__31, bpftraceParserT__32, bpftraceParserT__33, bpftraceParserT__34, bpftraceParserT__35, bpftraceParserT__36, bpftraceParserT__37, bpftraceParserT__38, bpftraceParserT__39, bpftraceParserT__40, bpftraceParserT__41, bpftraceParserT__42, bpftraceParserT__43, bpftraceParserT__44, bpftraceParserT__45, bpftraceParserT__46, bpftraceParserT__47, bpftraceParserT__48, bpftraceParserT__49, bpftraceParserT__50, bpftraceParserT__51, bpftraceParserT__52, bpftraceParserT__53, bpftraceParserT__54, bpftraceParserT__55, bpftraceParserT__56, bpftraceParserT__57, bpftraceParserT__58, bpftraceParserT__59, bpftraceParserT__60, bpftraceParserT__61, bpftraceParserCLEAR, bpftraceParserDELETE, bpftraceParserEXIT, bpftraceParserPRINT, bpftraceParserPRINTF:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(224)
			p.Builtin_name()
		}
		{
			p.SetState(225)
			p.Match(bpftraceParserT__10)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(227)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&9223372036853098624) != 0) || ((int64((_la-74)) & ^0x3f) == 0 && ((int64(1)<<(_la-74))&8522825759) != 0) {
			{
				p.SetState(226)
				p.Expr_list()
			}

		}
		{
			p.SetState(229)
			p.Match(bpftraceParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIf_statementContext is an interface to support dynamic dispatch.
type IIf_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IF() antlr.TerminalNode
	Expression() IExpressionContext
	AllBlock() []IBlockContext
	Block(i int) IBlockContext
	ELSE() antlr.TerminalNode

	// IsIf_statementContext differentiates from other interfaces.
	IsIf_statementContext()
}

type If_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIf_statementContext() *If_statementContext {
	var p = new(If_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_if_statement
	return p
}

func InitEmptyIf_statementContext(p *If_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_if_statement
}

func (*If_statementContext) IsIf_statementContext() {}

func NewIf_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *If_statementContext {
	var p = new(If_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_if_statement

	return p
}

func (s *If_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *If_statementContext) IF() antlr.TerminalNode {
	return s.GetToken(bpftraceParserIF, 0)
}

func (s *If_statementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *If_statementContext) AllBlock() []IBlockContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBlockContext); ok {
			len++
		}
	}

	tst := make([]IBlockContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBlockContext); ok {
			tst[i] = t.(IBlockContext)
			i++
		}
	}

	return tst
}

func (s *If_statementContext) Block(i int) IBlockContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *If_statementContext) ELSE() antlr.TerminalNode {
	return s.GetToken(bpftraceParserELSE, 0)
}

func (s *If_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *If_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *If_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterIf_statement(s)
	}
}

func (s *If_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitIf_statement(s)
	}
}

func (p *bpftraceParser) If_statement() (localctx IIf_statementContext) {
	localctx = NewIf_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, bpftraceParserRULE_if_statement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(233)
		p.Match(bpftraceParserIF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(234)
		p.Match(bpftraceParserT__10)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(235)
		p.Expression()
	}
	{
		p.SetState(236)
		p.Match(bpftraceParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(237)
		p.Block()
	}
	p.SetState(240)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == bpftraceParserELSE {
		{
			p.SetState(238)
			p.Match(bpftraceParserELSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(239)
			p.Block()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWhile_statementContext is an interface to support dynamic dispatch.
type IWhile_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHILE() antlr.TerminalNode
	Expression() IExpressionContext
	Block() IBlockContext

	// IsWhile_statementContext differentiates from other interfaces.
	IsWhile_statementContext()
}

type While_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhile_statementContext() *While_statementContext {
	var p = new(While_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_while_statement
	return p
}

func InitEmptyWhile_statementContext(p *While_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_while_statement
}

func (*While_statementContext) IsWhile_statementContext() {}

func NewWhile_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *While_statementContext {
	var p = new(While_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_while_statement

	return p
}

func (s *While_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *While_statementContext) WHILE() antlr.TerminalNode {
	return s.GetToken(bpftraceParserWHILE, 0)
}

func (s *While_statementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *While_statementContext) Block() IBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *While_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *While_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *While_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterWhile_statement(s)
	}
}

func (s *While_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitWhile_statement(s)
	}
}

func (p *bpftraceParser) While_statement() (localctx IWhile_statementContext) {
	localctx = NewWhile_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, bpftraceParserRULE_while_statement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(242)
		p.Match(bpftraceParserWHILE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(243)
		p.Match(bpftraceParserT__10)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(244)
		p.Expression()
	}
	{
		p.SetState(245)
		p.Match(bpftraceParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(246)
		p.Block()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFor_statementContext is an interface to support dynamic dispatch.
type IFor_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FOR() antlr.TerminalNode
	VARIABLE() antlr.TerminalNode
	IN() antlr.TerminalNode
	MAP_NAME() antlr.TerminalNode
	Block() IBlockContext
	AllAssignment() []IAssignmentContext
	Assignment(i int) IAssignmentContext
	Expression() IExpressionContext

	// IsFor_statementContext differentiates from other interfaces.
	IsFor_statementContext()
}

type For_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFor_statementContext() *For_statementContext {
	var p = new(For_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_for_statement
	return p
}

func InitEmptyFor_statementContext(p *For_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_for_statement
}

func (*For_statementContext) IsFor_statementContext() {}

func NewFor_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *For_statementContext {
	var p = new(For_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_for_statement

	return p
}

func (s *For_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *For_statementContext) FOR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserFOR, 0)
}

func (s *For_statementContext) VARIABLE() antlr.TerminalNode {
	return s.GetToken(bpftraceParserVARIABLE, 0)
}

func (s *For_statementContext) IN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserIN, 0)
}

func (s *For_statementContext) MAP_NAME() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMAP_NAME, 0)
}

func (s *For_statementContext) Block() IBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *For_statementContext) AllAssignment() []IAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAssignmentContext); ok {
			len++
		}
	}

	tst := make([]IAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAssignmentContext); ok {
			tst[i] = t.(IAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *For_statementContext) Assignment(i int) IAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignmentContext)
}

func (s *For_statementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *For_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *For_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *For_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterFor_statement(s)
	}
}

func (s *For_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitFor_statement(s)
	}
}

func (p *bpftraceParser) For_statement() (localctx IFor_statementContext) {
	localctx = NewFor_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, bpftraceParserRULE_for_statement)
	var _la int

	p.SetState(270)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(248)
			p.Match(bpftraceParserFOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(249)
			p.Match(bpftraceParserT__10)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(250)
			p.Match(bpftraceParserVARIABLE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(251)
			p.Match(bpftraceParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(252)
			p.Match(bpftraceParserMAP_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(253)
			p.Match(bpftraceParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(254)
			p.Block()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(255)
			p.Match(bpftraceParserFOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(256)
			p.Match(bpftraceParserT__10)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(258)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == bpftraceParserT__6 || _la == bpftraceParserVARIABLE || _la == bpftraceParserMAP_NAME {
			{
				p.SetState(257)
				p.Assignment()
			}

		}
		{
			p.SetState(260)
			p.Match(bpftraceParserT__3)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(262)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&9223372036853098624) != 0) || ((int64((_la-74)) & ^0x3f) == 0 && ((int64(1)<<(_la-74))&8522825759) != 0) {
			{
				p.SetState(261)
				p.Expression()
			}

		}
		{
			p.SetState(264)
			p.Match(bpftraceParserT__3)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(266)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == bpftraceParserT__6 || _la == bpftraceParserVARIABLE || _la == bpftraceParserMAP_NAME {
			{
				p.SetState(265)
				p.Assignment()
			}

		}
		{
			p.SetState(268)
			p.Match(bpftraceParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(269)
			p.Block()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IReturn_statementContext is an interface to support dynamic dispatch.
type IReturn_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RETURN() antlr.TerminalNode
	Expression() IExpressionContext

	// IsReturn_statementContext differentiates from other interfaces.
	IsReturn_statementContext()
}

type Return_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReturn_statementContext() *Return_statementContext {
	var p = new(Return_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_return_statement
	return p
}

func InitEmptyReturn_statementContext(p *Return_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_return_statement
}

func (*Return_statementContext) IsReturn_statementContext() {}

func NewReturn_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Return_statementContext {
	var p = new(Return_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_return_statement

	return p
}

func (s *Return_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *Return_statementContext) RETURN() antlr.TerminalNode {
	return s.GetToken(bpftraceParserRETURN, 0)
}

func (s *Return_statementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Return_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Return_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Return_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterReturn_statement(s)
	}
}

func (s *Return_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitReturn_statement(s)
	}
}

func (p *bpftraceParser) Return_statement() (localctx IReturn_statementContext) {
	localctx = NewReturn_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, bpftraceParserRULE_return_statement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(272)
		p.Match(bpftraceParserRETURN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(274)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(273)
			p.Expression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IClear_statementContext is an interface to support dynamic dispatch.
type IClear_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLEAR() antlr.TerminalNode
	MAP_NAME() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsClear_statementContext differentiates from other interfaces.
	IsClear_statementContext()
}

type Clear_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClear_statementContext() *Clear_statementContext {
	var p = new(Clear_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_clear_statement
	return p
}

func InitEmptyClear_statementContext(p *Clear_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_clear_statement
}

func (*Clear_statementContext) IsClear_statementContext() {}

func NewClear_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Clear_statementContext {
	var p = new(Clear_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_clear_statement

	return p
}

func (s *Clear_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *Clear_statementContext) CLEAR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserCLEAR, 0)
}

func (s *Clear_statementContext) MAP_NAME() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMAP_NAME, 0)
}

func (s *Clear_statementContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Clear_statementContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Clear_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Clear_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Clear_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterClear_statement(s)
	}
}

func (s *Clear_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitClear_statement(s)
	}
}

func (p *bpftraceParser) Clear_statement() (localctx IClear_statementContext) {
	localctx = NewClear_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, bpftraceParserRULE_clear_statement)
	var _la int

	p.SetState(293)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(276)
			p.Match(bpftraceParserCLEAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(277)
			p.Match(bpftraceParserMAP_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(278)
			p.Match(bpftraceParserCLEAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(279)
			p.Match(bpftraceParserT__6)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(291)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == bpftraceParserT__7 {
			{
				p.SetState(280)
				p.Match(bpftraceParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(281)
				p.Expression()
			}
			p.SetState(286)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == bpftraceParserT__8 {
				{
					p.SetState(282)
					p.Match(bpftraceParserT__8)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(283)
					p.Expression()
				}

				p.SetState(288)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(289)
				p.Match(bpftraceParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDelete_statementContext is an interface to support dynamic dispatch.
type IDelete_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DELETE() antlr.TerminalNode
	MAP_NAME() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsDelete_statementContext differentiates from other interfaces.
	IsDelete_statementContext()
}

type Delete_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDelete_statementContext() *Delete_statementContext {
	var p = new(Delete_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_delete_statement
	return p
}

func InitEmptyDelete_statementContext(p *Delete_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_delete_statement
}

func (*Delete_statementContext) IsDelete_statementContext() {}

func NewDelete_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Delete_statementContext {
	var p = new(Delete_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_delete_statement

	return p
}

func (s *Delete_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *Delete_statementContext) DELETE() antlr.TerminalNode {
	return s.GetToken(bpftraceParserDELETE, 0)
}

func (s *Delete_statementContext) MAP_NAME() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMAP_NAME, 0)
}

func (s *Delete_statementContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Delete_statementContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Delete_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Delete_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Delete_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterDelete_statement(s)
	}
}

func (s *Delete_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitDelete_statement(s)
	}
}

func (p *bpftraceParser) Delete_statement() (localctx IDelete_statementContext) {
	localctx = NewDelete_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, bpftraceParserRULE_delete_statement)
	var _la int

	p.SetState(321)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(295)
			p.Match(bpftraceParserDELETE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(296)
			p.Match(bpftraceParserMAP_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(297)
			p.Match(bpftraceParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(298)
			p.Expression()
		}
		p.SetState(303)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == bpftraceParserT__8 {
			{
				p.SetState(299)
				p.Match(bpftraceParserT__8)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(300)
				p.Expression()
			}

			p.SetState(305)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(306)
			p.Match(bpftraceParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(308)
			p.Match(bpftraceParserDELETE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(309)
			p.Match(bpftraceParserT__6)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(310)
			p.Match(bpftraceParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(311)
			p.Expression()
		}
		p.SetState(316)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == bpftraceParserT__8 {
			{
				p.SetState(312)
				p.Match(bpftraceParserT__8)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(313)
				p.Expression()
			}

			p.SetState(318)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(319)
			p.Match(bpftraceParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExit_statementContext is an interface to support dynamic dispatch.
type IExit_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXIT() antlr.TerminalNode
	Expression() IExpressionContext

	// IsExit_statementContext differentiates from other interfaces.
	IsExit_statementContext()
}

type Exit_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExit_statementContext() *Exit_statementContext {
	var p = new(Exit_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_exit_statement
	return p
}

func InitEmptyExit_statementContext(p *Exit_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_exit_statement
}

func (*Exit_statementContext) IsExit_statementContext() {}

func NewExit_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Exit_statementContext {
	var p = new(Exit_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_exit_statement

	return p
}

func (s *Exit_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *Exit_statementContext) EXIT() antlr.TerminalNode {
	return s.GetToken(bpftraceParserEXIT, 0)
}

func (s *Exit_statementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Exit_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Exit_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Exit_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterExit_statement(s)
	}
}

func (s *Exit_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitExit_statement(s)
	}
}

func (p *bpftraceParser) Exit_statement() (localctx IExit_statementContext) {
	localctx = NewExit_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, bpftraceParserRULE_exit_statement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(323)
		p.Match(bpftraceParserEXIT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(325)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(324)
			p.Expression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrint_statementContext is an interface to support dynamic dispatch.
type IPrint_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PRINT() antlr.TerminalNode
	Expression() IExpressionContext
	Output_redirection() IOutput_redirectionContext

	// IsPrint_statementContext differentiates from other interfaces.
	IsPrint_statementContext()
}

type Print_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrint_statementContext() *Print_statementContext {
	var p = new(Print_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_print_statement
	return p
}

func InitEmptyPrint_statementContext(p *Print_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_print_statement
}

func (*Print_statementContext) IsPrint_statementContext() {}

func NewPrint_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Print_statementContext {
	var p = new(Print_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_print_statement

	return p
}

func (s *Print_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *Print_statementContext) PRINT() antlr.TerminalNode {
	return s.GetToken(bpftraceParserPRINT, 0)
}

func (s *Print_statementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Print_statementContext) Output_redirection() IOutput_redirectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOutput_redirectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOutput_redirectionContext)
}

func (s *Print_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Print_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Print_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterPrint_statement(s)
	}
}

func (s *Print_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitPrint_statement(s)
	}
}

func (p *bpftraceParser) Print_statement() (localctx IPrint_statementContext) {
	localctx = NewPrint_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, bpftraceParserRULE_print_statement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(327)
		p.Match(bpftraceParserPRINT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(329)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(328)
			p.Expression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(332)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == bpftraceParserT__19 || _la == bpftraceParserGT || _la == bpftraceParserSHR {
		{
			p.SetState(331)
			p.Output_redirection()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrintf_statementContext is an interface to support dynamic dispatch.
type IPrintf_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PRINTF() antlr.TerminalNode
	STRING() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsPrintf_statementContext differentiates from other interfaces.
	IsPrintf_statementContext()
}

type Printf_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrintf_statementContext() *Printf_statementContext {
	var p = new(Printf_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_printf_statement
	return p
}

func InitEmptyPrintf_statementContext(p *Printf_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_printf_statement
}

func (*Printf_statementContext) IsPrintf_statementContext() {}

func NewPrintf_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Printf_statementContext {
	var p = new(Printf_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_printf_statement

	return p
}

func (s *Printf_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *Printf_statementContext) PRINTF() antlr.TerminalNode {
	return s.GetToken(bpftraceParserPRINTF, 0)
}

func (s *Printf_statementContext) STRING() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSTRING, 0)
}

func (s *Printf_statementContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Printf_statementContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Printf_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Printf_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Printf_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterPrintf_statement(s)
	}
}

func (s *Printf_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitPrintf_statement(s)
	}
}

func (p *bpftraceParser) Printf_statement() (localctx IPrintf_statementContext) {
	localctx = NewPrintf_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, bpftraceParserRULE_printf_statement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(334)
		p.Match(bpftraceParserPRINTF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(335)
		p.Match(bpftraceParserT__10)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(336)
		p.Match(bpftraceParserSTRING)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(341)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == bpftraceParserT__8 {
		{
			p.SetState(337)
			p.Match(bpftraceParserT__8)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(338)
			p.Expression()
		}

		p.SetState(343)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(344)
		p.Match(bpftraceParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Logical_or_expression() ILogical_or_expressionContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) Logical_or_expression() ILogical_or_expressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogical_or_expressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogical_or_expressionContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *bpftraceParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, bpftraceParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(346)
		p.Logical_or_expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILogical_or_expressionContext is an interface to support dynamic dispatch.
type ILogical_or_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLogical_and_expression() []ILogical_and_expressionContext
	Logical_and_expression(i int) ILogical_and_expressionContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsLogical_or_expressionContext differentiates from other interfaces.
	IsLogical_or_expressionContext()
}

type Logical_or_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogical_or_expressionContext() *Logical_or_expressionContext {
	var p = new(Logical_or_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_logical_or_expression
	return p
}

func InitEmptyLogical_or_expressionContext(p *Logical_or_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_logical_or_expression
}

func (*Logical_or_expressionContext) IsLogical_or_expressionContext() {}

func NewLogical_or_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Logical_or_expressionContext {
	var p = new(Logical_or_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_logical_or_expression

	return p
}

func (s *Logical_or_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Logical_or_expressionContext) AllLogical_and_expression() []ILogical_and_expressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILogical_and_expressionContext); ok {
			len++
		}
	}

	tst := make([]ILogical_and_expressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILogical_and_expressionContext); ok {
			tst[i] = t.(ILogical_and_expressionContext)
			i++
		}
	}

	return tst
}

func (s *Logical_or_expressionContext) Logical_and_expression(i int) ILogical_and_expressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogical_and_expressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogical_and_expressionContext)
}

func (s *Logical_or_expressionContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserOR)
}

func (s *Logical_or_expressionContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserOR, i)
}

func (s *Logical_or_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Logical_or_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Logical_or_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterLogical_or_expression(s)
	}
}

func (s *Logical_or_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitLogical_or_expression(s)
	}
}

func (p *bpftraceParser) Logical_or_expression() (localctx ILogical_or_expressionContext) {
	localctx = NewLogical_or_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, bpftraceParserRULE_logical_or_expression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(348)
		p.Logical_and_expression()
	}
	p.SetState(353)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == bpftraceParserOR {
		{
			p.SetState(349)
			p.Match(bpftraceParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(350)
			p.Logical_and_expression()
		}

		p.SetState(355)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILogical_and_expressionContext is an interface to support dynamic dispatch.
type ILogical_and_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllEquality_expression() []IEquality_expressionContext
	Equality_expression(i int) IEquality_expressionContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsLogical_and_expressionContext differentiates from other interfaces.
	IsLogical_and_expressionContext()
}

type Logical_and_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogical_and_expressionContext() *Logical_and_expressionContext {
	var p = new(Logical_and_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_logical_and_expression
	return p
}

func InitEmptyLogical_and_expressionContext(p *Logical_and_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_logical_and_expression
}

func (*Logical_and_expressionContext) IsLogical_and_expressionContext() {}

func NewLogical_and_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Logical_and_expressionContext {
	var p = new(Logical_and_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_logical_and_expression

	return p
}

func (s *Logical_and_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Logical_and_expressionContext) AllEquality_expression() []IEquality_expressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEquality_expressionContext); ok {
			len++
		}
	}

	tst := make([]IEquality_expressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEquality_expressionContext); ok {
			tst[i] = t.(IEquality_expressionContext)
			i++
		}
	}

	return tst
}

func (s *Logical_and_expressionContext) Equality_expression(i int) IEquality_expressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEquality_expressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEquality_expressionContext)
}

func (s *Logical_and_expressionContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserAND)
}

func (s *Logical_and_expressionContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserAND, i)
}

func (s *Logical_and_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Logical_and_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Logical_and_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterLogical_and_expression(s)
	}
}

func (s *Logical_and_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitLogical_and_expression(s)
	}
}

func (p *bpftraceParser) Logical_and_expression() (localctx ILogical_and_expressionContext) {
	localctx = NewLogical_and_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, bpftraceParserRULE_logical_and_expression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(356)
		p.Equality_expression()
	}
	p.SetState(361)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == bpftraceParserAND {
		{
			p.SetState(357)
			p.Match(bpftraceParserAND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(358)
			p.Equality_expression()
		}

		p.SetState(363)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEquality_expressionContext is an interface to support dynamic dispatch.
type IEquality_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRelational_expression() []IRelational_expressionContext
	Relational_expression(i int) IRelational_expressionContext
	AllEQ() []antlr.TerminalNode
	EQ(i int) antlr.TerminalNode
	AllNE() []antlr.TerminalNode
	NE(i int) antlr.TerminalNode

	// IsEquality_expressionContext differentiates from other interfaces.
	IsEquality_expressionContext()
}

type Equality_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEquality_expressionContext() *Equality_expressionContext {
	var p = new(Equality_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_equality_expression
	return p
}

func InitEmptyEquality_expressionContext(p *Equality_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_equality_expression
}

func (*Equality_expressionContext) IsEquality_expressionContext() {}

func NewEquality_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Equality_expressionContext {
	var p = new(Equality_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_equality_expression

	return p
}

func (s *Equality_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Equality_expressionContext) AllRelational_expression() []IRelational_expressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelational_expressionContext); ok {
			len++
		}
	}

	tst := make([]IRelational_expressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelational_expressionContext); ok {
			tst[i] = t.(IRelational_expressionContext)
			i++
		}
	}

	return tst
}

func (s *Equality_expressionContext) Relational_expression(i int) IRelational_expressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelational_expressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelational_expressionContext)
}

func (s *Equality_expressionContext) AllEQ() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserEQ)
}

func (s *Equality_expressionContext) EQ(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserEQ, i)
}

func (s *Equality_expressionContext) AllNE() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserNE)
}

func (s *Equality_expressionContext) NE(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserNE, i)
}

func (s *Equality_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Equality_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Equality_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterEquality_expression(s)
	}
}

func (s *Equality_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitEquality_expression(s)
	}
}

func (p *bpftraceParser) Equality_expression() (localctx IEquality_expressionContext) {
	localctx = NewEquality_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, bpftraceParserRULE_equality_expression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(364)
		p.Relational_expression()
	}
	p.SetState(369)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == bpftraceParserEQ || _la == bpftraceParserNE {
		{
			p.SetState(365)
			_la = p.GetTokenStream().LA(1)

			if !(_la == bpftraceParserEQ || _la == bpftraceParserNE) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(366)
			p.Relational_expression()
		}

		p.SetState(371)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRelational_expressionContext is an interface to support dynamic dispatch.
type IRelational_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllShift_expression() []IShift_expressionContext
	Shift_expression(i int) IShift_expressionContext
	AllLT() []antlr.TerminalNode
	LT(i int) antlr.TerminalNode
	AllGT() []antlr.TerminalNode
	GT(i int) antlr.TerminalNode
	AllLE() []antlr.TerminalNode
	LE(i int) antlr.TerminalNode
	AllGE() []antlr.TerminalNode
	GE(i int) antlr.TerminalNode

	// IsRelational_expressionContext differentiates from other interfaces.
	IsRelational_expressionContext()
}

type Relational_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelational_expressionContext() *Relational_expressionContext {
	var p = new(Relational_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_relational_expression
	return p
}

func InitEmptyRelational_expressionContext(p *Relational_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_relational_expression
}

func (*Relational_expressionContext) IsRelational_expressionContext() {}

func NewRelational_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Relational_expressionContext {
	var p = new(Relational_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_relational_expression

	return p
}

func (s *Relational_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Relational_expressionContext) AllShift_expression() []IShift_expressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IShift_expressionContext); ok {
			len++
		}
	}

	tst := make([]IShift_expressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IShift_expressionContext); ok {
			tst[i] = t.(IShift_expressionContext)
			i++
		}
	}

	return tst
}

func (s *Relational_expressionContext) Shift_expression(i int) IShift_expressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShift_expressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShift_expressionContext)
}

func (s *Relational_expressionContext) AllLT() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserLT)
}

func (s *Relational_expressionContext) LT(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserLT, i)
}

func (s *Relational_expressionContext) AllGT() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserGT)
}

func (s *Relational_expressionContext) GT(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserGT, i)
}

func (s *Relational_expressionContext) AllLE() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserLE)
}

func (s *Relational_expressionContext) LE(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserLE, i)
}

func (s *Relational_expressionContext) AllGE() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserGE)
}

func (s *Relational_expressionContext) GE(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserGE, i)
}

func (s *Relational_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Relational_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Relational_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterRelational_expression(s)
	}
}

func (s *Relational_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitRelational_expression(s)
	}
}

func (p *bpftraceParser) Relational_expression() (localctx IRelational_expressionContext) {
	localctx = NewRelational_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, bpftraceParserRULE_relational_expression)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(372)
		p.Shift_expression()
	}
	p.SetState(377)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(373)
				_la = p.GetTokenStream().LA(1)

				if !((int64((_la-82)) & ^0x3f) == 0 && ((int64(1)<<(_la-82))&15) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(374)
				p.Shift_expression()
			}

		}
		p.SetState(379)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IShift_expressionContext is an interface to support dynamic dispatch.
type IShift_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAdditive_expression() []IAdditive_expressionContext
	Additive_expression(i int) IAdditive_expressionContext
	AllSHL() []antlr.TerminalNode
	SHL(i int) antlr.TerminalNode
	AllSHR() []antlr.TerminalNode
	SHR(i int) antlr.TerminalNode

	// IsShift_expressionContext differentiates from other interfaces.
	IsShift_expressionContext()
}

type Shift_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShift_expressionContext() *Shift_expressionContext {
	var p = new(Shift_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_shift_expression
	return p
}

func InitEmptyShift_expressionContext(p *Shift_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_shift_expression
}

func (*Shift_expressionContext) IsShift_expressionContext() {}

func NewShift_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Shift_expressionContext {
	var p = new(Shift_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_shift_expression

	return p
}

func (s *Shift_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Shift_expressionContext) AllAdditive_expression() []IAdditive_expressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAdditive_expressionContext); ok {
			len++
		}
	}

	tst := make([]IAdditive_expressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAdditive_expressionContext); ok {
			tst[i] = t.(IAdditive_expressionContext)
			i++
		}
	}

	return tst
}

func (s *Shift_expressionContext) Additive_expression(i int) IAdditive_expressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAdditive_expressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAdditive_expressionContext)
}

func (s *Shift_expressionContext) AllSHL() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserSHL)
}

func (s *Shift_expressionContext) SHL(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHL, i)
}

func (s *Shift_expressionContext) AllSHR() []antlr.TerminalNode {
	return s.GetTokens(bpftraceParserSHR)
}

func (s *Shift_expressionContext) SHR(i int) antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHR, i)
}

func (s *Shift_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Shift_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Shift_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterShift_expression(s)
	}
}

func (s *Shift_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitShift_expression(s)
	}
}

func (p *bpftraceParser) Shift_expression() (localctx IShift_expressionContext) {
	localctx = NewShift_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, bpftraceParserRULE_shift_expression)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(380)
		p.Additive_expression()
	}
	p.SetState(385)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(381)
				_la = p.GetTokenStream().LA(1)

				if !(_la == bpftraceParserSHL || _la == bpftraceParserSHR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(382)
				p.Additive_expression()
			}

		}
		p.SetState(387)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAdditive_expressionContext is an interface to support dynamic dispatch.
type IAdditive_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllMultiplicative_expression() []IMultiplicative_expressionContext
	Multiplicative_expression(i int) IMultiplicative_expressionContext

	// IsAdditive_expressionContext differentiates from other interfaces.
	IsAdditive_expressionContext()
}

type Additive_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAdditive_expressionContext() *Additive_expressionContext {
	var p = new(Additive_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_additive_expression
	return p
}

func InitEmptyAdditive_expressionContext(p *Additive_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_additive_expression
}

func (*Additive_expressionContext) IsAdditive_expressionContext() {}

func NewAdditive_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Additive_expressionContext {
	var p = new(Additive_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_additive_expression

	return p
}

func (s *Additive_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Additive_expressionContext) AllMultiplicative_expression() []IMultiplicative_expressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMultiplicative_expressionContext); ok {
			len++
		}
	}

	tst := make([]IMultiplicative_expressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMultiplicative_expressionContext); ok {
			tst[i] = t.(IMultiplicative_expressionContext)
			i++
		}
	}

	return tst
}

func (s *Additive_expressionContext) Multiplicative_expression(i int) IMultiplicative_expressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMultiplicative_expressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMultiplicative_expressionContext)
}

func (s *Additive_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Additive_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Additive_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterAdditive_expression(s)
	}
}

func (s *Additive_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitAdditive_expression(s)
	}
}

func (p *bpftraceParser) Additive_expression() (localctx IAdditive_expressionContext) {
	localctx = NewAdditive_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, bpftraceParserRULE_additive_expression)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(388)
		p.Multiplicative_expression()
	}
	p.SetState(393)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(389)
				_la = p.GetTokenStream().LA(1)

				if !(_la == bpftraceParserT__12 || _la == bpftraceParserT__13) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(390)
				p.Multiplicative_expression()
			}

		}
		p.SetState(395)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMultiplicative_expressionContext is an interface to support dynamic dispatch.
type IMultiplicative_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllUnary_expression() []IUnary_expressionContext
	Unary_expression(i int) IUnary_expressionContext

	// IsMultiplicative_expressionContext differentiates from other interfaces.
	IsMultiplicative_expressionContext()
}

type Multiplicative_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMultiplicative_expressionContext() *Multiplicative_expressionContext {
	var p = new(Multiplicative_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_multiplicative_expression
	return p
}

func InitEmptyMultiplicative_expressionContext(p *Multiplicative_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_multiplicative_expression
}

func (*Multiplicative_expressionContext) IsMultiplicative_expressionContext() {}

func NewMultiplicative_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Multiplicative_expressionContext {
	var p = new(Multiplicative_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_multiplicative_expression

	return p
}

func (s *Multiplicative_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Multiplicative_expressionContext) AllUnary_expression() []IUnary_expressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUnary_expressionContext); ok {
			len++
		}
	}

	tst := make([]IUnary_expressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUnary_expressionContext); ok {
			tst[i] = t.(IUnary_expressionContext)
			i++
		}
	}

	return tst
}

func (s *Multiplicative_expressionContext) Unary_expression(i int) IUnary_expressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnary_expressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnary_expressionContext)
}

func (s *Multiplicative_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Multiplicative_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Multiplicative_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterMultiplicative_expression(s)
	}
}

func (s *Multiplicative_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitMultiplicative_expression(s)
	}
}

func (p *bpftraceParser) Multiplicative_expression() (localctx IMultiplicative_expressionContext) {
	localctx = NewMultiplicative_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, bpftraceParserRULE_multiplicative_expression)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(396)
		p.Unary_expression()
	}
	p.SetState(401)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(397)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&98308) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(398)
				p.Unary_expression()
			}

		}
		p.SetState(403)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnary_expressionContext is an interface to support dynamic dispatch.
type IUnary_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Unary_expression() IUnary_expressionContext
	Variable() IVariableContext
	INCR() antlr.TerminalNode
	DECR() antlr.TerminalNode
	Postfix_expression() IPostfix_expressionContext

	// IsUnary_expressionContext differentiates from other interfaces.
	IsUnary_expressionContext()
}

type Unary_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnary_expressionContext() *Unary_expressionContext {
	var p = new(Unary_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_unary_expression
	return p
}

func InitEmptyUnary_expressionContext(p *Unary_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_unary_expression
}

func (*Unary_expressionContext) IsUnary_expressionContext() {}

func NewUnary_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Unary_expressionContext {
	var p = new(Unary_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_unary_expression

	return p
}

func (s *Unary_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Unary_expressionContext) Unary_expression() IUnary_expressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnary_expressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnary_expressionContext)
}

func (s *Unary_expressionContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *Unary_expressionContext) INCR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserINCR, 0)
}

func (s *Unary_expressionContext) DECR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserDECR, 0)
}

func (s *Unary_expressionContext) Postfix_expression() IPostfix_expressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPostfix_expressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPostfix_expressionContext)
}

func (s *Unary_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Unary_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Unary_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterUnary_expression(s)
	}
}

func (s *Unary_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitUnary_expression(s)
	}
}

func (p *bpftraceParser) Unary_expression() (localctx IUnary_expressionContext) {
	localctx = NewUnary_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, bpftraceParserRULE_unary_expression)
	var _la int

	p.SetState(409)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case bpftraceParserT__12, bpftraceParserT__13, bpftraceParserT__16, bpftraceParserT__17:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(404)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&417792) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(405)
			p.Unary_expression()
		}

	case bpftraceParserINCR, bpftraceParserDECR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(406)
			_la = p.GetTokenStream().LA(1)

			if !(_la == bpftraceParserINCR || _la == bpftraceParserDECR) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(407)
			p.Variable()
		}

	case bpftraceParserT__6, bpftraceParserT__10, bpftraceParserT__20, bpftraceParserT__21, bpftraceParserT__22, bpftraceParserT__23, bpftraceParserT__24, bpftraceParserT__25, bpftraceParserT__26, bpftraceParserT__27, bpftraceParserT__28, bpftraceParserT__29, bpftraceParserT__30, bpftraceParserT__31, bpftraceParserT__32, bpftraceParserT__33, bpftraceParserT__34, bpftraceParserT__35, bpftraceParserT__36, bpftraceParserT__37, bpftraceParserT__38, bpftraceParserT__39, bpftraceParserT__40, bpftraceParserT__41, bpftraceParserT__42, bpftraceParserT__43, bpftraceParserT__44, bpftraceParserT__45, bpftraceParserT__46, bpftraceParserT__47, bpftraceParserT__48, bpftraceParserT__49, bpftraceParserT__50, bpftraceParserT__51, bpftraceParserT__52, bpftraceParserT__53, bpftraceParserT__54, bpftraceParserT__55, bpftraceParserT__56, bpftraceParserT__57, bpftraceParserT__58, bpftraceParserT__59, bpftraceParserT__60, bpftraceParserT__61, bpftraceParserCLEAR, bpftraceParserDELETE, bpftraceParserEXIT, bpftraceParserPRINT, bpftraceParserPRINTF, bpftraceParserNUMBER, bpftraceParserSTRING, bpftraceParserIDENTIFIER, bpftraceParserVARIABLE, bpftraceParserMAP_NAME:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(408)
			p.Postfix_expression()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPostfix_expressionContext is an interface to support dynamic dispatch.
type IPostfix_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primary_expression() IPrimary_expressionContext
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	IDENTIFIER() antlr.TerminalNode
	INCR() antlr.TerminalNode
	DECR() antlr.TerminalNode
	Expr_list() IExpr_listContext

	// IsPostfix_expressionContext differentiates from other interfaces.
	IsPostfix_expressionContext()
}

type Postfix_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPostfix_expressionContext() *Postfix_expressionContext {
	var p = new(Postfix_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_postfix_expression
	return p
}

func InitEmptyPostfix_expressionContext(p *Postfix_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_postfix_expression
}

func (*Postfix_expressionContext) IsPostfix_expressionContext() {}

func NewPostfix_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Postfix_expressionContext {
	var p = new(Postfix_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_postfix_expression

	return p
}

func (s *Postfix_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Postfix_expressionContext) Primary_expression() IPrimary_expressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimary_expressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimary_expressionContext)
}

func (s *Postfix_expressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Postfix_expressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Postfix_expressionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(bpftraceParserIDENTIFIER, 0)
}

func (s *Postfix_expressionContext) INCR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserINCR, 0)
}

func (s *Postfix_expressionContext) DECR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserDECR, 0)
}

func (s *Postfix_expressionContext) Expr_list() IExpr_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpr_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpr_listContext)
}

func (s *Postfix_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Postfix_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Postfix_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterPostfix_expression(s)
	}
}

func (s *Postfix_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitPostfix_expression(s)
	}
}

func (p *bpftraceParser) Postfix_expression() (localctx IPostfix_expressionContext) {
	localctx = NewPostfix_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, bpftraceParserRULE_postfix_expression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(411)
		p.Primary_expression()
	}
	p.SetState(431)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(412)
			_la = p.GetTokenStream().LA(1)

			if !(_la == bpftraceParserINCR || _la == bpftraceParserDECR) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) == 2 {
		{
			p.SetState(413)
			p.Match(bpftraceParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(414)
			p.Expression()
		}
		p.SetState(419)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == bpftraceParserT__8 {
			{
				p.SetState(415)
				p.Match(bpftraceParserT__8)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(416)
				p.Expression()
			}

			p.SetState(421)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(422)
			p.Match(bpftraceParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) == 3 {
		{
			p.SetState(424)
			p.Match(bpftraceParserT__10)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(426)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&9223372036853098624) != 0) || ((int64((_la-74)) & ^0x3f) == 0 && ((int64(1)<<(_la-74))&8522825759) != 0) {
			{
				p.SetState(425)
				p.Expr_list()
			}

		}
		{
			p.SetState(428)
			p.Match(bpftraceParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) == 4 {
		{
			p.SetState(429)
			p.Match(bpftraceParserT__18)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(430)
			p.Match(bpftraceParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimary_expressionContext is an interface to support dynamic dispatch.
type IPrimary_expressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode
	String_() IStringContext
	Variable() IVariableContext
	Map_access() IMap_accessContext
	Expression() IExpressionContext
	Function_call() IFunction_callContext

	// IsPrimary_expressionContext differentiates from other interfaces.
	IsPrimary_expressionContext()
}

type Primary_expressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimary_expressionContext() *Primary_expressionContext {
	var p = new(Primary_expressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_primary_expression
	return p
}

func InitEmptyPrimary_expressionContext(p *Primary_expressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_primary_expression
}

func (*Primary_expressionContext) IsPrimary_expressionContext() {}

func NewPrimary_expressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Primary_expressionContext {
	var p = new(Primary_expressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_primary_expression

	return p
}

func (s *Primary_expressionContext) GetParser() antlr.Parser { return s.parser }

func (s *Primary_expressionContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(bpftraceParserNUMBER, 0)
}

func (s *Primary_expressionContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *Primary_expressionContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *Primary_expressionContext) Map_access() IMap_accessContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMap_accessContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMap_accessContext)
}

func (s *Primary_expressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Primary_expressionContext) Function_call() IFunction_callContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunction_callContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunction_callContext)
}

func (s *Primary_expressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Primary_expressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Primary_expressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterPrimary_expression(s)
	}
}

func (s *Primary_expressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitPrimary_expression(s)
	}
}

func (p *bpftraceParser) Primary_expression() (localctx IPrimary_expressionContext) {
	localctx = NewPrimary_expressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, bpftraceParserRULE_primary_expression)
	p.SetState(442)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(433)
			p.Match(bpftraceParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(434)
			p.String_()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(435)
			p.Variable()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(436)
			p.Map_access()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(437)
			p.Match(bpftraceParserT__10)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(438)
			p.Expression()
		}
		{
			p.SetState(439)
			p.Match(bpftraceParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(441)
			p.Function_call()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVariableContext is an interface to support dynamic dispatch.
type IVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VARIABLE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsVariableContext differentiates from other interfaces.
	IsVariableContext()
}

type VariableContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableContext() *VariableContext {
	var p = new(VariableContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_variable
	return p
}

func InitEmptyVariableContext(p *VariableContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_variable
}

func (*VariableContext) IsVariableContext() {}

func NewVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableContext {
	var p = new(VariableContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_variable

	return p
}

func (s *VariableContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableContext) VARIABLE() antlr.TerminalNode {
	return s.GetToken(bpftraceParserVARIABLE, 0)
}

func (s *VariableContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(bpftraceParserIDENTIFIER, 0)
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitVariable(s)
	}
}

func (p *bpftraceParser) Variable() (localctx IVariableContext) {
	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, bpftraceParserRULE_variable)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(444)
		_la = p.GetTokenStream().LA(1)

		if !(_la == bpftraceParserIDENTIFIER || _la == bpftraceParserVARIABLE) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMap_accessContext is an interface to support dynamic dispatch.
type IMap_accessContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MAP_NAME() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsMap_accessContext differentiates from other interfaces.
	IsMap_accessContext()
}

type Map_accessContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMap_accessContext() *Map_accessContext {
	var p = new(Map_accessContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_map_access
	return p
}

func InitEmptyMap_accessContext(p *Map_accessContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_map_access
}

func (*Map_accessContext) IsMap_accessContext() {}

func NewMap_accessContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Map_accessContext {
	var p = new(Map_accessContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_map_access

	return p
}

func (s *Map_accessContext) GetParser() antlr.Parser { return s.parser }

func (s *Map_accessContext) MAP_NAME() antlr.TerminalNode {
	return s.GetToken(bpftraceParserMAP_NAME, 0)
}

func (s *Map_accessContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Map_accessContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Map_accessContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Map_accessContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Map_accessContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterMap_access(s)
	}
}

func (s *Map_accessContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitMap_access(s)
	}
}

func (p *bpftraceParser) Map_access() (localctx IMap_accessContext) {
	localctx = NewMap_accessContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, bpftraceParserRULE_map_access)
	var _la int

	p.SetState(472)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case bpftraceParserMAP_NAME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(446)
			p.Match(bpftraceParserMAP_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(447)
			p.Match(bpftraceParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(448)
			p.Expression()
		}
		p.SetState(453)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == bpftraceParserT__8 {
			{
				p.SetState(449)
				p.Match(bpftraceParserT__8)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(450)
				p.Expression()
			}

			p.SetState(455)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(456)
			p.Match(bpftraceParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case bpftraceParserT__6:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(458)
			p.Match(bpftraceParserT__6)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(470)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 52, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(459)
				p.Match(bpftraceParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(460)
				p.Expression()
			}
			p.SetState(465)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == bpftraceParserT__8 {
				{
					p.SetState(461)
					p.Match(bpftraceParserT__8)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(462)
					p.Expression()
				}

				p.SetState(467)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(468)
				p.Match(bpftraceParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpr_listContext is an interface to support dynamic dispatch.
type IExpr_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsExpr_listContext differentiates from other interfaces.
	IsExpr_listContext()
}

type Expr_listContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpr_listContext() *Expr_listContext {
	var p = new(Expr_listContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_expr_list
	return p
}

func InitEmptyExpr_listContext(p *Expr_listContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_expr_list
}

func (*Expr_listContext) IsExpr_listContext() {}

func NewExpr_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Expr_listContext {
	var p = new(Expr_listContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_expr_list

	return p
}

func (s *Expr_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Expr_listContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Expr_listContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Expr_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Expr_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Expr_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterExpr_list(s)
	}
}

func (s *Expr_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitExpr_list(s)
	}
}

func (p *bpftraceParser) Expr_list() (localctx IExpr_listContext) {
	localctx = NewExpr_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, bpftraceParserRULE_expr_list)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(474)
		p.Expression()
	}
	p.SetState(479)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == bpftraceParserT__8 {
		{
			p.SetState(475)
			p.Match(bpftraceParserT__8)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(476)
			p.Expression()
		}

		p.SetState(481)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOutput_redirectionContext is an interface to support dynamic dispatch.
type IOutput_redirectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GT() antlr.TerminalNode
	Expression() IExpressionContext
	SHR() antlr.TerminalNode

	// IsOutput_redirectionContext differentiates from other interfaces.
	IsOutput_redirectionContext()
}

type Output_redirectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOutput_redirectionContext() *Output_redirectionContext {
	var p = new(Output_redirectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_output_redirection
	return p
}

func InitEmptyOutput_redirectionContext(p *Output_redirectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_output_redirection
}

func (*Output_redirectionContext) IsOutput_redirectionContext() {}

func NewOutput_redirectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Output_redirectionContext {
	var p = new(Output_redirectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_output_redirection

	return p
}

func (s *Output_redirectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Output_redirectionContext) GT() antlr.TerminalNode {
	return s.GetToken(bpftraceParserGT, 0)
}

func (s *Output_redirectionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Output_redirectionContext) SHR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSHR, 0)
}

func (s *Output_redirectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Output_redirectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Output_redirectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterOutput_redirection(s)
	}
}

func (s *Output_redirectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitOutput_redirection(s)
	}
}

func (p *bpftraceParser) Output_redirection() (localctx IOutput_redirectionContext) {
	localctx = NewOutput_redirectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, bpftraceParserRULE_output_redirection)
	p.SetState(488)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case bpftraceParserGT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(482)
			p.Match(bpftraceParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(483)
			p.Expression()
		}

	case bpftraceParserSHR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(484)
			p.Match(bpftraceParserSHR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(485)
			p.Expression()
		}

	case bpftraceParserT__19:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(486)
			p.Match(bpftraceParserT__19)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(487)
			p.Expression()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunction_nameContext is an interface to support dynamic dispatch.
type IFunction_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode

	// IsFunction_nameContext differentiates from other interfaces.
	IsFunction_nameContext()
}

type Function_nameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunction_nameContext() *Function_nameContext {
	var p = new(Function_nameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_function_name
	return p
}

func InitEmptyFunction_nameContext(p *Function_nameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_function_name
}

func (*Function_nameContext) IsFunction_nameContext() {}

func NewFunction_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Function_nameContext {
	var p = new(Function_nameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_function_name

	return p
}

func (s *Function_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Function_nameContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(bpftraceParserIDENTIFIER, 0)
}

func (s *Function_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Function_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Function_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterFunction_name(s)
	}
}

func (s *Function_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitFunction_name(s)
	}
}

func (p *bpftraceParser) Function_name() (localctx IFunction_nameContext) {
	localctx = NewFunction_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, bpftraceParserRULE_function_name)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(490)
		p.Match(bpftraceParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBuiltin_nameContext is an interface to support dynamic dispatch.
type IBuiltin_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PRINTF() antlr.TerminalNode
	EXIT() antlr.TerminalNode
	DELETE() antlr.TerminalNode
	CLEAR() antlr.TerminalNode
	PRINT() antlr.TerminalNode

	// IsBuiltin_nameContext differentiates from other interfaces.
	IsBuiltin_nameContext()
}

type Builtin_nameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBuiltin_nameContext() *Builtin_nameContext {
	var p = new(Builtin_nameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_builtin_name
	return p
}

func InitEmptyBuiltin_nameContext(p *Builtin_nameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_builtin_name
}

func (*Builtin_nameContext) IsBuiltin_nameContext() {}

func NewBuiltin_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Builtin_nameContext {
	var p = new(Builtin_nameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_builtin_name

	return p
}

func (s *Builtin_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Builtin_nameContext) PRINTF() antlr.TerminalNode {
	return s.GetToken(bpftraceParserPRINTF, 0)
}

func (s *Builtin_nameContext) EXIT() antlr.TerminalNode {
	return s.GetToken(bpftraceParserEXIT, 0)
}

func (s *Builtin_nameContext) DELETE() antlr.TerminalNode {
	return s.GetToken(bpftraceParserDELETE, 0)
}

func (s *Builtin_nameContext) CLEAR() antlr.TerminalNode {
	return s.GetToken(bpftraceParserCLEAR, 0)
}

func (s *Builtin_nameContext) PRINT() antlr.TerminalNode {
	return s.GetToken(bpftraceParserPRINT, 0)
}

func (s *Builtin_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Builtin_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Builtin_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterBuiltin_name(s)
	}
}

func (s *Builtin_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitBuiltin_name(s)
	}
}

func (p *bpftraceParser) Builtin_name() (localctx IBuiltin_nameContext) {
	localctx = NewBuiltin_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, bpftraceParserRULE_builtin_name)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(492)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-21)) & ^0x3f) == 0 && ((int64(1)<<(_la-21))&279227574943481855) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommentContext is an interface to support dynamic dispatch.
type ICommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMENT() antlr.TerminalNode

	// IsCommentContext differentiates from other interfaces.
	IsCommentContext()
}

type CommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommentContext() *CommentContext {
	var p = new(CommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_comment
	return p
}

func InitEmptyCommentContext(p *CommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_comment
}

func (*CommentContext) IsCommentContext() {}

func NewCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommentContext {
	var p = new(CommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_comment

	return p
}

func (s *CommentContext) GetParser() antlr.Parser { return s.parser }

func (s *CommentContext) COMMENT() antlr.TerminalNode {
	return s.GetToken(bpftraceParserCOMMENT, 0)
}

func (s *CommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterComment(s)
	}
}

func (s *CommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitComment(s)
	}
}

func (p *bpftraceParser) Comment() (localctx ICommentContext) {
	localctx = NewCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, bpftraceParserRULE_comment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(494)
		p.Match(bpftraceParserCOMMENT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStringContext is an interface to support dynamic dispatch.
type IStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING() antlr.TerminalNode

	// IsStringContext differentiates from other interfaces.
	IsStringContext()
}

type StringContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringContext() *StringContext {
	var p = new(StringContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_string
	return p
}

func InitEmptyStringContext(p *StringContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = bpftraceParserRULE_string
}

func (*StringContext) IsStringContext() {}

func NewStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringContext {
	var p = new(StringContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = bpftraceParserRULE_string

	return p
}

func (s *StringContext) GetParser() antlr.Parser { return s.parser }

func (s *StringContext) STRING() antlr.TerminalNode {
	return s.GetToken(bpftraceParserSTRING, 0)
}

func (s *StringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.EnterString(s)
	}
}

func (s *StringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(bpftraceListener); ok {
		listenerT.ExitString(s)
	}
}

func (p *bpftraceParser) String_() (localctx IStringContext) {
	localctx = NewStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, bpftraceParserRULE_string)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(496)
		p.Match(bpftraceParserSTRING)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
