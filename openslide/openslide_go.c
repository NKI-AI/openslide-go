/*
 * openslide-go - Unofficial Golang bindings for the OpenSlide library
 *
 * Copyright (c) 2015 Carnegie Mellon University
 * Copyright (c) 2020 GitHub user jammy-dodgers
 * https://github.com/jammy-dodgers/gophenslide
 * Copyright (c) 2022 Jonas Teuwen, Netherlands Cancer Institute
 *
 * This file has been derived from the openslide-python bindings, under the same license
 * https://github.com/openslide/openslide-python/blob/3215e1ba96641bada0ec859b2632ba0d5f5b7168/openslide/_convert.c
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

#include <stdio.h>
#include <stdlib.h>
#include <openslide_go.h>

char * str_at(char ** p, int i) { return p[i]; }

typedef unsigned char u8;

void argb2rgba(uint32_t *buf, int len) {
    ssize_t cur;

    for (cur = 0; cur < len; cur++) {
        uint32_t val = buf[cur];
        u8 a = val >> 24;
        switch (a) {
        case 0:
            break;
        case 255:
            val = (val << 8) | a;
#ifndef WORDS_BIGENDIAN
            // compiler should optimize this to bswap
            val = (((val & 0x000000ff) << 24) |
                   ((val & 0x0000ff00) <<  8) |
                   ((val & 0x00ff0000) >>  8) |
                   ((val & 0xff000000) >> 24));
#endif
            buf[cur] = val;
            break;
        default:
            ; // label cannot point to a variable declaration
            u8 r = 255 * ((val >> 16) & 0xff) / a;
            u8 g = 255 * ((val >>  8) & 0xff) / a;
            u8 b = 255 * ((val >>  0) & 0xff) / a;
#ifdef WORDS_BIGENDIAN
            val = r << 24 | g << 16 | b << 8 | a;
#else
            val = a << 24 | b << 16 | g << 8 | r;
#endif
            buf[cur] = val;
            break;
        }
    }
}