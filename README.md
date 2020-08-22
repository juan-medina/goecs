# goecs
A simple Go [ECS](https://en.wikipedia.org/wiki/Entity_component_system)

[![License: Apache2](https://img.shields.io/badge/license-Apache%202-blue.svg)](/LICENSE)
[![go version](https://img.shields.io/github/go-mod/go-version/juan-medina/goecs)](https://pkg.go.dev/mod/github.com/juan-medina/goecs)
[![godoc](https://godoc.org/github.com/juan-medina/goecs?status.svg)](https://pkg.go.dev/mod/github.com/juan-medina/goecs)
[![Build Status](https://travis-ci.com/juan-medina/goecs.svg?branch=main)](https://travis-ci.com/juan-medina/goecs)
[![codecov](https://codecov.io/gh/juan-medina/goecs/branch/main/graph/badge.svg)](https://codecov.io/gh/juan-medina/goecs)

## Info
Entity–component–system (ECS) is an architectural patter that follows the composition over inheritance principle that allows greater flexibility in defining entities where every object in a world.

Every entity consists of one or more components which contains data or state. Therefore, the behavior of an entity can be changed at runtime by systems that add, remove or mutate components.

This eliminates the ambiguity problems of deep and wide inheritance hierarchies that are difficult to understand, maintain and extend.

Common ECS approaches are highly compatible and often combined with data-oriented design techniques.

For a more in deep read on this topic I could recommend this [article](https://medium.com/ingeniouslysimple/entities-components-and-systems-89c31464240d).

## License

```text
    Copyright (C) 2020 Juan Medina

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
```
