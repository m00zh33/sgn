package sgn

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/EgeBalci/keystone-go"
)

// REG structure for registers
type REG struct {
	Full     string
	Extended string
	High     string
	Low      string
	Arch     int
}

// Initialize the register values
func init() {

	// Setup x86 the register values
	REGS = make(map[int][]REG)
	REGS[32] = append(REGS[32], REG{Extended: "EAX", High: "AX", Low: "AL", Arch: 32})
	REGS[32] = append(REGS[32], REG{Extended: "EBX", High: "BX", Low: "BL", Arch: 32})
	REGS[32] = append(REGS[32], REG{Extended: "ECX", High: "CX", Low: "CL", Arch: 32})
	REGS[32] = append(REGS[32], REG{Extended: "EDX", High: "DX", Low: "DL", Arch: 32})
	// since there is no way to access 1 byte use above instead
	REGS[32] = append(REGS[32], REG{Extended: "ESI", High: "SI", Low: "AL", Arch: 32})
	REGS[32] = append(REGS[32], REG{Extended: "EDI", High: "DI", Low: "BL", Arch: 32})
	// Setup x64 the register values
	REGS[64] = append(REGS[64], REG{Full: "RAX", Extended: "EAX", High: "AX", Low: "AL", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "RBX", Extended: "EBX", High: "BX", Low: "BL", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "RCX", Extended: "ECX", High: "CX", Low: "CL", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "RDX", Extended: "EDX", High: "DX", Low: "DL", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "RSI", Extended: "ESI", High: "SI", Low: "SIL", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "RDI", Extended: "EDI", High: "DX", Low: "DIL", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R8", Extended: "R8D", High: "R8W", Low: "R8B", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R9", Extended: "R9D", High: "R9W", Low: "R9B", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R10", Extended: "R10D", High: "R10W", Low: "R10B", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R11", Extended: "R11D", High: "R11W", Low: "R11B", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R12", Extended: "R12D", High: "R12W", Low: "R12B", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R13", Extended: "R13D", High: "R13W", Low: "R13B", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R14", Extended: "R14D", High: "R14W", Low: "R14B", Arch: 64})
	REGS[64] = append(REGS[64], REG{Full: "R15", Extended: "R15D", High: "R15W", Low: "R15B", Arch: 64})

	// Set the decoder stubs
	STUB = make(map[int]string)
	STUB[32] = x86DecoderStub
	STUB[64] = x64DecoderStub

	// Set safe register prefix/suffix
	SafeRegisterPrefix = make(map[int]([]byte))
	SafeRegisterSuffix = make(map[int]([]byte))
	SafeRegisterPrefix[32] = safeX86Prefix
	SafeRegisterPrefix[64] = safeX64Prefix

	SafeRegisterSuffix[32] = safeX86Suffix
	SafeRegisterSuffix[64] = safeX64Suffix

	// Increase random garbage instruction generation pool
	addGarbageJumpMnemonics()

	// Set random seed
	rand.Seed(time.Now().UTC().UnixNano())

}

// SafeRegisterPrefix contains the instructions for saving registers to stack
var SafeRegisterPrefix map[int]([]byte)

// SafeRegisterSuffix contains the instructions for restoring registers from stack
var SafeRegisterSuffix map[int]([]byte)

// safeX86Prefix instructions for saving registers to stack
var safeX86Prefix = []byte{0x60, 0x9c} // PUSHAD, PUSHFD
// safeX86Suffix instructions for saving registers to stack
var safeX86Suffix = []byte{0x9d, 0x61} // POPFD, POPAD

// safeX64Prefix instructions for saving registers to stack
var safeX64Prefix = []byte{
	0x50, 0x53, 0x51, 0x52, // PUSH RAX,RBX,RCX,RDX
	0x56, 0x57, 0x55, 0x54, // PUSH RSI,RDI,RBP,RSP
	0x41, 0x50, 0x41, 0x51, // PUSH R8,R9
	0x41, 0x52, 0x41, 0x53, // PUSH R10,R11
	0x41, 0x54, 0x41, 0x55, // PUSH R12,R13
	0x41, 0x56, 0x41, 0x57, // PUSH R14,R15
}

// safeX64Suffix instructions for saving registers to stack
var safeX64Suffix = []byte{
	0x41, 0x5f, 0x41, 0x5e, // POP R15,R14
	0x41, 0x5d, 0x41, 0x5c, // POP R13,R12
	0x41, 0x5b, 0x41, 0x5a, // POP R11,R10
	0x41, 0x59, 0x41, 0x58, // POP R9,R8
	0x5c, 0x5d, 0x5f, 0x5e, // POP RSP,RBP,RDI,RSI
	0x5a, 0x59, 0x5b, 0x58, // POP RDX,RCX,RBX,RAX
}

// REGS contains 32/64 bit registers
var REGS map[int][]REG

// RandomRegister returns a random register name based on given size and architecture
func (encoder Encoder) RandomRegister(size int) string {

	switch size {
	case 1:
		return REGS[encoder.architecture][rand.Intn(len(REGS[encoder.architecture]))].Low
	case 2:
		return REGS[encoder.architecture][rand.Intn(len(REGS[encoder.architecture]))].High
	case 4:
		return REGS[encoder.architecture][rand.Intn(len(REGS[encoder.architecture]))].Extended
	case 8:
		return REGS[encoder.architecture][rand.Intn(len(REGS[encoder.architecture]))].Full
	default:
		panic("invalid register size")
	}

}

// SafeRandomRegister returns a random register amoung all (registers-excluded parameters) based on given size
func (encoder Encoder) SafeRandomRegister(size int, excludes ...string) string {

	for {
		r := REGS[encoder.architecture][rand.Intn(len(REGS[encoder.architecture]))]
		for i, exclude := range excludes {
			if r.Full != exclude && r.Extended != exclude && r.High != exclude && r.Low != exclude {
				if i == len(excludes)-1 {
					switch size {
					case 1:
						return r.Low
					case 2:
						return r.High
					case 4:
						return r.Extended
					case 8:
						return r.Full
					default:
						panic("invalid register size")
					}
				}
			} else {
				break
			}
		}
	}

}

// Assemble assembes the given instructions
// and return a byte array with a boolean value indicating wether the operation is successfull or not
func (encoder Encoder) Assemble(asm string) ([]byte, bool) {
	var mode keystone.Mode
	switch encoder.architecture {
	case 32:
		mode = keystone.MODE_32
	case 64:
		mode = keystone.MODE_64
	default:
		return nil, false
	}

	ks, err := keystone.New(keystone.ARCH_X86, mode)
	if err != nil {
		return nil, false
	}
	defer ks.Close()

	//err = ks.Option(keystone.OPT_SYNTAX, keystone.OPT_SYNTAX_INTEL)
	//err = ks.Option(keystone.OPT_SYNTAX, keystone.KS_OPT_SYNTAX_NASM)
	err = ks.Option(keystone.OPT_SYNTAX, keystone.OPT_SYNTAX_INTEL)
	if err != nil {
		return nil, false
	}
	//log.Println(asm)
	bin, _, ok := ks.Assemble(asm, 0)
	return bin, ok
}

// GenerateIPToStack function generates instructions series that pushes the instruction pointer to stack
func (encoder Encoder) GenerateIPToStack() []byte {

	callBin, ok := encoder.Assemble("call 5")
	if !ok {
		panic("call 5 assembly failed")
	}
	return callBin
}

// AddCallOver function adds a call instruction over the end of the given payload
// address of the payload will be pushed to the stack and execution will continiou after the end of payload
func (encoder Encoder) AddCallOver(payload []byte) ([]byte, error) {

	// Perform a shport call over the payload
	call := fmt.Sprintf("call 0x%x", len(payload)+5)
	callBin, ok := encoder.Assemble(call)
	if !ok {
		return nil, errors.New("call-over assembly failed")
	}
	payload = append(callBin, payload...)

	return payload, nil
}

// AddJmpOver function adds a jmp instruction over the end of the given payload
// execution will continiou after the end of payload
func (encoder Encoder) AddJmpOver(payload []byte) ([]byte, error) {
	// JMP 2 -> Jumps to next instruction
	// Perform a short call over the payload
	jmp := fmt.Sprintf("jmp 0x%x", len(payload)+2)
	jmpBin, ok := encoder.Assemble(jmp)
	if !ok {
		return nil, errors.New("jmp-over assembly failed")
	}
	payload = append(jmpBin, payload...)

	return payload, nil
}

// AddCondJmpOver function adds a jmp instruction over the end of the given payload
// execution will continiou after the end of payload
func (encoder Encoder) AddCondJmpOver(payload []byte) ([]byte, error) {
	// JZ 2 -> Jumps to next instruction
	// Perform a short call over the payload

	randomConditional := ConditionalJumpMnemonics[rand.Intn(len(ConditionalJumpMnemonics))]

	jmp := fmt.Sprintf("%s 0x%x", randomConditional, len(payload)+2)
	jmpBin, ok := encoder.Assemble(jmp)
	if !ok {
		return nil, errors.New("conditional call-over assembly failed")
	}
	payload = append(jmpBin, payload...)

	return payload, nil
}
