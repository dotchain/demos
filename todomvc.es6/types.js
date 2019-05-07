// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

"use strict";

import { Decoder, ValueStream, Transformer, Conn, Session } from "https://dotchain.github.io/dotjs/index.js";

import {
  StructDef,
  StructBase,
  ListDef,
  ListBase,
  Bool,
  Text,
  makeStreamClass
} from "https://dotchain.github.io/dotjs/types/index.js";

let todoDef = null;
let todoDefs = null;

export class Todo extends StructBase {
  constructor(done, description) {
    super();
    this.done = Boolean(done);
    this.description = "" + description;
  }

  static structDef() {
    if (todoDef == null) {
      todoDef = new StructDef("todomvc.Todo", Todo)
        .withField("done", "Done", Bool)
        .withField("description", "Description", Text);
    }
    return todoDef
  }

  static get Stream() {
    return TodoStream;
  }
}

export class Todos extends ListBase {
  static listDef() {
    return new ListDef("todomvc.Todos", Todos, Todo);
  }
  static get Stream() {
    return TodosStream;
  }
}

export const TodoStream = makeStreamClass(Todo);
export const TodosStream = makeStreamClass(Todos);

Decoder.registerValueClass(Todo);
Decoder.registerValueClass(Todos);

export function appSession(url) {
  const cache = null;
  const conn = new Transformer(new Conn(url, window.fetch), cache);
  const session = new Session().withLog({error: console.log});
  return {session, conn};
}
