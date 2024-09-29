package js

import (
	"encoding/json"
	"html"
	"strings"
)

const (
	defaultTransitionTime = 300
)

type Operation []any

func JS(ops ...Operation) string {
	jsonData, err := json.Marshal(ops)
	if err != nil {
		return ""
	}
	return html.EscapeString(string(jsonData))
}

type PushArgs struct {
	Target      string         `json:"target,omitempty"`
	Loading     string         `json:"loading,omitempty"`
	PageLoading string         `json:"page_loading,omitempty"`
	Value       map[string]any `json:"value,omitempty"`
}

type pushArgs struct {
	Event string `json:"event,omitempty"`
	*PushArgs
}

func Push(event string, args *PushArgs) Operation {
	a := &pushArgs{
		Event:    event,
		PushArgs: args,
	}

	if a.PushArgs == nil {
		a.PushArgs = &PushArgs{}
	}

	return Operation{"push", a}
}

type DispatchArgs struct {
	Event   string         `json:"event,omitempty"`
	To      string         `json:"to,omitempty"`
	Detail  map[string]any `json:"detail,omitempty"`
	Bubbles *bool          `json:"bubbles,omitempty"`
}

type dispatchArgs struct {
	Event string `json:"event,omitempty"`
	*DispatchArgs
}

func Dispatch(event string, args *DispatchArgs) Operation {
	a := &dispatchArgs{
		Event:        event,
		DispatchArgs: args,
	}

	if a.DispatchArgs == nil {
		a.DispatchArgs = &DispatchArgs{}
	}

	if a.Bubbles == nil {
		a.Bubbles = Bool(true)
	}
	return Operation{"dispatch", a}
}

type ToggleArgs struct {
	To       string    `json:"to,omitempty"`
	In       [3]string `json:"-"`
	Out      [3]string `json:"-"`
	Time     *int      `json:"time,omitempty"`
	Display  string    `json:"display,omitempty"`
	Blocking *bool     `json:"blocking,omitempty"`
}

type toggleArgs struct {
	Event string      `json:"event,omitempty"`
	In    [3][]string `json:"in,omitempty"`
	Out   [3][]string `json:"out,omitempty"`
	*ToggleArgs
}

func Toggle(event string, args *ToggleArgs) Operation {
	a := &toggleArgs{
		Event:      event,
		ToggleArgs: args,
	}

	if a.ToggleArgs == nil {
		a.ToggleArgs = &ToggleArgs{}
	}

	a.In = [3][]string{
		classNames(args.In[0]),
		classNames(args.In[1]),
		classNames(args.In[2]),
	}

	a.Out = [3][]string{
		classNames(args.Out[0]),
		classNames(args.Out[1]),
		classNames(args.Out[2]),
	}

	if a.Time == nil {
		a.Time = Int(defaultTransitionTime)
	}

	if a.Display == "" {
		a.Display = "block"
	}

	if a.Blocking == nil {
		a.Blocking = Bool(true)
	}

	return Operation{"toggle", a}
}

type ShowArgs struct {
	To         string    `json:"to,omitempty"`
	Transition [3]string `json:"-"`
	Time       *int      `json:"time,omitempty"`
	Blocking   *bool     `json:"blocking,omitempty"`
	Display    string    `json:"display,omitempty"`
}

type showArgs struct {
	Transition [3][]string `json:"transition,omitempty"`
	*ShowArgs
}

func Show(args *ShowArgs) Operation {
	a := &showArgs{
		ShowArgs: args,
	}

	if a.ShowArgs == nil {
		a.ShowArgs = &ShowArgs{}
	}

	a.Transition = [3][]string{
		classNames(args.Transition[0]),
		classNames(args.Transition[1]),
		classNames(args.Transition[2]),
	}

	if a.Time == nil {
		a.Time = Int(defaultTransitionTime)
	}

	if a.Blocking == nil {
		a.Blocking = Bool(true)
	}

	return Operation{"show", a}
}

type HideArgs struct {
	To         string    `json:"to,omitempty"`
	Transition [3]string `json:"-"`
	Time       *int      `json:"time,omitempty"`
	Blocking   *bool     `json:"blocking,omitempty"`
}

type hideArgs struct {
	Transition [3][]string `json:"transition,omitempty"`
	*HideArgs
}

func Hide(args *HideArgs) Operation {
	a := &hideArgs{
		HideArgs: args,
	}

	if a.HideArgs == nil {
		a.HideArgs = &HideArgs{}
	}

	a.Transition = [3][]string{
		classNames(args.Transition[0]),
		classNames(args.Transition[1]),
		classNames(args.Transition[2]),
	}

	if a.Time == nil {
		a.Time = Int(defaultTransitionTime)
	}

	if a.Blocking == nil {
		a.Blocking = Bool(true)
	}

	return Operation{"hide", a}
}

type AddClassArgs struct {
	To         string    `json:"to,omitempty"`
	Transition [3]string `json:"-"`
	Time       *int      `json:"time,omitempty"`
	Blocking   *bool     `json:"blocking,omitempty"`
}

type addClassArgs struct {
	Names      []string    `json:"names,omitempty"`
	Transition [3][]string `json:"transition,omitempty"`
	*AddClassArgs
}

func AddClass(names string, args *AddClassArgs) Operation {
	a := &addClassArgs{
		Names:        classNames(names),
		AddClassArgs: args,
	}

	if a.AddClassArgs == nil {
		a.AddClassArgs = &AddClassArgs{}
	}

	a.Transition = [3][]string{
		classNames(args.Transition[0]),
		classNames(args.Transition[1]),
		classNames(args.Transition[2]),
	}

	if a.Time == nil {
		a.Time = Int(defaultTransitionTime)
	}

	if a.Blocking == nil {
		a.Blocking = Bool(true)
	}

	return Operation{"add_class", a}
}

type ToggleClassArgs struct {
	To         string    `json:"to,omitempty"`
	Transition [3]string `json:"-"`
	Time       *int      `json:"time,omitempty"`
	Blocking   *bool     `json:"blocking,omitempty"`
}

type toggleClassArgs struct {
	Names      []string    `json:"names,omitempty"`
	Transition [3][]string `json:"transition,omitempty"`
	*ToggleClassArgs
}

func ToggleClass(names string, args *ToggleClassArgs) Operation {
	a := &toggleClassArgs{
		Names:           classNames(names),
		ToggleClassArgs: args,
	}

	if a.ToggleClassArgs == nil {
		a.ToggleClassArgs = &ToggleClassArgs{}
	}

	a.Transition = [3][]string{
		classNames(args.Transition[0]),
		classNames(args.Transition[1]),
		classNames(args.Transition[2]),
	}

	if a.Time == nil {
		a.Time = Int(defaultTransitionTime)
	}

	if a.Blocking == nil {
		a.Blocking = Bool(true)
	}

	return Operation{"toggle_class", a}
}

type RemoveClassArgs struct {
	To         string    `json:"to,omitempty"`
	Transition [3]string `json:"-"`
	Time       *int      `json:"time,omitempty"`
	Blocking   *bool     `json:"blocking,omitempty"`
}

type removeClassArgs struct {
	Names      []string    `json:"names,omitempty"`
	Transition [3][]string `json:"transition,omitempty"`
	*RemoveClassArgs
}

func RemoveClass(names string, args *RemoveClassArgs) Operation {
	a := &removeClassArgs{
		Names:           classNames(names),
		RemoveClassArgs: args,
	}

	if a.RemoveClassArgs == nil {
		a.RemoveClassArgs = &RemoveClassArgs{}
	}

	a.Transition = [3][]string{
		classNames(args.Transition[0]),
		classNames(args.Transition[1]),
		classNames(args.Transition[2]),
	}

	if a.Time == nil {
		args.Time = Int(defaultTransitionTime)
	}

	if a.Blocking == nil {
		a.Blocking = Bool(true)
	}

	return Operation{"remove_class", a}
}

type TransitionArgs struct {
	To       string `json:"to,omitempty"`
	Time     *int   `json:"time,omitempty"`
	Blocking *bool  `json:"blocking,omitempty"`
}

type transitionArgs struct {
	Transition [3][]string `json:"transition,omitempty"`
	*TransitionArgs
}

func Transition(transition [3]string, args *TransitionArgs) Operation {
	a := &transitionArgs{
		Transition: [3][]string{
			classNames(transition[0]),
			classNames(transition[1]),
			classNames(transition[2]),
		},
		TransitionArgs: args,
	}

	if a.TransitionArgs == nil {
		a.TransitionArgs = &TransitionArgs{}
	}

	if a.Time == nil {
		a.Time = Int(defaultTransitionTime)
	}

	if a.Blocking == nil {
		a.Blocking = Bool(true)
	}

	return Operation{"transition", a}
}

type SetAttrArgs struct {
	To string `json:"to,omitempty"`
}

type setAttrArgs struct {
	Attr [2]any `json:"attr,omitempty"`
	*SetAttrArgs
}

func SetAttr(key string, value any, args *SetAttrArgs) Operation {
	a := &setAttrArgs{
		Attr:        [2]any{key, value},
		SetAttrArgs: args,
	}

	if a.SetAttrArgs == nil {
		a.SetAttrArgs = &SetAttrArgs{}
	}

	return Operation{"set_attr", a}
}

type RemoveAttrArgs struct {
	To string `json:"to,omitempty"`
}

type removeAttrArgs struct {
	Attr string `json:"attr,omitempty"`
	*RemoveAttrArgs
}

func RemoveAttr(attr string, args *RemoveAttrArgs) Operation {
	a := &removeAttrArgs{
		Attr:           attr,
		RemoveAttrArgs: args,
	}

	if a.RemoveAttrArgs != nil {
		a.RemoveAttrArgs = &RemoveAttrArgs{}
	}

	return Operation{"remove_attr", a}
}

type ToggleAttrArgs struct {
	To string `json:"to,omitempty"`
}

type toggleAttrArgs struct {
	Attr [2]any `json:"attr,omitempty"`
	*ToggleAttrArgs
}

func ToggleAttr(attr string, val string, args *ToggleAttrArgs) Operation {
	a := &toggleAttrArgs{
		Attr:           [2]any{attr, val},
		ToggleAttrArgs: args,
	}

	if a.ToggleAttrArgs == nil {
		a.ToggleAttrArgs = &ToggleAttrArgs{}
	}

	return Operation{"toggle_attr", a}
}

type ToggleAttrsArgs struct {
	To string `json:"to,omitempty"`
}

type toggleAttrsArgs struct {
	Attrs [3]any `json:"attrs,omitempty"`
	*ToggleAttrsArgs
}

func ToggleAttrs(attr string, val1 string, val2 string, args *ToggleAttrsArgs) Operation {
	a := &toggleAttrsArgs{
		Attrs:           [3]any{attr, val1, val2},
		ToggleAttrsArgs: args,
	}

	if a.ToggleAttrsArgs == nil {
		a.ToggleAttrsArgs = &ToggleAttrsArgs{}
	}

	return Operation{"toggle_attr", a}
}

type FocusArgs struct {
	To string `json:"to,omitempty"`
}

func Focus(args *FocusArgs) Operation {
	if args == nil {
		args = &FocusArgs{}
	}
	return Operation{"focus", args}
}

type FocusFirstArgs struct {
	To string `json:"to,omitempty"`
}

func FocusFirst(args *FocusFirstArgs) Operation {
	if args == nil {
		args = &FocusFirstArgs{}
	}
	return Operation{"focus_first", args}
}

type PushFocusArgs struct {
	To string `json:"to,omitempty"`
}

func PushFocus(args *PushFocusArgs) Operation {
	if args == nil {
		args = &PushFocusArgs{}
	}
	return Operation{"push_focus", args}
}

func PopFocus() Operation {
	return Operation{"pop_focus", []any{}}
}

type NavigateArgs struct {
	Replace bool `json:"replace,omitempty"`
}

type navigateArgs struct {
	Href string `json:"href,omitempty"`
	*NavigateArgs
}

func Navigate(href string, args *NavigateArgs) Operation {
	a := &navigateArgs{
		Href:         href,
		NavigateArgs: args,
	}

	if a.NavigateArgs == nil {
		a.NavigateArgs = &NavigateArgs{}
	}

	return Operation{"navigate", a}
}

type PatchArgs struct {
	Replace bool `json:"replace,omitempty"`
}

type patchArgs struct {
	Href string `json:"href,omitempty"`
	*PatchArgs
}

func Patch(href string, args *PatchArgs) Operation {
	a := &patchArgs{
		Href:      href,
		PatchArgs: args,
	}

	if a.PatchArgs == nil {
		a.PatchArgs = &PatchArgs{}
	}

	return Operation{"patch", a}
}

type ExecArgs struct {
	To string `json:"to,omitempty"`
}

type execArgs struct {
	Attr string `json:"attr,omitempty"`
	*ExecArgs
}

func Exec(attr string, args *ExecArgs) Operation {
	a := &execArgs{
		Attr:     attr,
		ExecArgs: args,
	}

	if a.ExecArgs == nil {
		a.ExecArgs = &ExecArgs{}
	}

	return Operation{"exec", a}
}

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}

func classNames(s string) []string {
	return strings.Fields(s)
}
