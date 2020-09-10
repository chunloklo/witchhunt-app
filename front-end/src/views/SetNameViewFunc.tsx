import React, { useState } from 'react';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';

type SetNameViewProps = {
    ws: WebSocket,
    message: any
}

function SetNameViewFunc(props: SetNameViewProps) {

    const [inputName, setInputName] = useState('')


    function onSubmit() {
        if (inputName !== '') {
            props.ws.send(JSON.stringify({
                "type": "SET_NAME",
                "data": [inputName]
            }))
        }
    }

    return (
        <div>
            <h1>Enter Your Name</h1>
            <TextField
                id="input"
                value={inputName}
                onChange={(event) => { setInputName(event.target.value) }}
            />
            <Button onClick={onSubmit}>Submit</Button>

            <h2>Existing Players</h2>
            <List >
                {props.message['players'].map((player: any, index: number) => {
                    return (
                        <ListItem key={index}>{player.name}</ListItem>
                    );
                })}
            </List>
        </div>)
}

export default SetNameViewFunc