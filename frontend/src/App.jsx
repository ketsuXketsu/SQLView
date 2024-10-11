import {useEffect, useReducer, useState} from 'react';
import './App.css';
import {Query} from '../wailsjs/go/main/App'
import {FileDialog} from '../wailsjs/go/main/App'

function App() {

    const query = async (_query) => {
        localStorage.clear()
        console.log(_query.currQuery)
        await QueryBackend(_query.currQuery);
        setState(state + 1);
    };

    const [state, setState] = useState(0);
    const [currQuery, setCurrQuery] = useState("");

    return (
        <div id="App">
            <div className='navbar'>
                <p>SView</p>
            </div>
            <div className='mainContainer'>
                <div className='sidebarL'>
                    <button onClick={() => query({currQuery})}>Execute Query</button>
                    <button onClick={FileDialog}>Open DB File</button>
                    <button>placeholder</button>
                    <button>placeholder</button>
                    <button>placeholder</button>
                </div>
                <div className='queryPanel'> 
                    <textarea className='queryArea' spellCheck="false" onChange={(e) => setCurrQuery(e.target.value)}></textarea>
                </div>
                <div className='infoViewPanel'>
                    <DebugItem /> 
                </div>
            </div>
        </div>
    );
};

async function QueryBackend(_query) {
    await Query(_query).then((res) => {
        localStorage.setItem('debugItems', res)
    });
};

function DebugItem() {
    if (localStorage.getItem('debugItems') != null && localStorage.getItem('debugItems') != "") {
        try {
            var values = JSON.parse(localStorage.getItem('debugItems') || `"error":"empty"`)
            var valueList = values.map((val) => {
                return (
                    <div className='sqlObject'>
                        <p>{val.fields}</p>
                    </div>
                )
            })
            return (
                <div>{valueList}</div>
            ); 
        }
        catch {
            return (
                <div className='sqlObject' style={{backgroundColor : 'orange'}}>
                    <p>Error displaying info</p>
                </div>
            )
        }
    } else {
        return (
            <div className='sqlObject' style={{backgroundColor : '#2a475e'}}>
            <p>"Waiting..."</p>
            </div>
        )
    }
};


export default App


