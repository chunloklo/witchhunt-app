import React from 'react';
import SetNameView from './views/SetNameView';
import LobbyView from './views/LobbyView';
import ApprenticeStartView from './views/ApprenticeStartView';
import NightView from './views/NightView';
import DayView from './views/DayView';
import SetNameViewFunc from './views/SetNameViewFunc';

type MainViewProps = {
}

type MainViewState = {
    page: String,
    serverMessage: any,
}

class MainView extends React.Component<MainViewProps, MainViewState> {

    constructor(props: any) {
        super(props);
        this.state = {
            page: '',
            serverMessage: {}
        }
    }


    ws = new WebSocket('ws://localhost:8080/ws')

    componentDidMount() {
        this.ws.onopen = () => {
            // on connecting, do nothing but log it to the console
            console.log('connected')
        }

        this.ws.onmessage = evt => {
            // listen to data sent from the websocket server
            // console.log(evt.data)
            const message = JSON.parse(evt.data)
            this.setState({ serverMessage: message })
            console.log(message)
        }

        this.ws.onclose = () => {
            console.log('disconnected')
            // automatically try to reconnect on connection loss
        }
    }

    render() {
        if (this.state.serverMessage['action'] === 'SET_NAME') {
            return <SetNameViewFunc ws={this.ws} message={this.state.serverMessage} />
        }

        if (this.state.serverMessage['action'] === 'LOBBY') {
            return <LobbyView ws={this.ws} message={this.state.serverMessage} />
        }

        if (this.state.serverMessage['action'] === 'APPRENTICE_START') {
            return <ApprenticeStartView ws={this.ws} message={this.state.serverMessage} />
        }

        if (this.state.serverMessage['action'] === 'NIGHT') {
            return <NightView ws={this.ws} message={this.state.serverMessage} />
        }

        if (this.state.serverMessage['action'] === 'DAY') {
            return <DayView ws={this.ws} message={this.state.serverMessage} />
        }
        return <h1>Hello</h1>;
    }
}

export default MainView