import React from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Button from '@material-ui/core/Button';

type ApprenticeStartViewProps = {
    ws: WebSocket,
    message: any
}

type ApprenticeStartViewState = {
    inputName: string,
}

class ApprenticeStartView extends React.Component<ApprenticeStartViewProps, ApprenticeStartViewState> {

    roleSelect(role: string) {
        // console.warn(event.target.checked.toString());
        this.props.ws.send(JSON.stringify({
            'type': 'APPRENTICE_START_ROLE_SELECT',
            'data': [role],
        }))
    }

    finishedReading() {
        this.props.ws.send(JSON.stringify({
            'type': 'APPRENTICE_START_ROLE_SELECT',
            'data': [this.props.message.selectedRole, 'TRUE'],
        }))
    }


    render() {
        if (this.props.message.selectedRole === undefined) {
            return (
                <div>
                    <h1>Choose A Role that you want to Apprentice</h1>
                    <List >
                        <ListItem
                            id='judge'
                        >
                            <Button onClick={() => this.roleSelect('JUDGE')}>Judge</Button>
                        </ListItem>
                        <ListItem
                            id='gravedigger'
                        >
                            <Button onClick={() => this.roleSelect('GRAVEDIGGER')}>Gravedigger</Button>
                        </ListItem>
                    </List>
                </div>);
        } else {
            return (
                <div>
                    <h1>Apprentice Info</h1>
                    <div>
                        {this.props.message.selectedRole}
                    </div>
                    <div>
                        {this.props.message.rolePlayer.name}
                    </div>
                    <Button onClick={this.finishedReading.bind(this)}>Done Reading</Button>
                </div>);
        }

    }
}

export default ApprenticeStartView