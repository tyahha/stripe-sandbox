import React from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";

import "./App.css";

import { Home } from "./pages/Home";
import { Success } from "./pages/Success";
import { Cancel } from "./pages/Cancel";

function App() {
  return (
    <div className="App">
      <Router>
        <header>
          <Link to={"/"}>TOP</Link>
        </header>
        <Switch>
          <Route path="/success">
            <Success />
          </Route>
          <Route path="/cancel">
            <Cancel />
          </Route>
          <Route path="/">
            <Home />
          </Route>
        </Switch>
      </Router>
    </div>
  );
}

export default App;
