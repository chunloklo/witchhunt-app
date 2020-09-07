import React from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Button from '@material-ui/core/Button';

type NightViewProps = {
    ws: WebSocket,
    message: any
}

type NightViewState = {
    inputName: string,
}

class NightView extends React.Component<NightViewProps, NightViewState> {

    roleSelect(role: string) {
        // console.warn(event.target.checked.toString());
        this.props.ws.send(JSON.stringify({
            'type': 'APPRENTICE_START_ROLE_SELECT',
            'data': [role],
        }))
    }


    render() {
        return (
            <div>
                <h1>Night {this.props.message.number}</h1>
            </div>);
    }

}

export default NightView