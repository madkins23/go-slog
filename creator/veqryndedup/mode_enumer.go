// Code generated by "enumer -type=Mode"; DO NOT EDIT.

package veqryndedup

import (
	"fmt"
	"strings"
)

const _ModeName = "NoneAppendIgnoreIncrementOverwrite"

var _ModeIndex = [...]uint8{0, 4, 10, 16, 25, 34}

const _ModeLowerName = "noneappendignoreincrementoverwrite"

func (i Mode) String() string {
	if i >= Mode(len(_ModeIndex)-1) {
		return fmt.Sprintf("Mode(%d)", i)
	}
	return _ModeName[_ModeIndex[i]:_ModeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ModeNoOp() {
	var x [1]struct{}
	_ = x[None-(0)]
	_ = x[Append-(1)]
	_ = x[Ignore-(2)]
	_ = x[Increment-(3)]
	_ = x[Overwrite-(4)]
}

var _ModeValues = []Mode{None, Append, Ignore, Increment, Overwrite}

var _ModeNameToValueMap = map[string]Mode{
	_ModeName[0:4]:        None,
	_ModeLowerName[0:4]:   None,
	_ModeName[4:10]:       Append,
	_ModeLowerName[4:10]:  Append,
	_ModeName[10:16]:      Ignore,
	_ModeLowerName[10:16]: Ignore,
	_ModeName[16:25]:      Increment,
	_ModeLowerName[16:25]: Increment,
	_ModeName[25:34]:      Overwrite,
	_ModeLowerName[25:34]: Overwrite,
}

var _ModeNames = []string{
	_ModeName[0:4],
	_ModeName[4:10],
	_ModeName[10:16],
	_ModeName[16:25],
	_ModeName[25:34],
}

// ModeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ModeString(s string) (Mode, error) {
	if val, ok := _ModeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ModeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Mode values", s)
}

// ModeValues returns all values of the enum
func ModeValues() []Mode {
	return _ModeValues
}

// ModeStrings returns a slice of all String values of the enum
func ModeStrings() []string {
	strs := make([]string, len(_ModeNames))
	copy(strs, _ModeNames)
	return strs
}

// IsAMode returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Mode) IsAMode() bool {
	for _, v := range _ModeValues {
		if i == v {
			return true
		}
	}
	return false
}
