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

func easyjson202377feDecodeBackendServerAppApiModels(in *jlexer.Lexer, out *Board) {
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
		case "bid":
			out.BID = uint(in.Uint())
		case "tid":
			out.TID = uint(in.Uint())
		case "board_name":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "access_path":
			out.AccessPath = string(in.String())
		case "members":
			if in.IsNull() {
				in.Skip()
				out.Members = nil
			} else {
				in.Delim('[')
				if out.Members == nil {
					if !in.IsDelim(']') {
						out.Members = make([]PublicUserInfo, 0, 1)
					} else {
						out.Members = []PublicUserInfo{}
					}
				} else {
					out.Members = (out.Members)[:0]
				}
				for !in.IsDelim(']') {
					var v1 PublicUserInfo
					(v1).UnmarshalEasyJSON(in)
					out.Members = append(out.Members, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "card_lists":
			if in.IsNull() {
				in.Skip()
				out.CardLists = nil
			} else {
				in.Delim('[')
				if out.CardLists == nil {
					if !in.IsDelim(']') {
						out.CardLists = make([]CardList, 0, 0)
					} else {
						out.CardLists = []CardList{}
					}
				} else {
					out.CardLists = (out.CardLists)[:0]
				}
				for !in.IsDelim(']') {
					var v2 CardList
					(v2).UnmarshalEasyJSON(in)
					out.CardLists = append(out.CardLists, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson202377feEncodeBackendServerAppApiModels(out *jwriter.Writer, in Board) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"bid\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.BID))
	}
	{
		const prefix string = ",\"tid\":"
		out.RawString(prefix)
		out.Uint(uint(in.TID))
	}
	{
		const prefix string = ",\"board_name\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"access_path\":"
		out.RawString(prefix)
		out.String(string(in.AccessPath))
	}
	{
		const prefix string = ",\"members\":"
		out.RawString(prefix)
		if in.Members == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Members {
				if v3 > 0 {
					out.RawByte(',')
				}
				(v4).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"card_lists\":"
		out.RawString(prefix)
		if in.CardLists == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.CardLists {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Board) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeBackendServerAppApiModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Board) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeBackendServerAppApiModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Board) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeBackendServerAppApiModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Board) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeBackendServerAppApiModels(l, v)
}
