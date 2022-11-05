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

char * str_at(char ** p, int i);
void argb2rgba(uint32_t *buf, int len);