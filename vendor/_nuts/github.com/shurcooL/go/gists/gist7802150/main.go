package gist7802150

import (
	"strings"
	"sync"
)

type DepNode2I interface {
	Update()

	GetSources() []DepNode2I

	addSink(*DepNode2)
	getNeedToUpdate() bool
	markAllAsNeedToUpdate()
	markAsNotNeedToUpdate()
}

type DepNode2ManualI interface {
	DepNode2I
	manual() // Noop, just to separate it from automatic DepNode2I
}

// Used to make MakeUpdated resilient to concurrent access. Not really a good long term solution,
// but needed for now to prevent race conditions. Especially when called form http handler funcs.
var MakeUpdatedLock sync.Mutex

// Updates dependencies and itself, only if its dependencies have changed.
func MakeUpdated(this DepNode2I) {
	MakeUpdatedLock.Lock()
	makeUpdated(this)
	MakeUpdatedLock.Unlock()
}

func makeUpdated(this DepNode2I) {
	if !this.getNeedToUpdate() {
		return
	}
	for _, source := range this.GetSources() {
		makeUpdated(source)
	}
	this.Update()
	this.markAsNotNeedToUpdate()
}

// Updates dependencies and itself, regardless.
/*func ForceUpdated(this DepNode2I) {
	this.markAllAsNeedToUpdate()
	MakeUpdated(this)
}*/

// Updates only itself, regardless (skipping Update()).
func ExternallyUpdated(this DepNode2ManualI) {
	this.markAllAsNeedToUpdate()
	//this.markAsNotNeedToUpdate()
}

// ---

type DepNode2 struct {
	updated bool
	sources []DepNode2I
	sinks   []*DepNode2
}

func (this *DepNode2) GetSources() []DepNode2I {
	return this.sources
}

func (this *DepNode2) AddSources(sources ...DepNode2I) {
	this.updated = false
	this.sources = append(this.sources, sources...)
	for _, source := range sources {
		source.addSink(this)
	}
}

func (this *DepNode2) addSink(sink *DepNode2) {
	this.sinks = append(this.sinks, sink)
}

func (this *DepNode2) getNeedToUpdate() bool {
	return !this.updated
}

func (this *DepNode2) markAllAsNeedToUpdate() {
	this.updated = false
	for _, sink := range this.sinks {
		// TODO: See if this can be optimized away...
		sink.markAllAsNeedToUpdate()
	}
}

func (this *DepNode2) markAsNotNeedToUpdate() {
	this.updated = true
}

// ---

type DepNode2Manual struct {
	sinks []*DepNode2
}

func (this *DepNode2Manual) Update()                 { panic("") }
func (this *DepNode2Manual) GetSources() []DepNode2I { panic("") }
func (this *DepNode2Manual) addSink(sink *DepNode2) {
	this.sinks = append(this.sinks, sink)
}
func (this *DepNode2Manual) getNeedToUpdate() bool { return false }
func (this *DepNode2Manual) markAllAsNeedToUpdate() {
	for _, sink := range this.sinks {
		// TODO: See if this can be optimized away...
		sink.markAllAsNeedToUpdate()
	}
}
func (this *DepNode2Manual) markAsNotNeedToUpdate() { panic("") }
func (this *DepNode2Manual) manual()                { panic("") }

// Given there are two distinct DepNode2Manual structs, each having a pointer,
// merge takes other and merges it (along with its current sinks) into this.
// Afterwards, both pointers point to a single unified DepNode2Manual struct.
func (this *DepNode2Manual) merge(other **DepNode2Manual) {
	presentSinks := make(map[*DepNode2]struct{})
	for _, sink := range this.sinks {
		presentSinks[sink] = struct{}{}
	}

	for _, sink := range (*other).sinks {
		if _, present := presentSinks[sink]; !present {
			this.sinks = append(this.sinks, sink)
		}
	}

	*other = this
}

// ---

type DepNode2Func struct {
	UpdateFunc func(DepNode2I)
	DepNode2
}

func (this *DepNode2Func) Update() {
	this.UpdateFunc(this)
}

// =====

type ViewGroupI interface {
	SetSelf(string)

	AddAndSetViewGroup(ViewGroupI, string)
	RemoveView(ViewGroupI)

	GetUri() FileUri
	GetAllUris() []FileUri
	GetUriForProtocol(protocol string) (uri FileUri, ok bool)
	ContainsUri(FileUri) bool

	getViewGroup() *ViewGroup

	DepNode2ManualI
}

// FileUri represents a URI with a protocol notation.
//
// For example, "file:///tmp/foo" or "memory://???".
type FileUri string

// Path returns the path of URI, without the protocol notation.
func (u FileUri) Path() string {
	i := strings.Index(string(u), "://") + len("://")
	return string(u)[i:]
}

type ViewGroup struct {
	all *map[ViewGroupI]struct{}
	uri FileUri

	*DepNode2Manual
}

func (this *ViewGroup) getViewGroup() *ViewGroup {
	return this
}

// InitViewGroup must be called after creating a new ViewGroupI,
// before any other ViewGroup method or ViewGroupI func.
func (this *ViewGroup) InitViewGroup(self ViewGroupI, uri FileUri) {
	this.all = &map[ViewGroupI]struct{}{self: struct{}{}}
	this.uri = uri
	this.DepNode2Manual = &DepNode2Manual{}
}

// AddAndSetViewGroup adds another ViewGroupI and sets it to thisCurrent value, the current state of this ViewGroup.
func (this *ViewGroup) AddAndSetViewGroup(other ViewGroupI, thisCurrent string) {
	// Set other ViewGroup to thisCurrent
	for v := range *other.getViewGroup().all {
		v.SetSelf(thisCurrent)
	}
	ExternallyUpdated(other.getViewGroup().DepNode2Manual) // Notify whatever depended on the other ViewGroupI that it's been updated

	(*this.all)[other] = struct{}{}
	other.getViewGroup().all = this.all
	this.DepNode2Manual.merge(&other.getViewGroup().DepNode2Manual)
}

// RemoveView removes a single view from the ViewGroup.
func (this *ViewGroup) RemoveView(other ViewGroupI) {
	delete(*this.all, other)
	other.getViewGroup().InitViewGroup(other, other.GetUri())
}

func (this *ViewGroup) GetUri() FileUri {
	return this.uri
}
func (this *ViewGroup) GetAllUris() (uris []FileUri) {
	for v := range *this.all {
		uris = append(uris, v.GetUri())
	}
	return uris
}
func (this *ViewGroup) GetUriForProtocol(protocol string) (uri FileUri, ok bool) {
	for v := range *this.all {
		if strings.HasPrefix(string(v.GetUri()), protocol) {
			return v.GetUri(), true
		}
	}
	return "", false
}
func (this *ViewGroup) ContainsUri(uri FileUri) bool {
	for v := range *this.all {
		if uri == v.GetUri() {
			return true
		}
	}
	return false
}

func SetViewGroup(this ViewGroupI, s string) {
	for v := range *this.getViewGroup().all {
		v.SetSelf(s)
	}

	ExternallyUpdated(this)
}

func SetViewGroupOther(this ViewGroupI, s string) {
	for v := range *this.getViewGroup().all {
		if v != this {
			v.SetSelf(s)
		}
	}

	ExternallyUpdated(this)
}
