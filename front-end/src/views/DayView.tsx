import React from 'react';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';

type DayViewProps = {
    ws: WebSocket,
    message: any
}

type DayViewState = {
    voteName: string,
}

class DayView extends React.Component<DayViewProps, DayViewState> {

    constructor(props: any) {
        super(props);
        this.state = {
            voteName: 'NO VOTE'
        }

    }

    handleChange(event: React.ChangeEvent<HTMLInputElement>, value: string) {
        this.props.ws.send(JSON.stringify({
            'type': 'DAY_VOTE',
            'data': [event.target.value],
        }))
    }

    render() {
        return (
            <div>
                <h1>Day {this.props.message.number}</h1>
                <h2>Vote for who to hang</h2>
                <RadioGroup name="vote" value={this.state.voteName} onChange={this.handleChange.bind(this)}>
                    {this.props.message.players.map((player: any, index: number) => {
                        if (player.alive) {
                            return (
                                <FormControlLabel key={player.name} value={player.name} control={<Radio />} label={player.name} />
                            )
                        }
                    })}
                </RadioGroup>
            </div>);
    }

}

export default DayView