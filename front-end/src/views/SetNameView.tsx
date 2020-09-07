import React from 'react';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';

type SetNameViewProps = {
    ws: WebSocket,
    message: any
}

type SetNameViewState = {
    inputName: string,
}

class SetNameView extends React.Component<SetNameViewProps, SetNameViewState> {
    constructor(props: any) {
        // Required step: always call the parent class' constructor
        super(props);

        // Set the state directly. Use props if necessary.
        this.state = {
            inputName: '',
        }
    }


    onSubmit() {
        if (this.state.inputName !== '') {
            this.props.ws.send(JSON.stringify({
                "type": "SET_NAME",
                "data": [this.state.inputName]
            }))
        }
    }

    render() {
        return (
            <div>
                <h1>Enter Your Name</h1>
                <TextField
                    id="input"
                    value={this.state.inputName}
                    onChange={(event) => { this.setState({ inputName: event.target.value }) }}
                />
                <Button onClick={this.onSubmit.bind(this)}>Submit</Button>

                <h2>Existing Players</h2>
                <List >
                    {this.props.message['players'].map((player: any, index: number) => {
                        return (
                            <ListItem key={index}>{player.name}</ListItem>
                        );
                    })}
                </List>
            </div>);
    }
}

export default SetNameView