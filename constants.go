/*
 * Copyright (c) 2014 Mark Samman <https://github.com/marksamman/golinkshortener>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package main

var decodeArray = [...]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 35, -1, -1, 13, 49, 0, 56, 3, 62, 26, 63, 59, 11, -1, -1, -1, -1, -1, -1, -1, 43, 37, 54, 2, 32, 39, 33, 16, 4, 30, 40, 57, 31, 51, 15, 61, 20, 55, 29, 36, 45, 38, 17, 7, 28, 24, -1, -1, -1, -1, 44, -1, 6, 27, 25, 50, 10, 1, 9, 19, 12, 47, 22, 34, 52, 42, 5, 8, 58, 21, 18, 46, 53, 41, 48, 14, 23, 60}
var urlSafe = [...]byte{
	'2', 'f', 'D', '4', 'I', 'o', 'a', 'X',
	'p', 'g', 'e', '9', 'i', '0', 'x', 'O',
	'H', 'W', 's', 'h', 'Q', 'r', 'k', 'y',
	'Z', 'c', '6', 'b', 'Y', 'S', 'J', 'M',
	'E', 'G', 'l', '-', 'T', 'B', 'V', 'F',
	'K', 'v', 'n', 'A', '_', 'U', 't', 'j',
	'w', '1', 'd', 'N', 'm', 'u', 'C', 'R',
	'3', 'L', 'q', '8', 'z', 'P', '5', '7',
}
