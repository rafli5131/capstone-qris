// Package util provides an EMV/QRIS TLV payload parser.
package util

import (
	"strconv"
	"strings"
)

// QRISData holds the fields extracted from a raw QRIS TLV string.
type QRISData struct {
	MerchantID   string
	MerchantName string
	City         string
	MCC          string
	FixedAmount  float64
	TerminalID   string
	RawPayload   string
}

// ParseQRIS parses a raw QRIS/EMV TLV string and returns extracted fields.
func ParseQRIS(payload string) *QRISData {
	data := &QRISData{RawPayload: payload}
	tlv := parseTLV(payload)
	if v, ok := tlv["52"]; ok { data.MCC = v }
	if v, ok := tlv["54"]; ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil { data.FixedAmount = f }
	}
	if v, ok := tlv["59"]; ok { data.MerchantName = v }
	if v, ok := tlv["60"]; ok { data.City = v }
	for _, parentTag := range []string{"26", "51"} {
		if rawTemplate, ok := tlv[parentTag]; ok {
			subTLV := parseTLV(rawTemplate)
			if mid, ok := subTLV["02"]; ok && mid != "" { data.MerchantID = mid }
			if tid, ok := subTLV["03"]; ok && tid != "" { data.TerminalID = tid }
			if data.MerchantID != "" { break }
		}
	}
	return data
}

// parseTLV reads a flat EMV TLV string and returns a map[tag]value.
func parseTLV(s string) map[string]string {
	result := make(map[string]string)
	i := 0
	for i+4 <= len(s) {
		tag := s[i : i+2]
		lenStr := s[i+2 : i+4]
		length, err := strconv.Atoi(lenStr)
		if err != nil || i+4+length > len(s) { break }
		value := s[i+4 : i+4+length]
		result[tag] = strings.TrimSpace(value)
		i += 4 + length
	}
	return result
}
