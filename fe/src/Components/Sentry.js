import React from "react";

const sentryStyle = {
    cursor: 'pointer'
  };

class Sentry extends React.Component { 
    render() {
        return (
            <div>
                <div className="row justify-content-center mb-3">
                    <div className="col-md-4 mt-5">
                        <div style={sentryStyle} onClick={raiseError} className="card card-outline bg-secondary translate-up ripple ripple-dark">
                            <div className="card-body text-center">
                                <h5 className="font-weight-bold">Break the world</h5>
                            </div>
                        </div>
                    </div>        
                </div>
            </div>
        )
    }
}

export default Sentry;

function raiseError() {
    throw new Error("Hi there sentry!");
}