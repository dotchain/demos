// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

"use strict";

import { branch } from "https://dotchain.github.io/dotjs/streams/branch.js";

export class Items extends React.PureComponent {
  render() {
    const {onToggleAll, active, children} = this.props;
    const inputProps = {
      id: "toggle-all",
      className: "toggle-all",
      type: "checkbox",
      onChange: onToggleAll,
      checked: active === 0,
    };

    return React.createElement(
      "section",
      {className: "main"},
      React.createElement("input", inputProps),
      React.createElement("label", {htmlFor: "toggle-all"}),
      React.createElement("ul", {className: "todo-list"}, children)
    );
  }
}

export class Item extends React.Component {
  constructor(props) {
    super(props);

    this.state = {editing: null};
  }

  _onEdit() {
    const {constructor, value, stream} = this.props.todo;
    const editing = new constructor(value, branch(stream));
    this.setState({editing});
  }

  _onSubmit() {
    const {todo} =  this.props;
    const {editing} = this.state;

    if (editing) {
      const {value} = editing;
      const val = value.description.trim();
      if (val) {
        todo.description().replace(val);
        this.setState({editing: null});
      } else {
        todo.remove();
      }
    }
  }

  _onChange(e) {
    const {editing} = this.state;
    if(editing) {
      editing.description().replace(e.target.value);
      this.setState({editing: editing.latest()});
    }
  }

  _onKeyDown(e) {
    if (e.which == 27 /* escape */) {
      this.setState({editing: null});
    } else if (event.which === 13 /* enter */) {
      this._onSubmit();
    }
  }

  render() {
    const {todo} = this.props;
    let {done, description} = todo.value;
    const  classes = [];

    if (done) {
      classes.push("completed");
    }

    if (this.state.editing) {
      classes.push("editing");
      description = this.state.editing.value.description;
    }

    const cbProps = {
      className: "toggle",
      type: "checkbox",
      checked: done,
      onChange: e => todo.done().replace(e.target.checked)
    };

    const onEdit = () => this._onEdit();
    const onSubmit = () => this._onSubmit();
    const onChange = e => this._onChange(e);
    const onKeyDown = e => this._onKeyDown(e);
    const onDestroy = () => todo.remove();
    
    return React.createElement(
      "li",
      {className: classes.join(" ")},
      React.createElement(
        "div",
        {className: "view"},
        React.createElement("input", cbProps),
        React.createElement("label", {onDoubleClick: onEdit}, description),
        React.createElement("button", {className: "destroy", onClick: onDestroy})
      ),
      React.createElement(
        "input",
        {
          ref: "editField",
          className: "edit",
          value: description,
          onBlur: onSubmit,
          onChange: onChange,
          onKeyDown: onKeyDown
        }
      )
    );
  }
}
