package wiki

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/yundream/gowiki/plugin"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

const (
	HEAD = 1 << iota
	INSTRUCTION
	PROCESSOR_OPEN
	PROCESSOR_CLOSE
	TABLE_OPEN
	TABLE_CLOSE
)

type listInfo struct {
	Depth    int
	Type     string
	CloseTag string
}
type RequestOptions struct {
	Id          int
	PageName    string
	UserName    string
	SessionName string
	Theme       string
	R           *http.Request
	W           http.ResponseWriter
}

type Compiler struct {
	text            string
	listnum         int
	tablenum        int
	processorBuffer bytes.Buffer
	Linenum         int
	ProcessName     string
	listStack       *list.List
	isCompile       bool
	TextType        int
	opt             RequestOptions
	P               *plugin.PlugIns
	w               http.ResponseWriter
	r               *http.Request
}

func (c Compiler) NewIns(w http.ResponseWriter, r *http.Request) *Compiler {
	/*
		c.TextType = HEAD
		c.listStack = list.New()
	*/
	ins := &Compiler{
		P:         c.P,
		TextType:  HEAD,
		listStack: list.New(),
		w:         w,
		r:         r,
	}
	return ins
}

func (c *Compiler) Start(id string) *Compiler {
	c.Linenum = 0
	c.text = ""
	return c
}

func (c *Compiler) LoadPlugin() error {
	var err error
	c.P, err = plugin.Load()
	if err != nil {
		return err
	}
	return nil
}

func (c *Compiler) Instruction() *Compiler {
	if c.Linenum == 1 {
		idx := strings.Index(c.text, "#!")
		if idx == 0 {
		}
		c.text = ""
	}
	return c
}

func (c *Compiler) Head() *Compiler {
	re := regexp.MustCompile("^[ \t]*([=]+) (.*) ([=]+)")
	data := re.FindStringSubmatch(c.text)
	if len(data) > 0 {
		if len(data[1]) == len(data[3]) {
			depth := len(data[1])
			c.text = fmt.Sprintf("<div class=\"row\"><div class=\"small-12 columns\"><h%d class=\"head\"> %s </h%d></div></div>", depth, data[2], depth)
			c.isCompile = true
		}
		c.TextType = HEAD
	}
	return c
}

func (c *Compiler) Hr() *Compiler {
	if c.text == "----" {
		c.text = `<div class="row" style="margin-top:20px">
		<span style="margin:auto; display:table; ">
		<i class="fi-asterisk" style="margin-left:4em;font-size:8px"></i>
		<i class="fi-asterisk" style="margin-left:4em;font-size:8px"></i>
		<i class="fi-asterisk" style="margin-left:4em;font-size:8px"></i>
		</span>
		</div>`
	}
	return c

}

func (c *Compiler) Table() *Compiler {
	re, _ := regexp.MatchString("^[ \t]*\\|\\|.+\\|\\|$", c.text)
	if re == true {
		token := strings.Split(c.text, "||")

		c.text = "<tr>"
		if c.tablenum == 0 {
			c.text = "<table><thead>" + c.text
		}

		for i := 1; i < len(token)-1; i++ {
			align := "left"
			lefti := strings.Index(token[i], " ")
			righti := strings.LastIndex(token[i], " ") - (len(token[i]) - 1)
			if lefti+righti == 0 {
				align = "center"
			} else {
				if righti != 0 {
					align = "left"
				}
			}
			c.text += "<td align='" + align + "'>" + token[i] + "</td>"
		}
		c.text += "</tr>"
		if c.tablenum == 0 {
			c.text += "</tbody></thead><tbody>"
		}
		c.tablenum++
	} else {
		if c.tablenum > 0 {
			c.text = "</tbody></table>\n" + c.text
			c.tablenum = 0
		}

	}
	return c
}

func (c *Compiler) List() *Compiler {
	re := regexp.MustCompile("^[ \t]*(?P<first>\\*|1\\.)[ ]+(?P<second>.*)")
	match := re.FindStringSubmatch(c.text)
	var openTag string
	var closeTag string
	if len(match) == 3 {
		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			result[name] = match[i]
		}
		if result["first"] == "*" {
			openTag = "<ul class=\"item_list\">"
			closeTag = "</ul>"
		} else if result["first"] == "1." {
			openTag = "<ol class=\"item_list\">"
			closeTag = "</ol>"
		}

		item := re.ReplaceAllString(c.text, result["second"])
		depth := strings.Index(c.text, result["first"])
		if c.listnum == 0 {
			c.listStack.PushBack(
				listInfo{depth,
					result["first"],
					closeTag})
			c.text = openTag + "\n<li>" + item
			c.listnum++
		} else {
			current := 0
			if c.listStack.Len() > 0 {
				current = c.listStack.Back().Value.(listInfo).Depth
			}
			if depth > current {
				c.text = "\n" + openTag + "\n"
				c.listStack.PushBack(
					listInfo{depth,
						result["first"],
						closeTag})
			} else if depth < current {
				count := 0
				for i := c.listStack.Back(); i != nil; i = i.Prev() {
					list := i.Value.(listInfo)
					if list.Depth > depth {
						count++
						c.text = list.CloseTag + "</li>\n"
					} else {
						break
					}
				}
				for i := 0; i < count; i++ {
					c.listStack.Remove(c.listStack.Back())
				}
			} else {
				c.text = "</li>\n"
			}
			c.text += "<li>" + item
		}
	} else {
		if c.listnum > 0 {
			stackSize := c.listStack.Len()
			if c.listStack.Len() > 0 {
				closeListString := ""
				for i := 0; i < stackSize; i++ {
					tag := c.listStack.Back().Value.(listInfo).CloseTag
					closeListString += "</li>" + tag
				}
				c.text = closeListString + c.text
				c.listStack.Init()
				c.listnum = 0
			}
		}
	}
	return c
}

func (c *Compiler) EscapeString() *Compiler {
	if c.TextType == HEAD {
		return c
	}
	re := regexp.MustCompile("<")
	c.text = re.ReplaceAllString(c.text, "&lt;")
	re = regexp.MustCompile(">")
	c.text = re.ReplaceAllString(c.text, "&gt;")
	return c
}

func (c *Compiler) Body() *Compiler {
	if c.TextType == HEAD {
		return c
	}
	if c.text == "" {
		c.text = "<p></p>\n"
	} else {
		c.Comma().Decorator().OuterLink().InnerLink().ExecPlugin()
	}
	return c
}

func (c *Compiler) InnerLink() *Compiler {
	re := regexp.MustCompile("\\[(wiki:|/|\\.\\./)([^ ]+)[ ]+([^\\]]+)\\]")
	c.text = re.ReplaceAllStringFunc(c.text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		switch parts[1] {
		case "wiki:":
			return "<a href=\"/w/" + parts[2] + "\">" + parts[3] + "</a>"
		case "/":
			return "<a href=\"/w/" + c.opt.PageName + "/" + parts[2] + "\">" + parts[3] + "</a>"
		case "../":
			idx := strings.LastIndex(c.opt.PageName, "/")
			fmt.Println(c.opt.PageName[:idx], " ", idx)
			if idx > 0 {
				return "<a href=\"/w/" + c.opt.PageName[:idx] + "/" + parts[2] + "\">" + parts[3] + "</a>"
			}
		}
		return ""
	})
	return c
}

func (c *Compiler) Decorator() *Compiler {
	re := regexp.MustCompile("\\-\\-([^\\-]+)\\-\\-")
	c.text = re.ReplaceAllString(c.text, "<span class='linethrough'>$1</span>")
	re = regexp.MustCompile("__([^\\_]+)__")
	c.text = re.ReplaceAllString(c.text, "<span class='underline'>$1</span>")
	return c
}

func (c *Compiler) Processor() int {
	re := regexp.MustCompile("^[ \t]*}}}")
	if len(re.FindString(c.text)) > 0 {
		c.TextType = 0
		c.processorBuffer.Reset()
		return PROCESSOR_CLOSE
	} else {
		c.processorBuffer.WriteString(c.text + "\n")
	}
	return 0
}

func (c *Compiler) ExecPlugin() *Compiler {
	var buffer bytes.Buffer
	re := regexp.MustCompile("\\[\\[[a-zA-Z]+\\(*[^)]+\\)*\\]\\]")
	data := re.FindAllStringIndex(c.text, 1)
	if len(data) > 0 {
		for _, v := range data {
			name := c.text[v[0]+2 : v[1]-2]
			sidx := strings.Index(name, "(")
			var (
				parameter string
				pname     string
			)
			if sidx > 0 {
				parameter = name[sidx+1 : len(name)-1]
				pname = name[:sidx]
			} else {
				pname = name
			}

			rtv, err := c.P.Exec(pname, parameter, c.w, c.r)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("PLUGIN RUN : ", pname)
			}

			buffer.WriteString(c.text[0:v[0]])
			buffer.WriteString(rtv) // edit by yundream
			buffer.WriteString(c.text[v[1]:])
		}
		c.text = buffer.String()
	}
	return c
}

func (c *Compiler) SetOptions(opt RequestOptions) {
	c.opt = opt
	//c.localPlugin.SetOptions(c.opt)
}

func (c *Compiler) OuterLink() *Compiler {
	GoogleImg := struct {
		Id  string
		Alt string
	}{}
	thumNail := `
	<p>
	<a data-open="{{.Id}}">
    <img src="https://drive.google.com/uc?export=view&id={{.Id}}" width="50%" height="50%" alt="{{.Alt}}">
	</a>
	</p>

	<!-- This is the first modal -->
	<div class="large reveal" id="{{.Id}}" data-reveal>
	<img src="https://drive.google.com/uc?export=view&id={{.Id}}" width="100%" alt="{{.Alt}}">
	<button class="close-button" data-close aria-label="Close reveal" type="button">
	<span aria-hidden="true">&times;</span>
	</button>
	</div>`
	re := regexp.MustCompile("\\[(http[s]*://[^ \t]+)([^\\]]+)\\]")
	c.text = re.ReplaceAllStringFunc(c.text, func(str string) string {
		parts := re.FindStringSubmatch(str)
		if strings.Index(parts[0], "docs.google.com/drawings") > 0 {
			return "<img src=\"" + parts[1] + "\" alt=\"" + parts[2] + "\">"
		} else if strings.Index(parts[0], "drive.google.com/file/d") > 0 {
			items := strings.Split(parts[1], "/")
			GoogleImg.Id = items[5]
			GoogleImg.Alt = parts[2]
			t := template.New("Googl Images")
			t, err := t.Parse(thumNail)
			if err != nil {
				fmt.Println(err.Error())
			}
			var output bytes.Buffer
			err = t.Execute(&output, GoogleImg)
			if err != nil {
				fmt.Println(err.Error())
			}
			return output.String()
		} else {
			if strings.LastIndex(c.text, ".gif") > 0 {
				return "<img src=\"" + parts[1] + "\" alt=\"" + parts[2] + "\">"
			}
			if strings.LastIndex(c.text, ".png") > 0 {
				return "<img src=\"" + parts[1] + "\" alt=\"" + parts[2] + "\">"
			}
			if strings.LastIndex(c.text, ".jpg") > 0 {
				return "<img src=\"" + parts[1] + "\" alt=\"" + parts[2] + "\">"
			}
		}
		return "<a href=\"" + parts[1] + "\" class=\"external\">" + parts[2] + "</a>"
	})
	return c
}

func (c *Compiler) Comma() *Compiler {
	var buffer bytes.Buffer
	re := regexp.MustCompile("[']+")
	data := re.FindAllStringIndex(c.text, -1)
	if len(data) > 0 {
		offset := 0
		iloop := 0
		for _, v := range data {
			buffer.WriteString(c.text[offset:v[0]])
			quotationNum := v[1] - v[0]
			switch quotationNum {
			case 3:
				if iloop%2 == 0 {
					buffer.WriteString("<b>")
				} else {
					buffer.WriteString("</b>")
				}
			case 2:
				if iloop%2 == 0 {
					buffer.WriteString("<i>")
				} else {
					buffer.WriteString("</i>")
				}
			}
			offset = v[1]
			iloop++
		}
		buffer.WriteString(c.text[offset:])
		c.text = buffer.String()
	}
	return c
}

func (c *Compiler) Print() {
	fmt.Println(c.text)
}

func (c *Compiler) String() string {
	return c.text + "\n"
}

func (c *Compiler) Line(line string) *Compiler {
	c.text = line
	if c.TextType == HEAD {
		c.TextType = 0
	}

	c.Linenum++
	re := regexp.MustCompile("^[ \t]*{{{#!([a-z]+)")
	data := re.FindStringSubmatch(c.text)
	if len(data) > 0 {
		c.TextType = PROCESSOR_OPEN
		c.text = ""
		//c.processName = data[1]
	}
	return c
}
