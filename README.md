# README

## About

SQLView is a very simple SQLite Interface, it has almost nothing more than a simple CLI would have, it was made simply for visualization. It is very lightweight (13mb, single binary)
It is written in Go using the Wails Framework, with React-JS as the frontend. 
This is purely a student project, don't expect to use it for anything serious, as breaking the program is very buggy right now.

## Usage
1. Open sqlview.exe
2. Click 'Open DB File'
3. Select a SQLite .db file
4. Run a query

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
