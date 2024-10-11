import {useEffect, useReducer, useState} from 'react';
import './App.css';
import {Query} from '../wailsjs/go/main/App'
import {FileDialog} from '../wailsjs/go/main/App'
import {DatabaseButtonClicked} from '../wailsjs/go/main/App';

function App() {
    const query = async (_query) => {
        localStorage.removeItem('sqlItems')
        await QueryBackend(_query.currQuery);
        setState(state + 1);
    };
    const changeDb = async () => {
        await FrontFileDialog();
        setState(state + 1);
    }
    const clearRecentFiles = async () => {
        localStorage.removeItem('recentFiles');
        setState(state + 1);
    }
    const clearItems = async () => {
        localStorage.removeItem('sqlItems');
        setState(state + 1);
    }

    const [state, setState] = useState(0);
    const [currQuery, setCurrQuery] = useState("");
    const [infoViewFontSize, setInfoViewFontSize] = useState('20') // Not working 

    return (
        <div id="App">
            <div className='navbar'>
                <p>SView</p>
            </div>
            <div className='mainContainer'>
                <div className='sidebarL'>
                    <div className='sidebarLActions'>
                        <button onClick={() => query({currQuery})}>Execute Query</button>
                        <button onClick={changeDb}>Open DB File</button>
                    </div>
                    <div>
                        <button onClick={() => clearRecentFiles()}>Clear</button>
                        <div className='fileExplorer'>
                            <RecentFiles />
                        </div>
                    </div>
                </div>
                <div className='queryPanel'> 
                    <textarea className='queryArea' spellCheck="false" onChange={(e) => setCurrQuery(e.target.value)}></textarea>
                </div>
                <div className='infoViewPanel' style={{fontSize : `${infoViewFontSize}px`}}>
                    <div className='toolbox'>
                        <button className='clearBtn' onClick={() => clearItems()}>Clear Items</button>
                        <div>
                            <p>Filter results: </p>
                            <input></input>
                        </div>
                    </div>
                    <SQLHeaders />
                    <SQLItems /> 
                </div>
            </div>
        </div>
    );
};

async function DbButton(button) {
    DatabaseButtonClicked(button)
}

async function FrontFileDialog() {
    await FileDialog().then((res) => {
        if(res != "") {
            const items = JSON.parse(localStorage.getItem('recentFiles') || '[]')
            const item = JSON.stringify([...items, {name: res}])
            localStorage.setItem('recentFiles', item)
        }
        else {
            return
        }
    })
}

function RecentFiles() {
    try {
        const files = JSON.parse(localStorage.getItem('recentFiles') || '[]')
        const fileList = files.map((file) => {
            return (
                <button value={file.name} onClick={(e) => changeDb(e.target.value)}>{file.name}</button>
            )
        })
        return (
            <div>
                {fileList}
            </div>
        )
    }
    catch {
        return null
    }
}

async function QueryBackend(_query) { // Sends a query to the middleware, who then sends it to the backend, awaits a response and sends it back here
    await Query(_query).then((res) => {
        localStorage.setItem('sqlItems', res)
    });
};

function SQLHeaders() {
    try {
        const items = JSON.parse(localStorage.getItem('sqlItems') || '[]')
        if(items.length == 0) {
            return null
        }

        const hdrs = items[0].headers
        let hdrList = hdrs.map((hdr) => 
        {
            return (
                <p>{hdr} | </p>
            )
        })
        return (
            <div className='sqlHeader'>
                {hdrList}
            </div>
        )
    }
    catch {
        return null
    } 
}

function SQLItems() {
    var colorMemory = -1
    function orderColor() { // Gradient effect on the info panel items 
        const colors = ['#a4eb81', '#4a8574', '#32585c']
        colorMemory++
        if(colorMemory > colors.length-1) {
            colorMemory = 0
            return colors[colorMemory]
        }
        return colors[colorMemory]
    }
    
    if (localStorage.getItem('sqlItems') != null && localStorage.getItem('sqlItems') != "") {
        try {
            var values = JSON.parse(localStorage.getItem('sqlItems') || `"error":"empty"`)
            var valueList = values.map((val) => 
                {
                const style = {borderLeft: `4px solid ${orderColor()}`}
                var valList = val.fields.map((v) => 
                {
                    return (
                        <p>{v}</p>
                    )
                })
                return (
                    <div className='sqlObject' style={style}>
                        {valList}
                    </div>
                )
            })
            return (
                <div>{valueList}</div>
            ); 
        }
        catch 
        {
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