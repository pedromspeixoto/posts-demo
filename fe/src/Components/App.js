import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import '../sass/main.scss';

import Header from './Header';
import Home from "./Home";
import Posts from "./Posts";
import Sentry from './Sentry';

class App extends Component {
  render() {
    return (
      <Router>
        <div className="main">
          <Header/>
          
          <Switch>
            <div className="container mt-3">
              <Route exact path="/" component={Home} />
              <Route path="/posts" component={Posts} />
              <Route path="/sentry" component={Sentry} />
            </div> 
          </Switch>

        </div>
      </Router>
    )
  }
}

export default App;
