package ntlmssp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

const expMsgBodyLen = 40

type negotiateMessageFields struct {
	messageHeader
	NegotiateFlags NegotiateFlags

	Domain      varField
	Workstation varField

	Version
}

//NewNegotiateMessage creates a new NEGOTIATE message with the
//flags that this package supports.
func NewNegotiateMessage(domainName, workstationName string) ([]byte, error) {
	payloadOffset := expMsgBodyLen
	flags := defaultFlags

	if domainName != "" {
		flags |= negotiateFlagNTLMSSPNEGOTIATEOEMDOMAINSUPPLIED
	}

	if workstationName != "" {
		flags |= negotiateFlagNTLMSSPNEGOTIATEOEMWORKSTATIONSUPPLIED
	}

	version := EmptyVersion()
	if flags.Has(negotiateFlagNTLMSSPNEGOTIATEVERSION) {
		version = DefaultVersion()
	}

	msg := negotiateMessageFields{
		messageHeader:  newMessageHeader(1),
		NegotiateFlags: flags,
		Domain:         newVarField(&payloadOffset, len(domainName)),
		Workstation:    newVarField(&payloadOffset, len(workstationName)),
		Version:        version,
	}

	b := bytes.Buffer{}
	if err := binary.Write(&b, binary.LittleEndian, &msg); err != nil {
		return nil, err
	}
	if b.Len() != expMsgBodyLen {
		return nil, errors.New("incorrect body length")
	}

	payload := strings.ToUpper(domainName + workstationName)
	if _, err := b.WriteString(payload); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
