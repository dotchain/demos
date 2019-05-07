// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

"use strict";

import {Todo} from "./types.js";
import {Footer} from "./footer.js";
import {Item, Items} from "./items.js";

const all = {href: "#/", message: "All"};
const active = {href: "#/active", message: "Active"};
const completed = {href: "#/completed", message: "Completed"};

export class App extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      editing: null,
      newTodo: "",
    };
  }

  onChange(e) {
    this.setState({newTodo: e.target.value});
  }

  onKeyDown(e) {
    if (e.keyCode === 13 /* enter key */) {
      e.preventDefault();
      const val = this.state.newTodo.trim();
      if (val) {
        this.props.todos.push(new Todo(false, val));
        this.setState({newTodo: ''});
      }
    }
  }

  onToggleAll(e) {
    // use e.target.checked
  }
  
  onClearCompleted() {
    // what does this do?
  }

  _count(filter) {
    let result = 0;
    for (let todo of this.props.todos.value) {
      result += filter(todo) ? 1 : 0;
    }
    return result;
  }

  _selected() {
    const routes = [all, active, completed];
    for (let r of routes) {
      if (r.href === this.props.hash) {
        return r;
      }
    }
    return all;
  }

  _filtered(fn) {
    const {todos, hash} = this.props;
    const result = [];
    for (let kk = 0; kk < todos.value.length; kk ++) {
      const {done} = todos.value[kk];
      if (done && hash === active.href || !done && hash === completed.href) {
        continue;
      }
      const onDestroy = kk => todos.splice(kk, 1);
      result.push(fn(todos.item(kk), onDestroy.bind(null, kk)));
    }
    return result;
  }
  
  render() {
    const {todos} = this.props;
    const {newTodo} = this.state;
    const inputProps = {
      className: "new-todo",
      placeholder: "What needs to be done?",
      value: newTodo,
      onKeyDown: e => this.onKeyDown(e),
      onChange: e => this.onChange(e),
      autoFocus: true,
    };

    return React.createElement(
      "div",
      null,
      React.createElement(
        "header",
        {className: "header"},
        React.createElement("h1", null, "todos"),
	React.createElement("input", inputProps),
      ),

      React.createElement(
        Items,
        {onToggleAll: () => this._onToggleAll(), active: this._count(todo => !todo.done)},
        ...this._filtered((todo, onDestroy) => React.createElement(Item, {key: todo, todo, onDestroy}))
      ),
      
      React.createElement(
        Footer,
        {
          active: this._count(todo => !todo.done),
          completed: this._count(todo => todo.done),
          routes: [all, active, completed],
          selected: this._selected(),
          onClearCompleted: () => this.onClearCompleted()
        }
      )
    );
  }
}

