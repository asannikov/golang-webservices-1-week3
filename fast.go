package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

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

type User struct {
	Phone    string   `json:"phone"`
	Name     string   `json:"name"`
	Job      string   `json:"job"`
	Email    string   `json:"email"`
	Country  string   `json:"country"`
	Company  string   `json:"company"`
	Browsers []string `json:"browsers"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

	var userPool = sync.Pool{
		New: func() interface{} {
			return User{}
		},
	}

	for i, line := range lines {
		// fmt.Printf("%v %v\n", err, line)
		user := userPool.Get().(User)
		err := (&user).UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false
		browsers := user.Browsers

		for _, browserRaw := range browsers {
			if strings.TrimSpace(browserRaw) == "" {
				//log.Println("cant cast browser to string")
				continue
			}
			if strings.Contains(browserRaw, "Android") == true {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					//log.Printf("SLOW New browser: %s, first seen: %s", browserRaw, user.Name)
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range browsers {
			if strings.TrimSpace(browserRaw) == "" {
				//log.Println("cant cast browser to string")
				continue
			}
			if strings.Contains(browserRaw, "MSIE") == true {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					//log.Printf("SLOW New browser: %s, first seen: %s", browserRaw, user.Name)
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := r.ReplaceAllString(user.Email, " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

// =========== =========== =========== =========== =========== =========== =========== =========== ===========
// =========== =========== =========== =========== =========== =========== =========== =========== ===========

func easyjson9e1087fdDecodeUsersGoSrc(in *jlexer.Lexer, out *User) {
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
		case "phone":
			out.Phone = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "job":
			out.Job = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "country":
			out.Country = string(in.String())
		case "company":
			out.Company = string(in.String())
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
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
func easyjson9e1087fdEncodeUsersGoSrc(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"phone\":"
		out.RawString(prefix[1:])
		out.String(string(in.Phone))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"job\":"
		out.RawString(prefix)
		out.String(string(in.Job))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"country\":"
		out.RawString(prefix)
		out.String(string(in.Country))
	}
	{
		const prefix string = ",\"company\":"
		out.RawString(prefix)
		out.String(string(in.Company))
	}
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix)
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeUsersGoSrc(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeUsersGoSrc(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeUsersGoSrc(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeUsersGoSrc(l, v)
}
