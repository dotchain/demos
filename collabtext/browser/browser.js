"use strict"
// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

class CollabEditable {
    constructor(text, node) {
        this.node = node
        this.refresh(text, 0, 0)
        var done = (e) => this.setEditable(e);
        var refresh = (text, start, end) => this.refresh(text, start, end);
        window.NewEditable("http://localhost:5000/api/", done, refresh)
    }
    setEditable(editable) {
        this.editable = editable
    }
    Insert(ch) {
        this.editable.Insert(ch)
    }
    Delete() {
        this.editable.Delete()
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
            r.setEnd(this.node.firstChild, this.end);
        }

        if (s.rangeCount == 0) {
            s.addRange(r);
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
        if (e.code.slice(0, 3) == "Key" || e.code.slice(0, 5) == "Digit") {
            editable.Insert(e.key);
        } else if (e.code == "Backspace") {
            editable.Delete();
        } else {
            console.log("Unknown code", e.code);
        }
    });
}

init();
