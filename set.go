package omniglot

import (
	"io/ioutil"
	"path/filepath"

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
