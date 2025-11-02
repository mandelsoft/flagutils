package tree

import (
	"fmt"
	"github.com/mandelsoft/flagutils/utils/history"
	"github.com/mandelsoft/goutils/general"
	"strings"
)

type Object[T comparable] interface {
	history.HistoryProvider[T]
	IsNode() *T
}

type Typed interface {
	GetKind() string
}

type ValidTreeElement interface {
	IsValid() bool
}

type NodeCreator[T comparable] func(history.History[T], T) Object[T]

// TreeObject is an element enriched by a textual
// tree graph prefix line.
type TreeObject[T comparable] struct {
	Graph  string
	Object Object[T]
	Node   *TreeNode[T] // for synthesized nodes this entry is used if no object can be synthesized
}

func (t *TreeObject[T]) String() string {
	if t.Object != nil {
		return fmt.Sprintf("%s %s", t.Graph, t.Object)
	}
	return fmt.Sprintf("%s %s", t.Graph, t.Node.String())
}

type TreeNode[T comparable] struct {
	key      T
	History  history.History[T]
	CausedBy Object[T] // the object causing the synthesized node to be inserted
}

func (t *TreeNode[T]) String() string {
	if s, ok := any(t.key).(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", t.key)
}

var (
	vertical   = "│" + space[1:]
	horizontal = "─"
	corner     = "└" + horizontal
	fork       = "├" + horizontal
	space      = "   "
	node       = "⊗" // \u2297
)

// MapToTree maps a list of elements featuring a resolution history
// into a list of elements providing an ascii tree graph field
// Intermediate nodes are synthesized, so only leaf elements are required.
// If an element should act as explicit node, it must state to be a node,
// in this case the node will be tagged with the nodeSymbol. If this
// is not desired, pass an empty symbol string.
func MapToTree[T comparable](objs Objects[T], creator NodeCreator[T], symbols ...string) TreeObjects[T] {
	result := TreeObjects[T]{}
	nodeSym := general.OptionalDefaulted(node, symbols...)
	if nodeSym != "" && !strings.HasPrefix(nodeSym, " ") {
		nodeSym = " " + nodeSym
	}
	handleLevel(objs, "", nil, 0, creator, &result, nodeSym)
	return result
}

func handleLevel[T comparable](objs Objects[T], header string, prefix history.History[T], start int, creator NodeCreator[T], result *TreeObjects[T], nodeSym string) {
	var node *T
	lvl := len(prefix)
	for i := start; i < len(objs); {
		var next int
		h := objs[i].GetHistory()
		if !h.HasPrefix(prefix) {
			return
		}
		ftag := corner
		stag := space
		key := objs[i].IsNode()
		for next = i + 1; next < len(objs); next++ {
			if s := objs[next].GetHistory(); s.HasPrefix(prefix) {
				if len(s) > lvl && len(h) > lvl && h[lvl] == s[lvl] { // skip same sub level
					continue
				}
				if key != nil {
					if len(s) > lvl && *key == s[lvl] { // skip same sub level
						continue
					}
				}
				ftag = fork
				stag = vertical
			}
			break
		}
		if len(h) == lvl {
			node = objs[i].IsNode() // Element acts as dedicate node
			sym := ""
			if node != nil {
				if i < len(objs)-1 {
					sub := objs[i+1].GetHistory()
					if len(sub) > len(h) && sub.HasPrefix(append(h, *node)) {
						sym = nodeSym
					}
				}
			}
			if t, ok := objs[i].(Typed); ok {
				k := t.GetKind()
				if k != "" {
					sym += " " + k
				}
			}
			if valid, ok := objs[i].(ValidTreeElement); !ok || valid.IsValid() {
				*result = append(*result, &TreeObject[T]{
					Graph:  header + ftag + sym,
					Object: objs[i],
				})
			}
			i++
		} else {
			if node == nil || *node != h[lvl] {
				// synthesize node if only leafs or non-matching node has been issued before
				var o Object[T]
				var n *TreeNode[T]
				if creator != nil {
					o = creator(prefix, h[len(prefix)])
				}
				if o == nil {
					n = &TreeNode[T]{h[len(prefix)], prefix, objs[i]}
				}
				*result = append(*result, &TreeObject[T]{
					Graph:  header + ftag, // + " " + h[len(prefix)].String(),
					Object: o,
					Node:   n,
				})
			}
			handleLevel(objs, header+stag, h[:len(prefix)+1], i, creator, result, nodeSym)
			i = next
			node = nil
		}
	}
}
