package rend

import (
	"reflect"
)

func (oldRoot *Root) Diff(newRoot *Root) *Root {
	// if the root fingerprint changed, force a full render
	if oldRoot.Rend.Fingerprint != newRoot.Rend.Fingerprint {
		return newRoot
	}

	root := compareComponents(oldRoot, newRoot)
	root.Rend = compareRend(oldRoot.Rend, newRoot.Rend)

	if root.Components == nil && root.Rend == nil {
		return nil
	}

	return root
}

func compareComponents(oldRoot, newRoot *Root) *Root {
	root := &Root{}

	for key, newComponent := range newRoot.Components {
		oldComponent, exists := oldRoot.Components[key]
		if !exists {
			if root.Components == nil {
				root.Components = make(map[int64]*Rend)
			}
			root.Components[key] = newComponent
			continue
		}

		diff, changed := elementsEqual(oldComponent, newComponent)
		if changed {
			if root.Components == nil {
				root.Components = make(map[int64]*Rend)
			}
			rend, ok := diff.(*Rend)
			if ok {
				root.Components[key] = rend
			}
		}
	}

	return root
}

func compareComprehension(oldComp, newComp *Comprehension) *Comprehension {
	diff := &Comprehension{
		Stream: newComp.Stream,
	}

	if oldComp.Fingerprint != newComp.Fingerprint {
		return newComp
	}

	if len(oldComp.Dynamics) != len(newComp.Dynamics) {
		diff.Dynamics = newComp.Dynamics
		return diff
	}

	for i, newDynamicSlice := range newComp.Dynamics {
		oldDynamicSlice := oldComp.Dynamics[i]

		for j, newDynamicElement := range newDynamicSlice {
			oldDynamicElement := oldDynamicSlice[j]
			if _, different := elementsEqual(oldDynamicElement, newDynamicElement); different {
				// if any of the elements are different, force a full render
				diff.Dynamics = newComp.Dynamics
				return diff
			}
		}
	}

	if len(newComp.Stream) > 0 {
		return diff
	}

	return nil
}

func compareRend(oldRend, newRend *Rend) *Rend {
	if oldRend.Fingerprint != newRend.Fingerprint {
		return newRend
	}

	diff := &Rend{}

	for key, newDynamic := range newRend.Dynamic {
		oldDynamic, exists := oldRend.Dynamic[key]
		if !exists {
			if diff.Dynamic == nil {
				diff.Dynamic = make(map[string]any)
			}
			diff.Dynamic[key] = newDynamic
			continue
		}

		if !sameType(oldDynamic, newDynamic) {
			if diff.Dynamic == nil {
				diff.Dynamic = make(map[string]any)
			}
			diff.Dynamic[key] = newDynamic
			continue
		}

		newVal, changed := elementsEqual(oldDynamic, newDynamic)
		if changed {
			if diff.Dynamic == nil {
				diff.Dynamic = make(map[string]any)
			}
			diff.Dynamic[key] = newVal
		}
	}

	return diff
}

func elementsEqual(a, b any) (any, bool) {
	if a == nil || b == nil {
		return b, a != b
	}

	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return bVal, aVal != bVal
		}
	case string:
		if bVal, ok := b.(string); ok {
			return bVal, aVal != bVal
		}
	case *Rend:
		if bVal, ok := b.(*Rend); ok {
			diff := compareRend(aVal, bVal)
			if diff != nil {
				return diff, true
			}
			return nil, false
		}
	case *Comprehension:
		if bVal, ok := b.(*Comprehension); ok {
			diff := compareComprehension(aVal, bVal)
			if diff != nil {
				return diff, true
			}
			return nil, false
		}
	default:
		// Fall back to reflect.DeepEqual for unknown types
		areEqual := reflect.DeepEqual(a, b)
		if !areEqual {
			return b, true
		}
	}

	return nil, true
}

func sameType(a, b any) bool {
	switch a.(type) {
	case int64:
		_, ok := b.(int64)
		return ok
	case string:
		_, ok := b.(string)
		return ok
	case *Rend:
		_, ok := b.(*Rend)
		return ok
	case *Comprehension:
		_, ok := b.(*Comprehension)
		return ok
	default:
		// Fall back to reflection for unknown types
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	}
}
