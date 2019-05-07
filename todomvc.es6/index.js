// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

"use strict";

import {TodosStream, TodoStream, Todos, Todo, appSession} from "./types.js";
import {App} from "./app.js";

function main() {
  const {session, conn} = appSession("http://localhost:8080/dotjs/demo");
  let todos = new TodosStream(new Todos(), session.stream);

  sync();
  refresh();

  function refresh() {
    window.requestAnimationFrame(() => {
      todos = todos.latest();
      ReactDOM.render(
        React.createElement(App, {todos, hash: window.location.hash}, null),
        document.getElementById("root")
      );
  
      refresh();
    });
  }

  function sync(result) {
    setTimeout(push, 1000);
  }

  function push() {
    session.push(conn).then(pull, pull);
  }

  function pull() {
    session.pull(conn).then(sync, sync);
  }
}

main();
