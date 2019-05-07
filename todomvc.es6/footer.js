// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

"use strict";

export class Footer extends React.PureComponent {
  render() {
    const {active, completed, routes, selected, onClearCompleted} = this.props;

    if (active === 0 && completed === 0) {
      return null;
    }

    return React.createElement(
      "footer",
      {className: "footer"},
      React.createElement(
        "span",
        {className: "todo-count"},
        React.createElement("strong", null, active.toString()),
        " " + (active > 1 ? "items" : "item") + " left"
      ),

      React.createElement(
        "ul",
        {className: "filters"},
        ...routes.map(route => React.createElement(Route, {route, selected}))
      ),
      React.createElement(ClearButton, {completed, onClearCompleted})
    );
  }
}

class Route extends React.PureComponent {
  render() {
    const {route, selected} = this.props;
    const className = (selected === route) ? "selected" : "";
    const {href, message} = route;
    
    return React.createElement(
      "li",
      null,
      React.createElement("a", {href, className}, message)
    );
  }
}

class ClearButton extends React.PureComponent {
  render() {
    const {completed, onClearCompleted} = this.props;
    if (completed === 0) {
      return null;
    }

    const props = {className: "clear-completed", onClick: onClearCompleted};
    return React.createElement("button", props, "Clear completed");
  }
}
