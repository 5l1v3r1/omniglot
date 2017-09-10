package omniglot

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"sort"

	"github.com/unixpickle/essentials"
)

// Set is a set of augmented samples.
type Set []*AugSample

// ReadSet reads a sample set from a directory.
//
// No augmentation is performed; all of the samples will
// have rotation 0.
func ReadSet(dirPath string) (set Set, err error) {
	defer essentials.AddCtxTo("read omniglot set", &err)
	listing, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, item := range listing {
		if !item.IsDir() {
			continue
		}
		alphabetDir := filepath.Join(dirPath, item.Name())
		charItems, err := ioutil.ReadDir(alphabetDir)
		if err != nil {
			return nil, err
		}
		for _, charItem := range charItems {
			if !charItem.IsDir() {
				continue
			}
			charPath := filepath.Join(alphabetDir, charItem.Name())
			imageItems, err := ioutil.ReadDir(charPath)
			if err != nil {
				return nil, err
			}
			for _, imageItem := range imageItems {
				if filepath.Ext(imageItem.Name()) != ".png" {
					continue
				}
				imagePath := filepath.Join(charPath, imageItem.Name())
				set = append(set, &AugSample{
					Sample: &Sample{
						Alphabet: item.Name(),
						CharName: charItem.Name(),
						Path:     imagePath,
					},
				})
			}
		}
	}
	return
}

// Augment generates a sample set with rotated samples.
func (s Set) Augment() Set {
	var res Set
	for _, x := range s {
		for i := 0; i < 4; i++ {
			res = append(res, x.rotated(i))
		}
	}
	return res
}

// ByClass sorts the set into different subsets, one per
// class.
// A class is an alphabet-character-rotation triple.
//
// The resulting classes are ordered deterministically
// regardless of the order of s.
// Within each class, the samples are ordered the same way
// that they are ordered in s.
func (s Set) ByClass() []Set {
	classByName := map[string]Set{}
	classNames := []string{}
	for _, obj := range s {
		className := fmt.Sprintf("%s/%s/%d", obj.Sample.Alphabet,
			obj.Sample.CharName, obj.Rotation)
		if classByName[className] == nil {
			classNames = append(classNames, className)
		}
		classByName[className] = append(classByName[className], obj)
	}
	sort.Strings(classNames)
	var res []Set
	for _, name := range classNames {
		res = append(res, classByName[name])
	}
	return res
}

// Select selects a random n elements from the set.
func (s Set) Select(n int) Set {
	if n > len(s) {
		panic("size out of bounds")
	}
	p := rand.Perm(len(s))
	res := make(Set, n)
	for i, j := range p[:n] {
		res[i] = s[j]
	}
	return res
}
