"use strict"
// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

class CollabEditable {
    constructor(text, node) {
        this.node = node
        this.refresh(text, 0, 0)
        var done = (e) => this.setEditable(e);
        var refresh = (e) => this.setEditable(e);
        var url
        if (window.location.hostname == "localhost" && window.location.search != "?remote") {
            url = "http://localhost:5000/api/"
        } else {
            url = "https://etaqi6hpp8.execute-api.us-east-1.amazonaws.com/prod"
        }
        window.NewEditable(url, done, refresh)
    }

    setEditable(editable) {
        this.editable = editable
        var start = editable.Start(true);
        var end = editable.End(true);
        this.refresh(editable.Value(), start[0], end[0])
    }

    refresh(text, start, end) {
        console.log("Updating", text, start, end)
        this.text = text;
        this.start = start;
        this.end = end;

        this.node.innerText = text;
        var s = document.getSelection();
        var r = document.createRange();
        
        if (s.rangeCount > 0) {
            r = s.getRangeAt(0);
        }
        
        if (text == "") {
            r.setStart(this.node, 0);
            r.setEnd(this.node, 0);
        } else {
            r.setStart(this.node.firstChild, this.start);
            r.setEnd(this.node.firstChild, this.start);
        }

        if (s.rangeCount == 0) {
            s.addRange(r);
        }

        if (this.start !=  this.end) {
            s.extend(this.node.firstChild, this.end)
        }
    }
}

function init() {
    var node = document.querySelector(".editor");
    var editable = new CollabEditable("", node);

    node.addEventListener("keydown", function(e) {
        e.preventDefault();
    });
    node.addEventListener("keyup", function(e) {
        var code = e.code
        if (e.altKey) {
            code = "Alt" + e.code
        }
        if (e.shiftKey) {
            code = "Shift" + e.code
        }
        if (e.ctrlKey) {
            code = "Control" + e.code
        }
        if (e.code.slice(0, 3) == "Key" || e.code.slice(0, 5) == "Digit") {
            editable.editable.Insert(e.key);
        } else if (e.code == "Backspace") {
            editable.editable.Delete();
        } else if (e.code == "Space") {
            editable.editable.Insert(" ")
        } else if (editable.editable[code]) {
            editable.editable[code]()
        } else {
            console.log("Unknown code", e.code, code);
        }
    });
}

init();
