package srt

import (
	"fmt"
	"github.com/juxuny/data-utils/lib"
	"github.com/pkg/errors"
	"strings"
)

type NodeType int8

const (
	NodeTypeUnknown = NodeType(0)
	NodeTypeText    = NodeType(1)
	NodeTypeTag     = NodeType(2)
)

type Node struct {
	Type     NodeType
	TagName  []byte
	Content  []byte
	Attr     map[string]string
	Children []*Node
}

func NewNode(nodeType NodeType) *Node {
	return &Node{
		Type: nodeType,
		Attr: map[string]string{},
	}
}

func (t *Node) String() string {
	if t.Type == NodeTypeTag {
		ret := []byte{'<'}
		ret = append(ret, t.TagName...)

		kvs := make([]string, 0)
		// attr value
		for k, v := range t.Attr {
			kvs = append(kvs, fmt.Sprintf("%s=\"%s\"", k, v))
		}
		if len(kvs) > 0 {
			ret = append(ret, ' ')
			ret = append(ret, []byte(strings.Join(kvs, " "))...)
		}
		ret = append(ret, '>')

		// add children
		var children []string
		for _, child := range t.Children {
			children = append(children, child.String())
		}
		ret = append(ret, []byte(strings.Join(children, ""))...)

		// close tag
		ret = append(ret, []byte(fmt.Sprintf("</%s>", t.TagName))...)
		return string(ret)
	} else if t.Type == NodeTypeText {
		var ret []byte
		ret = append(ret, t.Content...)
		var children []string
		for _, child := range t.Children {
			children = append(children, child.String())
		}
		ret = append(ret, strings.Join(children, "")...)
		return string(ret)
	}
	return ""
}

type SegmentType int8

const (
	SegmentTypeTag     = SegmentType(1)
	SegmentTypeContent = SegmentType(2)
	SegmentTypeEndTag  = SegmentType(3)
)

type segment struct {
	Content   []byte
	Type      SegmentType
	Children  []*segment
	Depth     int
	closedTag bool
	parent    *segment
	hasTag    bool // 记录 SegmentTypeContent 类型的节点是否包含标签
}

func newSegment(parent *segment, t SegmentType, depth int) *segment {
	return &segment{Type: t, Depth: depth, parent: parent}
}

func closeTab(s *segment) {
	if s == nil {
		return
	}
	s.closedTag = true
	if s.Type == SegmentTypeTag {
		return
	}
	closeTab(s.parent)
}

func parseSegment(data []byte, startIndex int, s *segment) (offset int, err error) {
	i := startIndex
	for i < len(data) {
		if s.closedTag {
			return i - startIndex, nil
		}
		if data[i] == '<' {
			if i+1 < len(data) && data[i+1] == '/' && s.Type != SegmentTypeEndTag {
				endSeg := newSegment(s, SegmentTypeEndTag, s.Depth-1)
				if offset, err := parseSegment(data, i, endSeg); err != nil {
					return 0, errors.Wrap(err, "parse end tag segment failed")
				} else {
					i += offset
					closeTab(s)
					return i - startIndex, nil
				}
			}
			if s.Type == SegmentTypeContent {
				s.hasTag = true
				contentSegment := newSegment(s, SegmentTypeTag, s.Depth)
				if offset, err := parseSegment(data, i, contentSegment); err != nil {
					return 0, errors.Wrap(err, "parse content segment ")
				} else {
					i += offset
					s.Children = append(s.Children, contentSegment)
					continue
				}
			}
		}
		if data[i] == '>' && s.Type == SegmentTypeEndTag {
			s.Content = append(s.Content, '>')
			return i - startIndex + 1, nil
		}
		if data[i] == '>' && s.Type == SegmentTypeTag {
			s.Content = append(s.Content, '>')
			i += 1
			contentSegment := newSegment(s, SegmentTypeContent, s.Depth+1)
			if offset, err := parseSegment(data, i, contentSegment); err != nil {
				return 0, errors.Wrap(err, "parse '>' failed")
			} else {
				i += offset
				s.Children = append(s.Children, contentSegment)
				continue
			}
		}
		if s.Type == SegmentTypeContent && s.hasTag {
			seg := newSegment(s, SegmentTypeContent, s.Depth)
			if offset, err := parseSegment(data, i, seg); err != nil {
				return 0, errors.Wrap(err, "parse body failed")
			} else {
				i += offset
				s.Children = append(s.Children, seg)
				continue
			}
		}
		s.Content = append(s.Content, data[i])
		i += 1
	}
	return len(data) - 1, nil
}

func parseSegmentWrapper(data []byte) (ret []*segment, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == '<' {
			tag := newSegment(nil, SegmentTypeTag, 0)
			if offset, err := parseSegment(data, i, tag); err != nil {
				return nil, errors.Wrap(err, "parse tag failed")
			} else {
				ret = append(ret, tag)
				i += offset
			}
		}
	}
	return
}

func (t *segment) Parse() (ret *Node, err error) {
	ret = NewNode(NodeTypeUnknown)
	if t.Type == SegmentTypeTag || t.Type == SegmentTypeEndTag {
		ret.Type = NodeTypeTag
		list := splitTagRawData(t.Content)
		if len(list) == 0 {
			return nil, errors.Errorf("invalid tag data: %s", string(t.Content))
		}
		for i := range list {
			list[i] = Trim(list[i], " ")
		}
		ret.TagName = list[0]
		var keyName []byte
		for i := 0; i < len(list); i++ {
			if lib.IsQuoted(list[i]) {
				if len(keyName) == 0 {
					return nil, errors.Errorf("invalid tag attr: %s", string(t.Content))
				}
				ret.Attr[string(keyName)] = string(Trim(list[i], "\""))
			} else {
				keyName = list[i]
			}
		}
	} else if t.Type == SegmentTypeContent {
		ret.Type = NodeTypeText
		ret.Content = make([]byte, len(t.Content))
		copy(ret.Content, t.Content)
	} else {
		return nil, errors.Errorf("unkonwn segment type: %v", t.Type)
	}
	if len(t.Children) > 0 {
		for _, child := range t.Children {
			if childNode, err := child.Parse(); err != nil {
				return nil, errors.Wrap(err, "parse child failed")
			} else {
				ret.Children = append(ret.Children, childNode)
			}
		}
	}
	return
}

func ParseSubtitle(data []byte) (ret []*Node, err error) {
	seg, err := parseSegmentWrapper(data)
	if err != nil {
		return nil, errors.Wrap(err, "parse segment failed")
	}
	for _, s := range seg {
		if n, err := s.Parse(); err != nil {
			return nil, errors.Wrap(err, "parse segment failed")
		} else {
			ret = append(ret, n)
		}
	}
	return
}
