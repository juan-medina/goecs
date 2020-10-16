/*
 * Copyright (c) 2020 Juan Medina.
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

package goecs

import "testing"

var velocityCompType = NewComponentType()

type velocityComp struct {
	x float32
	y float32
}

func (v velocityComp) Type() ComponentType {
	return velocityCompType
}

var positionCompType = NewComponentType()

type positionComp struct {
	x float32
	y float32
}

func (p positionComp) Type() ComponentType {
	return positionCompType
}

func TestNewComponentType(t *testing.T) {
	vel := velocityComp{
		x: 1,
		y: 1,
	}

	pos := positionComp{
		x: 1,
		y: 1,
	}

	if vel.Type() == pos.Type() {
		t.Fatalf("vel and pos have the same type")
	}
}
