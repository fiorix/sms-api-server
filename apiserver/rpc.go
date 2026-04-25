// Copyright 2015 sms-api-server authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package apiserver

import (
	"fmt"
	"net/http"
	"net/rpc"
	"net/url"
	"strconv"

	"github.com/fiorix/go-smpp/v2/smpp"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdutext"
)

// parseUint8 parses a form value as uint8. Empty string returns 0.
func parseUint8(name, v string) (uint8, error) {
	if v == "" {
		return 0, nil
	}
	n, err := strconv.ParseUint(v, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid parameter %q=%q: %v", name, v, err)
	}
	return uint8(n), nil
}

// SM export its public methods to JSON RPC.
type SM struct {
	tx  *smpp.Transceiver
	rpc *rpc.Server
}

// NewSM creates and initializes a new SM, registering its own
// methods onto the given rpc server.
func NewSM(tx *smpp.Transceiver, rs *rpc.Server) *SM {
	sm := &SM{
		tx:  tx,
		rpc: rs,
	}
	sm.rpc.Register(sm) // hax more
	return sm
}

// ShortMessage contains the arguments of RPC call to SM.Submit.
type ShortMessage struct {
	Src      string `json:"src"`
	SrcTON   uint8  `json:"src_ton"`
	SrcNPI   uint8  `json:"src_npi"`
	Dst      string `json:"dst"`
	Text     string `json:"text"`
	Enc      string `json:"enc"`
	Register string `json:"register"`
}

// ShortMessageResp contains of RPC response from SM.Submit.
type ShortMessageResp struct {
	MessageID string `json:"message_id"`
}

// Submit sends a short message via RPC.
func (rpc *SM) Submit(args *ShortMessage, resp *ShortMessageResp) error {
	req := url.Values{
		"src":      {args.Src},
		"src_ton":  {strconv.FormatUint(uint64(args.SrcTON), 10)},
		"src_npi":  {strconv.FormatUint(uint64(args.SrcNPI), 10)},
		"dst":      {args.Dst},
		"text":     {args.Text},
		"enc":      {args.Enc},
		"register": {args.Register},
	}
	r, s, err := rpc.submit(req)
	if err != nil {
		return fmt.Errorf("%d %s: %v", s, http.StatusText(s), err)
	}
	*resp = *r
	return nil
}

func (rpc *SM) submit(req url.Values) (resp *ShortMessageResp, status int, err error) {
	sm := &smpp.ShortMessage{}
	var msg, enc, register, srcTON, srcNPI string
	f := form{
		{"src", "number of sender", false, nil, &sm.Src},
		{"src_ton", "type of number of sender", false, nil, &srcTON},
		{"src_npi", "numbering plan indicator of sender", false, nil, &srcNPI},
		{"dst", "number of recipient", true, nil, &sm.Dst},
		{"text", "text message", true, nil, &msg},
		{"enc", "text encoding", false, []string{"latin1", "ucs2"}, &enc},
		{"register", "registered delivery", false, []string{"final", "failure"}, &register},
	}
	if err := f.Validate(req); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if sm.SourceAddrTON, err = parseUint8("src_ton", srcTON); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if sm.SourceAddrNPI, err = parseUint8("src_npi", srcNPI); err != nil {
		return nil, http.StatusBadRequest, err
	}
	switch enc {
	case "":
		sm.Text = pdutext.Raw(msg)
	case "latin1", "latin-1":
		sm.Text = pdutext.Latin1(msg)
	case "ucs2", "ucs-2":
		sm.Text = pdutext.UCS2(msg)
	}
	switch register {
	case "final":
		sm.Register = pdufield.FinalDeliveryReceipt
	case "failure":
		sm.Register = pdufield.FailureDeliveryReceipt
	}
	sm, err = rpc.tx.Submit(sm)
	if err == smpp.ErrNotConnected {
		return nil, http.StatusServiceUnavailable, err
	}
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	resp = &ShortMessageResp{MessageID: sm.RespID()}
	return resp, http.StatusOK, nil
}

// QueryMessage contains the arguments of RPC call to SM.Query.
type QueryMessage struct {
	Src       string `json:"src"`
	SrcTON    uint8  `json:"src_ton"`
	SrcNPI    uint8  `json:"src_npi"`
	MessageID string `json:"message_id"`
}

// QueryMessageResp contains RPC response from SM.Query.
type QueryMessageResp struct {
	MsgState  string `json:"message_state"`
	FinalDate string `json:"final_date"`
	ErrCode   uint8  `json:"error_code"`
}

// Query queries the status of a short message via RPC.
func (rpc *SM) Query(args *QueryMessage, resp *QueryMessageResp) error {
	req := url.Values{
		"src":        {args.Src},
		"src_ton":    {strconv.FormatUint(uint64(args.SrcTON), 10)},
		"src_npi":    {strconv.FormatUint(uint64(args.SrcNPI), 10)},
		"message_id": {args.MessageID},
	}
	r, s, err := rpc.query(req)
	if err != nil {
		return fmt.Errorf("%d %s: %v", s, http.StatusText(s), err)
	}
	*resp = *r
	return nil
}

func (rpc *SM) query(req url.Values) (resp *QueryMessageResp, status int, err error) {
	var srcTON, srcNPI string
	f := form{
		{"src", "number of sender", false, nil, nil},
		{"src_ton", "type of number of sender", false, nil, &srcTON},
		{"src_npi", "numbering plan indicator of sender", false, nil, &srcNPI},
		{"message_id", "message id from send", true, nil, nil},
	}
	if err := f.Validate(req); err != nil {
		return nil, http.StatusBadRequest, err
	}
	ton, err := parseUint8("src_ton", srcTON)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	npi, err := parseUint8("src_npi", srcNPI)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	qr, err := rpc.tx.QuerySM(req.Get("src"), req.Get("message_id"), ton, npi)
	if err == smpp.ErrNotConnected {
		return nil, http.StatusServiceUnavailable, err
	}
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	resp = &QueryMessageResp{
		MsgState:  qr.MsgState,
		FinalDate: qr.FinalDate,
		ErrCode:   qr.ErrCode,
	}
	return resp, http.StatusOK, nil
}
