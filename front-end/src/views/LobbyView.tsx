import React from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Checkbox from '@material-ui/core/Checkbox';

type LobbyViewProps = {
    ws: WebSocket,
    message: any
}

type LobbyViewState = {
    inputName: string,
}

class LobbyView extends React.Component<LobbyViewProps, LobbyViewState> {

    onReadyChange(event: React.ChangeEvent<HTMLInputElement>, checked: boolean) {
        console.warn(event.target.checked.toString());
        this.props.ws.send(JSON.stringify({
            'type': 'READY',
            'data': [event.target.checked.toString()]
        }))
    }

    render() {
        return (
            <div>
                <h1>Lobby</h1>
                <h2>Name: {this.props.message.selfInfo.name}</h2>
                <List >
                    {this.props.message.players.map((player: any, index: number) => {
                        if (player.ready) {
                            return (
                                <ListItem key={index} style={{ backgroundColor: 'green' }}>{player.name}</ListItem>
                            );
                        } else {
                            return (
                                <ListItem key={index} >{player.name}</ListItem>
                            );
                        }
                    })}
                </List>
                <Checkbox value={this.props.message.selfInfo.ready} onChange={this.onReadyChange.bind(this)} />
            </div>);
    }
}

export default LobbyView