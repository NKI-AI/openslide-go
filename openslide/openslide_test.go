/*
 * openslide-go - Unofficial Golang bindings for the OpenSlide library
 *
 * Copyright (c) 2020 GitHub user jammy-dodgers
 * https://github.com/jammy-dodgers/gophenslide
 * Copyright (c) 2022 Jonas Teuwen, Netherlands Cancer Institute
 *
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

import (
	"image/png"
	"os"
	"testing"
)

const testTiff = "testdata/CMU-1.tiff"

func TestDetectVendor(t *testing.T) {
	vendor, err := DetectVendor(testTiff)
	if err != nil {
		t.Error("Failed to load image: ", err.Error())
	} else if err == nil && vendor == "" {
		t.Error("Err nil but vendor blank")
	}
	t.Log("Vendor: ", vendor)
}

func TestOpen(t *testing.T) {
	slide, err := Open(testTiff)
	defer slide.Close()
	if err != nil {
		t.Error("Failed to load image: ", err.Error())
	}
}

func TestLevels(t *testing.T) {
	slide, err := Open(testTiff)
	defer slide.Close()
	if err != nil {
		t.Error("Failed to load image: ", err.Error())
	}
	levels := slide.LevelCount()
	if levels == -1 {
		t.Error("Cannot parse levels")
	}
	largestDimensions := slide.LargestLevelDimensions()
	t.Log("Base lvl0 (", largestDimensions[0], ", ", largestDimensions[1], "): ", slide.LevelDownsample(0))
	for i := 1; i < levels; i++ {
		levelDimensions := slide.LevelDimensions(i)
		t.Log("Level ", i, " (", levelDimensions[0], ", ", levelDimensions[1], "): ", slide.LevelDownsample(i))
	}
}

func TestReadRegion(t *testing.T) {
	slide, err := Open(testTiff)
	defer slide.Close()
	if err != nil {
		t.Error("Failed to load image: ", err.Error())
	}
	region, err := slide.ReadRegion(10, 10, 6, 400, 400)
	if err != nil {
		t.Fatal(err.Error())
	}
	const testRawFilename = "testdata/region.png"
	if info, e := os.Stat(testRawFilename); os.IsExist(e) && !info.IsDir() {
		if remErr := os.Remove(testRawFilename); remErr != nil {
			t.Log("Could not remove file ", testRawFilename)
		}
	}

	f, _ := os.Create(testRawFilename)
	encodeErr := png.Encode(f, region)

	if encodeErr != nil {
		t.Fatal(encodeErr.Error())
	}
}

func TestProperties(t *testing.T) {
	slide, err := Open(testTiff)
	defer slide.Close()
	if err != nil {
		t.Error("Failed to load image: ", err.Error())
	}
	props := slide.PropertyNames()
	for i := 0; i < len(props); i++ {
		t.Log(props[i], "=", slide.PropertyValue(props[i]))
	}
}

func TestVersion(t *testing.T) {
	t.Log("Version", Version())
}
