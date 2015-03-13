package marcopolo

import (
	"bytes"
	"testing"
)

func _TestVersion(t *testing.T) {
	a1_1 := Version{1, 1}
	a1_12 := Version{1, 12}
	a20_1 := Version{20, 1}
	a20_12 := Version{20, 12}
	var tests = []struct {
		v1     Version
		v2     Version
		cmpRes int
	}{
		{a1_1, a1_1, 0},
		{a1_1, a1_12, -1},
		{a1_1, a20_1, -1},
		{a1_1, a20_12, -1},

		{a1_12, a1_1, 1},
		{a1_12, a1_12, 0},
		{a1_12, a20_1, -1},
		{a1_12, a20_12, -1},

		{a20_1, a1_1, 1},
		{a20_1, a1_12, 1},
		{a20_1, a20_1, 0},
		{a20_1, a20_12, -1},

		{a20_12, a1_1, 1},
		{a20_12, a1_12, 1},
		{a20_12, a20_1, 1},
		{a20_12, a20_12, 0},
	}
	for i, test := range tests {
		cmp := test.v1.Compare(test.v2)

		if cmp != test.cmpRes {
			t.Errorf("#%d: %v.Compare(%v) != %d (expect %d)", i+1, test.v1, test.v2, cmp, test.cmpRes)
		}

		if cmp == 0 && !test.v1.Equals(test.v2) {
			t.Errorf("#%d: cmp == 0 && !%v.Equals(%v)", i+1, test.v1, test.v2)
		}
		if cmp != 0 && test.v1.Equals(test.v2) {
			t.Errorf("#%d: cmp != 0 && %v.Equals(%v)", i+1, test.v1, test.v2)
		}

		if cmp == -1 && !test.v1.SmallerThan(test.v2) {
			t.Errorf("#%d: cmp == -1 && !%v.SmallerThan(%v)", i+1, test.v1, test.v2)
		}
		if cmp == -1 && test.v1.GreatherThan(test.v2) {
			t.Errorf("#%d: cmp == -1 && %v.GreatherThan(%v)", i+1, test.v1, test.v2)
		}

		if cmp == 1 && !test.v1.GreatherThan(test.v2) {
			t.Errorf("#%d: cmp == 1 && !%v.GreatherThan(%v)", i+1, test.v1, test.v2)
		}
		if cmp == 1 && test.v1.SmallerThan(test.v2) {
			t.Errorf("#%d: cmp == 1 && %v.SmallerThan(%v)", i+1, test.v1, test.v2)
		}

	}
}

func TestValidateCmd(t *testing.T) {

	var tests = []struct {
		in       string
		out0     string
		out1     string
		hasError bool
	}{
		{"", "", "", true},
		{"toto", "", "", true},
		{"toto|", "", "", true},
		{"marco.polo", "", "", true},
		{"marco.polo.", "", "", true},
		{"marco.polo.regapp", "", "", true},
		{"marco.polo.regapp|", "marco.polo.regapp", "", false},
		{"marco.polo.qryapp|somestuff", "marco.polo.qryapp", "somestuff", false},
	}

	for i, test := range tests {
		// prep in param
		var udpPacket UdpPacket
		udpPacket.data = []byte(test.in)

		// call validateCmd()
		var cmdMsgStrs [2][]byte
		cmdMsgStrs, err := validateCmd(udpPacket.data)
		hasError := (err != nil)

		// check out err
		if test.hasError != hasError {
			t.Errorf("#%d: validateCmd(%q), expected %v", i+1, udpPacket.data, test.hasError)
		}

		// check out cmdMsgStrs[0]
		if !bytes.Equal([]byte(test.out0), cmdMsgStrs[0]) {
			t.Errorf("#%d: validateCmd(%q), cmdMsgStrs[0]->%q, expect %q ",
				i+1, udpPacket.data, cmdMsgStrs[0], test.out0)
		}

		// check out cmdMsgStrs[1]
		if !bytes.Equal([]byte(test.out1), cmdMsgStrs[1]) {
			t.Errorf("#%d: validateCmd(%q), cmdMsgStrs[1]->%q, expect %q ",
				i+1, udpPacket.data, cmdMsgStrs[1], test.out1)
		}

	}

}
