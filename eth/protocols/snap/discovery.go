// Copyright 2020 The go-neworder Authors
// This file is part of the go-neworder library.
//
// The go-neworder library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-neworder library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-neworder library. If not, see <http://www.gnu.org/licenses/>.

package snap

import (
	"github.com/neworder-network/go-neworder/rlp"
)

// enrEntry is the ENR entry which advertises `snap` protocol on the discovery.
type enrEntry struct {
	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e enrEntry) ENRKey() string {
	return "snap"
}
