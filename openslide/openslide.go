/*
 * openslide-go - Unofficial Golang bindings for the OpenSlide library
 *
 * Copyright (c) 2020 GitHub user jammy-dodgers
 * https://github.com/jammy-dodgers/gophenslide
 * Copyright (c) 2022 Jonas Teuwen, Netherlands Cancer Institute
 *
 * The bindings have been modified from
 * https://github.com/jammy-dodgers/gophenslide
 *
 * This library is free software; you can redistribute it and/or modify it
 * under the terms of version 2.1 of the GNU Lesser General Public License
 * as published by the Free Software Foundation.
 *
 * This library is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public
 * License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this library; if not, write to the Free Software Foundation,
 * Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package openslide

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lopenslide
// #include <stdio.h>
// #include <stdlib.h>
// #include <stdint.h>
// #include <openslide.h>
// #include <openslide_go.h>
import "C"
import (
	"errors"
	"golang.org/x/image/draw"
	"image"
	"math"
	"strconv"
	"unsafe"
)

// Slide struct holding the OpenSlide pointer.
type Slide struct {
	ptr *C.openslide_t
}

// Open Get the Slide object referring to an OpenSlide image.
// Do not forget to defer opening the slide.
// This is an expensive operation, you will want to cache the result.
func Open(filename string) (Slide, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	slideData := C.openslide_open(cFilename)
	if slideData == nil {
		return Slide{nil}, errors.New("File " + filename + " unrecognized.")
	}
	return Slide{slideData}, nil
}

// Close Closes a slide.
func (slide Slide) Close() {
	C.openslide_close(slide.ptr)
}

// LevelCount Get the number of levels in the whole slide image.
func (slide Slide) LevelCount() int {
	return int(C.openslide_get_level_count(slide.ptr))
}

// LargestLevelDimensions Get the dimensions of level 0, the largest level (aka get_level0_dimensions).
func (slide Slide) LargestLevelDimensions() [2]int {
	var a, b C.int64_t
	C.openslide_get_level0_dimensions(slide.ptr, &a, &b)
	return [2]int{int(a), int(b)}
}

// LevelDimensions Get the dimensions of a level.
func (slide Slide) LevelDimensions(level int) [2]int {
	var a, b C.int64_t
	C.openslide_get_level_dimensions(slide.ptr, C.int32_t(level), &a, &b)
	return [2]int{int(a), int(b)}
}

// LevelDownsample Get the downsampling factor of the given level
func (slide Slide) LevelDownsample(level int) float64 {
	return float64(C.openslide_get_level_downsample(slide.ptr, C.int32_t(level)))
}

// LevelDownsamples Get the downsampling factors for all levels
func (slide Slide) LevelDownsamples() []float64 {
	downSamples := make([]float64, 0)
	for i := 0; i < int(slide.LevelCount()); i++ {
		levelDownsample := slide.LevelDownsample(i)
		downSamples = append(downSamples, levelDownsample)
	}
	return downSamples
}

// BestLevelForDownsample Get the best level to use for a particular downsampling factor
func (slide Slide) BestLevelForDownsample(downsample float64) int {
	return int(C.openslide_get_best_level_for_downsample(slide.ptr, C.double(downsample)))
}

// ReadRegion Read a region of the image as an RGBA image
func (slide Slide) ReadRegion(x, y int, level int, w, h int) (image.Image, error) {
	length := w * h * 4
	rawPtr := C.malloc(C.size_t(length))
	defer C.free(rawPtr)
	C.openslide_read_region(slide.ptr, (*C.uint32_t)(rawPtr), C.int64_t(x), C.int64_t(y), C.int32_t(level), C.int64_t(w), C.int64_t(h))
	if txt := C.openslide_get_error(slide.ptr); txt != nil {
		return nil, errors.New(C.GoString(txt))
	}

	// Convert ARGB to RGBA
	C.argb2rgba((*C.uint32_t)(rawPtr), C.int(length/4))
	rawArray := C.GoBytes(rawPtr, C.int(length))

	upLeft := image.Point{}
	lowRight := image.Point{X: w, Y: h}
	region := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	region.Pix = rawArray

	return region, nil
}

// AssociatedImageNames Get the names of the associated images
func (slide Slide) AssociatedImageNames() []string {
	cAssociatedImageNames := C.openslide_get_associated_image_names(slide.ptr)
	var strings []string
	for i := 0; C.str_at(cAssociatedImageNames, C.int(i)) != nil; i++ {
		strings = append(strings, C.GoString(C.str_at(cAssociatedImageNames, C.int(i))))
	}
	return strings
}

// AssociatedImageDimensions Get the dimensions of the associated images
func (slide Slide) AssociatedImageDimensions() map[string][2]int {
	associatedNames := slide.AssociatedImageNames()
	output := make(map[string][2]int)

	for _, associatedName := range associatedNames {
		var a, b C.int64_t
		C.openslide_get_associated_image_dimensions(slide.ptr, C.CString(associatedName), &a, &b)
		output[associatedName] = [2]int{int(a), int(b)}
	}
	return output
}

// ReadAssociatedImage Read an associated image as an RGBA image.
func (slide Slide) ReadAssociatedImage(associatedName string) (image.Image, error) {
	dimensions, ok := slide.AssociatedImageDimensions()[associatedName]
	if !ok {
		return nil, errors.New("associated image does not exist")
	}
	length := dimensions[0] * dimensions[1] * 4
	rawPtr := C.malloc(C.size_t(length))
	defer C.free(rawPtr)

	C.openslide_read_associated_image(slide.ptr, C.CString(associatedName), (*C.uint32_t)(rawPtr))
	if txt := C.openslide_get_error(slide.ptr); txt != nil {
		return nil, errors.New(C.GoString(txt))
	}

	// Convert ARGB to RGBA
	C.argb2rgba((*C.uint32_t)(rawPtr), C.int(length/4))
	rawArray := C.GoBytes(rawPtr, C.int(length))

	upLeft := image.Point{}
	lowRight := image.Point{X: dimensions[0], Y: dimensions[1]}
	region := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	region.Pix = rawArray

	return region, nil
}

// DetectVendor Quickly determine whether a whole slide image is recognized.
func DetectVendor(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	slideVendor := C.openslide_detect_vendor(cFilename)
	if slideVendor == nil {
		return "", errors.New("No vendor for " + filename)
	}
	return C.GoString(slideVendor), nil
}

// PropertyNames Get all property names available for this slide.
func (slide Slide) PropertyNames() []string {
	cPropNames := C.openslide_get_property_names(slide.ptr)
	var strings []string
	for i := 0; C.str_at(cPropNames, C.int(i)) != nil; i++ {
		strings = append(strings, C.GoString(C.str_at(cPropNames, C.int(i))))
	}
	return strings
}

// PropertyValue Get the value for a specific property.
func (slide Slide) PropertyValue(propName string) string {
	cPropName := C.CString(propName)
	defer C.free(unsafe.Pointer(cPropName))
	cPropValue := C.openslide_get_property_value(slide.ptr, cPropName)
	return C.GoString(cPropValue)
}

// Properties Get all properties as a map.
func (slide Slide) Properties() map[string]string {
	propertyNames := slide.PropertyNames()
	output := make(map[string]string)

	for _, propertyName := range propertyNames {
		if slide.PropertyValue(propertyName) != "" {
			output[propertyName] = slide.PropertyValue(propertyName)
		}
	}
	return output
}

// GetSpacing Get the spacing of the slide
func (slide Slide) GetSpacing() ([2]float64, error) {
	// TODO: For a TIFF different tags need to be read
	mppX := slide.PropertyValue(PropMPPX)
	mppY := slide.PropertyValue(PropMPPY)
	var output [2]float64

	if mppX == "" || mppY == "" {
		err := errors.New("mpp property not available")
		return output, err
	}
	mppXfloat, err0 := strconv.ParseFloat(mppX, 64)
	mppYfloat, err1 := strconv.ParseFloat(mppY, 64)

	if err0 != nil && err1 != nil {
		err := errors.New("cannot parse mpp values")
		return output, err
	}
	output = [2]float64{mppXfloat, mppYfloat}

	return output, nil
}

// GetThumbnail Get thumbnail of the image
func (slide Slide) GetThumbnail(size int) (image.Image, error) {
	var dimensions = slide.LargestLevelDimensions()
	var downsample float64 = 0
	for _, dim := range dimensions {
		currDownsample := float64(dim) / float64(size)
		if currDownsample >= downsample {
			downsample = currDownsample
		}
	}
	bestLevel := slide.BestLevelForDownsample(downsample)
	var thumbSize = slide.LevelDimensions(bestLevel)

	img, err := slide.ReadRegion(0, 0, bestLevel, thumbSize[0], thumbSize[1])

	// Compute the new size
	var outputSize [2]int
	var scaling float64
	if thumbSize[0] <= thumbSize[1] {
		outputSize[1] = size
		scaling = float64(thumbSize[1]) / float64(size)
		outputSize[0] = int(math.Floor(float64(thumbSize[0]) / scaling))
	} else {
		outputSize[0] = size
		scaling = float64(thumbSize[0]) / float64(size)
		outputSize[1] = int(math.Floor(float64(thumbSize[1]) / scaling))
	}
	outputImage := image.NewRGBA(image.Rect(0, 0, int(outputSize[0]), int(outputSize[1])))
	draw.BiLinear.Scale(outputImage, outputImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	return outputImage, err
}

// Version Get the current version of OpenSlide as a string
func Version() string {
	cVer := C.openslide_get_version()
	return C.GoString(cVer)
}

// PropBackgroundColor The name of the property containing a slide's background color, if any.
// It is represented as an RGB hex triplet.
const PropBackgroundColor = "openslide.background-color"

// PropBoundsHeight The name of the property containing the height of the rectangle bounding the non-empty region of the slide, if available.
const PropBoundsHeight = "openslide.bounds-height"

// PropBoundsWidth The name of the property containing the width of the rectangle bounding the non-empty region of the slide, if available.
const PropBoundsWidth = "openslide.bounds-width"

// PropBoundsX The name of the property containing the X coordinate of the rectangle bounding the non-empty region of the slide, if available.
const PropBoundsX = "openslide.bounds-x"

// PropBoundsY The name of the property containing the Y coordinate of the rectangle bounding the non-empty region of the slide, if available.
const PropBoundsY = "openslide.bounds-y"

// PropMPPX The name of the property containing the number of microns per pixel in the X dimension of level 0, if known.
const PropMPPX = "openslide.mpp-x"

// PropMPPY The name of the property containing the number of microns per pixel in the Y dimension of level 0, if known.
const PropMPPY = "openslide.mpp-y"

// PropObjectivePower The name of the property containing a slide's objective power, if known.
const PropObjectivePower = "openslide.objective-power"
