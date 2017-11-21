// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package nvm

import "C"

import (
	"regexp"
	"strings"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

var (
	pathRe = regexp.MustCompile("^\\.{0,2}/")
)

// Module
type Module struct {
	id         string
	source     string
	lineOffset int
}

type Modules map[string]*Module

// NewModules create new modules and return it.
func NewModules() Modules {
	return make(Modules, 1)
}

// NewModule create new module and return it.
func NewModule(id, source string, lineOffset int) *Module {
	paths := strings.Split(id, "/")
	if !pathRe.MatchString(id) {
		paths = append([]string{"lib"}, paths...)
	}

	id = strings.Join(paths, "/")

	return &Module{
		id:         id,
		source:     source,
		lineOffset: lineOffset,
	}
}

// Add add source to module.
func (ms Modules) Add(m *Module) {
	ms[m.id] = m
}

// Get get module from Modules by id.
func (ms Modules) Get(id string) *Module {
	return ms[id]
}

// RequireDelegateFunc delegate func for require.
//export RequireDelegateFunc
func RequireDelegateFunc(handler unsafe.Pointer, filename *C.char, lineOffset *C.size_t) *C.char {
	id := C.GoString(filename)
	log.Debugf("require load %s", id)

	e := getEngineByEngineHandler(handler)
	if e == nil {
		log.WithFields(log.Fields{
			"filename": id,
		}).Error("require delegate handler does not found.")
		return nil
	}

	module := e.modules.Get(id)
	if module == nil {
		return nil
	}

	*lineOffset = C.size_t(module.lineOffset)
	cSource := C.CString(module.source)
	return cSource
}