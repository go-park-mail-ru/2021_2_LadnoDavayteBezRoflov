// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonAe927906DecodeBackendServerAppApiModels(in *jlexer.Lexer, out *CardList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "clid":
			out.CLID = uint(in.Uint())
		case "bid":
			out.BID = uint(in.Uint())
		case "cid":
			out.CID = uint(in.Uint())
		case "pos":
			out.PositionOnBoard = uint(in.Uint())
		case "cards":
			if in.IsNull() {
				in.Skip()
				out.Cards = nil
			} else {
				in.Delim('[')
				if out.Cards == nil {
					if !in.IsDelim(']') {
						out.Cards = make([]Card, 0, 0)
					} else {
						out.Cards = []Card{}
					}
				} else {
					out.Cards = (out.Cards)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Card
					(v1).UnmarshalEasyJSON(in)
					out.Cards = append(out.Cards, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "cardList_name":
			out.Title = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonAe927906EncodeBackendServerAppApiModels(out *jwriter.Writer, in CardList) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"clid\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.CLID))
	}
	{
		const prefix string = ",\"bid\":"
		out.RawString(prefix)
		out.Uint(uint(in.BID))
	}
	{
		const prefix string = ",\"cid\":"
		out.RawString(prefix)
		out.Uint(uint(in.CID))
	}
	{
		const prefix string = ",\"pos\":"
		out.RawString(prefix)
		out.Uint(uint(in.PositionOnBoard))
	}
	{
		const prefix string = ",\"cards\":"
		out.RawString(prefix)
		if in.Cards == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Cards {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"cardList_name\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CardList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonAe927906EncodeBackendServerAppApiModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CardList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonAe927906EncodeBackendServerAppApiModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CardList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonAe927906DecodeBackendServerAppApiModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CardList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonAe927906DecodeBackendServerAppApiModels(l, v)
}
