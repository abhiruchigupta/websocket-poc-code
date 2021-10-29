import React, {useState} from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";
import { Container } from '@material-ui/core';
import Login from './Login';
import AddListing from './AddListing';
import ViewListing from './ViewListing';
import './App.css';

const SERVER = 'http://ec2-3-237-18-207.compute-1.amazonaws.com:8080';

interface Data {
  email?: string,
  address?: string,
  city?: string,
  state?: string,
  zip?: string,
};

function App() {
  const [data, setData] = useState<Data>({});

  const handleLogin = (loginData: Data) => {
    setData({
      ...data,
      ...loginData,
    })
  }

  const handleAddListing = async (listingData: Data) => {
    const newData = {
      ...data,
      ...listingData,
    };

    setData(newData)

    try {
      const resp = await fetch(`${SERVER}/mls`, {
        method: 'POST',
        body: JSON.stringify(newData),
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json'
        },
      });
      const respJson = await resp.json();
      console.log('RESPONSE:', respJson);
    } catch (e) {
      console.error('FAILED REQUEST:', e)
    }
  }

  return (
    <Router>
      <Container>
        <Switch>
          <Route path="/login">
            <Login submit={handleLogin}/>
          </Route>
          <Route path="/add">
            <AddListing submit={handleAddListing}/>
          </Route>
          <Route path="/view">
            <ViewListing data={data}/>
          </Route>
          <Route path="*">
            <Redirect to="/login" />
          </Route>
        </Switch>
      </Container>
    </Router>
  );
}

export default App;
